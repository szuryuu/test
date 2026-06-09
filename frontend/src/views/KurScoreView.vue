<script setup>
import { ref, onMounted } from "vue";
import { kurAPI } from "@/api/kur";
import { formatRupiah } from "@/utils/format";
import KurScoreGauge from "@/components/kur/KurScoreGauge.vue";
import KurRecommendations from "@/components/kur/KurRecommendations.vue";

const score = ref(null);
const isLoading = ref(false);

onMounted(() => fetchScore());

async function fetchScore() {
  isLoading.value = true;
  try {
    const res = await kurAPI.getScore();
    score.value = res.data.data;
  } finally {
    isLoading.value = false;
  }
}

async function recalculate() {
  isLoading.value = true;
  try {
    const res = await kurAPI.recalculate();
    score.value = res.data.data;
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
          Skor KUR
        </h1>
        <p class="text-sm m-0" style="color: var(--color-text-secondary)">
          Skor kesiapan kredit usaha Anda
        </p>
      </div>
      <button
        @click="recalculate"
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
        <i class="pi pi-refresh" :class="{ 'pi-spin': isLoading }" style="font-size: 14px" />
        {{ isLoading ? "Menghitung..." : "Hitung Ulang" }}
      </button>
    </div>

    <div v-if="!score && !isLoading" class="text-center py-6" style="color: var(--color-text-secondary)">
      <p>Belum ada data skor KUR. Tambahkan transaksi terlebih dahulu.</p>
    </div>

    <div v-else class="grid" style="grid-template-columns: 1fr 1fr; gap: 16px">
      <!-- Left column: Gauge + Details -->
      <div class="flex flex-column gap-3">
        <KurScoreGauge :score="score?.score || 0" />

        <div
          class="p-4"
          style="
            background: var(--color-pure-white);
            border-radius: var(--radius-md);
            box-shadow: 0 1px 3px rgba(0,0,0,0.04);
          "
        >
          <h3 class="text-base m-0 mb-3" style="font-weight: 600; color: var(--color-near-black)">
            Detail Skor
          </h3>
          <div class="flex flex-column gap-2" v-if="score">
            <div class="flex justify-content-between py-2" style="border-bottom: 1px solid var(--color-border-soft)">
              <span class="text-sm" style="color: var(--color-text-secondary)">Rata-rata Pemasukan/Bulan</span>
              <span class="text-sm" style="font-weight: 600">{{ formatRupiah(score.monthly_income_avg) }}</span>
            </div>
            <div class="flex justify-content-between py-2" style="border-bottom: 1px solid var(--color-border-soft)">
              <span class="text-sm" style="color: var(--color-text-secondary)">Rata-rata Pengeluaran/Bulan</span>
              <span class="text-sm" style="font-weight: 600">{{ formatRupiah(score.monthly_expense_avg) }}</span>
            </div>
            <div class="flex justify-content-between py-2" style="border-bottom: 1px solid var(--color-border-soft)">
              <span class="text-sm" style="color: var(--color-text-secondary)">Margin Laba</span>
              <span class="text-sm" style="font-weight: 600; color: #10b981">{{ score.profit_margin }}%</span>
            </div>
            <div class="flex justify-content-between py-2" style="border-bottom: 1px solid var(--color-border-soft)">
              <span class="text-sm" style="color: var(--color-text-secondary)">Konsistensi Pencatatan</span>
              <span class="text-sm" style="font-weight: 600">{{ score.consistency_score }}/100</span>
            </div>
            <div class="flex justify-content-between py-2">
              <span class="text-sm" style="color: var(--color-text-secondary)">Bulan Data Tersedia</span>
              <span class="text-sm" style="font-weight: 600">{{ score.months_of_data }} bulan</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Right column: Recommendations -->
      <KurRecommendations :recommendations="score?.recommendations || []" />
    </div>
  </div>
</template>
