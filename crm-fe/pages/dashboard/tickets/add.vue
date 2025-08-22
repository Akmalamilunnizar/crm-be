<script setup lang="ts">
import { ref, onMounted, computed, watch } from 'vue'
import { ticketsApi } from '@/api/tickets'
import { customerAdminApi } from '@/api/admin/customer'

const loading = ref(true)
const customers = ref<{id:string, name:string}[]>([])
const types = ref<{id:string, name?:string}[]>([])
const showNewType = ref(false)
const newTypeName = ref('')

const form = ref({
  customer_id: '',
  title: '',
  description: '' as string | undefined,
  type: ''
})

// derived GPS from selected customer
const gpsLat = ref<number | undefined>(undefined)
const gpsLng = ref<number | undefined>(undefined)
watch(() => form.value.customer_id, (id) => {
  const c = customers.value.find(c => c.id === id)
  gpsLat.value = c?.latitude
  gpsLng.value = c?.longitude
})

onMounted(async () => {
  // load customers (reuse admin API which already exists)
  try {
    const custRes: any = await customerAdminApi().getAllCustomers()
    customers.value = (custRes.data || custRes || []).map((c:any) => ({ id: c.id, name: c.name }))
  } catch (e) { console.error('load customers', e) }

  try {
    const tt: any = await ticketsApi().troubleTypes()
    types.value = tt.data || tt || []
  } catch (e) { console.error('load trouble types', e) }

  // preselect first options if empty
  if (!form.value.customer_id && customers.value.length) form.value.customer_id = customers.value[0].id
  if (!form.value.type && types.value.length) form.value.type = types.value[0].id
  showNewType.value = types.value.length === 0
  loading.value = false
})

async function submit(){
  await ticketsApi().create({
    customer_id: form.value.customer_id,
    title: form.value.title,
    description: form.value.description,
    type: form.value.type,
  })
  navigateTo('/dashboard/tickets')
}

function generateTypeId(): string {
  if (typeof crypto !== 'undefined' && 'randomUUID' in crypto) {
    return (crypto as any).randomUUID()
  }
  // fallback
  return 'tt_' + Math.random().toString(36).slice(2, 10) + Date.now().toString(36)
}

const saveNewType = async () => {
  const id = generateTypeId()
  await ticketsApi().createTroubleType(id, newTypeName.value || undefined)
  const tt: any = await ticketsApi().troubleTypes()
  types.value = tt.data || tt || []
  form.value.type = id
  newTypeName.value = ''
  showNewType.value = false
}
</script>

<template>
  <div class="max-w-2xl p-6 bg-white rounded-lg shadow border border-gray-100">
    <h1 class="text-2xl font-semibold mb-4">Add Trouble Ticket</h1>

    <div v-if="loading">Loading...</div>
    <form v-else class="space-y-4" @submit.prevent="submit">
      <div>
        <label class="block text-sm text-gray-600 mb-1">Customer</label>
        <select v-model="form.customer_id" class="w-full border rounded px-3 py-2">
          <option v-for="c in customers" :key="c.id" :value="c.id">{{ c.name }} ({{ c.id }})</option>
        </select>
      </div>
      <div>
        <label class="block text-sm text-gray-600 mb-1">Title</label>
        <input v-model="form.title" class="w-full border rounded px-3 py-2" />
      </div>
      <div>
        <label class="block text-sm text-gray-600 mb-1">Description</label>
        <textarea v-model="form.description as any" class="w-full border rounded px-3 py-2"></textarea>
      </div>
      <div>
        <label class="block text-sm text-gray-600 mb-1">Type</label>
        <div v-if="!showNewType" class="flex gap-2">
          <select v-model="form.type" class="w-full border rounded px-3 py-2">
            <option v-for="t in types" :key="t.id" :value="t.id">{{ t.name || t.id }}</option>
          </select>
          <button type="button" class="px-3 py-2 border rounded" @click="showNewType = true">New</button>
        </div>
        <div v-else class="space-y-2">
          <input v-model="newTypeName" placeholder="Display Name (optional)" class="w-full border rounded px-3 py-2" />
          <div class="flex gap-2">
            <button type="button" class="px-3 py-2 bg-emerald-600 text-white rounded" @click="saveNewType">Save Type</button>
            <button type="button" class="px-3 py-2 border rounded" @click="showNewType = false">Cancel</button>
          </div>
        </div>
      </div>
      <div class="grid grid-cols-2 gap-4">
        <div>
          <label class="block text-sm text-gray-600 mb-1">GPS Lat (from customer)</label>
          <input :value="gpsLat ?? ''" disabled class="w-full border rounded px-3 py-2 bg-gray-100" />
        </div>
        <div>
          <label class="block text-sm text-gray-600 mb-1">GPS Lng (from customer)</label>
          <input :value="gpsLng ?? ''" disabled class="w-full border rounded px-3 py-2 bg-gray-100" />
        </div>
      </div>

      <button class="px-4 py-2 bg-emerald-600 text-white rounded">Create Ticket</button>
    </form>
  </div>
</template>


