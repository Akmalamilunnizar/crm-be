<script setup lang="ts">
import { object, string, type InferType } from "yup";
import type { FormSubmitEvent } from "#ui/types";
import { companyAdminApi } from "@/api/admin/company";
const props = defineProps({
  isEdit: {
    type: Boolean,
    default: false,
  },
  id: {
    type: String,
    required: false,
  },
});

const state = reactive({
  name: "",
  url: "",
  email: "",
  phone: "",
  logo_url: "",
  description: "",
  npwp: "",
  address: "",
});

if (props.id) {
  await companyAdminApi()
    .getCompany(props.id)
    .then((response) => {
      Object.assign(state, response.data);
    })
    .catch((error) => {
      console.error("Error fetching companies:", error);
    });
}
const schema = object({
  name: string().required("Name is required"),
  // url: string().url("URL must be a valid URL").required("URL is required"),
  email: string()
    .required("Email is required")
    .email("Email must be a valid email"),
  phone: string()
    .matches(/^\d+$/, "Phone number must be digits")
    .required("Phone number is required"),
  npwp: string().required("NPWP is required"),
  address: string().required("Address is required"),
});

const emit = defineEmits(["success"]);

function onSuccess() {
  emit("success");
}

async function onSubmit() {
  if (props.id) {
    await companyAdminApi()
      .editCompany(props.id, state)
      .then((response) => {
        console.log("Success creating / editing company", response);
        onSuccess();
      })
      .catch((error) => {
        console.error("Error creating company:", error);
      });
  } else {
    await companyAdminApi()
      .createCompanies(state)
      .then((response) => {
        console.log("Success creating / editing company", response);
        onSuccess();
      })
      .catch((error) => {
        console.error("Error creating company:", error);
      });
  }
}
</script>

<template>
  <UModal>
    <div class="p-4">
      <div class="p-2 mb-4 text-2xl font-bold text-center">
        <h1>{{ props.isEdit ? "Edit" : "Add New" }} Company</h1>
      </div>
      <UForm
        :schema="schema"
        :state="state"
        class="space-y-4"
        @submit="onSubmit"
      >
        <UFormGroup label="Name" name="name">
          <UInput v-model="state.name" />
        </UFormGroup>
        <UFormGroup label="URL" name="url">
          <UInput v-model="state.url" />
        </UFormGroup>
        <UFormGroup label="Email" name="email">
          <UInput v-model="state.email" />
        </UFormGroup>
        <UFormGroup label="Phone" name="phone">
          <UInput v-model="state.phone" />
        </UFormGroup>
        <UFormGroup label="Logo URL" name="logo_url">
          <UInput v-model="state.logo_url" />
        </UFormGroup>
        <UFormGroup label="NPWP" name="npwp">
          <UInput v-model="state.npwp" />
        </UFormGroup>
        <UFormGroup label="Address" name="address">
          <UInput v-model="state.address" />
        </UFormGroup>
        <UFormGroup label="Description" name="description">
          <UInput v-model="state.description" />
        </UFormGroup>
        <UButton type="submit"> Submit </UButton>
      </UForm>
    </div>
  </UModal>
</template>
