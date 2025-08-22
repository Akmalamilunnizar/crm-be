<script setup lang="ts">
import { formatIDR } from "@/helper/currency";
const props = defineProps<{
  data: { name: string; amount: number }[]
}>()
const rows = ref(props.data);

watch(() => props.data, (newData) => {
  rows.value = newData;
}, { immediate: true });

type User = {
  id: number;
  username: string;
  fullName: string;
  Type: string;
};
const columns = [
  { key: "number", label: "Number" },
  { key: "account", label: "Account" },
  { key: "description", label: "Description" },
  { key: "amount", label: "Amount" },
  // {
  //   key: "actions",
  //   label: "Actions",
  // },
];
const q = ref("");
const page = 0;
const pageCount = 0;

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
  <UTable :rows="rows" :columns="columns">
    <template #actions-data="{ row }">
      <UDropdown :items="items(row)">
        <UButton
          color="gray"
          variant="ghost"
          icon="i-heroicons-ellipsis-horizontal-20-solid"
        />
      </UDropdown>
    </template>
    
     <template #account-data="{ row }">
      <p >{{row.account.name}}</p>
    </template>
     <template #amount-data="{ row }">
      <p >{{ formatIDR(row.amount) }}</p>
    </template>
  </UTable>

  <div
    class="flex justify-end px-3 py-3.5 border-t border-gray-200 dark:border-gray-700"
  >
    <!-- <UPagination
            v-model="page"
            :page-count="pageCount"
            :total="100"
          /> -->
  </div>
</template>
