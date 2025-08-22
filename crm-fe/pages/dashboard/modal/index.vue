<script setup lang="ts">
import Modal from './Modal.vue'

const toast = useToast()
const modal = useModal()
const count = ref(0)

function openModal() {
  count.value += 1
  modal.open(Modal, {
    count: count.value,
    onSuccess() {
      toast.add({
        title: 'Success !',
        id: 'modal-success'
      })
    }
  })
}
</script>

<template>
 
</template>
