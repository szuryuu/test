<script setup>
defineProps({
  recommendations: { type: Array, default: () => [] },
});
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
    <h3 class="text-base m-0 mb-3" style="font-weight: 600; color: var(--color-near-black)">
      Rekomendasi
    </h3>
    <div v-if="!recommendations.length" class="py-4 text-center" style="color: var(--color-text-secondary)">
      <p class="text-sm m-0">Belum ada rekomendasi</p>
    </div>
    <ul v-else class="list-none p-0 m-0 flex flex-column gap-3">
      <li
        v-for="(rec, i) in recommendations"
        :key="i"
        class="flex align-items-start gap-3"
      >
        <span
          class="flex align-items-center justify-content-center flex-shrink-0"
          style="
            width: 24px;
            height: 24px;
            border-radius: var(--radius-full);
            background: #eff6ff;
            color: var(--color-apple-blue);
            font-size: 12px;
            margin-top: 1px;
          "
        >
          {{ i + 1 }}
        </span>
        <p class="text-sm m-0" style="color: var(--color-near-black); line-height: 1.5">
          {{ rec }}
        </p>
      </li>
    </ul>
  </div>
</template>
