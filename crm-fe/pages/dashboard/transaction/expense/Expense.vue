<script setup lang="ts">
import { formatIDR } from "@/helper/currency";
import FormExpense from "./FormExpense.vue";
import { transactionAdminApi } from "@/api/admin/transaction";

const props = defineProps<{
  data: any[];
  onRefresh: () => Promise<void>;
  type?: any;
}>();
const rows = ref(props.data);

watch(
  () => props.data,
  (newData) => {
    rows.value = newData;
  },
  { immediate: true }
);

type User = {
  id: number;
  username: string;
  fullName: string;
  Type: string;
};
const columns = [
  { key: "number", label: "Number" },
  { key: "description", label: "Description" },
  { key: "amount", label: "Amount" },
  { key: "category", label: "Category" },
  { key: "method", label: "Method" },
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
const isOpen = ref(false);

function openAddExpenseModal() {
  modal.open(FormExpense, {
    type: props.type,
    onSuccess: handleSubmitExpense,
  });
}

function openEditExpenseModal(transactionId: string) {
  modal.open(FormExpense, {
    id: transactionId,
    type: props.type,
    onSuccess: handleSubmitExpense,
  });
}

async function handleSubmitExpense() {
  toast.add({
    title: "Success!",
    id: "modal-success",
  });

  // await fetchAllExpenses();
  modal.close();
  isOpen.value = false;
}

async function deleteExpense(transactionId: string) {
  const confirmed = window.confirm(
    "Are you sure you want to delete this company?"
  );
  if (!confirmed) return;

  try {
    await transactionAdminApi().deleteTransaction(transactionId);
    toast.add({
      title: "Success!",
      id: "modal-success",
    });
    await props.onRefresh();
    // await fetchAllExpenses();
    // props.function;
  } catch (error) {
    console.error("Error deleting company:", error);
  }
}
</script>

<template>
  <UButton label="Add Expense" @click="openAddExpenseModal" />
  <div class="flex px-3 py-3.5 border-b border-gray-200 dark:border-gray-700">
    <UInput v-model="q" placeholder="Search" />
  </div>
  <UTable :rows="rows" :columns="columns">
    <template #actions-data="{ row }">
      <UDropdown :items="items(row)">
        <UButton color="gray" variant="ghost" icon="i-heroicons-ellipsis-horizontal-20-solid" />
      </UDropdown>
    </template>
    <template #amount-data="{ row }">
      <p>{{ formatIDR(row.amount) }}</p>
    </template>
  </UTable>

  <div class="flex justify-end px-3 py-3.5 border-t border-gray-200 dark:border-gray-700">
    <!-- <UPagination
            v-model="page"
            :page-count="pageCount"
            :total="100"
          /> -->
  </div>
</template>
