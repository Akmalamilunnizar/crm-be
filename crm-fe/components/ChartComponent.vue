<template>
    <Line id="my-chart-id" :options="chartOptions" :data="chartData" />
  </template>
  
  <script setup lang="ts">
  import { Line } from 'vue-chartjs';
  import type { ChartLine } from "@/types/props";

  import {
    Chart as ChartJS,
    Title,
    Tooltip,
    Legend,
    CategoryScale,
    LinearScale,
    PointElement,
    LineElement,
  } from 'chart.js';
  
  // Registrasi elemen yang diperlukan untuk Chart.js
  ChartJS.register(Title, Tooltip, Legend, CategoryScale, LinearScale, PointElement, LineElement);
  
  // Definisi tipe data untuk props

  const props = defineProps<{ data: ChartLine }>();
  
  // Data chart
  const chartData = {
    labels: props.data.x,
    datasets: [
      {
        data: props.data.y,
        label: props.data.label,
        borderColor: 'rgb(75, 192, 192)',
        tension: 0.5,
      },
    ],
  };
  
  // Opsi konfigurasi chart
  const chartOptions = {
    responsive: true,
  };
  </script>
  