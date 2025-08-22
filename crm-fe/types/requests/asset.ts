export type CreateAssetRequest = {
  type: string,
  brand: string,
  model: string,
  serial_number: string,
  mac_address: string,
  date: Date,
  site: string,
  quantity: number,
  status: string,
  price: number,
  description: string,
  status_in_out: string,
  }