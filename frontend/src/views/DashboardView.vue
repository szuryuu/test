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
    <div class="mb-[24px]">
      <h1 class="m-0 mb-[4px] text-[24px] font-bold text-[var(--color-text)] tracking-[-0.02em]">
        Dashboard
      </h1>
      <p class="m-0 text-[14px] text-[var(--color-text-secondary)]">
        Ringkasan keuangan UMKM Anda
      </p>
    </div>

    <!-- Period selector -->
    <div
      class="flex gap-[4px] mb-[24px] bg-[var(--color-surface)] border border-[var(--color-border)] rounded-[10px] p-[4px] w-fit"
    >
      <button
        v-for="p in periods"
        :key="p.value"
        @click="changePeriod(p.value)"
        class="py-[7px] px-[16px] border-0 rounded-[7px] text-[13px] font-medium font-[inherit] cursor-pointer transition-all duration-[0.15s] ease"
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
      class="p-[12px_16px] bg-[var(--color-expense-bg)] border border-[rgba(239,68,68,0.2)] rounded-[10px] mb-[16px] text-[13px] text-[#dc2626]"
    >
      <i class="pi pi-exclamation-triangle mr-[8px]" />
      {{ store.error }}
    </div>

    <!-- Loading skeleton -->
    <div v-if="isLoading && !store.summary" class="flex flex-col gap-[20px]">
      <div class="grid grid-cols-3 gap-[16px]">
        <div v-for="i in 3" :key="i" class="skeleton h-[120px] rounded-[var(--radius-card)]" />
      </div>
      <div class="skeleton h-[340px] rounded-[var(--radius-card)]" />
    </div>

    <!-- Content -->
    <template v-else>
      <!-- KPI Cards -->
      <div class="kpi-grid grid grid-cols-3 gap-[16px] mb-[20px]">
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
      <div class="chart-recent-grid grid grid-cols-[1.6fr_1fr] gap-[16px]">
        <IncomeExpenseChart :chartData="store.summary?.chart_data || []" />
        <RecentTransactions :transactions="txStore.transactions" />
      </div>
    </template>
  </div>
</template>

<style scoped>
@media (max-width: 900px) {
  .kpi-grid {
    grid-template-columns: repeat(2, 1fr) !important;
  }
  .chart-recent-grid {
    grid-template-columns: 1fr !important;
  }
}
@media (max-width: 500px) {
  .kpi-grid {
    grid-template-columns: 1fr !important;
  }
}
</style>
