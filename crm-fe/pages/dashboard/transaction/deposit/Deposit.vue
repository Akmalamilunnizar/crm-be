<script setup lang="ts">
import { transactionAdminApi } from "@/api/admin/transaction";
import FormDeposit from "./FormDeposit.vue";
import FormRole from "./FormDeposit.vue";
import { formatIDR } from "@/helper/currency";
const props = defineProps<{
  data: any[];
  onRefresh: () => Promise<void>;
  type?: any;
}>();
const q = ref("");
const page = ref(1);
const pageCount = 5;

const deposites = ref<any[]>(
  props.data
); // Use ref for reactivity
const depositesList = computed(() => {
  if (!q.value) {
    return deposites.value.slice((page.value - 1) * pageCount, page.value * pageCount);
  }

  const newData = deposites.value.filter((transaction) => {
    return Object.values(transaction).some((value) => {
      return String(value).toLowerCase().includes(q.value.toLowerCase());
    });
  });

  return newData.slice((page.value - 1) * pageCount, page.value * pageCount);
})

type Deposit = {
  id: string;
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
];

function handleClick(row: { id: number }) {
  alert("clicked" + row);
}

const toast = useToast();
const modal = useModal();
const count = ref(0);

const items = (row: Deposit) => [
  [
    {
      label: "Edit",
      icon: "i-heroicons-pencil-square-20-solid",
      click: () => openEditDepositModal(row.id.toString()),
    },
  ],
  [
    {
      label: "Delete",
      icon: "i-heroicons-trash-20-solid",
      click: () => deleteDeposit(row.id.toString()),
    },
  ],
];

const isOpen = ref(false);

function openAddDepositModal() {
  modal.open(FormDeposit, {
    type: props.type,
    onSuccess: handleSubmitDeposit,
  });
}

function openEditDepositModal(transactionId: string) {
  modal.open(FormDeposit, {
    id: transactionId,
    type: props.type,
    onSuccess: handleSubmitDeposit,
  });
}

async function handleSubmitDeposit() {
  toast.add({
    title: "Success!",
    id: "modal-success",
  });

  // await fetchAllDeposits();
  modal.close();
  isOpen.value = false;
}

async function deleteDeposit(transactionId: string) {
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
    // await fetchAllDeposits();
    // props.function;
  } catch (error) {
    console.error("Error deleting company:", error);
  }
}

</script>

<template>
  <UButton label="Add Deposit" @click="openAddDepositModal" />
  <div class="flex px-3 py-3.5 border-b border-gray-200 dark:border-gray-700">
    <UInput v-model="q" placeholder="Search" />
  </div>
  <UTable :rows="depositesList" :columns="columns">
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
    <UPagination v-model="page" :page-count="pageCount" :total="depositesList.length" />
  </div>
</template>