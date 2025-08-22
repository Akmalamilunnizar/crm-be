<script setup lang="ts">
import { boolean, object, string, type InferType } from "yup";
import type { FormSubmitEvent } from "#ui/types";
import { userManagementAdminApi } from "@/api/admin/user-management";
import { customerAdminApi } from "@/api/admin/customer";
import { uploadFileAdminApi } from "@/api/admin/file-upload";

const props = defineProps({
  isEdit: {
    type: Boolean,
    required: false,
  },
  data: {
    type: Object,
    default: () => ({
      id: {
        type: string,
        default: "",
      },
      technician_id: {
        type: string,
        default: "",
      },
    }),
  },
});

const schema = object({
  technician_id: string().required("Customer is required"),
  date: string().required("Date is required"),
});

const state = reactive({
  customer_id: "",
  technician_id: "",
  date: "",
  description: "",
  image_ids: [] as string[],
  previews: [] as string[],
  selectedImage: "",
  showModal: false,
});

state.date =
  new Date().getFullYear() +
  "-" +
  (new Date().getMonth() + 1).toString().padStart(2, "0") +
  "-" +
  new Date().getDate().toString().padStart(2, "0");

watch(
  () => props.isEdit,
  (newValue) => {
    if (newValue) {
      state.customer_id = props.data.id;
    }
  },
  { immediate: true }
);

type Schema = InferType<typeof schema>;
async function onSubmit(event: FormSubmitEvent<Schema>) {
  await customerAdminApi()
    .createCustomerInstallation(state)
    .then((response) => {
      console.log("Success creating / editing company", response);
      onSuccess();
    })
    .catch((error) => {
      console.error("Error creating company:", error);
    });
  // Do something with event.data
  console.log(event.data);
}

const emit = defineEmits(["success"]);

function onSuccess() {
  emit("success");
}

const getFormattedMonth = () => {
  return (new Date().getMonth() + 1).toString().padStart(2, "0");
};

// Handle file upload with safe event handling
const handleFileUpload = async (event: FileList) => {
  const input = event;

  if (input) {
    const files = input;


    Array.from(files).forEach(async (file, index) => {
      if (file.type.startsWith("image/")) {
        await uploadFileAdminApi()
          .createUploadFile({
            name: `imageInstallation${index + 1}`,
            path: `${state.technician_id}/${props.data.id}`,
            file: file,
          })
          .then((response) => {
            state.image_ids[index] = response.data.id;
          });
        const reader = new FileReader();

        reader.onload = (e) => {
          if (e.target?.result) {
            state.previews.push(e.target.result as string); // Add base64 preview
          }
        };
        reader.readAsDataURL(file); // Convert to base64 for preview
      } else {
        console.warn(`File ${file.name} is not an image`);
      }
    });
  } else {
    console.error("No files selected or input element not found");
  }
};

const options = ref<any[]>([]);

async function getTechnicianData() {
  userManagementAdminApi()
    .getAllUsers({ query: { role: "TECHNICIAN" } })
    .then((response) => {
      options.value = [...response.data];
    });
}

await getTechnicianData();
</script>

<template>
  <UModal>
    <div class="p-4">
      <div class="p-2 mb-4 text-2xl font-bold text-center">
        <h1>Add New Customer Installation</h1>
      </div>
      <UForm
        :schema="schema"
        :state="state"
        class="space-y-4"
        @submit="onSubmit"
      >
        <UFormGroup label="Technician" name="technician_id">
          <USelectMenu
            v-model="state.technician_id"
            :options="options"
            placeholder="Select a person"
            searchable
            searchable-placeholder="Search by Technician name"
            option-attribute="name"
            value-attribute="id"
            :search-attributes="['name']"
          />
        </UFormGroup>
        <UFormGroup label="Description" name="description">
          <UTextarea autoresize v-model="state.description" />
        </UFormGroup>
        <UFormGroup label="Date Installation" name="date">
          <UInput v-model="state.date" type="date" />
        </UFormGroup>
        <UFormGroup
          label="Upload Photos"
          name="image"
          help="Select multiple images (PNG, JPG, etc.)"
        >
          <UInput
            :disabled="!state.technician_id"
            ref="fileInput"
            type="file"
            size="sm"
            icon="i-heroicons-folder"
            multiple
            accept="image/*"
            @change="handleFileUpload"
          />
        </UFormGroup>
        <div v-if="state.previews.length" class="grid grid-cols-3 gap-2 mt-4">
          <img
            v-for="(preview, index) in state.previews"
            :key="index"
            :src="preview"
            alt="Preview"
            class="object-cover w-full h-32 rounded-md"
          />
        </div>
        <UButton type="submit">Submit</UButton>
      </UForm>
    </div>
  </UModal>
</template>
