import type { CreateCompanyRequest } from "@/types/requests/company";
import type { CreateInternetPackageRequest } from "@/types/requests/internet-package";

export const internetPackageAdminApi = () => {
  const api = useApiHost();
  return {
    getInternetPacket: async (productId: string) => {
      const response = await fetch(`${api}/api/admin/product/${productId}`, {
        method: "GET",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${useCookie("token").value}`,
        },
      });
      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.message || 'Login failed');
      }
      return response.json();
    },

    getAllInternetPacket: async () => {
      const response = await fetch(`${api}/api/admin/product`, {
        method: "GET",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${useCookie("token").value}`,
        },
      });
      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.message || 'Login failed');
      }
      return response.json();
    },

    createInternetPackage: async (
      data: CreateInternetPackageRequest
    ) => {
      const response = await fetch(`${api}/api/admin/product`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${useCookie("token").value}`,
        },
        body: JSON.stringify(data),
      });
      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.message || 'Login failed');
      }
      return response.json();
    },

    editInternetPacket: async (
      productId: string,
      data: CreateInternetPackageRequest
    ) => {
      const response = await fetch(`${api}/api/admin/product/${productId}`, {
        method: "PUT",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${useCookie("token").value}`,
        },
        body: JSON.stringify(data),
      });
      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.message || 'Login failed');
      }
      return response.json();
    },

    deleteInternetPacket: async (productId: string) => {
      const response = await fetch(`${api}/api/admin/product/${productId}`, {
        method: "DELETE",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${useCookie("token").value}`,
        },
      });
      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.message || 'Login failed');
      }
      return response.json();
    },
  };
};
