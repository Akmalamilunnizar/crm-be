<script setup lang="ts">
import Role from './role/Role.vue';
import User from './user/User.vue';



const tab_items = [
  {
    label: "User Management",
    value: "user",
  },
  {
    label: "Role",
    value: "role",
  },
  // {
  //   label: "Report Customer Installation",
  //   value: "report-customer",
  // },
];
</script>

<template>
  <UTabs :items="tab_items" class="w-full">
    <template #item="{ item }">
      <div v-if="item.value == 'user'">
        <User />
      </div>

      <div v-if="item.value == 'role'">
        <Role />
      </div>

      <!-- <div v-if="item.value == 'report-customer'">
        <CustomerInstallation />

      </div> -->
    </template>
  </UTabs>
</template>
