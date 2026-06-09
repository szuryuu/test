<script setup>
import { ref, onMounted } from "vue";
import { useTransactionStore } from "@/stores/transaction";
import TransactionTable from "@/components/transaction/TransactionTable.vue";

const store = useTransactionStore();

onMounted(() => {
  store.fetchTransactions({ page: 1, limit: 20 });
});

const showCreate = ref(false);
const newTx = ref({
  amount: "",
  type: "income",
  description: "",
  transaction_date: new Date().toISOString().slice(0, 10),
});

async function handleCreate() {
  if (!newTx.value.amount || !newTx.value.description) return;
  try {
    await store.createTransaction({
      amount: parseInt(newTx.value.amount),
      type: newTx.value.type,
      description: newTx.value.description,
      transaction_date: newTx.value.transaction_date,
    });
    showCreate.value = false;
    newTx.value = {
      amount: "",
      type: "income",
      description: "",
      transaction_date: new Date().toISOString().slice(0, 10),
    };
    store.fetchTransactions({ page: 1, limit: 20 });
  } catch (e) {
    // handled by interceptor
  }
}

async function handleDelete(id) {
  if (!confirm("Hapus transaksi ini?")) return;
  await store.deleteTransaction(id);
  store.fetchTransactions({ page: 1, limit: 20 });
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
          Transaksi
        </h1>
        <p style="margin: 0; font-size: 14px; color: var(--color-text-secondary)">
          {{ store.total }} transaksi tercatat
        </p>
      </div>

      <div style="display: flex; gap: 8px">
        <button
          @click="showCreate = !showCreate"
          style="
            display: flex; align-items: center; gap: 8px;
            padding: 10px 20px; font-size: 14px; font-weight: 600;
            font-family: inherit; border: none; border-radius: 10px;
            cursor: pointer; transition: all 0.15s ease;
            background: linear-gradient(135deg, #10b981, #059669);
            color: white;
          "
          @mouseenter="$event.target.style.boxShadow = '0 4px 12px rgba(16, 185, 129, 0.4)'"
          @mouseleave="$event.target.style.boxShadow = 'none'"
        >
          <i class="pi pi-plus" style="font-size: 14px" />
          Tambah Transaksi
        </button>
      </div>
    </div>

    <!-- Create form -->
    <div
      v-if="showCreate"
      style="
        margin-bottom: 20px; padding: 20px;
        background: var(--color-surface);
        border: 1px solid var(--color-border);
        border-radius: var(--radius-card);
        box-shadow: var(--shadow-elevated);
      "
    >
      <div
        style="
          display: flex; justify-content: space-between;
          align-items: center; margin-bottom: 16px;
        "
      >
        <h3
          style="
            margin: 0; font-size: 15px; font-weight: 600;
            color: var(--color-text);
          "
        >
          <i class="pi pi-plus-circle" style="margin-right: 8px; color: var(--color-brand-500)" />
          Catat Transaksi Manual
        </h3>
        <button
          @click="showCreate = false"
          style="
            background: transparent; border: none; cursor: pointer;
            font-size: 18px; color: var(--color-text-tertiary);
            padding: 4px;
          "
        >
          <i class="pi pi-times" />
        </button>
      </div>

      <div
        style="
          display: grid;
          grid-template-columns: 1fr 1fr;
          gap: 12px;
          margin-bottom: 12px;
        "
      >
        <div style="display: flex; flex-direction: column; gap: 6px">
          <label style="font-size: 12px; font-weight: 500; color: var(--color-text-secondary)">
            Jumlah (Rp)
          </label>
          <div style="position: relative">
            <span
              style="
                position: absolute; left: 12px; top: 50%;
                transform: translateY(-50%);
                font-size: 13px; color: var(--color-text-tertiary);
              "
            >
              Rp
            </span>
            <input
              v-model="newTx.amount"
              type="number"
              placeholder="50000"
              style="
                width: 100%; padding: 10px 12px 10px 36px;
                font-size: 14px; font-family: inherit;
                border: 1px solid var(--color-border); border-radius: 8px;
                background: var(--color-bg); color: var(--color-text);
                transition: border-color 0.15s ease;
              "
              @focus="$event.target.style.borderColor = 'var(--color-brand-500)'"
              @blur="$event.target.style.borderColor = 'var(--color-border)'"
            />
          </div>
        </div>

        <div style="display: flex; flex-direction: column; gap: 6px">
          <label style="font-size: 12px; font-weight: 500; color: var(--color-text-secondary)">
            Tipe
          </label>
          <div style="display: flex; gap: 8px; height: 42px;">
            <button
              @click="newTx.type = 'income'"
              style="
                flex: 1; border-radius: 8px; font-size: 13px;
                font-weight: 500; font-family: inherit; cursor: pointer;
                transition: all 0.15s ease; border: 1px solid var(--color-border);
              "
              :style="{
                background: newTx.type === 'income' ? 'var(--color-income-bg)' : 'var(--color-surface)',
                color: newTx.type === 'income' ? 'var(--color-income)' : 'var(--color-text-secondary)',
                borderColor: newTx.type === 'income' ? 'var(--color-income)' : 'var(--color-border)',
              }"
            >
              + Pemasukan
            </button>
            <button
              @click="newTx.type = 'expense'"
              style="
                flex: 1; border-radius: 8px; font-size: 13px;
                font-weight: 500; font-family: inherit; cursor: pointer;
                transition: all 0.15s ease; border: 1px solid var(--color-border);
              "
              :style="{
                background: newTx.type === 'expense' ? 'var(--color-expense-bg)' : 'var(--color-surface)',
                color: newTx.type === 'expense' ? 'var(--color-expense)' : 'var(--color-text-secondary)',
                borderColor: newTx.type === 'expense' ? 'var(--color-expense)' : 'var(--color-border)',
              }"
            >
              − Pengeluaran
            </button>
          </div>
        </div>
      </div>

      <div
        style="
          display: grid;
          grid-template-columns: 2fr 1fr;
          gap: 12px;
          margin-bottom: 20px;
        "
      >
        <div style="display: flex; flex-direction: column; gap: 6px">
          <label style="font-size: 12px; font-weight: 500; color: var(--color-text-secondary)">
            Deskripsi
          </label>
          <input
            v-model="newTx.description"
            type="text"
            placeholder="Penjualan nasi uduk"
            style="
              width: 100%; padding: 10px 12px; font-size: 14px; font-family: inherit;
              border: 1px solid var(--color-border); border-radius: 8px;
              background: var(--color-bg); color: var(--color-text);
              transition: border-color 0.15s ease;
            "
            @focus="$event.target.style.borderColor = 'var(--color-brand-500)'"
            @blur="$event.target.style.borderColor = 'var(--color-border)'"
          />
        </div>

        <div style="display: flex; flex-direction: column; gap: 6px">
          <label style="font-size: 12px; font-weight: 500; color: var(--color-text-secondary)">
            Tanggal
          </label>
          <input
            v-model="newTx.transaction_date"
            type="date"
            style="
              width: 100%; padding: 10px 12px; font-size: 14px; font-family: inherit;
              border: 1px solid var(--color-border); border-radius: 8px;
              background: var(--color-bg); color: var(--color-text);
              transition: border-color 0.15s ease;
            "
            @focus="$event.target.style.borderColor = 'var(--color-brand-500)'"
            @blur="$event.target.style.borderColor = 'var(--color-border)'"
          />
        </div>
      </div>

      <div style="display: flex; gap: 8px">
        <button
          @click="handleCreate"
          style="
            padding: 10px 24px; font-size: 14px; font-weight: 600;
            font-family: inherit; border: none; border-radius: 8px;
            cursor: pointer; transition: all 0.15s ease;
            background: linear-gradient(135deg, #10b981, #059669);
            color: white;
          "
          @mouseenter="$event.target.style.boxShadow = '0 4px 12px rgba(16, 185, 129, 0.4)'"
          @mouseleave="$event.target.style.boxShadow = 'none'"
        >
          <i class="pi pi-check" style="margin-right: 6px" />
          Simpan
        </button>
        <button
          @click="showCreate = false"
          style="
            padding: 10px 20px; font-size: 14px; font-weight: 500;
            font-family: inherit; border: 1px solid var(--color-border);
            border-radius: 8px; cursor: pointer; transition: all 0.15s ease;
            background: transparent; color: var(--color-text-secondary);
          "
          @mouseenter="$event.target.style.background = 'var(--color-bg)'"
          @mouseleave="$event.target.style.background = 'transparent'"
        >
          Batal
        </button>
      </div>
    </div>

    <!-- Table -->
    <TransactionTable @delete="handleDelete" />
  </div>
</template>

<style scoped>
@media (max-width: 640px) {
  div[style*="grid-template-columns: 1fr 1fr"] {
    grid-template-columns: 1fr !important;
  }
  div[style*="grid-template-columns: 2fr 1fr"] {
    grid-template-columns: 1fr !important;
  }
}
</style>
