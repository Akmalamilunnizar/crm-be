<script setup lang="ts">
import type { FormSubmitEvent } from "#ui/types";
import {
  asset,
  assetSchema,
  status,
  status_in_out,
  type AssetSchema,
} from "./asset.model";
import { DatePicker as VCalendarDatePicker } from 'v-calendar'
import 'v-calendar/dist/style.css'
import { format } from 'date-fns'
import AddAreaForm from "../area/AddAreaForm.vue";
import { assetAdminApi } from "@/api/admin/asset";
const state = reactive(asset);

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
      type: {
        type: String,
        default: ""
      },
      brand: {
        type: String,
        default: ""
      },
      model: {
        type: String,
        default: ""
      },
      serial_number: {
        type: String,
        default: ""
      },
      mac_address: {
        type: String,
        default: ""
      },
      date: {
        type: String,
        default: new Date()
      },
      site: {
        type: String,
        default: ""
      },
      quantity: {
        type: String,
        default: ""
      },
      status: {
        type: String,
        default: ""
      },
      price: {
        type: String,
        default: ""
      },
      description: {
        type: String,
        default: ""
      },
      status_in_out: {
        type: String,
        default: ""
      },
    })
  }
})

watch(
  () => props.isEdit,
  (newValue) => {
    console.log(props.data)
    if (newValue) {
      state.type = props.data.type
      state.brand = props.data.brand
      state.model = props.data.model
      state.serial_number = props.data.serial_number
      state.mac_address = props.data.mac_address
      state.date = props.data.date
      state.site = props.data.site
      state.quantity = props.data.quantity
      state.status = props.data.status
      state.price = props.data.price
      state.description = props.data.description
      state.status_in_out = props.data.status_in_out
    }else{
      // clearState()
    }
  },
  { immediate: true }
)


const emit = defineEmits(["success"]);

function onSuccess() {
  emit("success");
}

function clearState() {
  state.type = ""
  state.brand = ""
  state.model = ""
  state.serial_number = ""
  state.mac_address = ""
  state.date = new Date()
  state.site = ""
  state.quantity = 0
  state.status = ""
  state.price = 0
  state.description = ""
  state.status_in_out = ""
}

async function onSubmit(event: FormSubmitEvent<AssetSchema>) {
  console.log(state)
  if (props.isEdit) {
    await assetAdminApi().editAsset(props.data.id, state).then((response) => {
      useToast().add({ title: response.message })
      onSuccess()
    }
    ).catch((err) => {
      useToast().add({ title: err, color: "red" })
      onSuccess()
    }
    )
  } else {
    await assetAdminApi().createAsset(state).then((response) => {
      useToast().add({ title: response.message })
      onSuccess()
    }
    ).catch((err) => {
      useToast().add({ title: err, color: "red" })
      onSuccess()
    }
    )

  }
}
</script>

<template>
  <UModal>
    <div class="p-4">


      <div class="p-2 mb-4 text-2xl font-bold text-center">
        <h1>{{ props.isEdit ? "Edit" : "Add New" }} Asset</h1>
      </div>
      <UForm :schema="assetSchema" :state="state" class="space-y-4" @submit="onSubmit">
        <UFormGroup label="Type" name="type">
          <UInput v-model="state.type" />
        </UFormGroup>

        <UFormGroup label="Brand" name="brand">
          <UInput v-model="state.brand" />
        </UFormGroup>

        <UFormGroup label="Model" name="model">
          <UInput v-model="state.model" />
        </UFormGroup>

        <UFormGroup label="Serial Number" name="serial_number">
          <UInput v-model="state.serial_number" />
        </UFormGroup>

        <UFormGroup label="Mac Address" name="mac_address">
          <UInput v-model="state.mac_address" />
        </UFormGroup>

        <UFormGroup label="Date" name="date">
          <UPopover :popper="{ placement: 'bottom-start' }">
            <UButton icon="i-heroicons-calendar-days-20-solid" :label="format(state.date, 'd MMM, yyy')" class="w-full"
              color="gray" />

            <template #panel="{ close }">
              <DatePickerComponent v-model="state.date" @close="close" />
            </template>
          </UPopover>
        </UFormGroup>

        <UFormGroup label="Site" name="site">
          <UInput v-model="state.site" />
        </UFormGroup>

        <UFormGroup label="Quantity" name="quantity">
          <UInput v-model="state.quantity" type="number" />
        </UFormGroup>

        <UFormGroup label="Status Asset IN/Out" name="status_in_out">
          <USelect v-model="state.status_in_out" :options="status_in_out"></USelect>
        </UFormGroup>

        <UFormGroup label="Status" name="status">
          <USelect v-model="state.status" :options="status"> </USelect>
        </UFormGroup>

        <UFormGroup label="Price" name="price">
          <UInput v-model="state.price" type="number" />
        </UFormGroup>

        <UFormGroup label="Description" name="description">
          <UTextarea v-model="state.description" />
        </UFormGroup>

        <UButton type="submit"> Submit </UButton>
      </UForm>
    </div>
  </UModal>
</template>
