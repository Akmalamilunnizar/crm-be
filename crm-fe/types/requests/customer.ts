export type CreateCustomerRequest = {
  type_of_service: String
  email: String
  name: String
  company_id: String
  gender: String
  card_identition: String
  no_identition: Number
  area_id: String
  phone: String
  address: String
  latitude: Number
  longitude: Number
  password: String
  product_id: String
  ip_static: String
  mac_address: String
  job: String
}

export type CreateCustomerInstallationRequest = {
  customer_id: String,
  technician_id: String,
  date: String,
  description:  String,
  image_ids: string[] ,
  previews: String[] ,
  selectedImage: String,
  showModal: boolean,
}