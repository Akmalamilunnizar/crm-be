<script setup lang="ts">
import { userManagementAdminApi } from '@/api/admin/user-management';
import { object, string } from 'yup';
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
});

const schema = object({
  name: string()
    .min(3, "Must be at least 3 characters")
    .required("Required"),
})

if (props.id) {
  await userManagementAdminApi()
    .getRole(props.id)
    .then((response) => {
      Object.assign(state, response.data);
    })
    .catch((error) => {
      console.error("Error fetching companies:", error);
    });
}

const emit = defineEmits(['success'])

function onSuccess() {
  emit('success')
}

interface SubmitEvent {
    data: {
        name: string;
    };
}

async function onSubmit(event: SubmitEvent) {
   if (props.id) {
    await userManagementAdminApi()
      .editRole(props.id, state)
      .then((response) => {
        console.log("Success creating / editing company", response);
        onSuccess();
      })
      .catch((error) => {
        console.error("Error creating company:", error);
      });
  } else {
    await userManagementAdminApi()
      .createRole(state)
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
            <h1>{{ props.isEdit ? "Edit" : "Add New" }} Role</h1>
        </div>
        <UForm
        :schema="schema"
        :state="state"
        class="space-y-4"
        @submit="onSubmit"
        >
        <UFormGroup label="Role name" name="name">
            <UInput v-model="state.name" />
        </UFormGroup>
        <!-- <UFormGroup label="Technician" name="technician">
            <UInput v-model="state.technician" />
        </UFormGroup>
        <UFormGroup label="Date Installation" name="date">
            <UInput v-model="state.date" type="date" />
        </UFormGroup> -->
        <UButton type="submit">Submit</UButton>
        </UForm>
    </div>
    </UModal>
</template>
