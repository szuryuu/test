<script setup>
import { computed } from "vue";
import { formatRupiah } from "@/utils/format";

const props = defineProps({
  title: { type: String, required: true },
  value: { type: Number, required: true },
  icon: { type: String, default: "" },
  trend: { type: String, default: "" },
  accent: { type: String, default: "neutral" },
});

const accentConfig = computed(() => {
  const configs = {
    income: {
      color: "var(--color-income)",
      bg: "var(--color-income-bg)",
      iconBg: "rgba(16, 185, 129, 0.12)",
    },
    expense: {
      color: "var(--color-expense)",
      bg: "var(--color-expense-bg)",
      iconBg: "rgba(239, 68, 68, 0.12)",
    },
    profit: {
      color: "var(--color-info)",
      bg: "var(--color-info-bg)",
      iconBg: "rgba(59, 130, 246, 0.12)",
    },
    neutral: {
      color: "var(--color-text)",
      bg: "transparent",
      iconBg: "var(--color-bg)",
    },
  };
  return configs[props.accent] || configs.neutral;
});
</script>

<template>
  <div
    class="summary-card"
    style="
      background: var(--color-surface);
      border: 1px solid var(--color-border);
      border-radius: var(--radius-card);
      padding: 20px;
      transition: all 0.2s ease;
    "
    @mouseenter="$event.target.style.boxShadow = 'var(--shadow-elevated)'; $event.target.style.borderColor = 'var(--color-border-hover)'"
    @mouseleave="$event.target.style.boxShadow = 'var(--shadow-card)'; $event.target.style.borderColor = 'var(--color-border)'"
  >
    <div style="display: flex; justify-content: space-between; align-items: flex-start; margin-bottom: 16px">
      <div
        style="
          width: 40px; height: 40px; border-radius: 10px;
          display: flex; align-items: center; justify-content: center;
          font-size: 18px;
        "
        :style="{ background: accentConfig.iconBg, color: accentConfig.color }"
      >
        <i v-if="icon" :class="icon" />
      </div>
      <span
        v-if="trend"
        style="
          font-size: 11px; font-weight: 500; padding: 2px 8px;
          border-radius: 999px;
        "
        :style="{ background: accentConfig.bg, color: accentConfig.color }"
      >
        {{ trend }}
      </span>
    </div>

    <p
      style="
        margin: 0 0 4px; font-size: 13px; font-weight: 500;
        color: var(--color-text-secondary);
      "
    >
      {{ title }}
    </p>

    <p
      style="
        margin: 0; font-size: 26px; font-weight: 700;
        letter-spacing: -0.03em;
      "
      :style="{ color: accentConfig.color }"
    >
      {{ formatRupiah(value) }}
    </p>
  </div>
</template>

<style scoped>
.summary-card {
  box-shadow: var(--shadow-card);
}
</style>
