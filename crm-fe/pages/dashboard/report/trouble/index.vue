<script setup lang="ts">
import { onMounted, ref, computed } from 'vue'
import { ticketsApi } from '@/api/tickets'

const rows = ref<any[]>([])
const seriesData = ref<any[]>([])
const loading = ref(true)
const troubleTypes = ref<any[]>([])
const typeNameMap = computed(() => {
  const map: Record<string, string> = {}
  for (const t of troubleTypes.value) map[t.id] = t.name || t.id
  return map
})
onMounted(async () => {
  const [list, byType] = await Promise.all([
    ticketsApi().list() as any,
    ticketsApi().byType() as any,
  ])
  rows.value = (list?.data || list) as any[]
  const s = (byType?.data || byType) as any[]
  seriesData.value = s.map((r:any)=>({ name:r.type, value:r.count }))
  loading.value = false
})

const pieOption = computed(()=>({
  tooltip: { trigger: 'item' },
  series: [{ type:'pie', radius:'60%', data: seriesData.value }]
}))

const hotspots = ref<any[]>([])
onMounted(async ()=>{
  const res:any = await ticketsApi().hotspots()
  hotspots.value = res.data || res
})
</script>

<template>
  <div class="space-y-6 text-gray-900">
    <h1 class="text-2xl font-semibold text-gray-900">Trouble Reports</h1>

    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <div class="p-4 bg-white rounded-lg shadow border border-gray-100">
        <h2 class="mb-3 font-semibold text-gray-800">Tickets by Type</h2>
        <ECharts :option="pieOption" style="height:320px" />
      </div>
      <div class="p-4 bg-white rounded-lg shadow border border-gray-100">
        <h2 class="mb-3 font-semibold text-gray-800">Hot Locations</h2>
        <div id="map" class="w-full h-80 rounded border"></div>
      </div>
    </div>

    <div class="p-4 bg-white rounded-lg shadow border border-gray-100">
      <h2 class="mb-4 font-semibold text-gray-800">All Trouble Tickets</h2>
      <div class="table-scroll-container">
        <div class="table-scroll-content">
          <table class="min-w-full text-sm text-gray-900">
            <thead class="bg-gray-100">
              <tr class="text-left border-b border-gray-200 uppercase text-xs tracking-wide text-gray-800">
                <th class="p-2">ID</th>
                <th class="p-2">Title</th>
                <th class="p-2">Type</th>
                <th class="p-2">Status</th>
                <th class="p-2">Assignee</th>
                <th class="p-2">Created</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="r in rows" :key="r.id" class="border-b border-gray-100 odd:bg-white even:bg-gray-50 hover:bg-gray-100/70">
                <td class="p-2">{{ r.id }}</td>
                <td class="p-2">{{ r.title }}</td>
                <td class="p-2 capitalize">{{ r.type }}</td>
                <td class="p-2 capitalize">{{ r.status }}</td>
                <td class="p-2 capitalize">{{ r.current_assignee_role }}</td>
                <td class="p-2">{{ r.created_at?.slice?.(0,10) }}</td>
              </tr>
            </tbody>
          </table>
        </div>
        <div class="table-scroll-footer">
          <span class="scroll-hint">↔ Scroll horizontally to see more columns | ↕ Scroll vertically for more rows</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
// Leaflet map attach using client-only hydration
export default {
  mounted() {
    if (process.client) {
      // defer import to client
      import('leaflet').then(async (Lmod:any)=>{
        const L = Lmod.default || Lmod
        const map = L.map('map').setView([-6.2,106.8], 11)
        L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', { attribution: '&copy; OSM' }).addTo(map)
        const res = await (ticketsApi() as any).hotspots()
        const data = (res.data || res) as any[]
        data.forEach((p:any)=>{
          if (p.gps_lat && p.gps_lng) {
            L.circleMarker([p.gps_lat, p.gps_lng], { radius: 4 + Math.min(p.count, 12), color:'#ef4444'}).addTo(map)
          }
        })
      })
    }
  }
}
</script>

<style scoped>
#map{ min-height: 20rem; }
</style>


