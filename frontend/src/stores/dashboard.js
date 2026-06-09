import { defineStore } from "pinia";
import { ref } from "vue";
import { dashboardAPI } from "@/api/dashboard";

export const useDashboardStore = defineStore("dashboard", () => {
  const summary = ref(null);
  const categories = ref([]);
  const isLoading = ref(false);

  async function fetchSummary(period = "monthly", date = "") {
    isLoading.value = true;
    try {
      const res = await dashboardAPI.summary(period, date);
      summary.value = res.data.data;
    } finally {
      isLoading.value = false;
    }
  }

  async function fetchCategories(period = "monthly", type = "expense") {
    const res = await dashboardAPI.categories(period, type);
    categories.value = res.data.data.categories;
  }

  return { summary, categories, isLoading, fetchSummary, fetchCategories };
});
