import type { CreateRoleRequest, CreateUserRequest } from "@/types/requests/user-management";

export const userManagementAdminApi = () => {
  const api = useApiHost();
  return {
    getUser: async (userId: string) => {
      const response = await fetch(`${api}/api/admin/user-management/${userId}`, {
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

    getAllUsers: async ({ query = { role: "" } }: { query?: { role?: string } }) => {
      let querySearch =""
      if(query.role != ""){
        querySearch = `?role=${query.role}`
      }

      const response = await fetch(`${api}/api/admin/user-management${querySearch}`, {
        method: "GET",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${useCookie("token").value}`,
        }
      });
      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.message || 'Login failed');
      }
      return response.json();
    },

    createUser: async (
      data: CreateUserRequest = {
        name: "",
        email: "",
        logo_url: "",
        role_id: "",
        password: "",
        password_confirm: "",
        phone: ""
      }
    ) => {
      const response = await fetch(`${api}/api/admin/user-management`, {
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

    editUser: async (
      userId: string,
      data: CreateUserRequest = {
        name: "",
        email: "",
        logo_url: "",
        role_id: "",
        password: "",
        password_confirm: "",
        phone: ""
      }
    ) => {
      const response = await fetch(`${api}/api/admin/user-management/${userId}`, {
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

    deleteUser: async (userId: string) => {
      const response = await fetch(`${api}/api/admin/user-management/${userId}`, {
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
    getRole: async (roleId: string) => {
      const response = await fetch(`${api}/api/admin/role/${roleId}`, {
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

    getAllRole: async () => {
      const response = await fetch(`${api}/api/admin/role`, {
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


    createRole: async (
      data: CreateRoleRequest = {
        name: "",
      }
    ) => {
      const response = await fetch(`${api}/api/admin/role`, {
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

    editRole: async (
      roleId: string,
      data: CreateRoleRequest = {
        name: "",
      }
    ) => {
      const response = await fetch(`${api}/api/admin/role/${roleId}`, {
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

    deleteRole: async (roleId: string) => {
      const response = await fetch(`${api}/api/admin/role/${roleId}`, {
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