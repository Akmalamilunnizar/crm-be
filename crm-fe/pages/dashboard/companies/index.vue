<script setup lang="ts">
import { companyAdminApi } from "@/api/admin/company";
import AddCompanyForm from "./AddCompanyForm.vue";
let people: any[] = [];

const page = ref(1);
const pageCount = 5;

const q = ref("");
async function fetchAllCompanies() {
    await companyAdminApi()
        .getAllCompanies()
        .then((response) => {
            response.data.forEach((company: any) => {
                company.number = response.data.indexOf(company) + 1;
                company.logo_url =
                    company.logo_url || "https://via.placeholder.com/150";
            });
            people = [...response.data];
            q.value = "changed";
            q.value = "";
        })
        .catch((error) => {
            console.error("Error fetching companies:", error);
        });
}

await fetchAllCompanies();

type Person = {
    id: number;
    name: string;
    email: string;
    phone: string;
    address: string;
    area_code: string;
    packet_internet: string;
    ip_static: string;
    mac_address: string;
};

const columns = [
    {
        key: "number",
        label: "No",
    },
    {
        key: "name",
        label: "Name",
    },
    {
        key: "email",
        label: "Email",
    },
    {
        key: "url",
        label: "Website",
    }, {
        key: "npwp",
        label: "NPWP",
    }, {
        key: "address",
        label: "Address",
    },
    {
        key: "phone",
        label: "Phone",
    },
    {
        key: "logo_url",
        label: "Image",
    },
    {
        key: "actions",
        label: "Actions",
    },
];


let peopleData = people;

const filteredRows = computed(() => {
    if (!q.value) {
        peopleData = people;
        return people.slice((page.value - 1) * pageCount, page.value * pageCount);
    }

    const newData = people.filter((person) => {
        return Object.values(person).some((value) => {
            return String(value).toLowerCase().includes(q.value.toLowerCase());
        });
    });
    peopleData = newData;
    return newData.slice((page.value - 1) * pageCount, page.value * pageCount);
});

const isOpen = ref(false);

function openAddCompanyModal() {
    modal.open(AddCompanyForm, {
        onSuccess: handleSubmitCompany,
    });
}

function openEditCompanyModal(companyId: string) {
    modal.open(AddCompanyForm, {
        isEdit: true,
        id: companyId,
        onSuccess: handleSubmitCompany,
    });
}

async function handleSubmitCompany() {
    toast.add({
        title: "Success!",
        id: "modal-success",
    });

    await fetchAllCompanies();
    modal.close();
    isOpen.value = false;
}

async function deleteCompany(companyId: string) {
    const confirmed = window.confirm("Are you sure you want to delete this company?");
    if (!confirmed) return;

    try {
        await companyAdminApi().deleteCompany(companyId);
        toast.add({
            title: "Success!",
            id: "modal-success",
        });

        await fetchAllCompanies();
    } catch (error) {
        console.error("Error deleting company:", error);
    }
}

const items = (row: Person) => [
    [
        {
            label: "Edit",
            icon: "i-heroicons-pencil-square-20-solid",
            click: () => openEditCompanyModal(row.id.toString()),
        },
    ],
    [
        {
            label: "Delete",
            icon: "i-heroicons-trash-20-solid",
            click: () => deleteCompany(row.id.toString()),
        },
    ],
];

const toast = useToast();
const modal = useModal();
</script>

<template>
    <UButton label="Add Companies" @click="() => openAddCompanyModal()" />
    <div class="flex px-3 py-3.5 border-b border-gray-200 dark:border-gray-700">
        <UInput v-model="q" placeholder="Search" />
    </div>

    <div class="table-scroll-container">
        <div class="table-scroll-content">
            <UTable :rows="filteredRows" :columns="columns" class="dashboard-table">
                <template #logo_url-data="{ row }">
                    <UAvatar :src="row.logo_url" size="xl" />
                </template>
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
        <UPagination v-model="page" :page-count="pageCount" :total="peopleData.length" />
    </div>

    <div>
        <UModal v-model="isOpen">
            <div class="p-4">
                <Placeholder class="h-48" />
                <AddCompanyForm />
            </div>
        </UModal>
    </div>
</template>
