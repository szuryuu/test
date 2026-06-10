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
  Filler,
} from "chart.js";

ChartJS.register(CategoryScale, LinearScale, BarElement, Tooltip, Legend, Filler);

const props = defineProps({
  chartData: { type: Array, default: () => [] },
});

const hasData = computed(() => props.chartData && props.chartData.length > 0);

const chartConfig = computed(() => ({
  responsive: true,
  maintainAspectRatio: false,
  data: {
    labels: props.chartData.map((d) => {
      const date = new Date(d.date);
      return date.toLocaleDateString("id-ID", {
        day: "numeric",
        month: "short",
      });
    }),
    datasets: [
      {
        label: "Pemasukan",
        data: props.chartData.map((d) => d.income),
        backgroundColor: "rgba(16, 185, 129, 0.85)",
        hoverBackgroundColor: "rgba(16, 185, 129, 1)",
        borderRadius: 6,
        borderSkipped: false,
        barPercentage: 0.55,
        categoryPercentage: 0.8,
      },
      {
        label: "Pengeluaran",
        data: props.chartData.map((d) => d.expense),
        backgroundColor: "rgba(239, 68, 68, 0.75)",
        hoverBackgroundColor: "rgba(239, 68, 68, 1)",
        borderRadius: 6,
        borderSkipped: false,
        barPercentage: 0.55,
        categoryPercentage: 0.8,
      },
    ],
  },
  options: {
    plugins: {
      legend: {
        position: "top",
        align: "end",
        labels: {
          usePointStyle: true,
          padding: 16,
          boxWidth: 8,
          boxHeight: 8,
          font: {
            family: "Inter, system-ui, sans-serif",
            size: 12,
          },
          color: "#64748b",
        },
      },
      tooltip: {
        backgroundColor: "#0f172a",
        titleFont: { family: "Inter, system-ui, sans-serif", size: 12 },
        bodyFont: { family: "Inter, system-ui, sans-serif", size: 13 },
        padding: 10,
        cornerRadius: 8,
        callbacks: {
          label: (ctx) => {
            const val = ctx.raw;
            return `${ctx.dataset.label}: Rp${val.toLocaleString("id-ID")}`;
          },
        },
      },
    },
    scales: {
      x: {
        grid: { display: false },
        ticks: {
          font: { family: "Inter, system-ui, sans-serif", size: 11 },
          color: "#94a3b8",
          maxRotation: 45,
        },
      },
      y: {
        grid: {
          color: "rgba(226, 232, 240, 0.6)",
          drawBorder: false,
        },
        border: { display: false },
        ticks: {
          font: { family: "Inter, system-ui, sans-serif", size: 11 },
          color: "#94a3b8",
          padding: 8,
          callback: (val) => {
            if (val >= 1000000) return `Rp${(val / 1000000).toFixed(1)}jt`;
            if (val >= 1000) return `Rp${(val / 1000).toFixed(0)}rb`;
            return `Rp${val}`;
          },
        },
      },
    },
  },
}));
</script>

<template>
  <div
    class="bg-[var(--color-surface)] border border-[var(--color-border)] rounded-[var(--radius-card)] p-[20px] shadow-[var(--shadow-card)]"
  >
    <div class="flex justify-between items-center mb-[16px]">
      <h3 class="m-0 text-[15px] font-semibold text-[var(--color-text)]">
        Pemasukan vs Pengeluaran
      </h3>
      <span class="text-[11px] text-[var(--color-text-tertiary)]">
        Per periode
      </span>
    </div>

    <div class="h-[300px] relative">
      <Bar
        v-if="hasData"
        :data="chartConfig.data"
        :options="chartConfig.options"
      />
      <div v-else class="flex items-center justify-center h-full">
        <div class="text-center">
          <div
            class="w-[56px] h-[56px] rounded-[12px] bg-[var(--color-bg)] mx-auto mb-[12px] flex items-center justify-center text-[24px] text-[var(--color-text-tertiary)]"
          >
            <i class="pi pi-chart-bar" />
          </div>
          <p class="text-[14px] text-[var(--color-text-secondary)] m-0">
            Belum ada data transaksi
          </p>
          <p class="text-[12px] text-[var(--color-text-tertiary)] mt-[4px]">
            Transaksi akan muncul setelah dicatat
          </p>
        </div>
      </div>
    </div>
  </div>
</template>
