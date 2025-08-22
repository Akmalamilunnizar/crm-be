import { computed } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { 
  hasPermission, 
  canAccessMenu, 
  getMenuForRole, 
  getRoleDisplayName, 
  getRoleDescription,
  normalizeRole,
  ROLES,
  ROLE_CONFIGS,
  MAIN_MENU,
  type MenuItem
} from '@/utilities/rolePermissions'

export function useRolePermissions() {
  const authStore = useAuthStore()

  // Current user's role (normalized)
  const userRole = computed(() => {
    const role = authStore.user?.role || ''
    return normalizeRole(role)
  })

  // Current user's menu items
  const userMenu = computed(() => {
    if (!userRole.value) return []
    return getMenuForRole(userRole.value)
  })

  // Check if user has specific permission
  const can = (permission: string) => {
    return hasPermission(userRole.value, permission)
  }

  // Check if user can access specific menu item
  const canAccess = (menuItem: MenuItem) => {
    return canAccessMenu(userRole.value, menuItem)
  }

  // Get role display name
  const getRoleName = (role?: string) => {
    return getRoleDisplayName(role || userRole.value)
  }

  // Get role description
  const getRoleDesc = (role?: string) => {
    return getRoleDescription(role || userRole.value)
  }

  // Check if user is admin
  const isAdmin = computed(() => userRole.value === ROLES.ADMIN)

  // Check if user is customer service
  const isCustomerService = computed(() => userRole.value === ROLES.CUSTOMER_SERVICE)

  // Check if user is NOC
  const isNOC = computed(() => userRole.value === ROLES.NOC)

  // Check if user is technician
  const isTechnician = computed(() => userRole.value === ROLES.TECHNICIAN)

  // Check if user is finance
  const isFinance = computed(() => userRole.value === ROLES.FINANCE)

  // Get all available roles
  const availableRoles = computed(() => Object.values(ROLES))

  // Get role configurations
  const roleConfigs = computed(() => ROLE_CONFIGS)

  return {
    // User state
    userRole,
    userMenu,
    
    // Permission checks
    can,
    canAccess,
    
    // Role information
    getRoleName,
    getRoleDesc,
    availableRoles,
    roleConfigs,
    
    // Role type checks
    isAdmin,
    isCustomerService,
    isNOC,
    isTechnician,
    isFinance,
    
    // Constants
    ROLES,
    MAIN_MENU
  }
}
