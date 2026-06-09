<script setup>
import { formatRupiah, formatDate } from "@/utils/format";

defineProps({
  transactions: { type: Array, default: () => [] },
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
      Transaksi Terbaru
    </h3>

    <div v-if="!transactions.length" class="py-4 text-center" style="color: var(--color-text-secondary)">
      <p class="text-sm m-0">Belum ada transaksi</p>
    </div>

    <div v-else class="flex flex-column gap-2">
      <div
        v-for="tx in transactions.slice(0, 5)"
        :key="tx.id"
        class="flex align-items-center justify-content-between py-2"
        style="border-bottom: 1px solid var(--color-border-soft)"
      >
        <div class="flex align-items-center gap-3">
          <div
            class="flex align-items-center justify-content-center"
            style="
              width: 36px;
              height: 36px;
              border-radius: var(--radius-sm);
              font-size: 14px;
            "
            :style="{
              background: tx.type === 'income' ? '#ecfdf5' : '#fef2f2',
              color: tx.type === 'income' ? '#10b981' : '#ef4444',
            }"
          >
            {{ tx.type === "income" ? "↓" : "↑" }}
          </div>
          <div>
            <p class="text-sm m-0" style="font-weight: 500; color: var(--color-near-black)">
              {{ tx.description }}
            </p>
            <p class="text-xs m-0" style="color: var(--color-text-secondary)">
              {{ formatDate(tx.transaction_date) }}
            </p>
          </div>
        </div>
        <span
          class="text-sm"
          :style="{
            fontWeight: 600,
            color: tx.type === 'income' ? '#10b981' : '#ef4444',
          }"
        >
          {{ tx.type === "income" ? "+" : "-" }}{{ formatRupiah(tx.amount) }}
        </span>
      </div>
    </div>
  </div>
</template>
