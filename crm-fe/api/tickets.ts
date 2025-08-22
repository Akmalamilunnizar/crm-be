import { useApiHost } from '@/composables/useApiHost'
import { useAuthStore } from '@/stores/auth'

const authHeader = () => {
  const authStore = useAuthStore()
  const token = authStore.getToken
  console.log('Auth token:', token) // Debug log
  return { Authorization: `Bearer ${token}` }
}

export const ticketsApi = () => {
  const base = useApiHost()
  return {
    list: () => $fetch(`${base}/api/tickets`, { headers: authHeader() }),
    create: (p:any) => $fetch(`${base}/api/tickets`, { method:'POST', body:p, headers: authHeader() }),
    sendToNOC: (id:number, note:string) => $fetch(`${base}/api/tickets/${id}/send-to-noc`, { method:'POST', body:{note}, headers: authHeader() }),
    sendToCS: (id:number, note:string) => $fetch(`${base}/api/tickets/${id}/send-to-cs`, { method:'POST', body:{note}, headers: authHeader() }),
    nocSolved: (id:number, note:string) => $fetch(`${base}/api/tickets/${id}/noc-solved`, { method:'POST', body:{note}, headers: authHeader() }),
    nocPhysical: (id:number, note:string) => $fetch(`${base}/api/tickets/${id}/noc-physical`, { method:'POST', body:{note}, headers: authHeader() }),
    assignTechnician: (id:number, technician_id:string) => $fetch(`${base}/api/tickets/${id}/assign-technician`, { method:'POST', body:{technician_id}, headers: authHeader() }),
    resolve: (id:number, note:string) => $fetch(`${base}/api/tickets/${id}/resolve`, { method:'POST', body:{note}, headers: authHeader() }),
    byType: () => $fetch(`${base}/api/tickets/reports/by-type`, { headers: authHeader() }),
    troubleTypes: () => $fetch(`${base}/api/tickets/lookups/trouble-types`, { headers: authHeader() }),
    createTroubleType: (id:string, name?:string) => $fetch(`${base}/api/tickets/lookups/trouble-types`, { method:'POST', body:{ id, name }, headers: authHeader() }),
    hotspots: () => $fetch(`${base}/api/tickets/reports/hotspots`, { headers: authHeader() }),
    debugRole: () => $fetch(`${base}/api/tickets/debug/role`, { headers: authHeader() }),
    updates: (since: string) => $fetch(`${base}/api/tickets/updates?since=${encodeURIComponent(since)}`, { headers: authHeader() }),
  }
}


