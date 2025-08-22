import { useApiHost } from "@/composables/useApiHost";


export const authCustomerApi = () => {
    const api = useApiHost();
    return {
        loginAuth: async (email: string, password: string) => {
            const response = await fetch(`${api}/api/auth/login-customer`, {
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
            const response = await fetch(`${api}/api/auth/verify-customer`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    "Authorization": `Bearer ${useCookie("token").value}`,
                },
            });

            return response.json();
        }
    }
}