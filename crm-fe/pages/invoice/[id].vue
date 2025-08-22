<script setup lang="ts">
import { ref } from "vue";
import html2canvas from "html2canvas";
import { invoiceAdminApi } from "@/api/admin/invoice";
import { formatDateToYMD } from "@/helper/date";
import { formatIDR } from "@/helper/currency";
const props = defineProps({
  invoice: {
    type: Object,
    required: true,
    default: () => ({
      id: "",
      status: "Paid",
      date: "",
      dueDate: "",
      customer: {
        name: "",
        attention: "",
        address: "",
        city: "",
        province: "",
        postalCode: "",
        country: "",
        phone: "",
        email: "",
      },
      items: [],
      subtotal: 0,
      tax: 0,
      total: 0,
      totalPaid: 0,
      amountDue: 0,
      transactions: [],
    }),
  },
});

const isLoading = ref(true);
let invoiceDetail = ref<any>({});

const route = useRoute();
const id = route.params.id as string;
async function getData() {
  isLoading.value = true;
  invoiceAdminApi()
    .getInvoice(id)
    .then((response) => {
      invoiceDetail.value = response.data;
    })
    .catch((err) => {
      useToast().add({
        title: err,
        color: "red",
      });
    }).finally(() => {
      isLoading.value = false;
    });
}
await getData();

function format(amount: number) {
  return amount.toLocaleString("id-ID", {
    minimumFractionDigits: 2,
    maximumFractionDigits: 2,
  });
}
const pdfContentRef = ref(null);
const invoiceNumber = computed(() => `#${props.invoice.id}`);
const generatePDF = async () => {
  try {
    const element = pdfContentRef.value;
    if (!element) {
      console.error("PDF content element not found");
      return;
    }

    // Render HTML to canvas
    const canvas = await html2canvas(element, { scale: 2 });

    // Dynamically import jsPDF (client-side only)
    const { jsPDF } = await import("jspdf");
    const pdf = new jsPDF({
      orientation: "portrait",
      unit: "mm",
      format: "a4",
    });

    const imgData = canvas.toDataURL("image/png");
    const imgProps = pdf.getImageProperties(imgData);
    const pdfWidth = 210; // A4 width in mm
    const pdfHeight = (imgProps.height * pdfWidth) / imgProps.width;

    pdf.addImage(imgData, "PNG", 0, 0, pdfWidth, pdfHeight);
    pdf.save("neraca_saldo.pdf");
  } catch (error) {
    console.error("Error generating PDF:", error);
  }
};
</script>

<template>
  <div v-if="!isLoading">
    <div class="p-8 text-gray-800 bg-white max-w-5xl mx-auto">
      <div ref="pdfContentRef" class="p-8 bg-grey-500">
        <client-only>
          <div class="flex justify-between items-start mb-6">
            <div>
              <h2 class="text-2xl font-semibold">INVOICE</h2>
              <p class="text-sm text-gray-500">{{ invoiceDetail.id }}</p>
              <div class="mt-2">
                <span class="bg-green-200 text-green-800 px-3 py-1 rounded capitalize">{{ invoiceDetail.status }}</span>
              </div>
              <div class="mt-4 text-sm">
                <p class="font-bold">Invoiced To</p>
                <p>{{ invoiceDetail.customer.name }}</p>
                <p>Phone: {{ invoiceDetail.customer.phone ?? "-" }}</p>
                <p>Email: {{ invoiceDetail.customer.email ?? "-" }}</p>
              </div>
            </div>
            <div class="text-right text-sm">
              <p class="font-bold">2024 ©Lilly Enterprise Billing</p>
              <p>Lilly Networks (PT JRNusa)</p>
              <p>Jl. Pratu Herman 34, Turen</p>
              <p>Kab. Malang Jawa Timur 65175</p>
              <p class="mt-4">
                Invoice Date: {{ formatDateToYMD(invoiceDetail.created_at) }}
              </p>
              <p>Due Date: {{ formatDateToYMD(invoiceDetail.created_at) }}</p>
              <p class="mt-4 font-bold">
                Invoice Total: Rp
                {{ props.invoice.total.toLocaleString("id-ID") }}
              </p>
              <p>
                Total Paid: Rp
                {{ props.invoice.totalPaid.toLocaleString("id-ID") }}
              </p>
              <!-- <p>
              Amount Due: Rp
              {{ props.invoice.amountDue.toLocaleString("id-ID") }}
            </p> -->
            </div>
          </div>

          <div class="table-scroll-container">
            <div class="table-scroll-content">
              <table class="w-full text-sm mb-6 border-t border-b border-gray-300 text-gray-900">
                <thead class="bg-gray-100">
                  <tr class="text-left uppercase text-xs tracking-wide text-gray-800">
                    <th class="py-2">#</th>
                    <th class="py-2">Item</th>
                    <th class="py-2">Price</th>
                    <th class="py-2">Qty</th>
                    <th class="py-2">Total</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="(item, i) in invoiceDetail.invoice_items" :key="i" class="odd:bg-white even:bg-gray-50">
                    <td class="py-2">{{ i + 1 }}</td>
                    <td class="py-2">{{ item.name }}</td>
                    <td class="py-2">{{ formatIDR(item.price) }}</td>
                    <td class="py-2">{{ item.qty }}</td>
                    <td class="py-2">{{ formatIDR(item.total) }}</td>
                  </tr>
                </tbody>
              </table>
            </div>
            <div class="table-scroll-footer">
              <span class="scroll-hint">↔ Scroll horizontally to see more columns | ↕ Scroll vertically for more rows</span>
            </div>
          </div>

          <div class="flex justify-end mb-6">
            <div class="text-sm w-1/3">
              <!-- <div class="flex justify-between">
              <span>Sub Total</span>
              <span
                >Rp {{ props.invoice.subtotal.toLocaleString("id-ID") }}</span
              >
            </div>
            <div class="flex justify-between">
              <span>TAX</span>
              <span>Rp {{ props.invoice.tax.toLocaleString("id-ID") }}</span>
            </div> -->
              <div class="flex justify-between font-bold">
                <span>Total</span>
                <span>Rp {{ props.invoice.total.toLocaleString("id-ID") }}</span>
              </div>
              <div class="flex justify-between">
                <span>Total Paid</span>
                <span>Rp {{ props.invoice.totalPaid.toLocaleString("id-ID") }}</span>
              </div>
              <!-- <div class="flex justify-between">
              <span>Amount Due</span>
              <span
                >Rp {{ props.invoice.amountDue.toLocaleString("id-ID") }}</span
              >
            </div> -->
            </div>
          </div>

          <!-- <h3 class="text-sm font-bold mb-2">Related Transactions</h3>
        <table class="w-full text-sm border-t border-gray-300">
          <thead>
            <tr class="text-left">
              <th class="py-2">Date</th>
              <th class="py-2">Account</th>
              <th class="py-2">Amount</th>
              <th class="py-2">Description</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(tx, i) in props.invoice.transactions" :key="i">
              <td class="py-2">{{ tx.date }}</td>
              <td class="py-2">{{ tx.account }}</td>
              <td class="py-2">Rp {{ tx.amount.toLocaleString("id-ID") }}</td>
              <td class="py-2">{{ tx.description }}</td>
            </tr>
          </tbody>
        </table> -->
        </client-only>
      </div>
      <button @click="generatePDF" class="px-4 py-2 mt-4 text-white bg-green-500 rounded">
        Download PDF
      </button>
    </div>
  </div>
</template>

<style scoped>
@media print {
  body * {
    visibility: hidden;
  }

  .print-area,
  .print-area * {
    visibility: visible;
  }

  .print-area {
    position: absolute;
    left: 0;
    top: 0;
  }
}
</style>
