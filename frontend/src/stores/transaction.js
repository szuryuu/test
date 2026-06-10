import { defineStore } from "pinia";
import { ref } from "vue";
import { transactionAPI } from "@/api/transaction";

export const useTransactionStore = defineStore("transaction", () => {
  const transactions = ref([]);
  const total = ref(0);
  const page = ref(1);
  const isLoading = ref(false);
  const error = ref(null);

  async function fetchTransactions(params = {}) {
    isLoading.value = true;
    error.value = null;
    try {
      const res = await transactionAPI.list(params);
      transactions.value = res.data.data.transactions;
      total.value = res.data.data.total;
      page.value = res.data.data.page;
    } catch (err) {
      error.value = err.response?.data?.message || "Gagal memuat transaksi";
    } finally {
      isLoading.value = false;
    }
  }

  async function createTransaction(data) {
    const res = await transactionAPI.create(data);
    return res.data.data.transaction;
  }

  async function deleteTransaction(id) {
    await transactionAPI.remove(id);
    transactions.value = transactions.value.filter((t) => t.id !== id);
  }

  return { transactions, total, page, isLoading, error, fetchTransactions, createTransaction, deleteTransaction };
});
