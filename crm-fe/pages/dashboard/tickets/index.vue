<script setup lang="ts">
import { onMounted, ref, watch, defineAsyncComponent, computed } from 'vue'
import { ticketsApi } from '@/api/tickets'
import { customerAdminApi } from '@/api/admin/customer'
import { userManagementAdminApi } from '@/api/admin/user-management'
import { useAuthStore } from '@/stores/auth'
import { useRolePermissions } from '@/composables/useRolePermissions'

const authStore = useAuthStore()
const { userRole, isAdmin, isCustomerService, isNOC, isTechnician } = useRolePermissions()

const rows = ref<any[]>([])
const loading = ref(true)
const note = ref('')
const techId = ref('')
const selectedId = ref<number | null>(null)
const activeTab = ref(0)
const isLoading = ref(true)

// Realtime notifications (simple polling)
const lastCheckedAt = ref<string>(new Date().toISOString())
const newUpdates = ref<any[]>([])
const showUpdatesBanner = ref(false)
let pollTimer: any = null

async function pollUpdates() {
  try {
    const res: any = await ticketsApi().updates(lastCheckedAt.value)
    const items = res.data || res || []
    if (items.length > 0) {
      newUpdates.value = items
      showUpdatesBanner.value = true
      // refresh table so user sees the latest entries
      await refresh()
    }
    lastCheckedAt.value = new Date().toISOString()
  } catch (e) {
    // silent
  }
}

function startPolling() {
  if (pollTimer) clearInterval(pollTimer)
  pollTimer = setInterval(pollUpdates, 15000) // 15s
}

function stopPolling() {
  if (pollTimer) clearInterval(pollTimer)
  pollTimer = null
}

// Lookups for modal and table rendering
const customers = ref<any[]>([])
const troubleTypes = ref<any[]>([])
const technicians = ref<any[]>([])
const loadingLookups = ref(true)
const showNewType = ref(false)
const newTypeName = ref('')
const showTechnicianModal = ref(false)
const showNOCNoteModal = ref(false)
const nocNote = ref('')
const typeNameMap = computed(() => {
  const map: Record<string, string> = {}
  for (const t of troubleTypes.value) map[t.id] = t.name || t.id
  return map
})

// hotspots state
const hotspots = ref<any[]>([])
async function loadHotspots() {
  try {
    const res: any = await ticketsApi().hotspots()
    hotspots.value = res.data || res || []
  } catch (e) { console.error('load hotspots', e) }
}

async function refresh() {
  try {
    console.log('Fetching tickets...')
    console.log('Auth token:', authStore.getToken)
    console.log('User role:', authStore.user?.role)

    const res: any = await ticketsApi().list()
    console.log('Tickets API response:', res)

    rows.value = res.data || res
    console.log('Processed tickets:', rows.value)

    loading.value = false
    isLoading.value = false
  } catch (error) {
    console.error('Error fetching tickets:', error)
    loading.value = false
    isLoading.value = false
  }
}

onMounted(async () => {
  // Test role extraction first
  try {
    const debugResponse = await ticketsApi().debugRole()
    console.log('Debug role response:', debugResponse)
  } catch (error) {
    console.error('Debug role error:', error)
  }

  await refresh();
  await loadLookups()
  startPolling()
})

function dismissUpdates() {
  showUpdatesBanner.value = false
  newUpdates.value = []
}

async function refreshUpdates() {
  await refresh()
  showUpdatesBanner.value = false
  newUpdates.value = []
}

function actPrepare(id: number) { selectedId.value = id; note.value = ''; techId.value = '' }

function actPrepareTechnician(id: number) {
  selectedId.value = id;
  techId.value = '';
  showTechnicianModal.value = true
}

function actPrepareNOC(id: number) {
  selectedId.value = id;
  nocNote.value = '';
  showNOCNoteModal.value = true
}

async function sendToNOC() { if (!selectedId.value) return; await ticketsApi().sendToNOC(selectedId.value, note.value); await refresh() }
async function nocSolved() { if (!selectedId.value) return; await ticketsApi().nocSolved(selectedId.value, note.value); await refresh() }
async function nocPhysical() { if (!selectedId.value) return; await ticketsApi().nocPhysical(selectedId.value, note.value); await refresh() }
async function assignTechnician() { if (!selectedId.value) return; await ticketsApi().assignTechnician(selectedId.value, techId.value); await refresh() }

async function assignTechnicianFromModal() {
  if (!selectedId.value || !techId.value) return;
  await ticketsApi().assignTechnician(selectedId.value, techId.value);
  showTechnicianModal.value = false;
  await refresh()
}

async function nocSolvedFromModal() {
  if (!selectedId.value) return;
  await ticketsApi().nocSolved(selectedId.value, nocNote.value);
  showNOCNoteModal.value = false;
  await refresh()
}

async function nocPhysicalFromModal() {
  if (!selectedId.value) return;
  await ticketsApi().nocPhysical(selectedId.value, nocNote.value);
  showNOCNoteModal.value = false;
  await refresh()
}

async function sendToCSFromModal() {
  if (!selectedId.value) return;
  await ticketsApi().sendToCS(selectedId.value, nocNote.value);
  showNOCNoteModal.value = false;
  await refresh()
}
async function resolve() { if (!selectedId.value) return; await ticketsApi().resolve(selectedId.value, note.value); await refresh() }

// Role-based action buttons with workflow awareness
const getTicketActions = (ticket: any) => {
  const actions: Array<{
    label: string
    color: string
    action: () => void
    show: boolean
    tooltip?: string
  }> = [
      {
        label: 'To NOC',
        color: 'bg-blue-600',
        action: () => { actPrepare(ticket.id); sendToNOC() },
        show: (isAdmin.value || isCustomerService.value) &&
          (ticket.current_assignee_name === 'CUSTOMER SERVICE' || ticket.current_assignee_name === 'CUSTOMER_SERVICE' || ticket.current_assignee_name === 'ADMIN') &&
          ticket.status !== 'finished',
        tooltip: 'Send ticket to Network Operations Center'
      },
      {
        label: 'To CS',
        color: 'bg-purple-600',
        action: () => { actPrepareNOC(ticket.id) },
        show: (isAdmin.value || isNOC.value) &&
          ticket.current_assignee_name === 'NOC' &&
          ticket.status !== 'finished',
        tooltip: 'Return ticket to Customer Service'
      },
      
      
      {
        label: 'Assign Tech',
        color: 'bg-cyan-600',
        action: () => { actPrepareTechnician(ticket.id) },
        show: (isAdmin.value || isCustomerService.value) &&
          (ticket.current_assignee_name === 'CUSTOMER SERVICE' || ticket.current_assignee_name === 'CUSTOMER_SERVICE' || ticket.current_assignee_name === 'ADMIN') &&
          ticket.status !== 'finished',
        tooltip: 'Assign to technician for field work'
      },
      {
        label: 'Resolve',
        color: 'bg-emerald-600',
        action: () => { actPrepare(ticket.id); resolve() },
        show: (isAdmin.value || isTechnician.value) &&
          (ticket.current_assignee_name === 'TECHNICIAN' || ticket.current_assignee_name === 'ADMIN') &&
          ticket.status !== 'finished',
        tooltip: 'Mark ticket as resolved'
      }
    ]

  return actions.filter(action => action.show)
}

// Add sendToCS function for NOC users
async function sendToCS() {
  if (!selectedId.value) return;
  // For now, we'll use the same API endpoint but with different logic
  // In a real implementation, you'd have a separate sendToCS endpoint
  await ticketsApi().sendToNOC(selectedId.value, note.value);
  await refresh()
}

const showAdd = ref(false)
const form = ref({ customer_id: '', title: '', description: '', type: '' })
async function loadLookups() {
  try {
    const cust: any = await customerAdminApi().getAllCustomers()
    customers.value = (cust.data || cust) || []
  } catch (e) { console.error('load customers', e) }
  try {
    const tt: any = await ticketsApi().troubleTypes()
    troubleTypes.value = tt.data || tt || []
  } catch (e) { console.error('load trouble types', e) }
  try {
    const tech: any = await userManagementAdminApi().getAllUsers({ query: { role: "TECHNICIAN" } })
    technicians.value = (tech.data || tech) || []
  } catch (e) { console.error('load technicians', e) }
  if (!form.value.customer_id && customers.value.length) form.value.customer_id = customers.value[0].id
  if (!form.value.type && troubleTypes.value.length) form.value.type = troubleTypes.value[0].id
  showNewType.value = troubleTypes.value.length === 0
  loadingLookups.value = false
}
function generateTypeId(): string {
  if (typeof crypto !== 'undefined' && 'randomUUID' in crypto) {
    return (crypto as any).randomUUID()
  }
  return 'tt_' + Math.random().toString(36).slice(2, 10) + Date.now().toString(36)
}

const saveNewType = async () => {
  const id = generateTypeId()
  await ticketsApi().createTroubleType(id, newTypeName.value || undefined)
  const tt: any = await ticketsApi().troubleTypes()
  troubleTypes.value = tt.data || tt || []
  form.value.type = id
  newTypeName.value = ''
  showNewType.value = false
}

// derived GPS for modal
const modalGpsLat = computed(() => {
  const c = customers.value.find(c => c.id === form.value.customer_id)
  return c?.latitude
})
const modalGpsLng = computed(() => {
  const c = customers.value.find(c => c.id === form.value.customer_id)
  return c?.longitude
})
async function createTicket() {
  try {
    console.log('Auth store token:', authStore.getToken) // Debug log
    console.log('Creating ticket with data:', form.value) // Debug log
    await ticketsApi().create({
      customer_id: String(form.value.customer_id),
      title: form.value.title,
      description: form.value.description,
      type: String(form.value.type),
    })
    showAdd.value = false
    form.value = { customer_id: customers.value[0]?.id || '', title: '', description: '', type: troubleTypes.value[0]?.id || '' }
    await refresh()
  } catch (error) {
    console.error('Error creating ticket:', error) // Debug log
  }
}
// Fetchers similar to transaction page
async function fetchAllTickets(params: any) {
  // params is kept for parity; current API does not filter server-side
  isLoading.value = true
  await ticketsApi()
    .list()
    .then((response: any) => {
      const data = response.data || response
      data.forEach((t: any, idx: number) => {
        t.number = idx + 1
      })
      rows.value = data
    })
    .catch((err: any) => {
      console.error('Error fetching tickets:', err)
    })
    .finally(() => {
      isLoading.value = false
    })
}

const tab_items = [
  { label: 'Tickets', value: 'tickets' },
  { label: 'Trouble Reports', value: 'trouble' },
]

async function fetchTicket() { await fetchAllTickets({}) }
async function fetchTrouble() { /* rendered component fetches itself */ }

watch(activeTab, (idx) => { if (idx === 0) fetchTicket(); else fetchTrouble() }, { immediate: true })

// Use alias consistent with tsconfig paths
const TroubleReport = defineAsyncComponent(() => import('@/pages/dashboard/report/trouble/index.vue'))
</script>
<template>
  <div class="space-y-4 text-gray-900">
    <UTabs :items="tab_items" class="w-full" v-model="activeTab" />

    <div v-if="activeTab === 0" class="space-y-4">
      <h1 class="text-2xl font-semibold text-gray-900">Trouble Tickets</h1>

      <div v-if="showUpdatesBanner"
        class="p-3 rounded bg-yellow-50 border border-yellow-200 text-yellow-800 flex items-center justify-between">
        <div>
          New updates available ({{ newUpdates.length }}). Some tickets may have been assigned to your role.
        </div>
        <div class="space-x-2">
          <button class="px-2 py-1 rounded bg-yellow-600 text-white" @click="dismissUpdates">Dismiss</button>
          <button class="px-2 py-1 rounded bg-yellow-700 text-white" @click="refreshUpdates">Refresh</button>
        </div>
      </div>

      <div class="p-4 bg-white rounded-lg shadow border border-gray-100 overflow-auto">
        <div class="flex items-center justify-between mb-3">
          <button v-if="isAdmin || isCustomerService" class="px-3 py-2 bg-emerald-600 text-white rounded"
            @click="showAdd = true">
            Add Ticket
          </button>
          <button class="px-3 py-2 bg-indigo-600 text-white rounded" @click="loadHotspots">Refresh Hotspots</button>
        </div>
        <div class="mb-3 text-xs text-gray-600" v-if="hotspots.length">
          <span class="font-semibold">Top Hotspots:</span>
          <span v-for="(h, idx) in hotspots.slice(0, 5)" :key="idx" class="ml-2">({{ h.gps_lat?.toFixed?.(5) }}, {{
            h.gps_lng?.toFixed?.(5) }}): {{ h.count }}</span>
        </div>
        <table class="min-w-full text-sm text-gray-700">
          <thead class="bg-gray-50">
            <tr class="text-left border-b border-gray-100 uppercase text-xs tracking-wide text-gray-600">
              <th class="p-2">ID</th>
              <th class="p-2">Title</th>
              <th class="p-2">Description</th>
              <th class="p-2">Type</th>
              <th class="p-2">Status</th>
              <th class="p-2">Assigned To Role</th>
              <th class="p-2">Notes</th>
              <th class="p-2">Actions</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="r in rows" :key="r.id" class="border-b border-gray-100 hover:bg-gray-50/60">
              <td class="p-2">{{ r.id }}</td>
              <td class="p-2">{{ r.title }}</td>
              <td class="p-2 max-w-xs truncate" :title="r.description">{{ r.description || '-' }}</td>
              <td class="p-2 capitalize">{{ typeNameMap[r.type] || r.type }}</td>
              <td class="p-2 capitalize">{{ r.status }}</td>
              <td class="p-2 capitalize">{{ r.current_assignee_name || r.current_assignee_role }}</td>
              <td class="p-2 capitalize">{{ r.notes }}</td>
              <td class="p-2 space-x-2">
                <button v-for="action in getTicketActions(r)" :key="action.label"
                  :class="['px-2 py-1 text-white rounded hover:opacity-80 transition-opacity', action.color]"
                  @click="action.action" :title="action.tooltip">
                  {{ action.label }}
                </button>
                <span v-if="getTicketActions(r).length === 0" class="text-gray-400 text-xs">
                  No actions available
                </span>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Modal Add Ticket -->
      <div v-if="showAdd" class="fixed inset-0 z-50 flex items-center justify-center">
        <div class="absolute inset-0 bg-black/60" @click="showAdd = false"></div>
        <div class="relative w-full max-w-2xl mx-4 rounded-xl shadow-xl bg-slate-900 text-slate-100 p-6">
          <div class="flex items-center justify-between mb-4">
            <h2 class="text-xl font-semibold">Add New Ticket</h2>
            <button class="text-slate-300 hover:text-white" @click="showAdd = false">✕</button>
          </div>
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4" v-if="!loadingLookups">
            <div>
              <label class="block text-sm text-slate-300 mb-1">Customer</label>
              <select v-model="form.customer_id"
                class="w-full rounded px-3 py-2 bg-slate-800 border border-slate-700 focus:outline-none">
                <option v-for="c in customers" :key="c.id" :value="c.id">{{ c.name }} ({{ c.id }})</option>
              </select>
            </div>
            <div>
              <label class="block text-sm text-slate-300 mb-1">Type</label>
              <div class="flex gap-2" v-if="!showNewType">
                <select v-model="form.type"
                  class="w-full rounded px-3 py-2 bg-slate-800 border border-slate-700 focus:outline-none">
                  <option v-for="t in troubleTypes" :key="t.id" :value="t.id">{{ t.name || t.id }}</option>
                </select>
                <button type="button" class="px-3 py-2 rounded bg-slate-700" @click="showNewType = true">New</button>
              </div>
              <div v-else class="space-y-2">
                <input v-model="newTypeName" placeholder="Display Name (optional)"
                  class="w-full rounded px-3 py-2 bg-slate-800 border border-slate-700" />
                <div class="flex gap-2">
                  <button type="button" class="px-3 py-2 bg-emerald-600 rounded" @click="saveNewType">Save Type</button>
                  <button type="button" class="px-3 py-2 bg-slate-700 rounded"
                    @click="showNewType = false">Cancel</button>
                </div>
              </div>
            </div>
            <div class="md:col-span-2">
              <label class="block text-sm text-slate-300 mb-1">Title</label>
              <input v-model="form.title"
                class="w-full rounded px-3 py-2 bg-slate-800 border border-slate-700 focus:outline-none" />
            </div>
            <div class="md:col-span-2">
              <label class="block text-sm text-slate-300 mb-1">Description</label>
              <textarea v-model="form.description"
                class="w-full rounded px-3 py-2 bg-slate-800 border border-slate-700 focus:outline-none"></textarea>
            </div>
            <div>
              <label class="block text-sm text-slate-300 mb-1">GPS Lat (from customer)</label>
              <input :value="modalGpsLat ?? ''" disabled
                class="w-full rounded px-3 py-2 bg-slate-800/50 border border-slate-700" />
            </div>
            <div>
              <label class="block text-sm text-slate-300 mb-1">GPS Lng (from customer)</label>
              <input :value="modalGpsLng ?? ''" disabled
                class="w-full rounded px-3 py-2 bg-slate-800/50 border border-slate-700" />
            </div>
          </div>
          <div v-else class="text-slate-300">Loading options...</div>
          <div class="mt-4 flex justify-end gap-2">
            <button class="px-4 py-2 rounded bg-gray-600 text-white" @click="showAdd = false">Cancel</button>
            <button class="px-4 py-2 rounded bg-emerald-600 text-white" @click="createTicket">Submit</button>
          </div>
        </div>
      </div>

      <!-- Modal Assign Technician -->
      <div v-if="showTechnicianModal" class="fixed inset-0 z-50 flex items-center justify-center">
        <div class="absolute inset-0 bg-black/60" @click="showTechnicianModal = false"></div>
        <div class="relative w-full max-w-md mx-4 rounded-xl shadow-xl bg-white p-6">
          <div class="flex items-center justify-between mb-4">
            <h2 class="text-xl font-semibold text-gray-900">Assign Technician</h2>
            <button class="text-gray-400 hover:text-gray-600" @click="showTechnicianModal = false">✕</button>
          </div>
          <div class="space-y-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Select Technician</label>
              <select v-model="techId"
                class="w-full rounded px-3 py-2 border border-gray-300 focus:outline-none focus:ring-2 focus:ring-blue-500">
                <option value="">Choose a technician...</option>
                <option v-for="tech in technicians" :key="tech.id" :value="tech.id">
                  {{ tech.name || tech.username }} ({{ tech.id }})
                </option>
              </select>
            </div>
          </div>
          <div class="mt-6 flex justify-end gap-2">
            <button class="px-4 py-2 rounded bg-gray-300 text-gray-700"
              @click="showTechnicianModal = false">Cancel</button>
            <button class="px-4 py-2 rounded bg-blue-600 text-white disabled:opacity-50"
              @click="assignTechnicianFromModal" :disabled="!techId">
              Assign Technician
            </button>
          </div>
        </div>
      </div>

      <!-- Modal NOC Note -->
      <div v-if="showNOCNoteModal" class="fixed inset-0 z-50 flex items-center justify-center">
        <div class="absolute inset-0 bg-black/60" @click="showNOCNoteModal = false"></div>
        <div class="relative w-full max-w-md mx-4 rounded-xl shadow-xl bg-white p-6">
          <div class="flex items-center justify-between mb-4">
            <h2 class="text-xl font-semibold text-gray-900">NOC Action</h2>
            <button class="text-gray-400 hover:text-gray-600" @click="showNOCNoteModal = false">✕</button>
          </div>
          <div class="space-y-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Notes (Optional)</label>
              <textarea v-model="nocNote" placeholder="Enter any notes about this action..."
                class="w-full rounded px-3 py-2 border border-gray-300 focus:outline-none focus:ring-2 focus:ring-blue-500 h-24 resize-none text-gray-900 bg-white"></textarea>
            </div>
          </div>
          <div class="mt-6 flex justify-end gap-2">
            <button class="px-4 py-2 rounded bg-gray-300 text-gray-700"
              @click="showNOCNoteModal = false">Cancel</button>
            <button class="px-4 py-2 rounded bg-purple-600 text-white" @click="sendToCSFromModal">
              To CS
            </button>
            
            
          </div>
        </div>
      </div>
    </div>

    <div v-else class="space-y-4">
      <ClientOnly>
        <component :is="TroubleReport" />
      </ClientOnly>
    </div>
  </div>
</template>