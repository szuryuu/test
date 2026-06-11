<script setup>
import { computed } from "vue";
import { Doughnut } from "vue-chartjs";
import { Chart as ChartJS, ArcElement, Tooltip, Legend } from "chart.js";

ChartJS.register(ArcElement, Tooltip, Legend);

const props = defineProps({
  score: { type: Number, default: 0 },
});

const levelConfig = computed(() => {
  if (props.score >= 80)
    return {
      color: "#10b981",
      label: "Sangat Baik",
      desc: "Siap mengajukan KUR",
    };
  if (props.score >= 60)
    return {
      color: "#3b82f6",
      label: "Baik",
      desc: "Siap dengan persiapan tambahan",
    };
  if (props.score >= 40)
    return {
      color: "#f59e0b",
      label: "Sedang",
      desc: "Tingkatkan konsistensi",
    };
  return {
    color: "#ef4444",
    label: "Rendah",
    desc: "Fokus pada pencatatan rutin",
  };
});

const chartConfig = computed(() => ({
  responsive: true,
  maintainAspectRatio: true,
  data: {
    datasets: [
      {
        data: [props.score, 100 - props.score, 100 / 3],
        backgroundColor: [
          levelConfig.value.color,
          "rgba(226, 232, 240, 0.5)",
          "transparent",
        ],
        borderWidth: 0,
        rotation: 225,
        cutout: "78%",
      },
    ],
  },
  options: {
    cutout: "78%",
    layout: {
      padding: 0,
    },
    plugins: {
      tooltip: { enabled: false },
      legend: { display: false },
    },
  },
}));
</script>

<template>
  <div
    class="bg-[var(--color-surface)] border border-[var(--color-border)] rounded-[var(--radius-card)] p-[24px] shadow-[var(--shadow-card)] flex flex-col items-center text-center"
  >
    <h3
      class="m-0 mb-[16px] text-[15px] font-semibold text-[var(--color-text)]"
    >
      Skor Kesiapan KUR
    </h3>

    <div class="relative w-[160px] h-[160px] mx-auto mb-[12px]">
      <Doughnut :data="chartConfig.data" :options="chartConfig.options" />
      <div
        class="absolute top-[50%] left-[50%] -translate-x-1/2 -translate-y-1/2 text-center w-full"
      >
        <p
          class="m-0 text-[36px] font-extrabold tracking-[-0.03em] leading-none"
          :style="{ color: levelConfig.color }"
        >
          {{ score }}
        </p>
        <p class="mt-[2px] text-[12px] text-[var(--color-text-tertiary)]">
          / 100
        </p>
      </div>
    </div>

    <p
      class="mt-[8px] mb-[4px] text-[16px] font-bold"
      :style="{ color: levelConfig.color }"
    >
      {{ levelConfig.label }}
    </p>
    <p class="m-0 text-[12px] text-[var(--color-text-secondary)]">
      {{ levelConfig.desc }}
    </p>
  </div>
</template>
