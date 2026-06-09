<script setup>
import { formatRupiah, formatDate } from "@/utils/format";

defineProps({
  transactions: { type: Array, default: () => [] },
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
        Transaksi Terbaru
      </h3>
      <span
        style="
          font-size: 11px; color: var(--color-text-tertiary);
        "
      >
        {{ transactions.length }} transaksi
      </span>
    </div>

    <!-- Empty state -->
    <div
      v-if="!transactions.length"
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
        <i class="pi pi-inbox" />
      </div>
      <p style="font-size: 14px; color: var(--color-text-secondary); margin: 0">
        Belum ada transaksi
      </p>
      <p style="font-size: 12px; color: var(--color-text-tertiary); margin: 4px 0 0">
        Kirim pesan via WhatsApp untuk memulai
      </p>
    </div>

    <!-- Transaction list -->
    <div v-else style="display: flex; flex-direction: column; gap: 2px">
      <div
        v-for="tx in transactions.slice(0, 5)"
        :key="tx.id"
        style="
          display: flex; align-items: center; gap: 12px;
          padding: 10px 8px; border-radius: 8px;
          transition: background 0.15s ease;
        "
        @mouseenter="$event.target.style.background = 'var(--color-bg)'"
        @mouseleave="$event.target.style.background = 'transparent'"
      >
        <!-- Type icon -->
        <div
          style="
            width: 36px; height: 36px; border-radius: 8px;
            display: flex; align-items: center; justify-content: center;
            font-size: 15px; flex-shrink: 0;
          "
          :style="{
            background: tx.type === 'income' ? 'var(--color-income-bg)' : 'var(--color-expense-bg)',
            color: tx.type === 'income' ? 'var(--color-income)' : 'var(--color-expense)',
          }"
        >
          <i :class="tx.type === 'income' ? 'pi pi-arrow-down' : 'pi pi-arrow-up'" />
        </div>

        <!-- Description + date -->
        <div style="flex: 1; min-width: 0">
          <p
            style="
              margin: 0; font-size: 13px; font-weight: 500;
              color: var(--color-text);
              overflow: hidden; text-overflow: ellipsis; white-space: nowrap;
            "
          >
            {{ tx.description }}
          </p>
          <p
            style="
              margin: 0; font-size: 11px; color: var(--color-text-tertiary);
            "
          >
            {{ formatDate(tx.transaction_date) }}
          </p>
        </div>

        <!-- Amount -->
        <span
          style="
            font-size: 13px; font-weight: 600; white-space: nowrap;
          "
          :style="{
            color: tx.type === 'income' ? 'var(--color-income)' : 'var(--color-expense)',
          }"
        >
          {{ tx.type === "income" ? "+" : "−" }}{{ formatRupiah(tx.amount) }}
        </span>
      </div>
    </div>
  </div>
</template>
