import type { CreateAssetRequest } from "@/types/requests/asset";
import { defaultStringValidation } from "@/utilities/constants";
import { object, string, type InferType } from "yup";

export const assetSchema = object({

});

export type AssetSchema = InferType<typeof assetSchema>;

export const asset :CreateAssetRequest= {
  type: "",
  brand: "",
  model: "",
  serial_number: "",
  mac_address: "",
  date: new Date(),
  site: "",
  quantity: 0,
  status: "",
  price: 0,
  description: "",
  status_in_out: "",
};

export const status_in_out = [
  { label: "IN", value: "in" },
  { label: "OUT", value: "out" },
];

export const status = [
  { label: "Broken", value: "broken" },
  { label: "Good", value: "good" },
  { label: "Bad", value: "bad" },
];

export const defaultAssets = ref<any[]>([])
