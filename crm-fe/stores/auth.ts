import { defineStore } from 'pinia'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    token: useCookie('token', { default: () => '' }).value,
    user: { 
      user_id: "", 
      role: useCookie('role_id', { default: () => '' }).value 
    },
  }),
  getters: {
    isLoggedIn: (state) => !!state.token && state.token !== '',
    getToken: (state) => {
      console.log('Auth store getToken called, token value:', state.token)
      return state.token
    },
  },
  actions: {
    login({token,role_id}:{token:string,role_id?:string}) {
      console.log('Auth store login called with token:', token, 'role_id:', role_id)
      
      const tokenCookie = useCookie('token', { 
        default: () => '',
        maxAge: 60 * 60 * 24 * 7, // 7 days
        // secure: true, // Commented out for development
        // sameSite: 'strict' // Commented out for development
      })
      const roleCookie = useCookie('role_id', { 
        default: () => '',
        maxAge: 60 * 60 * 24 * 7, // 7 days
        // secure: true, // Commented out for development
        // sameSite: 'strict' // Commented out for development
      })
      
      tokenCookie.value = token
      roleCookie.value = role_id || ''
      this.token = token
      this.user.role = role_id || ''
      
      console.log('Auth store after login - token:', this.token, 'role:', this.user.role)
    },
    logout() {
      const tokenCookie = useCookie('token')
      const roleCookie = useCookie('role_id')
      
      tokenCookie.value = ''
      roleCookie.value = ''
      this.token = ''
      this.user = { user_id: "", role: "" }
    },
  },
})