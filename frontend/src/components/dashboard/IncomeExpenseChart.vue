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
    style="
      background: var(--color-surface);
      border: 1px solid var(--color-border);
      border-radius: var(--radius-card);
      padding: 20px;
      box-shadow: var(--shadow-card);
    "
  >
    <div
      style="
        display: flex; justify-content: space-between;
        align-items: center; margin-bottom: 16px;
      "
    >
      <h3
        style="
          margin: 0; font-size: 15px; font-weight: 600;
          color: var(--color-text);
        "
      >
        Pemasukan vs Pengeluaran
      </h3>
      <span
        style="
          font-size: 11px; color: var(--color-text-tertiary);
        "
      >
        Per periode
      </span>
    </div>

    <div style="height: 300px; position: relative">
      <Bar
        v-if="hasData"
        :data="chartConfig.data"
        :options="chartConfig.options"
      />
      <div
        v-else
        style="
          display: flex; align-items: center; justify-content: center;
          height: 100%;
        "
      >
        <div style="text-align: center">
          <div
            style="
              width: 56px; height: 56px; border-radius: 12px;
              background: var(--color-bg); margin: 0 auto 12px;
              display: flex; align-items: center; justify-content: center;
              font-size: 24px; color: var(--color-text-tertiary);
            "
          >
            <i class="pi pi-chart-bar" />
          </div>
          <p style="font-size: 14px; color: var(--color-text-secondary); margin: 0">
            Belum ada data transaksi
          </p>
          <p style="font-size: 12px; color: var(--color-text-tertiary); margin: 4px 0 0">
            Transaksi akan muncul setelah dicatat
          </p>
        </div>
      </div>
    </div>
  </div>
</template>
