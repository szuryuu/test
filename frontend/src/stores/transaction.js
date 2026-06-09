import { defineStore } from "pinia";
import { ref } from "vue";
import { transactionAPI } from "@/api/transaction";

export const useTransactionStore = defineStore("transaction", () => {
  const transactions = ref([]);
  const total = ref(0);
  const page = ref(1);
  const isLoading = ref(false);

  async function fetchTransactions(params = {}) {
    isLoading.value = true;
    try {
      const res = await transactionAPI.list(params);
      transactions.value = res.data.data.transactions;
      total.value = res.data.data.total;
      page.value = res.data.data.page;
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

  return { transactions, total, page, isLoading, fetchTransactions, createTransaction, deleteTransaction };
});
