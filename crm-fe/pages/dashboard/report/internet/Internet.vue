<script setup lang="ts">
import { reportAdminApi } from '@/api/admin/report';
import FormInternet from './FormInternet.vue';


const dataList = ref([
])

const columns = [
  { key: "number", label: "Number" },
  { key: "customer", label: "Customer" },
  { key: "internet_packet", label: "Internet Packet" },
  { key: "date", label: "Date Installation" },
  // {
  //   key:"actions",
  //   label: "Actions",
  // }
];

const toast = useToast()
const modal = useModal()


async function getDataReport() {
  reportAdminApi().getReportInternet().then((res) => {
    dataList.value = res.data.map((item:any) => {
      return {
        number: res.data.indexOf(item) + 1,
        customer: item.name,
        internet_packet: item.product,
        date: item.installation_date.split('T')[0],
      }
    })
  }).catch((err) => {
    toast.add({
      title: err.message,
    })
  })
}

getDataReport()

function openModal() {
  modal.open(FormInternet, {
    onSuccess() {
      toast.add({
        title: 'Success !',
        id: 'modal-success'
      })
    }
  })
}
</script>

<template>
        <div
          class="flex px-3 py-3.5 border-b border-gray-200 dark:border-gray-700"
        >
        </div>
        <UTable :rows="dataList" :columns="columns"/>

        <div
          class="flex justify-end px-3 py-3.5 border-t border-gray-200 dark:border-gray-700"
        >
          <!-- <UPagination
            v-model="page"
            :page-count="pageCount"
            :total="peopleData.length"
          /> -->
        </div>
</template>