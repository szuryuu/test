<script setup>
import { ref } from "vue";
import { reportAPI } from "@/api/kur";
import { formatRupiah } from "@/utils/format";

const year = ref(new Date().getFullYear());
const month = ref(new Date().getMonth() + 1);
const isLoading = ref(false);
const report = ref(null);
const error = ref("");

const monthNames = [
  "Januari", "Februari", "Maret", "April", "Mei", "Juni",
  "Juli", "Agustus", "September", "Oktober", "November", "Desember",
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
        Laporan Bulanan
      </h1>
      <p style="margin: 0; font-size: 14px; color: var(--color-text-secondary)">
        Generate laporan keuangan bulanan dan kirim via WhatsApp
      </p>
    </div>

    <!-- Generator card -->
    <div
      style="
        background: var(--color-surface);
        border: 1px solid var(--color-border);
        border-radius: var(--radius-card);
        padding: 24px;
        box-shadow: var(--shadow-card);
        margin-bottom: 20px;
      "
    >
      <h3
        style="
          margin: 0 0 16px; font-size: 15px; font-weight: 600;
          color: var(--color-text);
        "
      >
        <i class="pi pi-calendar" style="margin-right: 8px; color: var(--color-brand-500)" />
        Pilih Periode Laporan
      </h3>

      <div
        style="
          display: flex; flex-wrap: wrap; align-items: flex-end;
          gap: 12px;
        "
      >
        <div style="display: flex; flex-direction: column; gap: 6px; min-width: 160px">
          <label style="font-size: 12px; font-weight: 500; color: var(--color-text-secondary)">
            Bulan
          </label>
          <select
            v-model="month"
            style="
              padding: 10px 12px; font-size: 14px; font-family: inherit;
              border: 1px solid var(--color-border); border-radius: 8px;
              background: var(--color-bg); color: var(--color-text);
              cursor: pointer;
            "
          >
            <option v-for="(name, i) in monthNames" :key="i" :value="i + 1">
              {{ name }}
            </option>
          </select>
        </div>

        <div style="display: flex; flex-direction: column; gap: 6px; min-width: 100px">
          <label style="font-size: 12px; font-weight: 500; color: var(--color-text-secondary)">
            Tahun
          </label>
          <input
            v-model="year"
            type="number"
            min="2024"
            max="2030"
            style="
              padding: 10px 12px; font-size: 14px; font-family: inherit;
              border: 1px solid var(--color-border); border-radius: 8px;
              background: var(--color-bg); color: var(--color-text);
            "
          />
        </div>

        <button
          @click="generateReport"
          :disabled="isLoading"
          style="
            display: flex; align-items: center; gap: 8px;
            padding: 10px 24px; font-size: 14px; font-weight: 600;
            font-family: inherit; border: none; border-radius: 8px;
            cursor: pointer; transition: all 0.15s ease;
            background: linear-gradient(135deg, #10b981, #059669);
            color: white;
          "
          :style="{ opacity: isLoading ? 0.7 : 1 }"
          @mouseenter="!isLoading && ($event.target.style.boxShadow = '0 4px 12px rgba(16, 185, 129, 0.4)')"
          @mouseleave="!isLoading && ($event.target.style.boxShadow = 'none')"
        >
          <i class="pi pi-file" style="font-size: 14px" />
          {{ isLoading ? "Membuat Laporan..." : "Buat Laporan" }}
        </button>
      </div>
    </div>

    <!-- Loading skeleton -->
    <div v-if="isLoading" style="display: flex; flex-direction: column; gap: 16px">
      <div class="skeleton" style="height: 200px; border-radius: var(--radius-card)" />
      <div class="skeleton" style="height: 80px; border-radius: var(--radius-card)" />
    </div>

    <!-- Error -->
    <div
      v-if="error"
      style="
        padding: 16px 20px; margin-bottom: 20px;
        background: var(--color-expense-bg);
        border: 1px solid rgba(239, 68, 68, 0.2);
        border-radius: var(--radius-card);
        display: flex; align-items: flex-start; gap: 12px;
      "
    >
      <i class="pi pi-exclamation-triangle" style="color: var(--color-expense); font-size: 18px; margin-top: 1px" />
      <div>
        <p style="margin: 0; font-size: 14px; font-weight: 500; color: #dc2626">
          {{ error }}
        </p>
      </div>
    </div>

    <!-- Report result -->
    <div
      v-if="report"
      style="
        background: var(--color-surface);
        border: 1px solid var(--color-border);
        border-radius: var(--radius-card);
        box-shadow: var(--shadow-card);
        overflow: hidden;
      "
    >
      <!-- Report header -->
      <div
        style="
          padding: 20px 24px;
          border-bottom: 1px solid var(--color-border);
          display: flex; justify-content: space-between;
          align-items: center; flex-wrap: wrap; gap: 8px;
        "
      >
        <div>
          <h3
            style="
              margin: 0; font-size: 16px; font-weight: 600;
              color: var(--color-text);
            "
          >
            Laporan {{ monthNames[month - 1] }} {{ year }}
          </h3>
          <p style="margin: 2px 0 0; font-size: 12px; color: var(--color-text-tertiary)">
            Ringkasan keuangan bulanan
          </p>
        </div>
        <span
          style="
            font-size: 11px; padding: 4px 10px;
            background: var(--color-income-bg); color: var(--color-income);
            border-radius: 999px; font-weight: 500;
          "
        >
          Siap dikirim via WhatsApp
        </span>
      </div>

      <!-- Summary numbers -->
      <div
        v-if="report.total_income"
        style="
          display: grid;
          grid-template-columns: repeat(3, 1fr);
          gap: 1px;
          background: var(--color-border);
        "
      >
        <div style="background: var(--color-surface); padding: 16px 20px; text-align: center">
          <p style="margin: 0 0 4px; font-size: 11px; font-weight: 600; color: var(--color-income); text-transform: uppercase; letter-spacing: 0.05em">
            Pemasukan
          </p>
          <p style="margin: 0; font-size: 18px; font-weight: 700; color: var(--color-income)">
            {{ formatRupiah(report.total_income) }}
          </p>
        </div>
        <div style="background: var(--color-surface); padding: 16px 20px; text-align: center">
          <p style="margin: 0 0 4px; font-size: 11px; font-weight: 600; color: var(--color-expense); text-transform: uppercase; letter-spacing: 0.05em">
            Pengeluaran
          </p>
          <p style="margin: 0; font-size: 18px; font-weight: 700; color: var(--color-expense)">
            {{ formatRupiah(report.total_expense) }}
          </p>
        </div>
        <div style="background: var(--color-surface); padding: 16px 20px; text-align: center">
          <p style="margin: 0 0 4px; font-size: 11px; font-weight: 600; color: var(--color-info); text-transform: uppercase; letter-spacing: 0.05em">
            Laba Bersih
          </p>
          <p style="margin: 0; font-size: 18px; font-weight: 700; color: var(--color-info)">
            {{ formatRupiah(report.net_profit) }}
          </p>
        </div>
      </div>

      <!-- Report text -->
      <div style="padding: 20px 24px">
        <p
          style="
            margin: 0 0 12px; font-size: 13px; font-weight: 600;
            color: var(--color-text);
          "
        >
          Ringkasan Teks
        </p>
        <div
          style="
            background: var(--color-bg);
            border: 1px solid var(--color-border);
            border-radius: 10px;
            padding: 16px 20px;
            font-size: 13px; line-height: 1.7;
            color: var(--color-text);
            white-space: pre-wrap;
          "
        >
          {{ report.report_text }}
        </div>
      </div>

      <!-- Top categories -->
      <div
        v-if="report.top_categories && report.top_categories.length"
        style="
          padding: 0 24px 20px;
        "
      >
        <p
          style="
            margin: 0 0 10px; font-size: 13px; font-weight: 600;
            color: var(--color-text);
          "
        >
          Kategori Teratas
        </p>
        <div style="display: flex; flex-wrap: wrap; gap: 8px">
          <span
            v-for="(cat, i) in report.top_categories"
            :key="i"
            style="
              padding: 6px 14px; border-radius: 999px;
              font-size: 12px; font-weight: 500;
              background: var(--color-bg);
              color: var(--color-text-secondary);
              border: 1px solid var(--color-border);
            "
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
  div[style*="grid-template-columns: repeat(3, 1fr)"] {
    grid-template-columns: 1fr !important;
  }
}
</style>
