import type { CreateCompanyRequest} from "@/types/requests/company";
import type { CreateTransactionRequest } from "@/types/requests/transaction";

export const transactionAdminApi = () => {
  const api = useApiHost();
  return {
    getTransaction: async (transactionId: string) => {
      const response = await fetch(`${api}/api/admin/transaction/${transactionId}`, {
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

    getAllTransactions: async (params:any) => {
      const queryString = new URLSearchParams(params).toString();
      const response = await fetch(`${api}/api/admin/transaction?${queryString}`, {
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

    createTransactions: async (
      data: CreateTransactionRequest = {
        type_cash: "",
        type_in_out: "",
        account_id: "",
        date: "",
        description: "",
        file: "",
        amount: "",
        category: "",
        tags: "",
        payer_id: "",
        method: "",
        ref: "",
      }
    ) => {
      const response = await fetch(`${api}/api/admin/transaction`, {
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

    editTransaction: async (
      transactionId: string,
      data: CreateCompanyRequest = {
        name: "",
        url: "",
        email: "",
        phone: "",
        logo_url: "",
        description: "",
    npwp: "",
    address: "",
      }
    ) => {
      const response = await fetch(`${api}/api/admin/transaction/${transactionId}`, {
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

    deleteTransaction: async (transactionId: string) => {
      const response = await fetch(`${api}/api/admin/transaction/${transactionId}`, {
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
