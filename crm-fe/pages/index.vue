// https://nuxt.com/docs/examples/advanced/jsx
// https://vuejs.org/guide/extras/render-function.html#jsx-tsx



<script setup>
    import { useAuthStore } from "@/stores/auth";
    const authStore = useAuthStore()
    function setToken(data) {
        authStore.token = data
    }
</script>

<template>

</template>