<script setup lang="ts">

const props = defineProps({
    image: {
        type: [String],
        required: true
    }
})

const emit = defineEmits(['close'])

</script>

<template>
    <UModal>
        <div class="p-4">
            <div class="p-2 mb-4 text-2xl font-bold text-center">
                <h1>Images Installation</h1>
            </div>

            <div class="flex flex-wrap justify-center gap-8">
                <img v-for="(item, index) in props.image" :key="index" :src="item" alt="Image Installation"
                    class="object-cover w-1/3 rounded-lg shadow-lg h-1/3" />
            </div>

            <div class="flex justify-end mt-4">
                <UButton @click="$emit('close')" color="primary">Close</UButton>
            </div>

        </div>
    </UModal>
</template>