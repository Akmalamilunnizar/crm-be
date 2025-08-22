import { useApiHost } from "@/composables/useApiHost";
import { useAuthStore } from "@/stores/auth";

export const authApi = () => {
    const api = useApiHost();
    const authStore = useAuthStore();
    
    return {
        loginAuth: async (email: string, password: string) => {
            const response = await fetch(`${api}/api/auth/login`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    email,
                    password,
                }),
            });
            if (!response.ok) {
                const errorData = await response.json();
                throw new Error(errorData.message || 'Login failed');
            }
            return response.json();
        },
        verifyAuth: async () => {
            const token = authStore.getToken;
            if (!token) {
                throw new Error('No token available');
            }
            
            const response = await fetch(`${api}/api/auth/verify`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    "Authorization": `Bearer ${token}`,
                },
            });

            return response.json();
        },
        verifyAuthCustomer: async () => {
            const token = authStore.getToken;
            if (!token) {
                throw new Error('No token available');
            }
            
            const response = await fetch(`${api}/api/auth/verify-customer`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    "Authorization": `Bearer ${token}`,
                },
            });

            return response.json();
        }
    }
}