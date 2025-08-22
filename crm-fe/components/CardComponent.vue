<script setup lang="ts">
import type { ChartLine, CardData } from '@/types/props';
const props = defineProps<{ dataChart: ChartLine, dataCard: CardData }>();
</script>

<template>
    <div class="p-4 bg-gray-200 rounded-lg">
        <div class="flex justify-between py-2">
            <div>
                {{ props.dataCard.name }}
            </div>
            <div>
                {{ props.dataCard.total }}
            </div>
        </div>
        <div class="p-2 bg-white rounded-lg">
            <ChartComponent :data="props.dataChart" />
        </div>
    </div>
</template>