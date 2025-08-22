import type { CreateAreaRequest } from "@/types/requests/area";
import type { CreateCompanyRequest } from "@/types/requests/company";

export const areaAdminApi = () => {
  const api = useApiHost();
  return {
    getArea: async (companyId: string) => {
      const response = await fetch(`${api}/api/admin/area/${companyId}`, {
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

    getAllAreas: async () => {
      const response = await fetch(`${api}/api/admin/area`, {
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

    createArea: async (
      data: CreateAreaRequest
    ) => {
      const response = await fetch(`${api}/api/admin/area`, {
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

    editArea: async (
      areaId: string,
      data: CreateAreaRequest
    ) => {
      const response = await fetch(`${api}/api/admin/area/${areaId}`, {
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

    deleteArea: async (areaId: string) => {
      const response = await fetch(`${api}/api/admin/area/${areaId}`, {
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
