<script setup lang="ts">
import { accountAdminApi } from "@/api/admin/account";
import FormAddComponent from "./FormAddComponent.vue";
import { transactionAdminApi } from "@/api/admin/transaction";
import { formatIDR } from "@/helper/currency";
import { formatDateToYMD } from "@/helper/date";

let transaction: any[] = [];

const page = ref(1);
const pageCount = 5;

const q = ref("");

let peopleData = transaction;

// const filteredRows = computed(() => {
//     if (!q.value) {

//         return transaction.slice((page.value - 1) * pageCount, page.value * pageCount);
//     }

//     const newData = transaction.filter((transaction) => {
//         return Object.values(transaction).some((value) => {
//             // person with paginate
//             return String(value).toLowerCase().includes(q.value.toLowerCase());
//         });
//     });

//     return newData.slice((page.value - 1) * pageCount, page.value * pageCount);
// });
const columns = [
  {
    key: "number",
    label: "Number",
  },
  {
    key: "createdAt",
    label: "date",
  },
  {
    key: "name",
    label: "Name",
  },
  {
    key: "type_in_out",
    label: "Type",
  },
  {
    key: "amount",
    label: "Amount",
  },
  {
    key: "description",
    label: "Description",
  },
  {
    key: "Dr",
    label: "Debit",
  },
  {
    key: "Cr",
    label: "Credit",
  },
  {
    key: "actions",
    label: "Actions",
  },
];
const isOpen = ref(false);

const toast = useToast();
const modal = useModal();
function openModal() {
  isOpen.value = true;
  modal.open(FormAddComponent, {
    // isOpen: true,
    onSuccess() {
      // Emit success event
      console.log("success");
      modal.close();
      toast.add({
        title: "Success !",
        id: "modal-success",
        ui: { position: "top-right" },
      });
    },
  });
  console.log("open modals");
}

async function fetchAllAccount() {
  await transactionAdminApi()
    .getAllTransactions({})
    .then((response) => {
      response.data.forEach((transaction: any) => {
        transaction.number = response.data.indexOf(transaction) + 1;
        transaction.createdAt = formatDateToYMD(transaction.createdAt);
        transaction.type_in_out = transaction.type_in_out === "debit" ? "Debit" : "Credit";
        transaction.Dr = transaction.type_in_out !== "Debit" ? "-" : formatIDR(transaction.amount);
        transaction.Cr = transaction.type_in_out !== "Credit" ? "-" : formatIDR(transaction.amount);
        transaction.balance = formatIDR(transaction.account.saldo);

      });
      transaction = [...response.data];
      q.value = "changed";
      q.value = "";
    })
    .catch((err) => {
      useToast().add({
        title: err,
        color: "red",
      });
    });
}
await fetchAllAccount();
const rows = computed(() => {
    if (!q.value) {

        return transaction.slice((page.value - 1) * pageCount, page.value * pageCount);
    }

    const newData = transaction.filter((transaction) => {
        return Object.values(transaction).some((value) => {
            // person with paginate
            return String(value).toLowerCase().includes(q.value.toLowerCase());
        });
    });

    return newData.slice((page.value - 1) * pageCount, page.value * pageCount);
});
</script>

<template>
  <!-- <UButton label="Add Report Cash Flow" @click="openModal" /> -->
  <div class="flex px-3 py-3.5 border-b border-gray-200 dark:border-gray-700">
    <UInput v-model="q" placeholder="Search" />
    <UButton
      label="Generate PDF"
      class="ml-2"
      @click="navigateTo('/dashboard/psf')"
    ></UButton>
  </div>
  <UTable :rows="rows" :columns="columns">
    <template #amount-data="{ row }">
      <p>{{ formatIDR(row.amount) }}</p>
    </template>
    <template #created_at-data="{ row }">
      <p>{{ formatDateToYMD(row.amount) }}</p>
    </template>
    <template #name-data="{ row }">
      <p>{{ row.account.name }}</p>
    </template>
     <template #balance-data="{ row }">
      <p>{{ formatIDR(row.account.saldo) }}</p>
    </template>
  </UTable>

  <div
    class="flex justify-end px-3 py-3.5 border-t border-gray-200 dark:border-gray-700"
  >
    <UPagination
      v-model="page"
      :page-count="pageCount"
      :total="peopleData.length"
    />
  </div>
</template>
