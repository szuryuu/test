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
    <div class="flex justify-content-between align-items-center mb-4">
      <div>
        <h1 class="text-2xl m-0 mb-1" style="font-weight: 700; color: var(--color-near-black)">
          Laporan Bulanan
        </h1>
        <p class="text-sm m-0" style="color: var(--color-text-secondary)">
          Generate laporan keuangan bulanan
        </p>
      </div>
    </div>

    <!-- Form -->
    <div
      class="p-4 mb-4"
      style="
        background: var(--color-pure-white);
        border-radius: var(--radius-md);
        box-shadow: 0 1px 3px rgba(0,0,0,0.04);
      "
    >
      <div class="flex align-items-end gap-3">
        <div class="flex flex-column gap-2">
          <label class="text-sm" style="font-weight: 500">Bulan</label>
          <select
            v-model="month"
            class="px-3 py-2 text-sm"
            style="
              border: 1px solid var(--color-border-soft);
              border-radius: var(--radius-sm);
              min-width: 140px;
            "
          >
            <option v-for="(name, i) in monthNames" :key="i" :value="i + 1">
              {{ name }}
            </option>
          </select>
        </div>
        <div class="flex flex-column gap-2">
          <label class="text-sm" style="font-weight: 500">Tahun</label>
          <input
            v-model="year"
            type="number"
            class="px-3 py-2 text-sm"
            style="
              border: 1px solid var(--color-border-soft);
              border-radius: var(--radius-sm);
              width: 100px;
            "
          />
        </div>
        <button
          @click="generateReport"
          :disabled="isLoading"
          class="flex align-items-center gap-2 px-4 py-3 text-sm cursor-pointer transition-all transition-duration-200"
          style="
            background: var(--color-apple-blue);
            color: var(--color-pure-white);
            border: none;
            border-radius: var(--radius-sm);
            font-weight: 600;
          "
        >
          <i class="pi pi-file" style="font-size: 14px" />
          {{ isLoading ? "Membuat..." : "Buat Laporan" }}
        </button>
      </div>
    </div>

    <!-- Error -->
    <div
      v-if="error"
      class="px-4 py-3 mb-4 text-sm"
      style="background: #fef2f2; color: #dc2626; border-radius: var(--radius-sm)"
    >
      {{ error }}
    </div>

    <!-- Result -->
    <div
      v-if="report"
      class="p-4"
      style="
        background: var(--color-pure-white);
        border-radius: var(--radius-md);
        box-shadow: 0 1px 3px rgba(0,0,0,0.04);
      "
    >
      <h3 class="text-base m-0 mb-3" style="font-weight: 600; color: var(--color-near-black)">
        Laporan {{ monthNames[month - 1] }} {{ year }}
      </h3>
      <pre
        class="text-sm p-4 m-0"
        style="
          background: var(--color-pale-gray);
          border-radius: var(--radius-sm);
          color: var(--color-near-black);
          white-space: pre-wrap;
          line-height: 1.6;
        "
      >{{ report.report_text }}</pre>

      <div class="grid mt-3" style="grid-template-columns: repeat(3, 1fr); gap: 12px" v-if="report.total_income">
        <div class="text-center p-3" style="background: #ecfdf5; border-radius: var(--radius-sm)">
          <p class="text-xs m-0 mb-1" style="color: #10b981; font-weight: 500">PEMASUKAN</p>
          <p class="text-lg m-0" style="font-weight: 700">{{ formatRupiah(report.total_income) }}</p>
        </div>
        <div class="text-center p-3" style="background: #fef2f2; border-radius: var(--radius-sm)">
          <p class="text-xs m-0 mb-1" style="color: #ef4444; font-weight: 500">PENGELUARAN</p>
          <p class="text-lg m-0" style="font-weight: 700">{{ formatRupiah(report.total_expense) }}</p>
        </div>
        <div class="text-center p-3" style="background: #eff6ff; border-radius: var(--radius-sm)">
          <p class="text-xs m-0 mb-1" style="color: var(--color-apple-blue); font-weight: 500">LABA BERSIH</p>
          <p class="text-lg m-0" style="font-weight: 700">{{ formatRupiah(report.net_profit) }}</p>
        </div>
      </div>
    </div>
  </div>
</template>
