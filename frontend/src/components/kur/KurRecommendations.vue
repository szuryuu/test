<script setup>
defineProps({
  recommendations: { type: Array, default: () => [] },
});
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
        Rekomendasi
      </h3>
      <span
        style="
          font-size: 11px; color: var(--color-text-tertiary);
          background: var(--color-bg); padding: 2px 8px;
          border-radius: 999px;
        "
      >
        {{ recommendations.length }} tips
      </span>
    </div>

    <!-- Empty state -->
    <div
      v-if="!recommendations.length"
      style="text-align: center; padding: 24px 0"
    >
      <div
        style="
          width: 48px; height: 48px; border-radius: 12px;
          background: var(--color-bg); margin: 0 auto 10px;
          display: flex; align-items: center; justify-content: center;
          font-size: 20px; color: var(--color-text-tertiary);
        "
      >
        <i class="pi pi-info-circle" />
      </div>
      <p style="font-size: 14px; color: var(--color-text-secondary); margin: 0">
        Belum ada rekomendasi
      </p>
      <p style="font-size: 12px; color: var(--color-text-tertiary); margin: 4px 0 0">
        Tambahkan lebih banyak transaksi
      </p>
    </div>

    <!-- Recommendation list -->
    <ul v-else style="list-style: none; padding: 0; margin: 0; display: flex; flex-direction: column; gap: 8px">
      <li
        v-for="(rec, i) in recommendations"
        :key="i"
        style="
          display: flex; align-items: flex-start; gap: 12px;
          padding: 12px; border-radius: 8px;
          background: var(--color-bg);
        "
      >
        <span
          style="
            width: 26px; height: 26px; border-radius: 8px;
            display: flex; align-items: center; justify-content: center;
            flex-shrink: 0; font-size: 12px; font-weight: 700;
            margin-top: 1px;
            background: linear-gradient(135deg, #10b981, #059669);
            color: white;
          "
        >
          {{ i + 1 }}
        </span>
        <p style="margin: 0; font-size: 13px; color: var(--color-text); line-height: 1.5">
          {{ rec }}
        </p>
      </li>
    </ul>
  </div>
</template>
