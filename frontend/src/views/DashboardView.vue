<script setup>
import { ref, onMounted, computed } from "vue";
import { useDashboardStore } from "@/stores/dashboard";
import { useTransactionStore } from "@/stores/transaction";
import SummaryCard from "@/components/dashboard/SummaryCard.vue";
import IncomeExpenseChart from "@/components/dashboard/IncomeExpenseChart.vue";
import RecentTransactions from "@/components/dashboard/RecentTransactions.vue";
import { formatRupiah } from "@/utils/format";

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

const marginDesc = computed(() => {
  if (!store.summary) return "";
  const m = store.summary.profit_margin;
  if (m >= 30) return "Margin sangat sehat";
  if (m >= 20) return "Margin baik";
  if (m >= 10) return "Margin cukup";
  return "Margin perlu ditingkatkan";
});
</script>

<template>
  <div>
    <div class="flex justify-content-between align-items-center mb-4">
      <div>
        <h1 class="text-2xl m-0 mb-1" style="font-weight: 700; color: var(--color-near-black)">
          Dashboard
        </h1>
        <p class="text-sm m-0" style="color: var(--color-text-secondary)">
          Ringkasan keuangan UMKM Anda
        </p>
      </div>
      <div class="flex gap-2">
        <button
          v-for="p in periods"
          :key="p.value"
          @click="changePeriod(p.value)"
          class="px-4 py-2 text-sm cursor-pointer transition-all transition-duration-200"
          :style="{
            background: period === p.value ? 'var(--color-near-black)' : 'var(--color-pure-white)',
            color: period === p.value ? 'var(--color-pure-white)' : 'var(--color-near-black)',
            border: period === p.value ? 'none' : '1px solid var(--color-border-soft)',
            borderRadius: 'var(--radius-sm)',
            fontWeight: 500,
          }"
        >
          {{ p.label }}
        </button>
      </div>
    </div>

    <!-- Summary Cards -->
    <div class="grid mb-4" style="grid-template-columns: repeat(3, 1fr); gap: 16px">
      <SummaryCard
        title="Pemasukan"
        :value="store.summary?.total_income || 0"
        icon="pi pi-arrow-down"
        accent="income"
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
        :trend="marginDesc"
      />
    </div>

    <!-- Chart + Recent -->
    <div class="grid" style="grid-template-columns: 2fr 1fr; gap: 16px">
      <IncomeExpenseChart :chartData="store.summary?.chart_data || []" />
      <RecentTransactions :transactions="txStore.transactions" />
    </div>
  </div>
</template>
