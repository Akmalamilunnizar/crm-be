<script setup lang="ts">
import FormAddComponent from './FormAddComponent.vue'
const people = [{
    id: 1,
    name: 'Lindsay Walton',
    email: 'lindsay.walton@example.com',
    role: 'Member'
}, {
    id: 2,
    name: 'Courtney Henry',
    email: 'courtney.henry@example.com',
    role: 'Admin'
}, {
    id: 3,
    name: 'Tom Cook',
    email: 'tom.cook@example.com',
    role: 'Member'
}, {
    id: 4,
    name: 'Whitney Francis',
    email: 'whitney.francis@example.com',
    role: 'Admin'
}, {
    id: 5,
    name: 'Leonard Krasner',
    email: 'leonard.krasner@example.com',
    role: 'Owner'
}, {
    id: 6,
    name: 'Floyd Miles',
    email: 'floyd.miles@example.com',
    role: 'Member'
}, {
    id: 7,
    name: 'Emily Selman',
    email: '',
    role: 'Admin'
}, {
    id: 8,
    name: 'Kristin Watson',
    email: '',
    role: 'Member'
}, {
    id: 9,
    name: 'Emma Watson',
    email: '',
    role: 'Member'
}, {
    id: 10,
    name: 'John Doe',
    email: '',
    role: 'Admin'
}, {
    id: 11,
    name: 'Jane Doe',
    email: '',
    role: 'Member'
}, {
    id: 12,
    name: 'John Smith',
    email: '',
    role: 'Admin'
}, {
    id: 13,
    name: 'Jane Smith',
    email: '',
    role: 'Owner'
}, {
    id: 14,
    name: 'Jane Smith',
    email: '',
    role: 'Owner'
}, {
    id: 15,
    name: 'Jane Smith',
    email: '',
    role: 'Owner'
}, {
    id: 16,
    name: 'Jane Smith',
    email: '',
    role: 'Owner'
}, {
    id: 17,
    name: 'Jane Smith',
    email: '',
    role: 'Owner'
}]

const page = ref(1)
const pageCount = 5

const rows = computed(() => {
    return people.slice((page.value - 1) * pageCount, (page.value) * pageCount)
})

const q = ref('')

let peopleData = people

const filteredRows = computed(() => {
    if (!q.value) {
        peopleData = people
        return people.slice((page.value - 1) * pageCount, (page.value) * pageCount)
    }

    const newData = people.filter((person) => {
        return Object.values(person).some((value) => {
            // person with paginate
            return String(value).toLowerCase().includes(q.value.toLowerCase())
        })
    })
    peopleData = newData
    return newData.slice((page.value - 1) * pageCount, (page.value) * pageCount)
})

const isOpen = ref(false)

</script>

<template>

    <UButton label="Add Customer" @click="isOpen = true" />
    <div class="flex px-3 py-3.5 border-b border-gray-200 dark:border-gray-700">
        <UInput v-model="q" placeholder="Search" />
    </div>
    <UTable :rows="filteredRows" />

    <div class="flex justify-end px-3 py-3.5 border-t border-gray-200 dark:border-gray-700">
        <UPagination v-model="page" :page-count="pageCount" :total="peopleData.length" />
    </div>

    <div>
        <UModal v-model="isOpen">
            <div class="p-4">
                <Placeholder class="h-48" />
                <FormAddComponent></FormAddComponent>
            </div>
        </UModal>
    </div>
</template>
