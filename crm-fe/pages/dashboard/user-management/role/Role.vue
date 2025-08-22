<script setup lang="ts">
import { userManagementAdminApi } from "@/api/admin/user-management";
import AddRoleForm from "./AddRoleForm.vue";
let role: any[] = [];
const page = ref(1);
const pageCount = 5;

const q = ref("");

async function fetchAllRole() {
  await userManagementAdminApi()
    .getAllRole()
    .then((response) => {
      response.data.forEach((user: any) => {
        user.number = response.data.indexOf(user) + 1;
      });
      role = [...response.data];
      q.value = "changed";
      q.value = "";
    })
    .catch((error) => {
      console.error("Error fetching user:", error);
    });
}
await fetchAllRole();

const columns = [
  { key: "number", label: "Number" },
  { key: "name", label: "Name" },
  {
    key: "actions",
    label: "Actions",
  },
];



type Role = {
  id: number;
  name: string;
  manahe: string;
};
const items = (row: Role) => [
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
      click: () => deleteRole(row.id.toString()),
    },
  ],
];


const toast = useToast();
const modal = useModal();
const count = ref(0);

const isOpen = ref(false);

function openAddUserModal() {
  modal.open(AddRoleForm, {
    onSuccess: handleSubmitUser,
  });
}

function openEditUserModal(companyId: string) {
  modal.open(AddRoleForm, {
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

  await fetchAllRole();
  modal.close();
  isOpen.value = false;
}
async function deleteRole(roleId: string) {
  const confirmed = window.confirm(
    "Are you sure you want to delete this company?"
  );
  if (!confirmed) return;

  try {
    await userManagementAdminApi().deleteRole(roleId);
    toast.add({
      title: "Success!",
      id: "modal-success",
    });

    await fetchAllRole();
  } catch (error) {
    console.error("Error deleting company:", error);
  }
}
</script>

<template>
  <UButton label="Add Role" @click="openAddUserModal" />
  <div class="flex px-3 py-3.5 border-b border-gray-200 dark:border-gray-700">
    <UInput v-model="q" placeholder="Search" />
  </div>
  <!-- <div
    class="flex px-3 py-3.5 border-b border-gray-200 dark:border-gray-700"
  ></div> -->
  <UTable :rows="role" :columns="columns">
    <template #actions-data="{ row }">
      <UDropdown :items="items(row)">
        <UButton
          color="gray"
          variant="ghost"
          icon="i-heroicons-ellipsis-horizontal-20-solid"
        />
      </UDropdown>
    </template>
  </UTable>

  <div
    class="flex justify-end px-3 py-3.5 border-t border-gray-200 dark:border-gray-700"
  >
    <UPagination
      v-model="page"
      :page-count="pageCount"
      :total="role.length"
    />
  </div>
</template>
