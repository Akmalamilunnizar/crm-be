<script setup lang="ts">
import { object, string } from 'yup';

const state = reactive({

  customer: "",
  technician: "",
  date: "",
});

const schema = object({
  customer: string()
    .min(3, "Must be at least 3 characters")
    .required("Required"),
  technician: string()
    .min(3, "Must be at least 3 characters")
    .required("Required"),
  date: string().required("Required"),
})



const emit = defineEmits(['success'])

function onSuccess() {
  emit('success')
}

interface SubmitEvent {
    data: {
        customer: string;
        technician: string;
        date: string;
    };
}

function onSubmit(event: SubmitEvent) {
    // Do something with event.data
    console.log(event.data);
    onSuccess();
}
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
        <UFormGroup label="Customer" name="customer">
            <UInput v-model="state.customer" />
        </UFormGroup>
        <UFormGroup label="Technician" name="technician">
            <UInput v-model="state.technician" />
        </UFormGroup>
        <UFormGroup label="Date Installation" name="date">
            <UInput v-model="state.date" type="date" />
        </UFormGroup>
        <UButton type="submit">Submit</UButton>
        </UForm>
    </div>
    </UModal>
</template>
