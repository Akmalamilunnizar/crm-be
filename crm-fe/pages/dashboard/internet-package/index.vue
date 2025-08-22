<script setup lang="ts">
import FormAddComponent from './FormAddComponent.vue'
import { internetPackageAdminApi } from '@/api/admin/internet-package'
const internetPackages = ref<any[]>([]); // Use ref for reactivity
const internetPackageList = computed(() => {
    if (!q.value) {
        return internetPackages.value.slice((page.value - 1) * pageCount, page.value * pageCount);
    }

    const newData = internetPackages.value.filter((transaction) => {
        return Object.values(transaction).some((value) => {
            return String(value).toLowerCase().includes(q.value.toLowerCase());
        });
    });

    return newData.slice((page.value - 1) * pageCount, page.value * pageCount);
})
const column = [
    {
        key: "number",
        label: "Number"
    },
    {
        key: "name",
        label: "Name"
    },
    {
        key: "price",
        label: "Price"
    },
    {
        key: "description",
        label: "Description"
    },
    {
        key: "actions",
        label: "Action"
    },
]

async function fetchData() {
    await internetPackageAdminApi().getAllInternetPacket().then((response) => {
        internetPackages.value = response.data.map((product: any, index: number) => ({
            ...product,
            number: index + 1,
        }));
        console.log(internetPackages)
    }).catch((error) => {
        console.error("Error fetching companies:", error);
    })
}

async function deletePacket(id: string) {
    await internetPackageAdminApi().deleteInternetPacket(id).then((response) => {
        useToast().add({
            title: "Success Delete Product"
        })

        fetchData();
    }).catch((error) => {
        useToast().add({
            title: "Failed Delete Product",
            color: "red"
        })
    })

}

await fetchData()

const page = ref(1)
const pageCount = 5
const q = ref('')

const isOpen = ref(false)
const modal = useModal()

function openModal(isEdit: boolean, data: any) {
    modal.open(FormAddComponent, {
        isEdit: isEdit,
        data: data,
        async onSuccess() {
            try {
                await fetchData();
                console.log('Data fetched successfully:');
                modal.close();
            } catch (error) {
                console.error('Failed to fetch data:', error);
                // Optionally show an error message to the user
                // For example, update FormAddComponent state to display error
                modal.close(); // Close modal or keep open based on your needs
            }
        }
    })
}
const items = (row: any) => [
    [
        {
            label: "Edit",
            icon: "i-heroicons-pencil-square-20-solid",
            click: () => openModal(true, row),
        },
    ],
    [
        {
            label: "Delete",
            icon: "i-heroicons-trash-20-solid",
            click: () => deletePacket(row.id.toString()),
        },
    ],
];

</script>

<template>


    <UButton label="Add Internet Package" @click="openModal(false, null)" />
    <!-- <UButton label="Push" @click="pushData" /> -->
    <div class="flex px-3 py-3.5 border-b border-gray-200 dark:border-gray-700">
        <UInput v-model="q" placeholder="Search" />
    </div>
    <UTable :rows="internetPackageList" :columns="column">
        <template #actions-data="{ row }">
            <UDropdown :items="items(row)">
                <UButton color="gray" variant="ghost" icon="i-heroicons-ellipsis-horizontal-20-solid" />
            </UDropdown>
        </template>
    </UTable>

    <div class="flex justify-end px-3 py-3.5 border-t border-gray-200 dark:border-gray-700">
        <UPagination v-model="page" :page-count="pageCount" :total="internetPackages.length" />
    </div>

    <!-- <div>
        <UModal v-model="isOpen">
            <div class="p-4">
                <Placeholder class="h-48" />
                <FormAddComponent></FormAddComponent>
            </div>
        </UModal>
    </div> -->
</template>
