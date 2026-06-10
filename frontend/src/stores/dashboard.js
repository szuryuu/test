import { defineStore } from "pinia";
import { ref } from "vue";
import { dashboardAPI } from "@/api/dashboard";

export const useDashboardStore = defineStore("dashboard", () => {
  const summary = ref(null);
  const categories = ref([]);
  const isLoading = ref(false);
  const error = ref(null);

  async function fetchSummary(period = "monthly", date = "") {
    isLoading.value = true;
    error.value = null;
    try {
      const res = await dashboardAPI.summary(period, date);
      summary.value = res.data.data;
    } catch (err) {
      error.value = err.response?.data?.message || "Gagal memuat ringkasan";
    } finally {
      isLoading.value = false;
    }
  }

  async function fetchCategories(period = "monthly", type = "expense") {
    try {
      const res = await dashboardAPI.categories(period, type);
      categories.value = res.data.data.categories;
    } catch (err) {
      error.value = err.response?.data?.message || "Gagal memuat kategori";
    }
  }

  return { summary, categories, isLoading, error, fetchSummary, fetchCategories };
});
