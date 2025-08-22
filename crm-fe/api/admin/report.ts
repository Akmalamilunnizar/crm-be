import type { CreateAreaRequest } from "@/types/requests/area";
import type { CreateCompanyRequest } from "@/types/requests/company";

export const reportAdminApi = () => {
  const api = useApiHost();
  return {
    getReportInternet: async () => {
      const response = await fetch(`${api}/api/admin/report/internet/`, {
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
