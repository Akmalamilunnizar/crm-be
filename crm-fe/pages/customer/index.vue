<script setup lang="ts">
import { formatIDR } from "@/helper/currency";
import { dashboardCustomerApi } from "../../api/customer/dashboard";
import { formatDateToYMD } from "@/helper/date";
import { use } from "echarts/core";

interface Customer {
  id: string;
  address: string;
  area_id: string;
  card_identition: string;
  company_id: string;
  email: string;
  gender: string;
  product_id: string;
  ip_static: string;
  job: string;
  latitude: number;
  longitude: number;
  mac_address: string;
  name: string;
  no_identition: number;
  phone: string;
  type_of_service: string;
  created_at: string;
  updated_at: string;
  installation_date: string;
  next_payment_date: string;
}

interface Product {
  id: string;
  name: string;
  price: number;
  description: string;
  createdAt: string;
  updatedAt: string;
}

interface Invoice {
  number: String;
  id: string;
  amount: number;
  customer_id: string;
  link: string;
  status: "unpaid" | "paid" | "pending"; // sesuaikan jika status lain ada
  created_at: string;
  updated_at: string;
}

interface DashboardResponse {
  customer: Customer;
  product: Product;
  invoice: Invoice[];
}

interface Column {
  key: keyof Invoice;
  label: string;
}

const columns: Column[] = [
  { key: "number", label: "Number" },
  { key: "created_at", label: "Date" },
  { key: "amount", label: "Amount" },
  { key: "status", label: "Status" },
  // { key: 'plan', label: 'Plan' },
];

let dashboardResponse = ref<DashboardResponse>({
  customer: {} as Customer,
  product: {} as Product,
  invoice: [] as Invoice[],
});

function handlePayment() {
  dashboardCustomerApi()
    .createPaymentDashboard({
      invoice_id: dashboardResponse.value.invoice[0].id,
    })
    .then((response) => {
      window.open(`${response.data.redirect_url}`, "_blank");
    })
    .catch((error) => {});
}
async function getData() {
  await dashboardCustomerApi()
    .getCustomerDashboard()
    .then((response) => {
      console.log(response.data);
      response.data.invoice.forEach((invoice: any) => {
        invoice.number = response.data.invoice.indexOf(invoice) + 1;
      });
      dashboardResponse.value = { ...response.data };

      // customer.value = [...response.data]
    })
    .catch((err) => {
      useToast().add({
        title: err,
        color: "red",
      });
    });
}
await getData();

const router = useRouter();
const status =
  new Date(dashboardResponse.value.customer.next_payment_date) >= new Date()
    ? "active"
    : "non active";
const logout = async () => {
  await useAuthStore().logout();
  router.push("/login");
};
</script>

<template>
  <div class="min-h-screen bg-gray-100">
    <div
      class="flex items-center justify-between p-4 text-white bg-green-600 shadow-md"
    >
      <div class="flex space-x-4">
        <div>
          <p class="text-sm">Hello, Welcome</p>
          <p class="text-2xl font-bold">
            {{ dashboardResponse.customer.name }}
          </p>
        </div>
        <UButton @click="logout" color="white" variant="outline"
          >Logout</UButton
        >
      </div>

      <div class="text-2xl font-bold">Lilly Apps</div>
    </div>

    <div class="max-w-4xl p-4 mx-auto space-y-6">
      <div class="p-6 bg-white border border-gray-200 shadow-lg rounded-xl">
        <div class="mb-4 text-center">
          <h2 class="text-xl font-bold text-gray-800">Next Payment</h2>
          <p class="text-gray-600">
            {{ formatDateToYMD(dashboardResponse.customer.next_payment_date) }}
          </p>
        </div>
        <div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
          <div class="p-2">
            <p class="text-gray-600">Package Chosen</p>
            <p class="text-xl font-bold text-gray-800">
              {{ dashboardResponse.product.name }}
            </p>
          </div>
          <div class="p-2">
            <p class="text-gray-600">Package Price</p>
            <p class="text-xl font-bold text-gray-800">
              {{ formatIDR(dashboardResponse.product.price) }}
            </p>
          </div>
        </div>
      </div>

      <div class="flex flex-col items-center gap-4">
        <div
          class="flex flex-col items-center justify-center w-full max-w-md p-6 transition-all duration-300 rounded-lg"
          :class="status === 'active' ? 'bg-green-500' : 'bg-red-500'"
        >
          <p class="text-lg font-bold text-white">Internet Status</p>
          <p class="text-white capitalize">{{ status }}</p>
        </div>
        <UButton
          :disabled="
            dashboardResponse.invoice.length === 0 ||
            dashboardResponse.invoice[0].status === 'paid'
          "
          class="w-full max-w-md py-3 font-semibold text-white transition-colors bg-green-600 rounded-lg hover:bg-green-700"
          size="xl"
          @click="handlePayment"
        >
          Payment Now
        </UButton>
      </div>

      <div class="p-6 bg-white border border-gray-200 shadow-lg rounded-xl">
        <h2 class="mb-4 text-lg font-bold text-gray-800">Recent Invoices</h2>
        <UTable
          :rows="dashboardResponse.invoice"
          :columns="columns"
          class="w-full"
          :class="{
            'text-gray-800': true,
            '[&>tbody>tr:nth-child(odd)]:bg-gray-50': true,
            '[&>tbody>tr:hover]:bg-gray-100': true,
          }"
        >
          <template #created_at-data="{ row }">
            <p>{{ row.created_at.split("T")[0] }}</p>
          </template>
          <template #amount-data="{ row }">
            <p>{{ formatIDR(row.amount) }}</p>
          </template>
          <template #status-data="{ row }">
            <UBadge
              size="sm"
              class="capitalize"
              :color="
                row.status === 'paid'
                  ? 'green'
                  : row.status === 'unpaid'
                  ? 'red'
                  : 'yellow'
              "
              :label="row.status"
            />
          </template>
        </UTable>
      </div>
    </div>
  </div>
</template>
