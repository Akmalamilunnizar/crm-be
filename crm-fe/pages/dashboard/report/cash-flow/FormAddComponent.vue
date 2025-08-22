<script setup lang="ts">
import { boolean, object, string, type InferType } from 'yup'
import type { FormSubmitEvent } from '#ui/types'

const schema = object({
    // email: string().email('Invalid email').required('Required'),
    // password: string()
    //     .min(8, 'Must be at least 8 characters')
    //     .required('Required'),
    // name: string()
    //     .min(3, 'Must be at least 3 characters')
    //     .required('Required'),
    // phone: string()
    //     .min(3, 'Must be at least 3 characters')
    //     .required('Required'),
    // address: string()
    //     .min(3, 'Must be at least 3 characters')
    //     .required('Required'),
    // packet_internet: string()
    //     .min(3, 'Must be at least 3 characters')
    //     .required('Required')
})

type Schema = InferType<typeof schema>

const state = reactive({
    type_cash: '',
    nominal: 0,
    description: '',
    type_use_of_money: '',
    type_transaction: '',
})
const type_cash = [
    { label: 'Cash', value: 'cash' },
    { label: 'Credit Card', value: 'credit_card' },
    { label: 'Debit Card', value: 'debit_card' },
]

const type_use_of_money = [
    { label: 'Operational', value: 'operational' },
    { label: 'Asset', value: 'asset' },
    { label: 'Investment', value: 'investment' },
    { label: 'Other', value: 'other' },
]
const type_transaction = [
    { label: 'Cash In', value: 'cash_in' },
    { label: 'Cash Out', value: 'cash_out' },
]

defineProps({
    // isOpen: { type: Boolean, required: true }
})
const emit = defineEmits(['success', 'close'])
async function onSubmit(event: FormSubmitEvent<Schema>) {
    // Do something with event.data
    // Emit an event to the parent component
    // Close the modal
    clearState()
    // Emit success event
    emit('success')
}
function clearState() {
    state.type_cash = ''
    state.nominal = 0
    state.description = ''
}
</script>

<!-- <template>
    <UModal>
        <UCard>
            <div class="space-y-2">
                <p>This modal was opened programmatically !</p>
                <p>Count:</p>
                <UButton @click="onSubmit">
                    Click to emit a success event
                </UButton>
            </div>
        </UCard>
    </UModal>
</template> -->
<template>
    <UModal>
        <div class="p-4">


            <div class="p-2 mb-4 text-2xl font-bold text-center">
                <h1>Add New Report Cash Flow</h1>
            </div>
            <UForm :schema="schema" :state="state" class="space-y-4" @submit="onSubmit">
                <UFormGroup label="Type Cash" name="type_cash">
                    <USelectMenu v-model="state.type_cash" :options="type_cash" value-attribute="value"
                        option-attribute="label" />
                </UFormGroup>
                <UFormGroup label="Type Transaction" name="type_transaction">
                    <USelectMenu v-model="state.type_transaction" :options="type_transaction" value-attribute="value"
                        option-attribute="label" />
                </UFormGroup>
                <UFormGroup label="Type Used Of Money" name="type_use_of_money">
                    <USelectMenu v-model="state.type_use_of_money" :options="type_use_of_money" value-attribute="value"
                        option-attribute="label" />
                </UFormGroup>
                <UFormGroup label="Nominal" name="nominal">
                    <UInput type="number" v-model="state.nominal" />
                </UFormGroup>
                <UFormGroup label="Description" name="description">
                    <UTextarea v-model="state.description" />
                </UFormGroup>

                <UButton type="submit">
                    Submit
                </UButton>
            </UForm>
        </div>
    </UModal>
</template>
