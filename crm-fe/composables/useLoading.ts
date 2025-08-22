import { ref } from 'vue'

const isLoading = ref(false)

export function useLoading() {
  const show = () => isLoading.value = true
  const hide = () => isLoading.value = false

  return {
    isLoading,
    show,
    hide
  }
}
