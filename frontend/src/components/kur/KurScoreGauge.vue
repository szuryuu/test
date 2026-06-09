<script setup>
import { computed } from "vue";
import { Doughnut } from "vue-chartjs";
import { Chart as ChartJS, ArcElement, Tooltip } from "chart.js";

ChartJS.register(ArcElement, Tooltip);

const props = defineProps({
  score: { type: Number, default: 0 },
});

const levelColor = computed(() => {
  if (props.score >= 80) return "#10b981";
  if (props.score >= 60) return "#0071e3";
  if (props.score >= 40) return "#f59e0b";
  return "#ef4444";
});

const chartConfig = computed(() => ({
  responsive: true,
  maintainAspectRatio: true,
  data: {
    datasets: [
      {
        data: [props.score, 100 - props.score],
        backgroundColor: [levelColor.value, "#f0f0f0"],
        borderWidth: 0,
        circumference: 270,
        rotation: 225,
        cutout: "75%",
      },
    ],
  },
  options: {
    plugins: {
      tooltip: { enabled: false },
    },
  },
}));

const levelLabel = computed(() => {
  if (props.score >= 80) return "Sangat Baik";
  if (props.score >= 60) return "Baik";
  if (props.score >= 40) return "Sedang";
  return "Rendah";
});
</script>

<template>
  <div
    class="text-center p-4"
    style="
      background: var(--color-pure-white);
      border-radius: var(--radius-md);
      box-shadow: 0 1px 3px rgba(0,0,0,0.04);
    "
  >
    <h3 class="text-base m-0 mb-4" style="font-weight: 600; color: var(--color-near-black)">
      Skor Kesiapan KUR
    </h3>
    <div style="position: relative; width: 200px; height: 140px; margin: 0 auto">
      <Doughnut :data="chartConfig.data" :options="chartConfig.options" />
      <div
        style="
          position: absolute;
          top: 55%;
          left: 50%;
          transform: translate(-50%, -50%);
          text-align: center;
        "
      >
        <p
          class="text-3xl m-0"
          :style="{ fontWeight: 700, color: levelColor }"
        >
          {{ score }}
        </p>
        <p class="text-xs m-0" style="color: var(--color-text-secondary)">
          / 100
        </p>
      </div>
    </div>
    <p
      class="text-sm mt-2 mb-0"
      :style="{ fontWeight: 600, color: levelColor }"
    >
      {{ levelLabel }}
    </p>
  </div>
</template>
