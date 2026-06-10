<script setup>
import { ref, onMounted, computed } from "vue";
import { useDashboardStore } from "@/stores/dashboard";
import { useTransactionStore } from "@/stores/transaction";
import SummaryCard from "@/components/dashboard/SummaryCard.vue";
import IncomeExpenseChart from "@/components/dashboard/IncomeExpenseChart.vue";
import RecentTransactions from "@/components/dashboard/RecentTransactions.vue";

const store = useDashboardStore();
const txStore = useTransactionStore();
const period = ref("monthly");

onMounted(() => {
  store.fetchSummary(period.value);
  txStore.fetchTransactions({ page: 1, limit: 5 });
});

function changePeriod(p) {
  period.value = p;
  store.fetchSummary(p);
}

const periods = [
  { label: "Harian", value: "daily" },
  { label: "Mingguan", value: "weekly" },
  { label: "Bulanan", value: "monthly" },
  { label: "Tahunan", value: "yearly" },
];

const marginTrend = computed(() => {
  if (!store.summary) return "";
  const m = store.summary.profit_margin;
  if (m >= 30) return "Sangat sehat";
  if (m >= 20) return "Baik";
  if (m >= 10) return "Cukup";
  if (m > 0) return "Perlu ditingkatkan";
  return "";
});

const isLoading = computed(() => store.isLoading);
</script>

<template>
  <div>
    <!-- Page header -->
    <div style="margin-bottom: 24px">
      <h1
        style="
          margin: 0 0 4px; font-size: 24px; font-weight: 700;
          color: var(--color-text); letter-spacing: -0.02em;
        "
      >
        Dashboard
      </h1>
      <p style="margin: 0; font-size: 14px; color: var(--color-text-secondary)">
        Ringkasan keuangan UMKM Anda
      </p>
    </div>

    <!-- Period selector -->
    <div
      style="
        display: flex; gap: 4px; margin-bottom: 24px;
        background: var(--color-surface);
        border: 1px solid var(--color-border);
        border-radius: 10px;
        padding: 4px;
        width: fit-content;
      "
    >
      <button
        v-for="p in periods"
        :key="p.value"
        @click="changePeriod(p.value)"
        style="
          padding: 7px 16px; border: none; border-radius: 7px;
          font-size: 13px; font-weight: 500; font-family: inherit;
          cursor: pointer; transition: all 0.15s ease;
        "
        :style="{
          background: period === p.value
            ? 'linear-gradient(135deg, #10b981, #059669)'
            : 'transparent',
          color: period === p.value ? 'white' : 'var(--color-text-secondary)',
        }"
        @mouseenter="(period !== p.value) && ($event.target.style.background = 'var(--color-bg)')"
        @mouseleave="(period !== p.value) && ($event.target.style.background = 'transparent')"
      >
        {{ p.label }}
      </button>
    </div>

    <!-- Error state -->
    <div
      v-if="store.error"
      style="padding: 12px 16px; background: var(--color-expense-bg); border: 1px solid rgba(239,68,68,0.2); border-radius: 10px; margin-bottom: 16px; font-size: 13px; color: #dc2626"
    >
      <i class="pi pi-exclamation-triangle" style="margin-right: 8px" />
      {{ store.error }}
    </div>

    <!-- Loading skeleton -->
    <div v-if="isLoading && !store.summary" style="display: flex; flex-direction: column; gap: 20px">
      <div style="display: grid; grid-template-columns: repeat(3, 1fr); gap: 16px">
        <div v-for="i in 3" :key="i" class="skeleton" style="height: 120px; border-radius: var(--radius-card)" />
      </div>
      <div class="skeleton" style="height: 340px; border-radius: var(--radius-card)" />
    </div>

    <!-- Content -->
    <template v-else>
      <!-- KPI Cards -->
      <div
        style="
          display: grid;
          grid-template-columns: repeat(3, 1fr);
          gap: 16px;
          margin-bottom: 20px;
        "
      >
        <SummaryCard
          title="Pemasukan"
          :value="store.summary?.total_income || 0"
          icon="pi pi-arrow-down"
          accent="income"
          :trend="store.summary?.profit_margin >= 0 ? marginTrend : ''"
        />
        <SummaryCard
          title="Pengeluaran"
          :value="store.summary?.total_expense || 0"
          icon="pi pi-arrow-up"
          accent="expense"
        />
        <SummaryCard
          title="Laba Bersih"
          :value="store.summary?.net_profit || 0"
          icon="pi pi-chart-line"
          accent="profit"
          :trend="store.summary?.net_profit > 0 ? 'Positif' : store.summary?.net_profit === 0 ? '' : 'Negatif'"
        />
      </div>

      <!-- Chart + Recent Transactions -->
      <div
        style="
          display: grid;
          grid-template-columns: 1.6fr 1fr;
          gap: 16px;
        "
      >
        <IncomeExpenseChart :chartData="store.summary?.chart_data || []" />
        <RecentTransactions :transactions="txStore.transactions" />
      </div>
    </template>
  </div>
</template>

<style scoped>
@media (max-width: 900px) {
  div > div[style*="grid-template-columns: repeat(3, 1fr)"] {
    grid-template-columns: repeat(2, 1fr) !important;
  }
  div > div[style*="grid-template-columns: 1.6fr 1fr"] {
    grid-template-columns: 1fr !important;
  }
}
@media (max-width: 500px) {
  div > div[style*="grid-template-columns: repeat(3, 1fr)"] {
    grid-template-columns: 1fr !important;
  }
}
</style>
