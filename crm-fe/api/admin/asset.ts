import type { CreateAreaRequest } from "@/types/requests/area";
import type { CreateAssetRequest } from "@/types/requests/asset";

export const assetAdminApi = () => {
  const api = useApiHost();
  return {
    getAsset: async (companyId: string) => {
      const response = await fetch(`${api}/api/admin/asset/${companyId}`, {
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

    getAllAssets: async () => {
      const response = await fetch(`${api}/api/admin/asset`, {
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

    createAsset: async (
      data: CreateAssetRequest
    ) => {
      const response = await fetch(`${api}/api/admin/asset`, {
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

    editAsset: async (
      assetId: string,
      data: CreateAssetRequest
    ) => {
      const response = await fetch(`${api}/api/admin/asset/${assetId}`, {
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

    deleteAsset: async (assetId: string) => {
      const response = await fetch(`${api}/api/admin/asset/${assetId}`, {
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
