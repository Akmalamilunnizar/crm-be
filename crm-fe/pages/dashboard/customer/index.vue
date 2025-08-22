<script setup lang="ts">
import FormCustomerInstallation from './FormCustomerInstallation.vue'
import FormAddComponent from './FormAddComponent.vue'
import { customerAdminApi } from '@/api/admin/customer'
let customer = ref<any[]>([])

type Customer = {
    id: string
    name: string
    email: string
    phone: string
    address: string
    area_code: string
    gmaps_link: string
    packet_internet: string
    ip_static: string
    mac_address: string
}


async function getData() {
    await customerAdminApi().getAllCustomers().then((response) => {
        response.data.forEach((customer: any) => {
            customer.number = response.data.indexOf(customer) + 1;
            customer.area_name = customer.area.name_city + "-" + customer.area.name_subdistrict + "-" + customer.area.name_village
            customer.product_name = customer.product.name
            customer.gmaps_link = "https://www.google.com/maps/place/" + customer.latitude + "," + customer.longitude
        })

        customer.value = [...response.data]
    }).catch((err) => {
        useToast().add({
            title: err,
            color: "red"
        })
    })
}

async function deleteData(id: string) {
    await customerAdminApi().deleteCustomer(id).then((response) => {
        getData()
        useToast().add({
            title: response.message,
        })
    }).catch((err) => {
        useToast().add({
            title: err,
            color: "red"
        })
    })
}

await getData();

const columns = [
    {
        key: 'number',
        label: 'Number'
    }, {
        key: 'name',
        label: 'Name'
    }, {
        key: 'email',
        label: 'Email'
    }, {
        key: 'phone',
        label: 'Phone'
    }, {
        key: 'address',
        label: 'Address'
    }, {
        key: 'area_name',
        label: 'Area Code'
    }, {
        key: 'product_name',
        label: 'Packet Internet'
    }, {
        key: 'ip_static',
        label: 'IP Static'
    }, {
        key: 'mac_address',
        label: 'Mac Address'
    }, {
        key: 'actions',
        label: 'Actions'
    }
]

const page = ref(1)
const pageCount = 5

const q = ref('')
const rows = computed(() => {
    if (!q.value) {
        return customer.value.slice((page.value - 1) * pageCount, (page.value) * pageCount)
    }
    const newData = customerData.value.filter((customer) => {
        return Object.values(customer).some((value) => {
            return String(value).toLowerCase().includes(q.value.toLowerCase())
        })
    })

    return newData.slice((page.value - 1) * pageCount, (page.value) * pageCount)
})


let customerData = customer

// const filteredRows = computed(() => {
//     if (!q.value) {
//         customerData = customer
//         return customer.value.slice((page.value - 1) * pageCount, (page.value) * pageCount)
//     }

//     const newData = customer.value.filter((person) => {
//         return Object.values(person).some((value) => {
//             // person with paginate
//             return String(value).toLowerCase().includes(q.value.toLowerCase())
//         })
//     })
//     customerData.value = newData
//     return newData.slice((page.value - 1) * pageCount, (page.value) * pageCount)
// })

const isOpen = ref(false)

const items = (row: Customer) => [
    [{
        label: 'Edit',
        icon: 'i-heroicons-pencil-square-20-solid',
        click: () => OpenModalAddCustomer(true, row)
    }], [{
        label: 'Add Report Installation',
        icon: 'i-heroicons-archive-box-20-solid',
        click: () => OpenModalReportInstallation(true, row)
    }, {
        label: 'View Maps',
        icon: 'i-heroicons-arrow-right-circle-20-solid',
        click: () => window.open(row.gmaps_link, '_blank')
    }], [{
        label: 'Delete',
        icon: 'i-heroicons-trash-20-solid',
        click: () => deleteData(row.id)
    }]
]

const toast = useToast()
const modal = useModal()

function OpenModalAddCustomer(isEdit: boolean, data: any) {
    modal.open(FormAddComponent, {
        isEdit,
        data,
        async onSuccess() {
            await getData()
            modal.close()
        }
    })
}

function OpenModalReportInstallation(isEdit: boolean, data: any) {
    modal.open(FormCustomerInstallation, {
        isEdit,
        data,
        async onSuccess() {
            await getData()
            toast.add({
                title: 'Success !',
                id: 'modal-success'
            })
            modal.close()
        }
    })
}
</script>


<template>

    <UButton label="Add Customer" @click="OpenModalAddCustomer(false, null)" />
    <div class="flex px-3 py-3.5 border-b border-gray-200 dark:border-gray-700">
        <UInput v-model="q" placeholder="Filter customer..." />
    </div>

    <div class="table-scroll-container">
        <div class="table-scroll-content">
            <UTable :rows="rows" :columns="columns" class="dashboard-table">

                <template #actions-data="{ row }">
                    <UDropdown :items="items(row)">
                        <UButton color="gray" variant="ghost" icon="i-heroicons-ellipsis-horizontal-20-solid" />
                    </UDropdown>
                </template>
            </UTable>
        </div>
        <div class="table-scroll-footer">
          <span class="scroll-hint">↔ Scroll horizontally to see more columns | ↕ Scroll vertically for more rows</span>
        </div>
    </div>

    <div class="flex justify-end px-3 py-3.5 border-t border-gray-200 dark:border-gray-700">
        <UPagination v-model="page" :page-count="pageCount" :total="customer.length" />
    </div>

</template>
