import type { CreateAreaRequest } from "@/types/requests/area";
import type { CreateCompanyRequest } from "@/types/requests/company";

export const archiveInstallationAdminApi = () => {
  const api = useApiHost();
  return {
    getArchiveInstallation: async (id: string) => {
      const response = await fetch(`${api}/api/admin/customer-installation/${id}`, {
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
    getAllArchiveInstallation: async () => {
      const response = await fetch(`${api}/api/admin/customer-installation`, {
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
  };
};
