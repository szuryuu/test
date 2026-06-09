<script setup>
import { computed } from "vue";
import { Bar } from "vue-chartjs";
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  BarElement,
  Tooltip,
  Legend,
} from "chart.js";

ChartJS.register(CategoryScale, LinearScale, BarElement, Tooltip, Legend);

const props = defineProps({
  chartData: { type: Array, default: () => [] },
});

const chartConfig = computed(() => ({
  responsive: true,
  maintainAspectRatio: false,
  data: {
    labels: props.chartData.map((d) => {
      const date = new Date(d.date);
      return date.toLocaleDateString("id-ID", { day: "numeric", month: "short" });
    }),
    datasets: [
      {
        label: "Pemasukan",
        data: props.chartData.map((d) => d.income),
        backgroundColor: "#10b981",
        borderRadius: 4,
        barPercentage: 0.6,
      },
      {
        label: "Pengeluaran",
        data: props.chartData.map((d) => d.expense),
        backgroundColor: "#ef4444",
        borderRadius: 4,
        barPercentage: 0.6,
      },
    ],
  },
  options: {
    plugins: {
      legend: {
        position: "bottom",
        labels: {
          usePointStyle: true,
          padding: 20,
          font: { family: "Inter, system-ui, sans-serif", size: 12 },
        },
      },
      tooltip: {
        callbacks: {
          label: (ctx) => {
            const val = ctx.raw;
            return `Rp${val.toLocaleString("id-ID")}`;
          },
        },
      },
    },
    scales: {
      x: {
        grid: { display: false },
        ticks: { font: { family: "Inter, system-ui, sans-serif", size: 11 } },
      },
      y: {
        grid: { color: "#f0f0f0" },
        ticks: {
          font: { family: "Inter, system-ui, sans-serif", size: 11 },
          callback: (val) => (val >= 1000 ? `${val / 1000}rb` : val),
        },
      },
    },
  },
}));
</script>

<template>
  <div
    class="p-4"
    style="
      background: var(--color-pure-white);
      border-radius: var(--radius-md);
      box-shadow: 0 1px 3px rgba(0,0,0,0.04);
    "
  >
    <h3 class="text-base m-0 mb-3" style="font-weight: 600; color: var(--color-near-black)">
      Pemasukan vs Pengeluaran
    </h3>
    <div style="height: 280px">
      <Bar v-if="chartData.length" :data="chartConfig.data" :options="chartConfig.options" />
      <div
        v-else
        class="flex align-items-center justify-content-center h-full"
        style="color: var(--color-text-secondary)"
      >
        <p class="text-sm">Belum ada data transaksi</p>
      </div>
    </div>
  </div>
</template>
