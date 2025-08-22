<script setup lang="ts">
import { object, string, type InferType } from "yup";
import type { FormSubmitEvent } from "#ui/types";
import { customerAdminApi } from "@/api/admin/customer";
import { areaAdminApi } from "@/api/admin/area";
import { companyAdminApi } from "@/api/admin/company";
import { internetPackageAdminApi } from "@/api/admin/internet-package";

const props = defineProps({
  isEdit: {
    type: Boolean,
    required: false
  },
  data: {
    type: Object,
    default: () => ({
      id: {
        type: String,
        default: "",
      },
      type_of_service: {
        type: String,
        default: ""
      },
      email: {
        type: String,
        default: ""
      },
      name: {
        type: String,
        default: ""
      },
      company_id: {
        type: String,
        default: ""
      },
      gender: {
        type: String,
        default: ""
      },
      card_identition: {
        type: String,
        default: ""
      },
      no_identition: {
        type: Number,
        default: 0
      },
      area_id: {
        type: String,
        default: ""
      },
      phone: {
        type: String,
        default: ""
      },
      address: {
        type: String,
        default: ""
      },
      latitude: {
        type: Number,
        default: 0,
      },
      longitude: {
        type: Number,
        default: 0
      },
      password: {
        type: String,
        default: ""
      },
      internet_package: {
        type: String,
        default: ""
      },
      ip_static: {
        type: String,
        default: ""
      },
      mac_address: {
        type: String,
        default: ""
      },
      job: {
        type: String,
        default: ""
      },

    })
  }
})



const schema = object({
  type_of_service: string().required(),
  // email: string().required(),
  // name: string().required(),
  // company: string().required(),
  // gender: string().required(),
  // card_identition: string().required(),
  // no_identition: string().required(),
  // area_code: string().required(),
  // phone: string().required(),
  // address: string().required(),
  // latitude: string().required(),
  // longitude: string().required(),
  // password: string().required(),
  // internet_package: string().required(),
  // ip_static: string().required(),
  // mac_address: string().required(),
  // job: string().required(),
});

type Schema = InferType<typeof schema>;

const state = reactive({
  type_of_service: "",
  email: "",
  name: "",
  company_id: "",
  gender: "",
  card_identition: "",
  no_identition: 0,
  area_id: "",
  phone: "",
  address: "",
  latitude: 0,
  longitude: 0,
  password: "",
  product_id: "",
  ip_static: "",
  mac_address: "",
  job: "",
});

watch(
  () => props.isEdit,
  (newValue) => {
    if (newValue) {
      state.type_of_service = props.data.type_of_service,
        state.email = props.data.email,
        state.name = props.data.name,
        state.company_id = props.data.company_id,
        state.gender = props.data.gender,
        state.card_identition = props.data.card_identition,
        state.no_identition = props.data.no_identition,
        state.area_id = props.data.area_id,
        state.phone = props.data.phone,
        state.address = props.data.address,
        state.latitude = props.data.latitude,
        state.longitude = props.data.longitude,
        state.password = props.data.password,
        state.product_id = props.data.product_id,
        state.ip_static = props.data.ip_static,
        state.mac_address = props.data.mac_address,
        state.job = props.data.job
    }
  },
  { immediate: true }
)
const emit = defineEmits(["success"]);

function onSuccess() {
  emit("success");
}

async function onSubmit(event: FormSubmitEvent<Schema>) {
  // Do something with event.data
  if (props.isEdit) {
    customerAdminApi().editCustomer(props.data.id, state).then((response) => {
      useToast().add({ title: response.message })
      onSuccess()
    }).catch((error) => {

    })
  } else {
    customerAdminApi().createCustomer(state).then((response) => {
      useToast().add({ title: response.message })
      onSuccess()
    }).catch((error) => {

    })
  }

}

function onMarkerDrag(e: any) {
  const latlng = e.target.getLatLng();
  state.latitude = latlng.lat;
  state.longitude = latlng.lng;
  reverseGeocode(latlng.lat, latlng.lng);
}

async function reverseGeocode(lat: number, lng: number) {
  try {
    const response = await fetch(
      `https://nominatim.openstreetmap.org/reverse?format=json&lat=${lat}&lon=${lng}`
    );
    const data = await response.json();
    state.address = data.display_name || "Address not found";
  } catch (error) {
    console.error("Reverse geocoding failed:", error);
    state.address = "Error fetching address";
  }
}

async function moveToMyLocation() {
  if (navigator.geolocation) {
    navigator.geolocation.getCurrentPosition(
      (position) => {
        const { latitude, longitude } = position.coords;
        state.latitude = latitude;
        state.longitude = longitude;
        reverseGeocode(latitude, longitude);
      },
      (error) => {
        console.error("Error getting location:", error);
      }
    );
  } else {
    console.error("Geolocation is not supported by this browser.");
  }
}

const genders = [
  {
    label: "Male",
    value: "male",
  },
  {
    label: "Female",
    value: "female",
  },
];

const card_identitions = [
  {
    label: "KTP",
    value: "ktp",
  },
  {
    label: "SIM",
    value: "sim",
  },
  {
    label: "PASPOR",
    value: "paspor",
  },
];



// const submission_types = [
//   {
//     label: "New",
//     value: "new",
//   },
//   {
//     label: "Upgrade",
//     value: "upgrade",
//   },
// ];

const type_of_services = [
  {
    label: "Internet",
    value: "internet",
  },
  {
    label: "IPTV",
    value: "iptv",
  },
];

const companies = ref([]);
const internet_packages = ref([]);
const areas = ref([])

async function getDataOptions() {
  companyAdminApi().getAllCompanies().then((response) => {
    companies.value = response.data.map((value: any, index: number) => ({
      label: value.name,
      value: value.id
    }))
  })

  internetPackageAdminApi().getAllInternetPacket().then((response) => {
    internet_packages.value = response.data.map((value: any, index: number) => ({
      label: value.name,
      value: value.id
    }))
  })

  areaAdminApi().getAllAreas().then((response) => {
    areas.value = response.data.map((value: any, index: number) => ({
      label: value.name_city + "-" + value.name_subdistrict + "-" + value.name_village,
      value: value.id
    }))
  })
}
await getDataOptions()


</script>

<template>
  <UModal>
    <div class="w-full p-4">


      <div class="p-2 mb-4 text-2xl font-bold text-center">
        <h1>{{ props.isEdit ? "Edit" : "Add New" }} Customer</h1>
      </div>
      <UForm :schema="schema" :state="state" class="space-y-4" @submit="onSubmit">
        <!-- <URadioGroup v-model="state.submission_type" legend="Submission Type" :options="submission_types" /> -->
        <div class="flex gap-4 flex-row-2">

          <div class="w-full">

            <URadioGroup v-model="state.type_of_service" legend="Type Of Service" :options="type_of_services" />

            <UFormGroup label="Email" name="email">
              <UInput v-model="state.email" />
            </UFormGroup>
            <UFormGroup label="Name" name="name">
              <UInput v-model="state.name" />
            </UFormGroup>
            <UFormGroup label="Company" name="company">
              <USelectMenu v-model="state.company_id" :options="companies" value-attribute="value"
                option-attribute="label" />
            </UFormGroup>
            <UFormGroup label="Gender" name="gender">
              <USelectMenu v-model="state.gender" :options="genders" value-attribute="value" option-attribute="label" />
            </UFormGroup>
            <UFormGroup label="Card Identition" name="card_identition">
              <USelectMenu v-model="state.card_identition" :options="card_identitions" value-attribute="value"
                option-attribute="label" />
            </UFormGroup>
            <UFormGroup label="No Indetition" name="no_identition">
              <UInput v-model="state.no_identition" type="number" />
            </UFormGroup>
            <UFormGroup label="Area Code" name="area_code">
              <USelectMenu v-model="state.area_id" :options="areas" value-attribute="value" />
            </UFormGroup>
            <UFormGroup label="Phone" name="phone">
              <UInput v-model="state.phone" />
            </UFormGroup>
            <UFormGroup label="Ip Static" name="ip_static">
              <UInput v-model="state.ip_static" />
            </UFormGroup>
            <UFormGroup label="Mac Address" name="mac_address">
              <UInput v-model="state.mac_address" />
            </UFormGroup>
          </div>
          <div class="w-full">

            <LMap style="height: 300px" :zoom="6" :center="[state.latitude, state.longitude]"
              :use-global-leaflet="false">
              <LTileLayer url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png" />
              <LMarker ref="map" :lat-lng="[state.latitude, state.longitude]" draggable @dragend="onMarkerDrag" />
              <LControl position="bottomleft">
                <UButton @click="moveToMyLocation">My Position</UButton>
              </LControl>
            </LMap>
            <UFormGroup label="Address" name="address">
              <UInput v-model="state.address" />
            </UFormGroup>
            <UFormGroup label="Map" name="coordinates">
              <div class="flex justify-around gap-4 flex-row-2">
                <div class="w-full text-sm">
                  <span>Latitude</span>
                  <UInput v-model="state.latitude" class="w-full" placeholder="Latitude" type="number" />
                </div>
                <div class="w-full text-sm">
                  Longitude
                  <UInput v-model="state.longitude" class="w-full" placeholder="Longitude" type="number" />
                </div>
              </div>
            </UFormGroup>
            <UFormGroup label="Packet Internet" name="internet_package">
              <USelectMenu v-model="state.product_id" :options="internet_packages" value-attribute="value"
                option-attribute="label" />
            </UFormGroup>

            <UFormGroup label="Job" name="job">
              <UInput v-model="state.job" />
            </UFormGroup>

            <div class="flex justify-end mt-4">
              <UButton type="submit"> Submit </UButton>
            </div>
          </div>
        </div>
      </UForm>
    </div>
  </UModal>
</template>
