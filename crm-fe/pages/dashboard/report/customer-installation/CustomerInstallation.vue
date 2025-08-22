<script setup lang="ts">
import { archiveInstallationAdminApi } from '@/api/admin/archive-installation';
import FormCustomerInstallation from './FormCustomerInstallation.vue';
import ImageViewComponent from './ImageViewComponent.vue';


const dataList = ref([]);
const columns = [
  { key: "customer", label: "Customer" },
  { key: "technician", label: "Technician" },
  { key: "date", label: "Date Installation" },
  {
    key:"actions",
    label: "Actions",
  }
];
const q = ref("");
const page = ref(1);
const pageCount = 5;

const modal = useModal()
function handleClick(row:any) {
  modal.open(ImageViewComponent, {
    image: row.images.map((item:any) => useApiHost()+"/"+item.full_path),
    onClose: () => {
      modal.close();
    },
  });
}


async function getData() {
  archiveInstallationAdminApi().getAllArchiveInstallation().then((res) => {
    res.data.map((item:any) => {
      item.customer = item.customer.name;
      item.technician = item.technician.name;
      item.date = item.date.split("T")[0];
    });
    dataList.value = res.data;
  }).catch((err) => {
    console.log(err);
  });
}

const rows = computed(() => {
    if (!q.value) {
        return dataList.value.slice((page.value - 1) * pageCount, page.value * pageCount);
    }

    const newData = dataList.value.filter((data) => {
        return Object.values(data).some((value) => {
            return String(value).toLowerCase().includes(q.value.toLowerCase());
        });
    });

    return newData.slice((page.value - 1) * pageCount, page.value * pageCount);
});

getData();
</script>

<template>
        <div
          class="flex px-3 py-3.5 border-b border-gray-200 dark:border-gray-700"
        >
          <UInput v-model="q" placeholder="Search" />
        </div>
        <UTable :rows="rows" :columns="columns">

            <template #actions-data="{ row }">
                <UButton @click="handleClick(row)">View Image</UButton>
            <!-- <UDropdown>
                <UButton color="gray" variant="ghost" icon="i-heroicons-ellipsis-horizontal-20-solid" />
            </UDropdown> -->
        </template>
        </UTable>

        <div
          class="flex justify-end px-3 py-3.5 border-t border-gray-200 dark:border-gray-700"
        >
          <UPagination
            v-model="page"
            :page-count="pageCount"
            :total="dataList.length"
          />
        </div>
</template>