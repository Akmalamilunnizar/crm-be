<template>
  <div v-if="show" class="fixed inset-0 z-[9999] flex items-center justify-center bg-white dark:bg-gray-900 transition-opacity">
    <div class="flex flex-col items-center">
      <svg class="animate-spin h-12 w-12 text-emerald-500" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
        <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
        <path class="opacity-75" fill="currentColor"
              d="M4 12a8 8 0 018-8v4a4 4 0 00-4 4H4z"/>
      </svg>
      <p class="mt-4 text-gray-500 dark:text-gray-300">Loading...</p>
    </div>
  </div>
</template>

<script setup lang="ts">
defineProps({
  show: {
    type: Boolean,
    default: false
  }
})
</script>

<style scoped>
/* optional: fade effect */
</style>
