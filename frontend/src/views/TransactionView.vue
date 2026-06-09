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
    newTx.value = { amount: "", type: "income", description: "", transaction_date: new Date().toISOString().slice(0, 10) };
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
    <div class="flex justify-content-between align-items-center mb-4">
      <div>
        <h1 class="text-2xl m-0 mb-1" style="font-weight: 700; color: var(--color-near-black)">
          Transaksi
        </h1>
        <p class="text-sm m-0" style="color: var(--color-text-secondary)">
          {{ store.total }} transaksi tercatat
        </p>
      </div>
      <button
        @click="showCreate = !showCreate"
        class="flex align-items-center gap-2 px-4 py-3 text-sm cursor-pointer transition-all transition-duration-200"
        style="
          background: var(--color-apple-blue);
          color: var(--color-pure-white);
          border: none;
          border-radius: var(--radius-sm);
          font-weight: 600;
        "
      >
        <i class="pi pi-plus" style="font-size: 14px" />
        Tambah Transaksi
      </button>
    </div>

    <!-- Create modal -->
    <div
      v-if="showCreate"
      class="mb-4 p-4"
      style="
        background: var(--color-pure-white);
        border-radius: var(--radius-md);
        box-shadow: 0 1px 3px rgba(0,0,0,0.04);
      "
    >
      <h3 class="text-base m-0 mb-3" style="font-weight: 600; color: var(--color-near-black)">
        Catat Transaksi Manual
      </h3>
      <div class="grid mb-3" style="grid-template-columns: 1fr 1fr; gap: 12px">
        <div class="flex flex-column gap-2">
          <label class="text-sm" style="font-weight: 500">Jumlah (Rp)</label>
          <input
            v-model="newTx.amount"
            type="number"
            placeholder="50000"
            class="px-3 py-2 text-sm"
            style="
              border: 1px solid var(--color-border-soft);
              border-radius: var(--radius-sm);
              background: var(--color-pale-gray);
            "
          />
        </div>
        <div class="flex flex-column gap-2">
          <label class="text-sm" style="font-weight: 500">Tipe</label>
          <select
            v-model="newTx.type"
            class="px-3 py-2 text-sm"
            style="
              border: 1px solid var(--color-border-soft);
              border-radius: var(--radius-sm);
              background: var(--color-pale-gray);
            "
          >
            <option value="income">Pemasukan</option>
            <option value="expense">Pengeluaran</option>
          </select>
        </div>
      </div>
      <div class="grid mb-3" style="grid-template-columns: 2fr 1fr; gap: 12px">
        <div class="flex flex-column gap-2">
          <label class="text-sm" style="font-weight: 500">Deskripsi</label>
          <input
            v-model="newTx.description"
            type="text"
            placeholder="Penjualan nasi uduk"
            class="px-3 py-2 text-sm"
            style="
              border: 1px solid var(--color-border-soft);
              border-radius: var(--radius-sm);
              background: var(--color-pale-gray);
            "
          />
        </div>
        <div class="flex flex-column gap-2">
          <label class="text-sm" style="font-weight: 500">Tanggal</label>
          <input
            v-model="newTx.transaction_date"
            type="date"
            class="px-3 py-2 text-sm"
            style="
              border: 1px solid var(--color-border-soft);
              border-radius: var(--radius-sm);
              background: var(--color-pale-gray);
            "
          />
        </div>
      </div>
      <div class="flex gap-2">
        <button
          @click="handleCreate"
          class="px-4 py-2 text-sm cursor-pointer"
          style="
            background: var(--color-apple-blue);
            color: var(--color-pure-white);
            border: none;
            border-radius: var(--radius-sm);
            font-weight: 600;
          "
        >
          Simpan
        </button>
        <button
          @click="showCreate = false"
          class="px-4 py-2 text-sm cursor-pointer"
          style="
            background: transparent;
            border: 1px solid var(--color-border-soft);
            border-radius: var(--radius-sm);
            color: var(--color-near-black);
          "
        >
          Batal
        </button>
      </div>
    </div>

    <TransactionTable @delete="handleDelete" />
  </div>
</template>
