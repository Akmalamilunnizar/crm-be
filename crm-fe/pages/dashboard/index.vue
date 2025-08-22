<script setup lang="ts">
import { dashboardAdminApi } from "@/api/admin/dashboard";
import { invoiceAdminApi } from "@/api/admin/invoice";
import { formatIDR } from "@/helper/currency";
import { formatDateToYMD } from "@/helper/date";


let invoices = ref<any[]>([]);
const latestDeposites = ref<any[]>([]);
const latestExpenses = ref<any[]>([]);
let totalIncome = ref<any>(0);
let totalExpenses = ref<any>(0);
let totalNetWorth = ref<any>(0);
let totalSales = ref<any>(0);
const optionCardCustomer = ref();
const optionCardPacketPopular = ref();
const optionCardArea = ref();
const optionCardReportCash = ref();


definePageMeta({
  layout: false,
});

const data = ref<{ y: number[]; x: string[]; label: string }[]>([]);
let cards = ref<{ name: string; total: number }[]>([]);
const CardList = [
  {
    name: "Customer",
    total: 0,
  },
  {
    name: "Packet Popular",
    total: 0,
  },
  {
    name: "Area",
    total: 0,
  },
  {
    name: "Report Cash",
    total: 0,
  },
];

for (const card of CardList) {
  const option = {
    title: {
      text: ''
    },
    tooltip: {},
    legend: {
      show: true,
      top: '5%'        // ✅ Pastikan legend tidak di luar viewport
    },
    xAxis: {
      data: [],
    },
    yAxis: {},
    series: [
      {
        name: '',
        type: '',
        smooth: true,
        label: { show: true },
        data: []
      }
    ]
  }
  if (card.name === "Customer") {
    try {
      const response = await dashboardAdminApi().cardCustomerDashboard();
      const graph = response.data.graph_customer;

      (option.title as any).text = card.name

      if (Array.isArray(option.series) && option.series[0]) {
        (option.series[0] as any).data = graph.map((item: any) => item.count);
        (option.series[0] as any).name = card.name;
        (option.series[0] as any).type = "line";
      }

      (option.xAxis as any).data =
        graph.map((item: any) => {
          const date = new Date(item.date);
          return date.toLocaleDateString('en-US', { month: 'short', day: '2-digit' }
          );
        })

      optionCardCustomer.value = option


    } catch (err: any) {
      useToast().add({
        title: err.message || "Failed to load customer data",
        color: "red",
      });
    }
  }

  if (card.name === "Packet Popular") {
    try {
      const response = await dashboardAdminApi().cardPacketPopularDashboard();
      const graph = response.data.graph_packet_popular;

      (option.title as any).text = card.name

      if (Array.isArray(option.series) && option.series[0]) {
        (option.series[0] as any).data = graph.map((item: any) => item.count);
        (option.series[0] as any).name = card.name;
        (option.series[0] as any).type = "bar";
      }

      (option.xAxis as any).data =
        graph.map((item: any) => {
          return item.name;
        })

      optionCardPacketPopular.value = option


    } catch (err: any) {
      useToast().add({
        title: err.message || "Failed to load customer data",
        color: "red",
      });
    }
  }

  if (card.name === "Area") {
    try {
      const response = await dashboardAdminApi().cardAreaPopularDashboard();
      const graph = response.data.graph_area_popular;

      (option.title as any).text = card.name

      if (Array.isArray(option.series) && option.series[0]) {
        (option.series[0] as any).data = graph.map((item: any) => item.count);
        (option.series[0] as any).name = card.name;
        (option.series[0] as any).type = "bar";
      }

      (option.xAxis as any).data =
        graph.map((item: any) => {
          return item.name_city+ " - " + item.name_subdistrict + " - " + item.name_village;
        })

      optionCardArea.value = option


    } catch (err: any) {
      useToast().add({
        title: err.message || "Failed to load customer data",
        color: "red",
      });
    }
  }

  if (card.name === "Report Cash") {
    try {
      const response = await dashboardAdminApi().cardReportCashDashboard();
      const graph = response.data.graph_report_cash;

      (option.title as any).text = card.name

      if (Array.isArray(option.series) && option.series[0]) {
        (option.series[0] as any).data = graph.map((item: any) => item.count);
        (option.series[0] as any).name = card.name;
        (option.series[0] as any).type = "bar";
      }

      (option.xAxis as any).data =
        graph.map((item: any) => {
          const date = new Date(item.date);
          return date.toLocaleDateString('en-US', { month: 'short', day: '2-digit' }
          );
        })

      optionCardReportCash.value = option


    } catch (err: any) {
      useToast().add({
        title: err.message || "Failed to load customer data",
        color: "red",
      });
    }
  }
}

const columns = [
  {
    key: "number",
    label: "Number",
  },
  {
    key: "customer.name",
    label: "Name",
  },
  {
    key: "created_at",
    label: "Date",
  },
  {
    key: "amount",
    label: "amount",
  }, {
    key: "status",
    label: "Status",
  },
  {
    key: "customer.product.name",
    label: "Plan",
  },
];
async function getData() {
  invoiceAdminApi()
    .getAllInvoices()
    .then((response) => {
      response.data
        .forEach((customer: any) => {
          customer.number = response.data.indexOf(customer) + 1;
        });

      invoices.value = [...response.data];
    })
    .catch((err) => {
      useToast().add({
        title: err,
        color: "red",
      });
    });

  transactionAdminApi()
    .getAllTransactions(null).then((response) => {
      latestDeposites.value = response.data.filter((item: any) => item.type_in_out == "debit").map((transaction: any,index: number) => {
        return {
          number: index + 1,
          type_cash: transaction.type_cash.split("_")[0] + " " + (transaction.type_cash.split("_")[1] ? transaction.type_cash.split("_")[1] : ""),
          amount: transaction.amount,
          description: transaction.description
        }
      }).slice((1 - 1) * 5, (1) * 5);
      latestExpenses.value = response.data.filter((item: any) => item.type_in_out == "credit").map((transaction: any,index: number) => {
        return {
          number: index + 1,
          type_cash: transaction.type_cash.split("_")[0] + " " + (transaction.type_cash.split("_")[1] ? transaction.type_cash.split("_")[1] : ""),
          amount: transaction.amount,
          description: transaction.description
        }
      }).slice((1 - 1) * 5, (1) * 5);
    }).catch((err) => {
      useToast().add({
        title: err,
        color: "red",
      });
    });

  dashboardAdminApi()
    .totalIncomeDashboard()
    .then((response) => {
      totalIncome.value = response.data.total_income;
    })
    .catch((err) => {
      useToast().add({
        title: err,
        color: "red",
      });
    });

  dashboardAdminApi()
    .totalExpensesDashboard()
    .then((response) => {
      totalExpenses.value = response.data.total_expenses;
    })
    .catch((err) => {
      useToast().add({
        title: err,
        color: "red",
      });
    });

  dashboardAdminApi()
    .totalNetWorthDashboard()
    .then((response) => {
      totalNetWorth.value = response.data.total_net_worth;
    })
    .catch((err) => {
      useToast().add({
        title: err,
        color: "red",
      });
    });

  dashboardAdminApi()
    .totalSalesDashboard()
    .then((response) => {
      totalSales.value = response.data.total_sales;
    })
    .catch((err) => {
      useToast().add({
        title: err,
        color: "red",
      });
    });


}
getData();


import { useLoading } from '@/composables/useLoading'
import { transactionAdminApi } from "@/api/admin/transaction";
import { number } from "yup";

const { show, hide } = useLoading()

onMounted(async () => {
  show()
})


</script>

<template>
  <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6 py-6">

    <div
      class="w-full p-6 bg-gradient-to-br from-green-400 via-emerald-500 to-teal-500 text-white rounded-2xl shadow-xl transition-transform hover:scale-[1.03] duration-300">
      <div class="flex items-center justify-between mb-6">
        <h1 class="text-base font-medium uppercase tracking-wider opacity-90">Total Income</h1>
        <span class="text-lg font-semibold">$</span>
      </div>
      <div class="text-center">
        <h1 class="text-3xl font-bold">{{ formatIDR(totalIncome) }}</h1>
      </div>
    </div>

    <div
      class="w-full p-6 bg-gradient-to-br from-green-400 via-emerald-500 to-teal-500 text-white rounded-2xl shadow-xl transition-transform hover:scale-[1.03] duration-300">
      <div class="flex items-center justify-between mb-6">
        <h1 class="text-base font-medium uppercase tracking-wider opacity-90">Total Expenses</h1>
        <span class="text-lg font-semibold">$</span>
      </div>
      <div class="text-center">
        <h1 class="text-3xl font-bold">{{ formatIDR(totalExpenses) }}</h1>
      </div>
    </div>

    <div
      class="w-full p-6 bg-gradient-to-br from-green-400 via-emerald-500 to-teal-500 text-white rounded-2xl shadow-xl transition-transform hover:scale-[1.03] duration-300">
      <div class="flex items-center justify-between mb-6">
        <h1 class="text-base font-medium uppercase tracking-wider opacity-90">Net Worth</h1>
        <span class="text-lg font-semibold">$</span>
      </div>
      <div class="text-center">
        <h1 class="text-3xl font-bold">{{ formatIDR(totalNetWorth) }}</h1>
      </div>
    </div>

    <div
      class="w-full p-6 bg-gradient-to-br from-green-400 via-emerald-500 to-teal-500 text-white rounded-2xl shadow-xl transition-transform hover:scale-[1.03] duration-300">
      <div class="flex items-center justify-between mb-6">
        <h1 class="text-base font-medium uppercase tracking-wider opacity-90">Sales</h1>
      </div>
      <div class="text-center">
        <h1 class="text-3xl font-bold">{{ totalSales }}</h1>
      </div>
    </div>
  </div>

  <!-- <div class="grid gap-6 md:grid-cols-4 sm:grid-cols-2 mb-10">
    <div v-for="(card, index) in cards" :key="index">
      <CardComponent :dataChart="data[index]" :dataCard="card" />
    </div>
  </div> -->
  <div class="grid gap-6 md:grid-cols-2 sm:grid-cols-1 mb-10">
    <VChart :option="optionCardCustomer" autoresize style="height: 400px;" />
    <VChart :option="optionCardPacketPopular" autoresize style="height: 400px;" />
    <VChart :option="optionCardArea" autoresize style="height: 400px;" />
    <VChart :option="optionCardReportCash" autoresize style="height: 400px;" />
  </div>

  <div class="grid gap-6 md:grid-cols-2 sm:grid-cols-1 mb-10">
    <div class="p-6 bg-white border border-slate-200 rounded-2xl shadow-lg">
      <h1 class="text-xl font-semibold text-slate-800 mb-4">Latest Deposites</h1>
      <UTable v-if="latestDeposites.length > 0" :rows="latestDeposites" :page-size="5"></UTable>
    </div>
    <div class="p-6 bg-white border border-slate-200 rounded-2xl shadow-lg">
      <h1 class="text-xl font-semibold text-slate-800 mb-4">Latest Expenses</h1>
      <UTable v-if="latestExpenses.length > 0" :rows="latestExpenses" :page-size="5"></UTable>
    </div>
  </div>


  <div class="p-6 bg-white border border-slate-200 rounded-2xl shadow-lg">
    <h1 class="text-xl font-semibold text-slate-800 mb-4">Recent Invoices</h1>
    <UTable :columns="columns" :rows="invoices" :page-size="5">
      <template #amount-data="{ row }">
        <p class=" font-medium">{{ formatIDR(row.amount) }}</p>
      </template>
      <template #created_at-data="{ row }">
        <p class="text-slate-500">{{ formatDateToYMD(row.created_at) }}</p>
      </template>
    </UTable>
  </div>
</template>
