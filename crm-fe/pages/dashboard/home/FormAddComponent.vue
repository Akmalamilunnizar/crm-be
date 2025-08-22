<script setup lang="ts">
import { object, string, type InferType } from 'yup'
import type { FormSubmitEvent } from '#ui/types'

const schema = object({
    email: string().email('Invalid email').required('Required'),
    password: string()
        .min(8, 'Must be at least 8 characters')
        .required('Required'),
    name: string()
        .min(3, 'Must be at least 3 characters')
        .required('Required'),
    phone: string()
        .min(3, 'Must be at least 3 characters')
        .required('Required'),
    address: string()
        .min(3, 'Must be at least 3 characters')
        .required('Required'),
    packet_internet: string()
        .min(3, 'Must be at least 3 characters')
        .required('Required')
})

type Schema = InferType<typeof schema>

const state = reactive({
    email: undefined,
    password: undefined,
    name: undefined,
    phone: undefined,
    address: undefined,
    packet_internet: undefined
})

async function onSubmit(event: FormSubmitEvent<Schema>) {
    // Do something with event.data
    console.log(event.data)
}
</script>

<template>
    <div class="p-2 mb-4 text-2xl font-bold text-center">
        <h1>Add New Customer</h1>
    </div>
    <UForm :schema="schema" :state="state" class="space-y-4" @submit="onSubmit">
        <UFormGroup label="Email" name="email">
            <UInput v-model="state.email" />
        </UFormGroup>
        <UFormGroup label="Name" name="name">
            <UInput v-model="state.name" />
        </UFormGroup>
        <UFormGroup label="Phone" name="phone">
            <UInput v-model="state.phone" />
        </UFormGroup>
        <UFormGroup label="Address" name="address">
            <UInput v-model="state.address" />
        </UFormGroup>
        <UFormGroup label="Packet Internet" name="packet_internet">
            <UInput v-model="state.packet_internet" />
        </UFormGroup>

        <UButton type="submit">
            Submit
        </UButton>
    </UForm>
</template>
