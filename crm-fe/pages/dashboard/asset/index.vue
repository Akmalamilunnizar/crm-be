<script setup lang="ts">
import { assetAdminApi } from '@/api/admin/asset'
import FormAddComponent from './AddAssetForm.vue'
import { defaultAssets } from './asset.model'
import { tableColumns } from './table'

const page = ref(1)
const pageCount = 5

const q = ref('')

let assetData: any = defaultAssets

const filteredRows = computed(() => {
    if (!q.value) {
        return defaultAssets.value.slice((page.value - 1) * pageCount, (page.value) * pageCount)
    }

    const newData = defaultAssets.value.filter((person) => {
        return Object.values(person).some((value) => {
            return String(value).toLowerCase().includes(q.value.toLowerCase())
        })
    })
    return newData.slice((page.value - 1) * pageCount, (page.value) * pageCount)
})

async function getData() {
    await assetAdminApi().getAllAssets().then((response) => {
        response.data.forEach((asset: any) => {
            asset.number = response.data.indexOf(asset) + 1;
        });

        defaultAssets.value = [...response.data]
    }).catch()
}

async function deleteAsset(assetId: string) {
    await assetAdminApi().deleteAsset(assetId).then((response) => {
        useToast().add({
            title: response.message
        })
        getData()


    }).catch((err) => {
        useToast().add({
            title: err,
            color: "red"
        })
        getData()

    })

}

await getData()
const isOpen = ref(false)

const toast = useToast()
const modal = useModal()

function OpenModalAddAsset(isEdit: boolean, data: any) {
    modal.open(FormAddComponent, {
        isEdit,
        data,
        async onSuccess() {
            await getData()
            modal.close()
        }
    })
}

const items = (row: any) => [
    [{
        label: 'Edit',
        icon: 'i-heroicons-pencil-square-20-solid',
        click: () => OpenModalAddAsset(true, row)
    }, {
        label: 'Delete',
        icon: 'i-heroicons-trash-20-solid',
        click: () => deleteAsset(row.id)
    }]
]

</script>

<template>
    <UButton label="Add Asset" @click="OpenModalAddAsset(false, null)" />
    <div class="flex px-3 py-3.5 border-b border-gray-200 dark:border-gray-700">
        <UInput v-model="q" placeholder="Filter asset..." />
    </div>
    <div class="table-scroll-container">
        <div class="table-scroll-content">
            <UTable :columns="tableColumns" :rows="filteredRows" class="dashboard-table">
                <template #actions-data="{ row }">
                    <UDropdown :items="items(row)">
                        <UButton color="gray" variant="ghost" icon="i-heroicons-ellipsis-horizontal-20-solid" />
                    </UDropdown>
                </template>
            </UTable>
        </div>
        <div class="table-scroll-footer">
          <span class="scroll-hint">↔ Scroll horizontally to see more columns | ↕ Scroll vertically for more rows</span>
        </div>
    </div>


    <div class="flex justify-end px-3 py-3.5 border-t border-gray-200 dark:border-gray-700">
        <UPagination v-model="page" :page-count="pageCount" :total="assetData.length" />
    </div>

</template>
