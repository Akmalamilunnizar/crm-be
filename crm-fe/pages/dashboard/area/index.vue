<script setup lang="ts">
import { defaultAreaData, type Area } from "./area.model";
import AddAreaFrom from "./AddAreaForm.vue";
import { areaAdminApi } from "@/api/admin/area";
import AddAreaForm from "./AddAreaForm.vue";

const area = ref<any[]>([]);
const areaList = computed(()=>{
    if (!q.value) {
    areaData = area;
    return areaData.value.slice((page.value - 1) * pageCount, page.value * pageCount);
  }

  const newData = area.value.filter((oneArea) => {
    return Object.values(oneArea).some((value) => {
      return String(value).toLowerCase().includes(q.value.toLowerCase());
    });
  });

  return newData.slice((page.value - 1) * pageCount, page.value * pageCount);
})
const columns = [
  {
    key: "name_city",
    label: "City",
  },
  {
    key: "name_subdistrict",
    label: "Subdistrict",
  },
  {
    key: "name_village",
    label: "Village",
  },
  {
    key: "actions",
    label: "Manage",
  },
];

async function fetchData() {
  await areaAdminApi().getAllAreas().then((response) => {
    area.value = response.data
  }).catch((error) => { })
}

await fetchData()
const page = ref(1);
const pageCount = 5;

const q = ref("");

let areaData: any;

const filteredRows = computed(() => {
  if (!q.value) {
    areaData = area;
    return areaData.value.slice((page.value - 1) * pageCount, page.value * pageCount);
  }

  const newData = area.value.filter((oneArea) => {
    return Object.values(oneArea).some((value) => {
      return String(value).toLowerCase().includes(q.value.toLowerCase());
    });
  });
  // areaData = newData;
  return newData.slice((page.value - 1) * pageCount, page.value * pageCount);
});

const isOpen = ref(false);

function openModal(isEdit: boolean, data: any) {
  useModal().open(AddAreaForm, {
    isEdit,
    data,
    async onSuccess() {
      await fetchData()
      useModal().close()
    }
  })
}

async function deleteData(id: string){
  await areaAdminApi().deleteArea(id).then((response) => {
    fetchData()
  }).catch((error) => { })
  
}
</script>

<template>
  <div class="flex gap-2">
    <UButton icon="fa-plus" color="blue" label="Add New Area" @click="openModal(false,null)" />
  </div>
  <div class="flex px-3 py-3.5 border-b border-gray-200 dark:border-gray-700">
    <UInput v-model="q" placeholder="Filter Area..." />
  </div>

  <UTable :rows="filteredRows" :columns="columns">
    <template #actions-data="{ row }">
      <div class="flex gap-2">
        <UButton icon="fa-pencil" color="yellow" label="Edit" @click="openModal(true,row)" />
        <UButton icon="fa-trash" color="red" label="Delete" @click="deleteData(row.id)" />
      </div>
    </template>
  </UTable>

  <div class="flex justify-end px-3 py-3.5 border-t border-gray-200 dark:border-gray-700">
    <UPagination v-model="page" :page-count="pageCount" :total="area.length" />
  </div>

  <div>
    <UModal v-model="isOpen">
      <div class="p-4">
        <Placeholder class="h-48" />
        <AddAreaFrom />
      </div>
    </UModal>
  </div>
</template>
