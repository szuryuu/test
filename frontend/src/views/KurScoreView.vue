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
    <div class="flex justify-between items-start mb-[24px] flex-wrap gap-[12px]">
      <div>
        <h1 class="m-0 mb-[4px] text-[24px] font-bold text-[var(--color-text)] tracking-[-0.02em]">
          Skor KUR
        </h1>
        <p class="m-0 text-[14px] text-[var(--color-text-secondary)]">
          Skor kesiapan kredit usaha Anda
        </p>
      </div>

      <button
        @click="recalculate"
        :disabled="isLoading"
        class="flex items-center gap-[8px] py-[10px] px-[20px] text-[14px] font-semibold font-[inherit] border-0 rounded-[10px] cursor-pointer transition-all duration-[0.15s] ease bg-[linear-gradient(135deg,#10b981,#059669)] text-white"
        :style="{ opacity: isLoading ? 0.7 : 1 }"
        @mouseenter="!isLoading && ($event.target.style.boxShadow = '0 4px 12px rgba(16, 185, 129, 0.4)')"
        @mouseleave="!isLoading && ($event.target.style.boxShadow = 'none')"
      >
        <i class="pi pi-refresh text-[14px]" :class="{ 'pi-spin': isLoading }" />
        {{ isLoading ? "Menghitung..." : "Hitung Ulang" }}
      </button>
    </div>

    <!-- Loading state -->
    <div v-if="isLoading && !score" class="grid grid-cols-2 gap-[16px]">
      <div class="skeleton h-[300px] rounded-[var(--radius-card)]" />
      <div class="skeleton h-[300px] rounded-[var(--radius-card)]" />
    </div>

    <!-- Error state -->
    <div
      v-if="error"
      class="p-[16px_20px] bg-[var(--color-expense-bg)] border border-[rgba(239,68,68,0.2)] rounded-[var(--radius-card)] mb-[16px] flex items-start gap-[12px]"
    >
      <i class="pi pi-exclamation-triangle text-[var(--color-expense)] text-[18px] mt-[1px]" />
      <div>
        <p class="m-0 text-[14px] font-medium text-[#dc2626]">
          {{ error }}
        </p>
      </div>
    </div>

    <!-- Empty state -->
    <div
      v-if="!score && !isLoading && !error"
      class="bg-[var(--color-surface)] border border-[var(--color-border)] rounded-[var(--radius-card)] shadow-[var(--shadow-card)] py-[48px] px-[24px] text-center"
    >
      <div
        class="w-[64px] h-[64px] rounded-[16px] bg-[var(--color-bg)] mx-auto mb-[16px] flex items-center justify-center text-[28px] text-[var(--color-text-tertiary)]"
      >
        <i class="pi pi-star" />
      </div>
      <h3 class="m-0 mb-[8px] text-[18px] font-semibold text-[var(--color-text)]">
        Belum Ada Data Skor
      </h3>
      <p class="mx-auto mb-[20px] text-[14px] text-[var(--color-text-secondary)] max-w-[360px] leading-[1.6]">
        Tambahkan transaksi terlebih dahulu melalui WhatsApp atau halaman Transaksi, lalu hitung ulang skor KUR Anda.
      </p>
      <router-link
        to="/transactions"
        class="inline-flex items-center gap-[8px] py-[10px] px-[20px] text-[14px] font-semibold font-[inherit] border-0 rounded-[10px] cursor-pointer no-underline transition-all duration-[0.15s] ease bg-[linear-gradient(135deg,#10b981,#059669)] text-white"
      >
        <i class="pi pi-arrow-right text-[14px]" />
        Ke Halaman Transaksi
      </router-link>
    </div>

    <!-- Score data -->
    <div v-if="score" class="score-grid grid grid-cols-2 gap-[16px]">
      <!-- Staleness warning -->
      <div
        v-if="score.is_stale"
        class="col-span-2 p-[12px_16px] bg-[rgba(245,158,11,0.08)] border border-[rgba(245,158,11,0.25)] rounded-[var(--radius-card)] flex items-center gap-[10px] text-[13px] text-[#92400e]"
      >
        <i class="pi pi-clock text-[16px] shrink-0" />
        <span>Data mungkin sudah tidak akurat. Klik <strong>"Hitung Ulang"</strong> untuk memperbarui.</span>
      </div>

      <!-- Left column -->
      <div class="flex flex-col gap-[16px]">
        <KurScoreGauge :score="score.score || 0" />

        <div
          class="bg-[var(--color-surface)] border border-[var(--color-border)] rounded-[var(--radius-card)] p-[20px] shadow-[var(--shadow-card)]"
        >
          <h3 class="m-0 mb-[16px] text-[15px] font-semibold text-[var(--color-text)]">
            Detail Skor
          </h3>

          <div class="flex flex-col">
            <div
              v-for="(item, i) in [
                { label: 'Rata-rata Pemasukan/Bulan', value: formatRupiah(score.monthly_income_avg) },
                { label: 'Rata-rata Pengeluaran/Bulan', value: formatRupiah(score.monthly_expense_avg) },
                { label: 'Margin Laba', value: `${score.profit_margin}%`, color: 'var(--color-income)' },
                { label: 'Konsistensi Pencatatan', value: `${score.consistency_score}/100` },
                { label: 'Bulan Data Tersedia', value: `${score.months_of_data} bulan` },
              ]"
              :key="i"
              class="flex justify-between items-center py-[10px]"
              :style="{ borderBottom: i < 4 ? '1px solid var(--color-border-light)' : 'none' }"
            >
              <span class="text-[13px] text-[var(--color-text-secondary)]">
                {{ item.label }}
              </span>
              <span
                class="text-[13px] font-semibold"
                :style="{ color: item.color || 'var(--color-text)' }"
              >
                {{ item.value }}
              </span>
            </div>
          </div>

          <div
            v-if="score.calculated_at"
            class="mt-[12px] pt-[12px] border-t border-[var(--color-border-light)]"
          >
            <p class="m-0 text-[11px] text-[var(--color-text-tertiary)]">
              Terakhir dihitung: {{ new Date(score.calculated_at).toLocaleDateString('id-ID', { day: 'numeric', month: 'long', year: 'numeric', hour: '2-digit', minute: '2-digit' }) }}
            </p>
          </div>
        </div>
      </div>

      <!-- Right column -->
      <div class="flex flex-col gap-[16px]">
        <KurRecommendations :recommendations="score.recommendations || []" />
      </div>
    </div>
  </div>
</template>

<style scoped>
@media (max-width: 768px) {
  .score-grid {
    grid-template-columns: 1fr !important;
  }
}
</style>
