<script setup lang="ts">
import { object, string, type InferType } from "yup";
import type { FormSubmitEvent } from "#ui/types";
import { userManagementAdminApi } from "@/api/admin/user-management";
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
const roles = ref<any[]>([]);

const state = reactive({
  name: "",
  email: "",
  logo_url: "",
  role_id: "",
  password: "",
  password_confirm: "",
  phone: "",
});
onMounted(async () => {
  await userManagementAdminApi()
    .getAllRole()
    .then((response) => {
      roles.value = response.data.map((role: any) => ({
        label: role.name,
        value: role.id,
      }));
    })
    .catch((error) => {
      console.error("Error fetching user:", error);
    });
  if (props.id) {
    await userManagementAdminApi()
      .getUser(props.id)
      .then((response) => {
        Object.assign(state,{...response.data,role_id:response.data.role?.id});
      })
      .catch((error) => {
        console.error("Error fetching companies:", error);
      });
  }
});

const schema = object({
  name: string().required("Name is required"),
  email: string()
    .email("Email must be a valid email")
    .required("Email is required"),
  role_id: string().required("Role is required"),
  password: string()
    .min(6, "Password must be at least 6 characters")
    .required("Password is required"),
  password_confirm: string()
    .required("Confirm Password is required")
    .test("passwords-match", "Passwords must match", function (value) {
      const { password } = this.parent; // Access parent object, which includes the password
      return password === value; // Compare password and confirmPassword
    }),
  phone: string()
    .matches(/^\d+$/, "Phone number must be digits")
    .required("Phone number is required"),
  // logo_url: string()
  //   .url("Logo URL must be a valid URL")
  //   .required("Logo URL is required"),
});

const emit = defineEmits(["success"]);

function onSuccess() {
  emit("success");
}

async function onSubmit() {
  console.log("schema", schema);
  console.log("state", state);
  if (props.id) {
    await userManagementAdminApi()
      .editUser(props.id, state)
      .then((response) => {
        console.log("Success creating / editing company", response);
        onSuccess();
      })
      .catch((error) => {
        console.error("Error creating company:", error);
      });
  } else {
    await userManagementAdminApi()
      .createUser(state)
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
    <div class="p-2 mb-4 text-2xl font-bold text-center">
      <h1>{{ props.isEdit ? "Edit" : "Add New" }} User</h1>
    </div>
    <div class="p-4">
      <UForm
        :schema="schema"
        :state="state"
        class="space-y-4"
        @submit="onSubmit"
      >
        <UFormGroup label="Username" name="name">
          <UInput v-model="state.name" />
        </UFormGroup>
        <UFormGroup label="Email" name="email">
          <UInput v-model="state.email" />
        </UFormGroup>
        <UFormGroup label="Phone number" name="phone">
          <UInput v-model="state.phone" />
        </UFormGroup>
        <div v-if="roles.length === 0">
          <span class="text-gray-500">Loading roles...</span>
        </div>
        <URadioGroup
          v-else
          v-model="state.role_id"
          :options="roles"
          name="role_id"
          legend="Role"
        />
        <UFormGroup label="Password" name="password">
          <UInput type="password" v-model="state.password" />
        </UFormGroup>
        <UFormGroup label="Confirm Password" name="password_confirm">
          <UInput type="password" v-model="state.password_confirm" />
        </UFormGroup>
        <UFormGroup label="Logo URL" name="logo_url">
          <UInput v-model="state.logo_url" />
        </UFormGroup>
        <UButton type="submit">Submit</UButton>
      </UForm>
    </div>
  </UModal>
</template>
