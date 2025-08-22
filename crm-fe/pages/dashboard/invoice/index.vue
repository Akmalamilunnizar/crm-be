<script setup lang="ts">
import { invoiceAdminApi } from "@/api/admin/invoice";
import FormAddComponent from "./FormAddInvoice.vue";
import * as currency from "@/helper/currency";
import type { UpdateStatusInvoiceRequest } from "@/types/requests/invoice";
import { WhatsappApi } from "@/api/admin/wa";
let customer = ref<any[]>([]);

const router = useRouter();
type Customer = {
  id: string;
  name: string;
  email: string;
  phone: string;
  address: string;
  area_code: string;
  gmaps_link: string;
  packet_internet: string;
  ip_static: string;
  mac_address: string;
};

async function getData() {
  invoiceAdminApi()
    .getAllInvoices()
    .then((response) => {
      response.data.forEach((customer: any) => {
        customer.number = response.data.indexOf(customer) + 1;
        customer.created_at = customer.created_at.split("T")[0];
      });

      customer.value = [...response.data];
    })
    .catch((err) => {
      useToast().add({
        title: err,
        color: "red",
      });
    });
}

async function updateStatus(id: string, status: string) {
  invoiceAdminApi()
    .updateStatusInvoice(id, { status })
    .then((response) => {
      useToast().add({
        title: response.message,
        color: "green",
      });
    })
    .catch((err) => {
      useToast().add({
        title: err,
        color: "red",
      });
    });
}

async function sendWhatsapp(number: string, id: string) {
  WhatsappApi()
    .sendWhatsapp({
      number,
      message:
        `berikut invoice yang harus anda bayarkan sekarang \n\nKami berikan Link untuk melihat invoice \n\nhttps://skripsi.rtrsite.com/invoice/${id} \n\nSilahkan menuju dashboard login customer kami https://skripsi.rtrsite.com/login \n\nTerimakasih`,
    })
    .then((response) => {
      useToast().add({
        title: response.message,
        color: "green",
      });
    })
    .catch((err) => {
      useToast().add({
        title: err,
        color: "red",
      });
    });
}

async function deleteData(id: string) {
  await invoiceAdminApi()
    .deleteInvoice(id)
    .then((response) => {
      getData();
      useToast().add({
        title: response.message,
      });
    })
    .catch((err) => {
      useToast().add({
        title: err,
        color: "red",
      });
    });
}

await getData();

const columns = [
  {
    key: "number",
    label: "Number",
  },
  {
    key: "customer.name",
    label: "Customer",
  },
  {
    key: "amount",
    label: "Amount",
  },
  {
    key: "status",
    label: "Status",
  },
  {
    key: "created_at",
    label: "Date",
  },
  {
    key: "actions",
    label: "Actions",
  },
];

const page = ref(1);
const pageCount = 5;

const rows = computed(() => {
  if (!q.value) {
    return customer.value.slice(
      (page.value - 1) * pageCount,
      page.value * pageCount
    );
  }

  const newData = customer.value.filter((transaction) => {
    return Object.values(transaction).some((value) => {
      return String(value).toLowerCase().includes(q.value.toLowerCase());
    });
  });

  return newData.slice((page.value - 1) * pageCount, page.value * pageCount);
});

const q = ref("");

const filteredRows = computed(() => {
    if (!q.value) {
        return customer.value.slice((page.value - 1) * pageCount, (page.value) * pageCount)
    }

    const newData = customer.value.filter((person) => {
        return Object.values(person).some((value) => {
            // person with paginate
            return String(value).toLowerCase().includes(q.value.toLowerCase())
        })
    })
    return newData.slice((page.value - 1) * pageCount, (page.value) * pageCount)
})

const items = (row: any) => [
  [
    {
      label: "Send Whatsapp",
      icon: "i-heroicons-chat-bubble-left-ellipsis-20-solid",
      click: () => sendWhatsapp(row.customer.phone, row.id),
    },
    {
      label: "Download PDF",
      icon: "i-heroicons-arrow-down-on-square-20-solid",
      click: () => navigateTo(`/invoice/${row.id}`),
    },
  ],
  [
    {
      label: "Delete",
      icon: "i-heroicons-trash-20-solid",
      click: () => deleteData(row.id),
    },
  ],
];

const toast = useToast();
const modal = useModal();

function OpenModalAddCustomer(isEdit: boolean, data: any) {
  console.log("Open Modal");
  modal.open(FormAddComponent, {
    isEdit,
    data,
    async onSuccess() {
      await getData();
      modal.close();
    },
  });
}
const sort = ref({
  column: "amount",
  direction: "asc" as "asc",
});
</script>

<template>
  <UButton label="Add Invoice" @click="OpenModalAddCustomer(false, null)" />
  <div class="flex px-3 py-3.5 border-b border-gray-200 dark:border-gray-700">
    <UInput v-model="q" placeholder="Filter customer..." />
  </div>

  <UTable :rows="filteredRows" :columns="columns">
    <template #actions-data="{ row }">
      <UDropdown :items="items(row)">
        <UButton
          color="gray"
          variant="ghost"
          icon="i-heroicons-ellipsis-horizontal-20-solid"
        />
      </UDropdown>
    </template>
    <template #amount-data="{ row }">
      <p>{{ currency.formatIDR(row.amount) }}</p>
    </template>
    <template #status-data="{ row }">
      <USelectMenu
        v-model="row.status"
        :options="[
          { label: 'Pending', value: 'pending' },
          { label: 'Paid', value: 'paid' },
          { label: 'Unpaid', value: 'unpaid' },
        ]"
        @change="updateStatus(row.id, row.status)"
        value-attribute="value"
        option-attribute="label"
      />
    </template>
  </UTable>

  <div
    class="flex justify-end px-3 py-3.5 border-t border-gray-200 dark:border-gray-700"
  >
    <UPagination
      v-model="page"
      :page-count="pageCount"
      :total="customer.length"
    />
  </div>
</template>
