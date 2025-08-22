<script setup lang="ts">
import { formatIDR } from '@/helper/currency';

const props = defineProps<{
  data: any[];
}>();
const q = ref("");
const page = ref(1);
const pageCount = 5;

const dataTable = ref<any[]>(
  props.data.map((data) => {
    data.account = data.account.name;
    return data;
  })
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

// watch(
//   () => props.data,
//   (newData) => {
//     rows.value = newData;
//   },
//   { immediate: true }
// );

type User = {
  id: number;
  username: string;
  fullName: string;
  Type: string;
};
const columns = [
  { key: "number", label: "Number" },
  { key: "date", label: "Date" },
  { key: "account", label: "Account" },
  { key: "type", label: "Type" },
  { key: "amount", label: "Amount" },
  { key: "description", label: "Description" },
  // { key: "dr", label: "Dr." },
  // { key: "cr", label: "CR." },
  // { key: "balance", label: "Balance" },
  // { key: "actions", label: "Actions" },
];

function handleClick(row: { id: number }) {
  alert("clicked" + row);
}

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
// function openModal() {
//   count.value += 1;
//   modal.open(FormRole, {
//     count: count.value,
//     onSuccess() {
//       toast.add({
//         title: "Success !",
//         id: "modal-success",
//       });
//     },
//   });
// }
</script>

<template>
  <!-- <UButton label="Add Expense" @click="openModal" /> -->
  <div class="flex px-3 py-3.5 border-b border-gray-200 dark:border-gray-700">
    <UInput v-model="q" placeholder="Search" />
  </div>
  <UTable :rows="dataTableList" :columns="columns">
    <template #actions-data="{ row }">
      <UDropdown :items="items(row)">
        <UButton color="gray" variant="ghost" icon="i-heroicons-ellipsis-horizontal-20-solid" />
      </UDropdown>
    </template>
    <template #date-data="{ row }">
      <p>{{ row.date.split("T")[0] }}</p>
    </template>
    <template #type-data="{ row }">
      <p class="capitalize">{{ row.type }}</p>
    </template>
    <template #amount-data="{ row }">
      <p>{{ formatIDR(row.amount) }}</p>
    </template>
    <template #balance-data="{ row }">
      <p>{{ formatIDR(row?.account?.saldo ?? 0) }}</p>
    </template>
  </UTable>

  <div class="flex justify-end px-3 py-3.5 border-t border-gray-200 dark:border-gray-700">
    <UPagination v-model="page" :page-count="pageCount" :total="dataTableList.length" />
  </div>
</template>
