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

const accentColors = {
  income: "#10b981",
  expense: "#ef4444",
  profit: "#0071e3",
  neutral: "var(--color-near-black)",
};

const accentColor = computed(() => accentColors[props.accent] || accentColors.neutral);
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
    <div class="flex justify-content-between align-items-start mb-3">
      <p class="text-sm m-0" style="color: var(--color-text-secondary); font-weight: 500">
        {{ title }}
      </p>
      <i
        v-if="icon"
        :class="icon"
        style="font-size: 20px; color: var(--color-border-mid)"
      />
    </div>
    <p
      class="text-2xl m-0 mb-1"
      :style="{ fontWeight: 700, color: accentColor }"
    >
      {{ formatRupiah(value) }}
    </p>
    <small
      v-if="trend"
      style="color: var(--color-text-secondary); font-size: 12px"
    >
      {{ trend }}
    </small>
  </div>
</template>
