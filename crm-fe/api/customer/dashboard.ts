

export const dashboardCustomerApi = () => {
  const api = useApiHost();
  return {
    getCustomerDashboard: async () => {
      const response = await fetch(`${api}/api/customer/dashboard`, {
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
     createPaymentDashboard: async (
      data: {invoice_id: string}
    ) => {
      const response = await fetch(`${api}/api/customer/dashboard`, {
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
  };
};
