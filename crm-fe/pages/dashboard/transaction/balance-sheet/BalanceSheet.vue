<script setup lang="ts">
const props = defineProps<{
  data: any[];
  onRefresh: () => Promise<void>;
}>();
const q = ref("");
const page = ref(1);
const pageCount = 5;

const dataTable = ref<any[]>(
  props.data
);
const dataTableList = computed(() => {
  if (!q.value) {
    return dataTable.value.slice((page.value - 1) * pageCount, page.value * pageCount);
  }

  const newData = dataTable.value.filter((transaction) => {
    return Object.values(transaction).some((value) => {
      return String(value).toLowerCase().includes(q.value.toLowerCase());
    });
  });

  return newData.slice((page.value - 1) * pageCount, page.value * pageCount);
})
type User = {
  id: number;
  username: string;
  fullName: string;
  Type: string;
};
const columns = [
  { key: "number", label: "Number" },
  { key: "name", label: "Account" },
  { key: "saldo", label: "Balance" },
];

const toast = useToast();
const modal = useModal();
const count = ref(0);

const items = (row: User) => [
  [
    {
      label: "Edit",
      icon: "i-heroicons-pencil-square-20-solid",
      click: () => console.log("Edit", row.id),
    },
  ],
  [
    {
      label: "Delete",
      icon: "i-heroicons-trash-20-solid",
    },
  ],
];

</script>

<template>
  <!-- <UButton label="Add Expense" @click="openModal" /> -->
  <div class="flex px-3 py-3.5 border-b border-gray-200 dark:border-gray-700">
    <UInput v-model="q" placeholder="Search" />
  </div>
  <UTable :rows="dataTableList" :columns="columns">
    <template #actions-data="{ row }">
      <UDropdown :items="items(row)">
        <UButton
          color="gray"
          variant="ghost"
          icon="i-heroicons-ellipsis-horizontal-20-solid"
        />
      </UDropdown>
    </template>
  </UTable>

  <div
    class="flex justify-end px-3 py-3.5 border-t border-gray-200 dark:border-gray-700"
  >
    <UPagination
            v-model="page"
            :page-count="pageCount"
            :total="dataTableList.length"
          />
  </div>
</template>
