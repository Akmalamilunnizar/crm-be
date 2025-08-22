<script setup lang="ts">
import { object, string, type InferType } from "yup";
import type { FormSubmitEvent } from "#ui/types";

const schema = object({
  email: string().email("Invalid email").required("Required"),
  password: string()
    .min(8, "Must be at least 8 characters")
    .required("Required"),
  name: string().min(3, "Must be at least 3 characters").required("Required"),
  phone: string().min(3, "Must be at least 3 characters").required("Required"),
  address: string()
    .min(3, "Must be at least 3 characters")
    .required("Required"),
  packet_internet: string()
    .min(3, "Must be at least 3 characters")
    .required("Required"),
});

type Schema = InferType<typeof schema>;

const state = reactive({
  email: undefined,
  password: undefined,
  name: undefined,
  phone: undefined,
  address: undefined,
  packet_internet: undefined,
  mac_address: undefined,
  quantity: undefined,
  role: undefined,
});

async function onSubmit(event: FormSubmitEvent<Schema>) {
  // Do something with event.data
  console.log(event.data);
}

const roles = [
  { label: "Admin", value: "admin" },
  { label: "Member", value: "member" },
  { label: "Owner", value: "owner" },
];
</script>

<template>
  <div class="p-2 mb-4 text-2xl font-bold text-center">
    <h1>Add User Access</h1>
  </div>
  <UForm :schema="schema" :state="state" class="space-y-4" @submit="onSubmit">
    <UFormGroup label="Name" name="name">
      <UInput v-model="state.name" />
    </UFormGroup>
    <UFormGroup label="Email" name="email">
      <UInput v-model="state.email" />
    </UFormGroup>
    <UFormGroup label="Password" name="password">
      <UInput v-model="state.password" />
    </UFormGroup>

    <UFormGroup label="Role" name="role">
      <USelectMenu
        v-model="state.role"
        :options="roles"
        value-attribute="value"
        option-attribute="label"
      />
    </UFormGroup>

    <UButton type="submit"> Submit </UButton>
  </UForm>
</template>
