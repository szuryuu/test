<script setup>
import { computed } from "vue";
import { Doughnut } from "vue-chartjs";
import { Chart as ChartJS, ArcElement, Tooltip } from "chart.js";

ChartJS.register(ArcElement, Tooltip);

const props = defineProps({
  score: { type: Number, default: 0 },
});

const levelConfig = computed(() => {
  if (props.score >= 80) return { color: "#10b981", label: "Sangat Baik", desc: "Siap mengajukan KUR" };
  if (props.score >= 60) return { color: "#3b82f6", label: "Baik", desc: "Siap dengan persiapan tambahan" };
  if (props.score >= 40) return { color: "#f59e0b", label: "Sedang", desc: "Tingkatkan konsistensi" };
  return { color: "#ef4444", label: "Rendah", desc: "Fokus pada pencatatan rutin" };
});

const chartConfig = computed(() => ({
  responsive: true,
  maintainAspectRatio: true,
  data: {
    datasets: [
      {
        data: [props.score, 100 - props.score],
        backgroundColor: [levelConfig.value.color, "rgba(226, 232, 240, 0.5)"],
        borderWidth: 0,
        circumference: 270,
        rotation: 225,
        cutout: "78%",
      },
    ],
  },
  options: {
    cutout: "78%",
    plugins: {
      tooltip: { enabled: false },
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
      padding: 24px;
      box-shadow: var(--shadow-card);
      text-align: center;
    "
  >
    <h3
      style="
        margin: 0 0 16px; font-size: 15px; font-weight: 600;
        color: var(--color-text);
      "
    >
      Skor Kesiapan KUR
    </h3>

    <div style="position: relative; width: 180px; height: 130px; margin: 0 auto">
      <Doughnut :data="chartConfig.data" :options="chartConfig.options" />
      <div
        style="
          position: absolute; top: 55%; left: 50%;
          transform: translate(-50%, -50%);
          text-align: center;
        "
      >
        <p
          style="
            margin: 0; font-size: 36px; font-weight: 800;
            letter-spacing: -0.03em; line-height: 1;
          "
          :style="{ color: levelConfig.color }"
        >
          {{ score }}
        </p>
        <p
          style="
            margin: 2px 0 0; font-size: 12px;
            color: var(--color-text-tertiary);
          "
        >
          / 100
        </p>
      </div>
    </div>

    <p
      style="
        margin: 12px 0 4px; font-size: 16px; font-weight: 700;
      "
      :style="{ color: levelConfig.color }"
    >
      {{ levelConfig.label }}
    </p>
    <p
      style="
        margin: 0; font-size: 12px; color: var(--color-text-secondary);
      "
    >
      {{ levelConfig.desc }}
    </p>
  </div>
</template>
