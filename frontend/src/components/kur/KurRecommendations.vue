<script setup>
defineProps({
  recommendations: { type: Array, default: () => [] },
});
</script>

<template>
  <div
    class="bg-[var(--color-surface)] border border-[var(--color-border)] rounded-[var(--radius-card)] p-[20px] shadow-[var(--shadow-card)]"
  >
    <div class="flex justify-between items-center mb-[16px]">
      <h3 class="m-0 text-[15px] font-semibold text-[var(--color-text)]">
        Rekomendasi
      </h3>
      <span
        class="text-[11px] text-[var(--color-text-tertiary)] bg-[var(--color-bg)] py-[2px] px-[8px] rounded-[999px]"
      >
        {{ recommendations.length }} tips
      </span>
    </div>

    <!-- Empty state -->
    <div
      v-if="!recommendations.length"
      class="text-center py-[24px]"
    >
      <div
        class="w-[48px] h-[48px] rounded-[12px] bg-[var(--color-bg)] mx-auto mb-[10px] flex items-center justify-center text-[20px] text-[var(--color-text-tertiary)]"
      >
        <i class="pi pi-info-circle" />
      </div>
      <p class="text-[14px] text-[var(--color-text-secondary)] m-0">
        Belum ada rekomendasi
      </p>
      <p class="text-[12px] text-[var(--color-text-tertiary)] mt-[4px]">
        Tambahkan lebih banyak transaksi
      </p>
    </div>

    <!-- Recommendation list -->
    <ul v-else class="list-none p-0 m-0 flex flex-col gap-[8px]">
      <li
        v-for="(rec, i) in recommendations"
        :key="i"
        class="flex items-start gap-[12px] p-[12px] rounded-[8px] bg-[var(--color-bg)]"
      >
        <span
          class="w-[26px] h-[26px] rounded-[8px] flex items-center justify-center shrink-0 text-[12px] font-bold mt-[1px] bg-[linear-gradient(135deg,#10b981,#059669)] text-white"
        >
          {{ i + 1 }}
        </span>
        <p class="m-0 text-[13px] text-[var(--color-text)] leading-[1.5]">
          {{ rec }}
        </p>
      </li>
    </ul>
  </div>
</template>
