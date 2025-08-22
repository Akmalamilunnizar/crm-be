<template>
  <div class="space-y-6">
    <div class="bg-white rounded-lg shadow p-6">
      <h1 class="text-2xl font-bold text-gray-900 mb-4">Role-Based Access Control Example</h1>
      
      <!-- User Information -->
      <div class="mb-6 p-4 bg-gray-50 rounded-lg">
        <h2 class="text-lg font-semibold mb-2">Current User</h2>
        <p><strong>Role:</strong> {{ getRoleName() }} ({{ userRole }})</p>
        <p><strong>Description:</strong> {{ getRoleDesc() }}</p>
      </div>

      <!-- Permission Tests -->
      <div class="mb-6">
        <h2 class="text-lg font-semibold mb-4">Permission Tests</h2>
        <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
          <div class="p-3 border rounded" :class="can('dashboard:view') ? 'bg-green-50 border-green-200' : 'bg-red-50 border-red-200'">
            <p class="font-medium">Dashboard View</p>
            <p class="text-sm">{{ can('dashboard:view') ? '✅ Allowed' : '❌ Denied' }}</p>
          </div>
          
          <div class="p-3 border rounded" :class="can('customer:manage') ? 'bg-green-50 border-green-200' : 'bg-red-50 border-red-200'">
            <p class="font-medium">Customer Manage</p>
            <p class="text-sm">{{ can('customer:manage') ? '✅ Allowed' : '❌ Denied' }}</p>
          </div>
          
          <div class="p-3 border rounded" :class="can('invoice:manage') ? 'bg-green-50 border-green-200' : 'bg-red-50 border-red-200'">
            <p class="font-medium">Invoice Manage</p>
            <p class="text-sm">{{ can('invoice:manage') ? '✅ Allowed' : '❌ Denied' }}</p>
          </div>
          
          <div class="p-3 border rounded" :class="can('user:manage') ? 'bg-green-50 border-green-200' : 'bg-red-50 border-red-200'">
            <p class="font-medium">User Manage</p>
            <p class="text-sm">{{ can('user:manage') ? '✅ Allowed' : '❌ Denied' }}</p>
          </div>
        </div>
      </div>

      <!-- Role-Based Content -->
      <div class="space-y-4">
        <h2 class="text-lg font-semibold">Role-Based Content</h2>
        
        <!-- Admin Only Content -->
        <div v-if="isAdmin" class="p-4 bg-blue-50 border border-blue-200 rounded-lg">
          <h3 class="font-semibold text-blue-900">Admin Only Section</h3>
          <p class="text-blue-700">This content is only visible to administrators.</p>
          <button class="mt-2 px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700">
            Admin Action
          </button>
        </div>

        <!-- Customer Service Content -->
        <div v-if="isCustomerService" class="p-4 bg-green-50 border border-green-200 rounded-lg">
          <h3 class="font-semibold text-green-900">Customer Service Section</h3>
          <p class="text-green-700">This content is only visible to customer service representatives.</p>
          <button class="mt-2 px-4 py-2 bg-green-600 text-white rounded hover:bg-green-700">
            Customer Service Action
          </button>
        </div>

        <!-- Finance Content -->
        <div v-if="isFinance" class="p-4 bg-purple-50 border border-purple-200 rounded-lg">
          <h3 class="font-semibold text-purple-900">Finance Section</h3>
          <p class="text-purple-700">This content is only visible to finance personnel.</p>
          <button class="mt-2 px-4 py-2 bg-purple-600 text-white rounded hover:bg-purple-700">
            Finance Action
          </button>
        </div>

        <!-- Technician Content -->
        <div v-if="isTechnician" class="p-4 bg-orange-50 border border-orange-200 rounded-lg">
          <h3 class="font-semibold text-orange-900">Technician Section</h3>
          <p class="text-orange-700">This content is only visible to technicians.</p>
          <button class="mt-2 px-4 py-2 bg-orange-600 text-white rounded hover:bg-orange-700">
            Technician Action
          </button>
        </div>

        <!-- NOC Content -->
        <div v-if="isNOC" class="p-4 bg-indigo-50 border border-indigo-200 rounded-lg">
          <h3 class="font-semibold text-indigo-900">NOC Section</h3>
          <p class="text-indigo-700">This content is only visible to NOC personnel.</p>
          <button class="mt-2 px-4 py-2 bg-indigo-600 text-white rounded hover:bg-indigo-700">
            NOC Action
          </button>
        </div>
      </div>

      <!-- Available Menu Items -->
      <div class="mt-6">
        <h2 class="text-lg font-semibold mb-4">Your Available Menu Items</h2>
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          <div 
            v-for="item in userMenu" 
            :key="item.label"
            class="p-4 border rounded-lg hover:bg-gray-50 cursor-pointer"
            @click="navigateTo(item.link)"
          >
            <div class="flex items-center gap-3">
              <UIcon :name="item.icon" class="w-6 h-6 text-gray-600" />
              <div>
                <h3 class="font-medium">{{ item.label }}</h3>
                <p class="text-sm text-gray-500">{{ item.description }}</p>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Role Information -->
      <div class="mt-6">
        <h2 class="text-lg font-semibold mb-4">All Available Roles</h2>
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          <div 
            v-for="(config, role) in roleConfigs" 
            :key="role"
            class="p-4 border rounded-lg"
            :class="role === userRole ? 'bg-blue-50 border-blue-200' : 'bg-gray-50'"
          >
            <h3 class="font-semibold">{{ config.displayName }}</h3>
            <p class="text-sm text-gray-600 mb-2">{{ config.description }}</p>
            <div class="text-xs text-gray-500">
              <p><strong>Permissions:</strong> {{ config.permissions.length }}</p>
              <p v-if="role === userRole" class="text-blue-600 font-medium">← Current Role</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRolePermissions } from '@/composables/useRolePermissions'

const { 
  userRole, 
  userMenu, 
  can, 
  getRoleName, 
  getRoleDesc, 
  isAdmin, 
  isCustomerService, 
  isFinance, 
  isTechnician, 
  isNOC, 
  roleConfigs 
} = useRolePermissions()

const navigateTo = (link: string) => {
  // In a real app, you'd use router.push(link)
  console.log('Navigating to:', link)
}
</script>
