export type CreateInvoiceRequest = {
    customer_id: string,
    amount: number,
    invoice_items: Array<{name: string, qty: number, price: number,total: number}>
  }

  export type UpdateStatusInvoiceRequest = {
    status: string
  }