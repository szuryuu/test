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
    <div class="flex justify-between items-start mb-[24px] flex-wrap gap-[12px]">
      <div>
        <h1 class="m-0 mb-[4px] text-[24px] font-bold text-[var(--color-text)] tracking-[-0.02em]">
          Transaksi
        </h1>
        <p class="m-0 text-[14px] text-[var(--color-text-secondary)]">
          {{ store.total }} transaksi tercatat
        </p>
      </div>

      <div class="flex gap-[8px]">
        <button
          @click="showCreate = !showCreate"
          class="flex items-center gap-[8px] py-[10px] px-[20px] text-[14px] font-semibold font-[inherit] border-0 rounded-[10px] cursor-pointer transition-all duration-[0.15s] ease bg-[linear-gradient(135deg,#10b981,#059669)] text-white"
          @mouseenter="$event.target.style.boxShadow = '0 4px 12px rgba(16, 185, 129, 0.4)'"
          @mouseleave="$event.target.style.boxShadow = 'none'"
        >
          <i class="pi pi-plus text-[14px]" />
          Tambah Transaksi
        </button>
      </div>
    </div>

    <!-- Error state -->
    <div
      v-if="store.error"
      class="p-[12px_16px] bg-[var(--color-expense-bg)] border border-[rgba(239,68,68,0.2)] rounded-[10px] mb-[16px] text-[13px] text-[#dc2626]"
    >
      <i class="pi pi-exclamation-triangle mr-[8px]" />
      {{ store.error }}
    </div>

    <!-- Create form -->
    <div
      v-if="showCreate"
      class="mb-[20px] p-[20px] bg-[var(--color-surface)] border border-[var(--color-border)] rounded-[var(--radius-card)] shadow-[var(--shadow-elevated)]"
    >
      <div class="flex justify-between items-center mb-[16px]">
        <h3 class="m-0 text-[15px] font-semibold text-[var(--color-text)]">
          <i class="pi pi-plus-circle mr-[8px] text-[var(--color-brand-500)]" />
          Catat Transaksi Manual
        </h3>
        <button
          @click="showCreate = false"
          class="bg-transparent border-0 cursor-pointer text-[18px] text-[var(--color-text-tertiary)] p-[4px]"
        >
          <i class="pi pi-times" />
        </button>
      </div>

      <div class="create-form-grid-type grid grid-cols-2 gap-[12px] mb-[12px]">
        <div class="flex flex-col gap-[6px]">
          <label class="text-[12px] font-medium text-[var(--color-text-secondary)]">
            Jumlah (Rp)
          </label>
          <div class="relative">
            <span
              class="absolute left-[12px] top-[50%] -translate-y-1/2 text-[13px] text-[var(--color-text-tertiary)]"
            >
              Rp
            </span>
            <input
              v-model="newTx.amount"
              type="number"
              placeholder="50000"
              class="w-full pt-[10px] pb-[10px] pr-[12px] pl-[36px] text-[14px] font-[inherit] border border-[var(--color-border)] rounded-[8px] bg-[var(--color-bg)] text-[var(--color-text)] transition-[border-color] duration-[0.15s] ease"
              @focus="$event.target.style.borderColor = 'var(--color-brand-500)'"
              @blur="$event.target.style.borderColor = 'var(--color-border)'"
            />
          </div>
        </div>

        <div class="flex flex-col gap-[6px]">
          <label class="text-[12px] font-medium text-[var(--color-text-secondary)]">
            Tipe
          </label>
          <div class="flex gap-[8px] h-[42px]">
            <button
              @click="newTx.type = 'income'"
              class="flex-1 rounded-[8px] text-[13px] font-medium font-[inherit] cursor-pointer transition-all duration-[0.15s] ease border"
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
              class="flex-1 rounded-[8px] text-[13px] font-medium font-[inherit] cursor-pointer transition-all duration-[0.15s] ease border"
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

      <div class="create-form-grid-details grid grid-cols-[2fr_1fr] gap-[12px] mb-[20px]">
        <div class="flex flex-col gap-[6px]">
          <label class="text-[12px] font-medium text-[var(--color-text-secondary)]">
            Deskripsi
          </label>
          <input
            v-model="newTx.description"
            type="text"
            placeholder="Penjualan nasi uduk"
            class="w-full py-[10px] px-[12px] text-[14px] font-[inherit] border border-[var(--color-border)] rounded-[8px] bg-[var(--color-bg)] text-[var(--color-text)] transition-[border-color] duration-[0.15s] ease"
            @focus="$event.target.style.borderColor = 'var(--color-brand-500)'"
            @blur="$event.target.style.borderColor = 'var(--color-border)'"
          />
        </div>

        <div class="flex flex-col gap-[6px]">
          <label class="text-[12px] font-medium text-[var(--color-text-secondary)]">
            Tanggal
          </label>
          <input
            v-model="newTx.transaction_date"
            type="date"
            class="w-full py-[10px] px-[12px] text-[14px] font-[inherit] border border-[var(--color-border)] rounded-[8px] bg-[var(--color-bg)] text-[var(--color-text)] transition-[border-color] duration-[0.15s] ease"
            @focus="$event.target.style.borderColor = 'var(--color-brand-500)'"
            @blur="$event.target.style.borderColor = 'var(--color-border)'"
          />
        </div>
      </div>

      <div class="flex gap-[8px]">
        <button
          @click="handleCreate"
          class="py-[10px] px-[24px] text-[14px] font-semibold font-[inherit] border-0 rounded-[8px] cursor-pointer transition-all duration-[0.15s] ease bg-[linear-gradient(135deg,#10b981,#059669)] text-white"
          @mouseenter="$event.target.style.boxShadow = '0 4px 12px rgba(16, 185, 129, 0.4)'"
          @mouseleave="$event.target.style.boxShadow = 'none'"
        >
          <i class="pi pi-check mr-[6px]" />
          Simpan
        </button>
        <button
          @click="showCreate = false"
          class="py-[10px] px-[20px] text-[14px] font-medium font-[inherit] border border-[var(--color-border)] rounded-[8px] cursor-pointer transition-all duration-[0.15s] ease bg-transparent text-[var(--color-text-secondary)]"
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
  .create-form-grid-type {
    grid-template-columns: 1fr !important;
  }
  .create-form-grid-details {
    grid-template-columns: 1fr !important;
  }
}
</style>
