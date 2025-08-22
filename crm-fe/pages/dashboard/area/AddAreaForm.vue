<script setup lang="ts">
import { object, string, type InferType } from "yup";
import type { FormSubmitEvent } from "#ui/types";
import { areaAdminApi } from "@/api/admin/area";

const schema = object({
  name_city: string().required("Name is required").min(3, "min 3 words"),
  name_subdistrict: string().required("Name is required").min(4, "min 4 words"),
  name_village: string().required("Name is required").min(4, "min 4 words"),
});

type Schema = InferType<typeof schema>;

const props = defineProps({
  isEdit: {
    type: Boolean,
    default: false,
  },
  data: {
    type: Object,
    default: () => ({
      id: {
        type: String,
        default: ''
      },
      name_city: {
        type: String,
        default: ''
      },
      name_subdistrict: {
        type: String,
        default: ''
      },
      name_village: {
        type: String,
        default: ''
      }
    })
  }
})
const state = reactive({
  name_city: "",
  name_subdistrict: "",
  name_village: "",
});

watch(
  () => props.isEdit,
  (newValue) => {
    console.log("newValue", props.isEdit, newValue)
    if (newValue) {
      state.name_city = props.data.name_city
      state.name_subdistrict = props.data.name_subdistrict
      state.name_village = props.data.name_village
    } else {
      clearState()
    }
  },
  { immediate: true }
)

function clearState() {
  state.name_city = ""
}

const emit = defineEmits(["success"]);

function onSuccess() {
  emit("success");
}
async function onSubmit(event: FormSubmitEvent<Schema>) {
  // Do something with event.data
  if (props.isEdit) {
    areaAdminApi().editArea(props.data.id, state).then((response) => {
      useToast().add({ title: response.message, color: "green" })
      onSuccess()
    }).catch((error) => {
      useToast().add({ title: "Failed Create Area", color: "red" })
    })
  } else {

    areaAdminApi().createArea(state).then((response) => {
      useToast().add({ title: response.message, color: "green" })
      onSuccess()
    }).catch((error) => {
      useToast().add({ title: "Failed Create Area", color: "red" })
    })
  }
}
</script>

<template>
  <UModal>
    <div class="p-4">
      <div class="p-2 mb-4 text-2xl font-bold text-center">
        <h1>{{ props.isEdit ? "Edit" : "Add New" }} Area</h1>
      </div>

      <UForm :schema="schema" :state="state" class="space-y-4" @submit="onSubmit">
        <UFormGroup label="Name City" name="name_city">
          <UInput v-model="state.name_city" />
        </UFormGroup>
        <UFormGroup label="Name Subdistrict" name="name_subdistrict">
          <UInput v-model="state.name_subdistrict" />
        </UFormGroup>
        <UFormGroup label="Name Village" name="name_village">
          <UInput v-model="state.name_village" />
        </UFormGroup>

        <UButton type="submit" color="blue"> Submit </UButton>
      </UForm>
    </div>
  </UModal>
</template>
