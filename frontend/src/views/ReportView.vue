<script setup>
import { ref } from "vue";
import { reportAPI } from "@/api/kur";
import { formatRupiah, parseWhatsAppText } from "@/utils/format";

const year = ref(new Date().getFullYear());
const month = ref(new Date().getMonth() + 1);
const isLoading = ref(false);
const report = ref(null);
const error = ref("");

const monthNames = [
  "Januari",
  "Februari",
  "Maret",
  "April",
  "Mei",
  "Juni",
  "Juli",
  "Agustus",
  "September",
  "Oktober",
  "November",
  "Desember",
];

async function generateReport() {
  isLoading.value = true;
  error.value = "";
  report.value = null;
  try {
    const res = await reportAPI.generate(year.value, month.value);
    report.value = res.data.data;
  } catch (err) {
    error.value = err.response?.data?.message || "Gagal membuat laporan";
  } finally {
    isLoading.value = false;
  }
}

const handlePrint = () => {
  if (typeof window !== "undefined") {
    window.print();
  }
};
</script>

<template>
  <div>
    <!-- Page header -->
    <div class="mb-[24px]">
      <h1
        class="m-0 mb-[4px] text-[24px] font-bold text-[var(--color-text)] tracking-[-0.02em]"
      >
        Laporan Bulanan
      </h1>
      <p class="m-0 text-[14px] text-[var(--color-text-secondary)]">
        Generate laporan keuangan bulanan dan kirim via WhatsApp
      </p>
    </div>

    <!-- Generator card -->
    <div
      class="no-print bg-[var(--color-surface)] border border-[var(--color-border)] rounded-[var(--radius-card)] p-[24px] shadow-[var(--shadow-card)] mb-[20px]"
    >
      <h3
        class="m-0 mb-[16px] text-[15px] font-semibold text-[var(--color-text)]"
      >
        <i class="pi pi-calendar mr-[8px] text-[var(--color-brand-500)]" />
        Pilih Periode Laporan
      </h3>

      <div class="flex flex-wrap items-end gap-[12px]">
        <div class="flex flex-col gap-[6px] min-w-[160px]">
          <label
            class="text-[12px] font-medium text-[var(--color-text-secondary)]"
          >
            Bulan
          </label>
          <select
            v-model="month"
            class="py-[10px] px-[12px] text-[14px] font-[inherit] border border-[var(--color-border)] rounded-[8px] bg-[var(--color-bg)] text-[var(--color-text)] cursor-pointer"
          >
            <option v-for="(name, i) in monthNames" :key="i" :value="i + 1">
              {{ name }}
            </option>
          </select>
        </div>

        <div class="flex flex-col gap-[6px] min-w-[100px]">
          <label
            class="text-[12px] font-medium text-[var(--color-text-secondary)]"
          >
            Tahun
          </label>
          <input
            v-model="year"
            type="number"
            min="2024"
            max="2030"
            class="py-[10px] px-[12px] text-[14px] font-[inherit] border border-[var(--color-border)] rounded-[8px] bg-[var(--color-bg)] text-[var(--color-text)]"
          />
        </div>

        <button
          @click="generateReport"
          :disabled="isLoading"
          class="flex items-center gap-[8px] py-[10px] px-[24px] text-[14px] font-semibold font-[inherit] border-0 rounded-[8px] cursor-pointer transition-all duration-[0.15s] ease bg-[linear-gradient(135deg,#10b981,#059669)] text-white disabled:cursor-default"
          :style="{ opacity: isLoading ? 0.7 : 1 }"
          @mouseenter="
            !isLoading &&
            ($event.target.style.boxShadow =
              '0 4px 12px rgba(16, 185, 129, 0.4)')
          "
          @mouseleave="!isLoading && ($event.target.style.boxShadow = 'none')"
        >
          <i class="pi pi-file text-[14px]" />
          {{ isLoading ? "Membuat Laporan..." : "Buat Laporan" }}
        </button>
      </div>
    </div>

    <!-- Loading skeleton -->
    <div v-if="isLoading" class="flex flex-col gap-[16px]">
      <div class="skeleton h-[200px] rounded-[var(--radius-card)]" />
      <div class="skeleton h-[80px] rounded-[var(--radius-card)]" />
    </div>

    <!-- Error -->
    <div
      v-if="error"
      class="py-[16px] px-[20px] mb-[20px] bg-[var(--color-expense-bg)] border border-[rgba(239,68,68,0.2)] rounded-[var(--radius-card)] flex items-start gap-[12px]"
    >
      <i
        class="pi pi-exclamation-triangle text-[var(--color-expense)] text-[18px] mt-[1px]"
      />
      <div>
        <p class="m-0 text-[14px] font-medium text-[#dc2626]">
          {{ error }}
        </p>
      </div>
    </div>

    <!-- Report result -->
    <div
      v-if="report"
      class="bg-[var(--color-surface)] border border-[var(--color-border)] rounded-[var(--radius-card)] shadow-[var(--shadow-card)] overflow-hidden"
    >
      <!-- Report header -->
      <div
        class="py-[20px] px-[24px] border-b border-[var(--color-border)] flex justify-between items-center flex-wrap gap-[8px]"
      >
        <div>
          <h3 class="m-0 text-[16px] font-semibold text-[var(--color-text)]">
            Laporan {{ monthNames[month - 1] }} {{ year }}
          </h3>
          <p
            class="mt-[2px] mb-0 text-[12px] text-[var(--color-text-tertiary)]"
          >
            Ringkasan keuangan bulanan
          </p>
        </div>
        <span
          class="text-[11px] py-[4px] px-[10px] bg-[var(--color-income-bg)] text-[var(--color-income)] rounded-[999px] font-medium"
        >
          Siap dikirim via WhatsApp
        </span>
        <button
          class="no-print flex items-center gap-[8px] py-[8px] px-[16px] text-[13px] font-semibold font-[inherit] border border-[var(--color-border)] rounded-[8px] cursor-pointer transition-all duration-[0.15s] ease bg-[var(--color-surface)] text-[var(--color-text-secondary)]"
          @click="handlePrint"
          @mouseenter="
            $event.target.style.borderColor = '#94a3b8';
            $event.target.style.color = 'var(--color-text)';
          "
          @mouseleave="
            $event.target.style.borderColor = 'var(--color-border)';
            $event.target.style.color = 'var(--color-text-secondary)';
          "
        >
          <i class="pi pi-print text-[14px]" />
          Cetak Laporan
        </button>
      </div>

      <!-- Summary numbers -->
      <div
        v-if="report.total_income"
        class="summary-grid grid grid-cols-3 gap-[1px] bg-[var(--color-border)]"
      >
        <div class="bg-[var(--color-surface)] py-[16px] px-[20px] text-center">
          <p
            class="m-0 mb-[4px] text-[11px] font-semibold text-[var(--color-income)] uppercase tracking-[0.05em]"
          >
            Pemasukan
          </p>
          <p class="m-0 text-[18px] font-bold text-[var(--color-income)]">
            {{ formatRupiah(report.total_income) }}
          </p>
        </div>
        <div class="bg-[var(--color-surface)] py-[16px] px-[20px] text-center">
          <p
            class="m-0 mb-[4px] text-[11px] font-semibold text-[var(--color-expense)] uppercase tracking-[0.05em]"
          >
            Pengeluaran
          </p>
          <p class="m-0 text-[18px] font-bold text-[var(--color-expense)]">
            {{ formatRupiah(report.total_expense) }}
          </p>
        </div>
        <div class="bg-[var(--color-surface)] py-[16px] px-[20px] text-center">
          <p
            class="m-0 mb-[4px] text-[11px] font-semibold text-[var(--color-info)] uppercase tracking-[0.05em]"
          >
            Laba Bersih
          </p>
          <p class="m-0 text-[18px] font-bold text-[var(--color-info)]">
            {{ formatRupiah(report.net_profit) }}
          </p>
        </div>
      </div>

      <!-- Report text -->
      <div class="py-[20px] px-[24px]">
        <p
          class="m-0 mb-[12px] text-[13px] font-semibold text-[var(--color-text)]"
        >
          Ringkasan Teks
        </p>
        <div
          class="bg-[var(--color-bg)] border border-[var(--color-border)] rounded-[10px] py-[16px] px-[20px] text-[13px] leading-[1.7] text-[var(--color-text)]"
          v-html="parseWhatsAppText(report.report_text)"
        />
      </div>

      <!-- Top categories -->
      <div
        v-if="report.top_categories && report.top_categories.length"
        class="px-[24px] pb-[20px]"
      >
        <p
          class="m-0 mb-[10px] text-[13px] font-semibold text-[var(--color-text)]"
        >
          Kategori Teratas
        </p>
        <div class="flex flex-wrap gap-[8px]">
          <span
            v-for="(cat, i) in report.top_categories"
            :key="i"
            class="py-[6px] px-[14px] rounded-[999px] text-[12px] font-medium bg-[var(--color-bg)] text-[var(--color-text-secondary)] border border-[var(--color-border)]"
          >
            {{ cat }}
          </span>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
@media (max-width: 640px) {
  .summary-grid {
    grid-template-columns: 1fr !important;
  }
}
</style>
