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
    class="summary-card bg-[var(--color-surface)] border border-[var(--color-border)] rounded-[var(--radius-card)] p-[20px] transition-all duration-[0.2s] ease"
    @mouseenter="$event.target.style.boxShadow = 'var(--shadow-elevated)'; $event.target.style.borderColor = 'var(--color-border-hover)'"
    @mouseleave="$event.target.style.boxShadow = 'var(--shadow-card)'; $event.target.style.borderColor = 'var(--color-border)'"
  >
    <div class="flex justify-between items-start mb-[16px]">
      <div
        class="w-[40px] h-[40px] rounded-[10px] flex items-center justify-center text-[18px]"
        :style="{ background: accentConfig.iconBg, color: accentConfig.color }"
      >
        <i v-if="icon" :class="icon" />
      </div>
      <span
        v-if="trend"
        class="text-[11px] font-medium py-[2px] px-[8px] rounded-[999px]"
        :style="{ background: accentConfig.bg, color: accentConfig.color }"
      >
        {{ trend }}
      </span>
    </div>

    <p class="m-0 mb-[4px] text-[13px] font-medium text-[var(--color-text-secondary)]">
      {{ title }}
    </p>

    <p
      class="m-0 text-[26px] font-bold tracking-[-0.03em]"
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
