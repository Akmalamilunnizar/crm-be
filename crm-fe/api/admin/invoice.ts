import type { CreateInvoiceRequest, UpdateStatusInvoiceRequest } from "@/types/requests/invoice";


export const invoiceAdminApi = () => {
  const api = useApiHost();
  return {
    getInvoice: async (invoiceId: string) => {
      const response = await fetch(`${api}/api/admin/invoice/${invoiceId}`, {
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

    getAllInvoices: async () => {
      const response = await fetch(`${api}/api/admin/invoice`, {
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

    createInvoice: async (
      data: CreateInvoiceRequest
    ) => {
      const response = await fetch(`${api}/api/admin/invoice`, {
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

    editInvoice: async (
      invoiceId: string,
      data: CreateInvoiceRequest
    ) => {
      const response = await fetch(`${api}/api/admin/invoice/${invoiceId}`, {
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

    deleteInvoice: async (invoiceId: string) => {
      const response = await fetch(`${api}/api/admin/invoice/${invoiceId}`, {
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

    updateStatusInvoice: async (
      invoiceId: string,
      data: UpdateStatusInvoiceRequest
    ) => {
      const response = await fetch(`${api}/api/admin/invoice/${invoiceId}/status`, {
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
  };
};
