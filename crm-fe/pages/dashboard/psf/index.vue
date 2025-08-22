<template>
  <div class="container p-6 mx-auto">
    <h1 class="mb-4 text-2xl font-bold text-center">Generate PDF</h1>
    <div class="flex items-center justify-center">

      <client-only>
        <div class="bg-white p-6 w-[800px] text-sm text-black" ref="pdfContentRef">
          <!-- Header -->
          <div class="p-3 text-white bg-green-800">
            <h1 class="font-bold">Lilly Apps</h1>
            <p class="text-xs">Saldo Pemasukan dan Pengeluaran</p>
            <p class="text-xs">Per Akhir <span class="float-right">{{ new Date() }}</span></p>
          </div>

          <!-- Section: Aktiva -->
          <div class="mt-4">
            <h2 class="font-semibold text-green-700">Aktiva Lancar</h2>
            <table class="w-full mt-1">
              <tbody>
                <tr v-for="item in aktivaLancar" :key="item.name">
                  <td class="pl-2">{{ item.name }}</td>
                  <td class="pr-2 text-right">{{ format(item.amount) }}</td>
                </tr>
                <tr class="font-bold text-green-700 border-t">
                  <td class="pl-2">Total Aktiva Lancar</td>
                  <td class="pr-2 text-right">{{ format(total(aktivaLancar)) }}</td>
                </tr>
              </tbody>
            </table>

            <h2 class="mt-4 font-semibold text-green-700">Aktiva Tidak Lancar</h2>
            <table class="w-full mt-1">
              <tbody>
                <tr v-for="item in aktivaTidakLancar" :key="item.name">
                  <td class="pl-2">{{ item.name }}</td>
                  <td class="pr-2 text-right">{{ format(item.amount) }}</td>
                </tr>
                <tr class="font-bold text-green-700 border-t">
                  <td class="pl-2">Total Aktiva Tidak Lancar</td>
                  <td class="pr-2 text-right">{{ format(total(aktivaTidakLancar)) }}</td>
                </tr>
              </tbody>
            </table>

            <div class="py-1 mt-3 font-bold text-right text-green-900 border-t border-b">
              Total Assets: {{ format(total(aktivaLancar) + total(aktivaTidakLancar)) }}
            </div>
          </div>

          <!-- Section: Kewajiban -->
          <div class="mt-6">
            <h2 class="font-semibold text-green-700">Kewajiban Lancar</h2>
            <table class="w-full mt-1">
              <tbody>
                <tr v-for="item in kewajibanLancar" :key="item.name">
                  <td class="pl-2">{{ item.name }}</td>
                  <td class="pr-2 text-right">{{ format(item.amount) }}</td>
                </tr>
                <tr class="font-bold text-green-700 border-t">
                  <td class="pl-2">Total Kewajiban Lancar</td>
                  <td class="pr-2 text-right">{{ format(total(kewajibanLancar)) }}</td>
                </tr>
              </tbody>
            </table>

            <h2 class="mt-4 font-semibold text-green-700">Kewajiban Tidak Lancar</h2>
            <table class="w-full mt-1">
              <tbody>
                <tr v-for="item in kewajibanTidakLancar" :key="item.name">
                  <td class="pl-2">{{ item.name }}</td>
                  <td class="pr-2 text-right">{{ format(item.amount) }}</td>
                </tr>
                <tr class="font-bold text-green-700 border-t">
                  <td class="pl-2">Total Kewajiban Tidak Lancar</td>
                  <td class="pr-2 text-right">{{ format(total(kewajibanTidakLancar)) }}</td>
                </tr>
              </tbody>
            </table>

            <!-- Footer -->
            <div class="mt-4 text-xs text-left">
              17 Mei 2025, 07:55 WIB
            </div>
          </div>
        </div>
      </client-only>
    </div>
    
  </div>
  <button @click="generatePDF" class="px-4 py-2 mt-4 text-white bg-green-500 rounded">Download PDF</button>
</template>

<script setup>
import { ref } from 'vue';
import html2canvas from 'html2canvas';

// Sample data
const aktivaLancar = ref([
  { name: 'Kas', amount: 5000000 },
  { name: 'Piutang', amount: 3000000 },
]);
const aktivaTidakLancar = ref([{ name: 'Peralatan', amount: 10000000 }]);
const kewajibanLancar = ref([{ name: 'Utang Dagang', amount: 4000000 }]);
const kewajibanTidakLancar = ref([{ name: 'Utang Jangka Panjang', amount: 8000000 }]);
const pdfContentRef = ref(null);

const format = (amount) => amount.toLocaleString('id-ID');
const total = (items) => items.reduce((sum, item) => sum + item.amount, 0);

const generatePDF = async () => {
  try {
    const element = pdfContentRef.value;
    if (!element) {
      console.error('PDF content element not found');
      return;
    }

    // Render HTML to canvas
    const canvas = await html2canvas(element, { scale: 2 });

    // Dynamically import jsPDF (client-side only)
    const { jsPDF } = await import('jspdf');
    const pdf = new jsPDF({
      orientation: 'portrait',
      unit: 'mm',
      format: 'a4',
    });

    const imgData = canvas.toDataURL('image/png');
    const imgProps = pdf.getImageProperties(imgData);
    const pdfWidth = 210; // A4 width in mm
    const pdfHeight = (imgProps.height * pdfWidth) / imgProps.width;

    pdf.addImage(imgData, 'PNG', 0, 0, pdfWidth, pdfHeight);
    pdf.save('neraca_saldo.pdf');
  } catch (error) {
    console.error('Error generating PDF:', error);
  }
};
</script>

<style scoped>
/* Add any additional custom styles if needed */
</style>