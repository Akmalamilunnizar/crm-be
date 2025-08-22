import type { CreateInvoiceRequest, UpdateStatusInvoiceRequest } from "@/types/requests/invoice";


export const WhatsappApi = () => {
  const wa = useWaHost();
  return {
    sendWhatsapp: async (
      data: {
        number:string,
        message:string
      }
    ) => {
      const response = await fetch(`${wa}/send-message`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
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
