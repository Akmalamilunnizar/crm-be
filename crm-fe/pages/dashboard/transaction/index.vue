<script setup lang="ts">
import { transactionAdminApi } from "@/api/admin/transaction";
import BalanceSheet from "./balance-sheet/BalanceSheet.vue";
import Deposit from "./deposit/Deposit.vue";
import Expense from "./expense/Expense.vue";
import Tranfer from "./tranfer/Tranfer.vue";
import ViewTransaction from "./view-transaction/ViewTransaction.vue";
import { accountAdminApi } from "@/api/admin/account";
const activeTab = ref(0);

let transaction: any[] = [];
let accounts: any[] = [];
const q = ref("");
const isLoading = ref(true); // Flag untuk menandakan apakah data sedang dimuat

async function fetchAllTransaction(params: any) {
  await transactionAdminApi()
    .getAllTransactions(params)
    .then((response) => {
      response.data.forEach((company: any) => {
        company.number = response.data.indexOf(company) + 1;
        company.logo_url =
          company.logo_url || "https://via.placeholder.com/150";
      });
      transaction = [...response.data];
      q.value = "changed";
      q.value = "";
    })
    .catch((error) => {
      console.error("Error fetching companies:", error);
    })
    .finally(() => {
      isLoading.value = false;
    });
}
async function fetchAllAccount() {
  await accountAdminApi()
    .getAllAccount()
    .then((response) => {
      response.data.forEach((company: any) => {
        company.number = response.data.indexOf(company) + 1;
        company.logo_url =
          company.logo_url || "https://via.placeholder.com/150";
      });
      accounts = [...response.data];
      q.value = "changed";
      q.value = "";
    })
    .catch((error) => {
      console.error("Error fetching companies:", error);
    })
    .finally(() => {
      console.log("succss");
      isLoading.value = false;
    });
}
const queries: any = {
  0: { type: "debit", type_cash: "cash_flow" },
  1: { type: "credit", type_cash: "cash_flow" },
  // 2: { type_cash: "cash_flow" },
  2: {},
  3: {},
};

watch(
  activeTab,
  (tab) => {
    const query = queries[tab]; // aman karena pakai Record<string, {type: string}>
    if (tab != 3) {
      isLoading.value = true; // Menandakan data sedang dimuat
      fetchAllTransaction(query);
    } else {
      isLoading.value = true; // Menandakan data sedang dimuat
      fetchAllAccount();
    }
  },
  { immediate: true }
);

const tab_items = [
  {
    label: "Deposit",
    value: "deposit",
  },
  {
    label: "Expense",
    value: "expense",
  },
  // {
  //   label: "Tranfer",
  //   value: "tranfer",
  // },
  {
    label: "View Transaction",
    value: "view-transaction",
  },
  {
    label: "Balance Sheet",
    value: "balance-sheet",
  },
];
</script>

<template>
  <UTabs :items="tab_items" class="w-full" v-model="activeTab">
    <template #item="{ item }">
      <div v-if="item.value == 'deposit'">
        <Deposit
          :type="queries[activeTab]"
          v-if="!isLoading"
          :data="transaction"
          :onRefresh="() => fetchAllTransaction(queries[activeTab])"
        />
      </div>
      <div v-if="item.value == 'expense'">
        <Expense
          :type="queries[activeTab]"
          :data="transaction"
          v-if="!isLoading"
          :onRefresh="() => fetchAllTransaction(queries[activeTab])"
        />
      </div>
      <!-- <div v-if="item.value == 'tranfer'">
        <Tranfer
          :data="transaction"
          v-if="!isLoading"
          :onRefresh="() => fetchAllTransaction(queries[activeTab])"
        />
      </div> -->
      <div v-if="item.value == 'view-transaction'">
        <ViewTransaction
          :data="transaction"
          v-if="!isLoading"
          :onRefresh="() => fetchAllTransaction(queries[activeTab])"
        />
      </div>
      <div v-if="item.value == 'balance-sheet'">
        <BalanceSheet
          v-if="!isLoading"
          :data="accounts"
          :onRefresh="() => fetchAllAccount()"
        />
      </div>
      <!-- <div v-if="item.value == 'report-customer'">
        <CustomerInstallation />

      </div> -->
    </template>
  </UTabs>
</template>
