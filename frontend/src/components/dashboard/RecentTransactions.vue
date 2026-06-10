<script setup>
import { formatRupiah, formatDate } from "@/utils/format";

defineProps({
  transactions: { type: Array, default: () => [] },
});
</script>

<template>
  <div
    class="bg-[var(--color-surface)] border border-[var(--color-border)] rounded-[var(--radius-card)] p-[20px] shadow-[var(--shadow-card)]"
  >
    <div class="flex justify-between items-center mb-[16px]">
      <h3 class="m-0 text-[15px] font-semibold text-[var(--color-text)]">
        Transaksi Terbaru
      </h3>
      <span class="text-[11px] text-[var(--color-text-tertiary)]">
        {{ transactions.length }} transaksi
      </span>
    </div>

    <!-- Empty state -->
    <div v-if="!transactions.length" class="text-center py-[24px]">
      <div
        class="w-[48px] h-[48px] rounded-[12px] bg-[var(--color-bg)] mx-auto mb-[10px] flex items-center justify-center text-[20px] text-[var(--color-text-tertiary)]"
      >
        <i class="pi pi-inbox" />
      </div>
      <p class="text-[14px] text-[var(--color-text-secondary)] m-0">
        Belum ada transaksi
      </p>
      <p class="text-[12px] text-[var(--color-text-tertiary)] mt-[4px]">
        Kirim pesan via WhatsApp untuk memulai
      </p>
    </div>

    <!-- Transaction list -->
    <div v-else class="flex flex-col gap-[2px]">
      <div
        v-for="tx in transactions.slice(0, 5)"
        :key="tx.id"
        class="flex items-center gap-[12px] py-[10px] px-[8px] rounded-[8px] transition-[background] duration-[0.15s] ease"
        @mouseenter="$event.target.style.background = 'var(--color-bg)'"
        @mouseleave="$event.target.style.background = 'transparent'"
      >
        <!-- Type icon -->
        <div
          class="w-[36px] h-[36px] rounded-[8px] flex items-center justify-center text-[15px] shrink-0"
          :style="{
            background: tx.type === 'income' ? 'var(--color-income-bg)' : 'var(--color-expense-bg)',
            color: tx.type === 'income' ? 'var(--color-income)' : 'var(--color-expense)',
          }"
        >
          <i :class="tx.type === 'income' ? 'pi pi-arrow-down' : 'pi pi-arrow-up'" />
        </div>

        <!-- Description + date -->
        <div class="flex-1 min-w-0">
          <p class="m-0 text-[13px] font-medium text-[var(--color-text)] truncate">
            {{ tx.description }}
          </p>
          <p class="m-0 text-[11px] text-[var(--color-text-tertiary)]">
            {{ formatDate(tx.transaction_date) }}
          </p>
        </div>

        <!-- Amount -->
        <span
          class="text-[13px] font-semibold whitespace-nowrap"
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
