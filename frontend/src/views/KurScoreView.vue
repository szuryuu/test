<script setup>
import { ref, onMounted } from "vue";
import { kurAPI } from "@/api/kur";
import { formatRupiah } from "@/utils/format";
import KurScoreGauge from "@/components/kur/KurScoreGauge.vue";
import KurRecommendations from "@/components/kur/KurRecommendations.vue";

const score = ref(null);
const isLoading = ref(false);
const error = ref("");

onMounted(() => fetchScore());

async function fetchScore() {
  isLoading.value = true;
  error.value = "";
  try {
    const res = await kurAPI.getScore();
    score.value = res.data.data;
  } catch (err) {
    if (err.response?.status === 404) {
      score.value = null;
    } else {
      error.value = err.response?.data?.message || "Gagal memuat skor KUR";
    }
  } finally {
    isLoading.value = false;
  }
}

async function recalculate() {
  isLoading.value = true;
  error.value = "";
  try {
    const res = await kurAPI.recalculate();
    score.value = res.data.data;
  } catch (err) {
    error.value = err.response?.data?.message || "Gagal menghitung ulang skor";
  } finally {
    isLoading.value = false;
  }
}
</script>

<template>
  <div>
    <!-- Page header -->
    <div
      style="
        display: flex; justify-content: space-between;
        align-items: flex-start; margin-bottom: 24px;
        flex-wrap: wrap; gap: 12px;
      "
    >
      <div>
        <h1
          style="
            margin: 0 0 4px; font-size: 24px; font-weight: 700;
            color: var(--color-text); letter-spacing: -0.02em;
          "
        >
          Skor KUR
        </h1>
        <p style="margin: 0; font-size: 14px; color: var(--color-text-secondary)">
          Skor kesiapan kredit usaha Anda
        </p>
      </div>

      <button
        @click="recalculate"
        :disabled="isLoading"
        style="
          display: flex; align-items: center; gap: 8px;
          padding: 10px 20px; font-size: 14px; font-weight: 600;
          font-family: inherit; border: none; border-radius: 10px;
          cursor: pointer; transition: all 0.15s ease;
          background: linear-gradient(135deg, #10b981, #059669);
          color: white; opacity: 1;
        "
        :style="{ opacity: isLoading ? 0.7 : 1 }"
        @mouseenter="!isLoading && ($event.target.style.boxShadow = '0 4px 12px rgba(16, 185, 129, 0.4)')"
        @mouseleave="!isLoading && ($event.target.style.boxShadow = 'none')"
      >
        <i
          class="pi pi-refresh"
          :class="{ 'pi-spin': isLoading }"
          style="font-size: 14px"
        />
        {{ isLoading ? "Menghitung..." : "Hitung Ulang" }}
      </button>
    </div>

    <!-- Loading state -->
    <div
      v-if="isLoading && !score"
      style="display: grid; grid-template-columns: 1fr 1fr; gap: 16px"
    >
      <div class="skeleton" style="height: 300px; border-radius: var(--radius-card)" />
      <div class="skeleton" style="height: 300px; border-radius: var(--radius-card)" />
    </div>

    <!-- Error state -->
    <div
      v-if="error"
      style="
        padding: 16px 20px; background: var(--color-expense-bg);
        border: 1px solid rgba(239, 68, 68, 0.2);
        border-radius: var(--radius-card); margin-bottom: 16px;
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

    <!-- Empty state -->
    <div
      v-if="!score && !isLoading && !error"
      style="
        background: var(--color-surface);
        border: 1px solid var(--color-border);
        border-radius: var(--radius-card);
        box-shadow: var(--shadow-card);
        padding: 48px 24px; text-align: center;
      "
    >
      <div
        style="
          width: 64px; height: 64px; border-radius: 16px;
          background: var(--color-bg); margin: 0 auto 16px;
          display: flex; align-items: center; justify-content: center;
          font-size: 28px; color: var(--color-text-tertiary);
        "
      >
        <i class="pi pi-star" />
      </div>
      <h3
        style="
          margin: 0 0 8px; font-size: 18px; font-weight: 600;
          color: var(--color-text);
        "
      >
        Belum Ada Data Skor
      </h3>
      <p
        style="
          margin: 0 auto 20px; font-size: 14px; color: var(--color-text-secondary);
          max-width: 360px; line-height: 1.6;
        "
      >
        Tambahkan transaksi terlebih dahulu melalui WhatsApp atau halaman Transaksi, lalu hitung ulang skor KUR Anda.
      </p>
      <router-link
        to="/transactions"
        style="
          display: inline-flex; align-items: center; gap: 8px;
          padding: 10px 20px; font-size: 14px; font-weight: 600;
          font-family: inherit; border: none; border-radius: 10px;
          cursor: pointer; text-decoration: none; transition: all 0.15s ease;
          background: linear-gradient(135deg, #10b981, #059669);
          color: white;
        "
      >
        <i class="pi pi-arrow-right" style="font-size: 14px" />
        Ke Halaman Transaksi
      </router-link>
    </div>

    <!-- Score data -->
    <div
      v-if="score"
      style="display: grid; grid-template-columns: 1fr 1fr; gap: 16px"
    >
      <!-- Left column -->
      <div style="display: flex; flex-direction: column; gap: 16px">
        <KurScoreGauge :score="score.score || 0" />

        <div
          style="
            background: var(--color-surface);
            border: 1px solid var(--color-border);
            border-radius: var(--radius-card);
            padding: 20px;
            box-shadow: var(--shadow-card);
          "
        >
          <h3
            style="
              margin: 0 0 16px; font-size: 15px; font-weight: 600;
              color: var(--color-text);
            "
          >
            Detail Skor
          </h3>

          <div style="display: flex; flex-direction: column; gap: 0">
            <div
              v-for="(item, i) in [
                { label: 'Rata-rata Pemasukan/Bulan', value: formatRupiah(score.monthly_income_avg) },
                { label: 'Rata-rata Pengeluaran/Bulan', value: formatRupiah(score.monthly_expense_avg) },
                { label: 'Margin Laba', value: `${score.profit_margin}%`, color: 'var(--color-income)' },
                { label: 'Konsistensi Pencatatan', value: `${score.consistency_score}/100` },
                { label: 'Bulan Data Tersedia', value: `${score.months_of_data} bulan` },
              ]"
              :key="i"
              style="
                display: flex; justify-content: space-between;
                align-items: center; padding: 10px 0;
              "
              :style="{ borderBottom: i < 4 ? '1px solid var(--color-border-light)' : 'none' }"
            >
              <span style="font-size: 13px; color: var(--color-text-secondary)">
                {{ item.label }}
              </span>
              <span
                style="font-size: 13px; font-weight: 600;"
                :style="{ color: item.color || 'var(--color-text)' }"
              >
                {{ item.value }}
              </span>
            </div>
          </div>

          <div
            v-if="score.calculated_at"
            style="margin-top: 12px; padding-top: 12px; border-top: 1px solid var(--color-border-light)"
          >
            <p style="margin: 0; font-size: 11px; color: var(--color-text-tertiary)">
              Terakhir dihitung: {{ new Date(score.calculated_at).toLocaleDateString('id-ID', { day: 'numeric', month: 'long', year: 'numeric', hour: '2-digit', minute: '2-digit' }) }}
            </p>
          </div>
        </div>
      </div>

      <!-- Right column -->
      <div style="display: flex; flex-direction: column; gap: 16px">
        <KurRecommendations :recommendations="score.recommendations || []" />
      </div>
    </div>
  </div>
</template>

<style scoped>
@media (max-width: 768px) {
  div[style*="grid-template-columns: 1fr 1fr"] {
    grid-template-columns: 1fr !important;
  }
}
</style>
