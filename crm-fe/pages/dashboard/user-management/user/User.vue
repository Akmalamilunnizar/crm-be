<script setup lang="ts">
import { userManagementAdminApi } from "@/api/admin/user-management";
import AddUserForm from "./AddUserForm.vue";
let user: any[] = [];
const page = ref(1);
const pageCount = 5;

const q = ref("");

async function fetchAllUsers() {
  await userManagementAdminApi()
    .getAllUsers({ query: { role: "" } })
    .then((response) => {
      response.data.forEach((user: any) => {
        user.number = response.data.indexOf(user) + 1;
        user.avatar = user.logo_url || "https://via.placeholder.com/150";
        user.role = user.role.name;
      });
      user = [...response.data];
      q.value = "changed";
      q.value = "";
    })
    .catch((error) => {
      console.error("Error fetching user:", error);
    });
}


await fetchAllUsers();

type User = {
  id: number;
  username: string;
  fullName: string;
  Type: string;
};
const columns = [
  { key: "number", label: "Number" },
  { key: "avatar", label: "Avatar" },
  { key: "name", label: "Username" },
  { key: "email", label: "Full Name" },
  { key: "role", label: "Type" },
  {
    key: "actions",
    label: "Actions",
  },
];



let userData = user;
const filteredRows = computed(() => {
  if (!q.value) {
    userData = user;
    return user.slice((page.value - 1) * pageCount, page.value * pageCount);
  }

  const newData = user.filter((person) => {
    return Object.values(person).some((value) => {
      return String(value).toLowerCase().includes(q.value.toLowerCase());
    });
  });
  userData = newData;
  return newData.slice((page.value - 1) * pageCount, page.value * pageCount);
});


function handleClick(row: { id: number }) {
  alert("clicked" + row);
}

const toast = useToast();
const modal = useModal();
const count = ref(0);

const items = (row: User) => [
  [
    {
      label: "Edit",
      icon: "i-heroicons-pencil-square-20-solid",
      click: () => openEditUserModal(row.id.toString()),
    },
  ],
  [
    {
      label: "Delete",
      icon: "i-heroicons-trash-20-solid",
      click: () => deleteUser(row.id.toString()),
    },
  ],
];

const isOpen = ref(false);

function openAddUserModal() {
  modal.open(AddUserForm, {
    onSuccess: handleSubmitUser,
  });
}

function openEditUserModal(companyId: string) {
  modal.open(AddUserForm, {
    isEdit: true,
    id: companyId,
    onSuccess: handleSubmitUser,
  });
}


async function handleSubmitUser() {
  toast.add({
    title: "Success!",
    id: "modal-success",
  });

  await fetchAllUsers();
  modal.close();
  isOpen.value = false;
}

async function deleteUser(companyId: string) {
  const confirmed = window.confirm("Are you sure you want to delete this company?");
  if (!confirmed) return;

  try {
    await userManagementAdminApi().deleteUser(companyId);
    toast.add({
      title: "Success!",
      id: "modal-success",
    });

    await fetchAllUsers();
  } catch (error) {
    console.error("Error deleting company:", error);
  }
}

</script>

<template>
  <UButton label="Add User" @click="openAddUserModal" />
  <div class="flex px-3 py-3.5 border-b border-gray-200 dark:border-gray-700">
    <UInput v-model="q" placeholder="Filter user..." />
  </div>
  <UTable :rows="filteredRows" :columns="columns">
    <template #avatar-data="{ row }">
      <UAvatar :src="row.avatar" size="xl" />
    </template>
    <template #actions-data="{ row }">
      <UDropdown :items="items(row)">
        <UButton color="gray" variant="ghost" icon="i-heroicons-ellipsis-horizontal-20-solid" />
      </UDropdown>
    </template>
  </UTable>

  <div class="flex justify-end px-3 py-3.5 border-t border-gray-200 dark:border-gray-700">
    <UPagination v-model="page" :page-count="pageCount" :total="userData.length" />
  </div>
</template>
