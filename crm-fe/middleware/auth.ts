import { authApi } from "@/api/auth";
import { useAuthStore } from "../stores/auth";
import type { RouteLocationNormalizedGeneric } from "vue-router";

export default defineNuxtRouteMiddleware((to, from) => {
  const authStore = useAuthStore();
  if (to.path.startsWith('/dashboard') ) {
      checkAuth();
  }

  if (to.path.startsWith('/customer')) {
      checkAuthCustomer();
  }

  if (to.path.startsWith('/dashboard')) {
    const check = checkPermission(to);
    if (check == "disallow") {
      return navigateTo("/dashboard");
    }
    if (!authStore.isLoggedIn && to.path !== "/login") {
      authStore.logout();
      return navigateTo("/login");
    }
  }
});

async function checkAuth() {
  const authStore = useAuthStore();
  
  // Check if we have a token first
  if (!authStore.isLoggedIn) {
    authStore.logout();
    return;
  }
  
  try {
    const response = await authApi().verifyAuth();
    if (response.success == false) {
      authStore.logout();
    } else {
      authStore.user = response.data;
    }
  } catch (error) {
    console.error("authStore cek auth error", error);
    authStore.logout();
  }
}

async function checkAuthCustomer() {
  const authStore = useAuthStore();
  
  // Check if we have a token first
  if (!authStore.isLoggedIn) {
    authStore.logout();
    return;
  }
  
  try {
    const response = await authApi().verifyAuthCustomer();
    if (response.success == false) {
      authStore.logout();
    } else {
      authStore.user = response.data;
    }
  } catch (error) {
    console.error("authStore cek auth error", error);
    authStore.logout();
  }
}

function checkPermission(to: RouteLocationNormalizedGeneric) {
  const authStore = useAuthStore();

  const restrictedForAdmins = [""];
  const restrictedForTechnicians = [
    "/dashboard/user-management",
    "/dashboard/report",
    "/dashboard/invoice",
    "/dashboard/transaction",
    "/dashboard/companies",
  ];
  const restrictedForFinances = [
    "/dashboard/user-management",
    "/dashboard/companies",
    "/dashboard/asset",
    "/dashboard/area",
    "/dashboard/internet-package",
  ];
  if (
    authStore.user.role === "ADMIN" &&
    restrictedForAdmins.includes(to.path)
  ) {
    return "disallow";
  }
  if (
    authStore.user.role === "TECHNICIAN" &&
    restrictedForTechnicians.includes(to.path)
  ) {
    return "disallow";
  }
  if (
    authStore.user.role === "FINANCE" &&
    restrictedForFinances.includes(to.path)
  ) {
    return "disallow";
  }
}
