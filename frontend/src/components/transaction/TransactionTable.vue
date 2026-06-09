<script setup>
import { ref, watch } from "vue";
import { formatRupiah, formatDate } from "@/utils/format";
import { useTransactionStore } from "@/stores/transaction";

const store = useTransactionStore();

const emit = defineEmits(["delete"]);

const filters = ref({
  type: "",
  start_date: "",
  end_date: "",
});

const page = ref(1);
const limit = 20;

watch(
  filters,
  () => {
    page.value = 1;
    loadData();
  },
  { deep: true }
);

function loadData() {
  store.fetchTransactions({
    page: page.value,
    limit,
    ...filters.value,
  });
}

function onPageChange(p) {
  page.value = p;
  loadData();
}

function handleDelete(id) {
  emit("delete", id);
}

loadData();

const typeOptions = [
  { label: "Semua", value: "" },
  { label: "Pemasukan", value: "income" },
  { label: "Pengeluaran", value: "expense" },
];

const totalPages = computed(() => Math.ceil(store.total / limit));

import { computed } from "vue";
</script>

<template>
  <div>
    <!-- Filters -->
    <div
      class="flex flex-wrap gap-3 mb-4 p-3"
      style="
        background: var(--color-pure-white);
        border-radius: var(--radius-md);
        box-shadow: 0 1px 3px rgba(0,0,0,0.04);
      "
    >
      <select
        v-model="filters.type"
        class="px-3 py-2 text-sm"
        style="
          border: 1px solid var(--color-border-soft);
          border-radius: var(--radius-sm);
          background: var(--color-pure-white);
          color: var(--color-near-black);
        "
      >
        <option v-for="opt in typeOptions" :key="opt.value" :value="opt.value">
          {{ opt.label }}
        </option>
      </select>

      <input
        v-model="filters.start_date"
        type="date"
        class="px-3 py-2 text-sm"
        style="
          border: 1px solid var(--color-border-soft);
          border-radius: var(--radius-sm);
          background: var(--color-pure-white);
          color: var(--color-near-black);
        "
        placeholder="Dari"
      />

      <input
        v-model="filters.end_date"
        type="date"
        class="px-3 py-2 text-sm"
        style="
          border: 1px solid var(--color-border-soft);
          border-radius: var(--radius-sm);
          background: var(--color-pure-white);
          color: var(--color-near-black);
        "
        placeholder="Sampai"
      />
    </div>

    <!-- Table -->
    <div
      style="
        background: var(--color-pure-white);
        border-radius: var(--radius-md);
        box-shadow: 0 1px 3px rgba(0,0,0,0.04);
        overflow: hidden;
      "
    >
      <table style="width: 100%; border-collapse: collapse">
        <thead>
          <tr style="border-bottom: 1px solid var(--color-border-soft)">
            <th class="text-left px-4 py-3 text-sm" style="font-weight: 600; color: var(--color-text-secondary)">Tanggal</th>
            <th class="text-left px-4 py-3 text-sm" style="font-weight: 600; color: var(--color-text-secondary)">Deskripsi</th>
            <th class="text-left px-4 py-3 text-sm" style="font-weight: 600; color: var(--color-text-secondary)">Tipe</th>
            <th class="text-right px-4 py-3 text-sm" style="font-weight: 600; color: var(--color-text-secondary)">Jumlah</th>
            <th class="text-center px-4 py-3 text-sm" style="font-weight: 600; color: var(--color-text-secondary)">Aksi</th>
          </tr>
        </thead>
        <tbody>
          <tr
            v-for="tx in store.transactions"
            :key="tx.id"
            style="border-bottom: 1px solid var(--color-border-soft)"
          >
            <td class="px-4 py-3 text-sm" style="color: var(--color-near-black)">
              {{ formatDate(tx.transaction_date) }}
            </td>
            <td class="px-4 py-3 text-sm" style="color: var(--color-near-black); max-width: 200px">
              <div class="overflow-hidden white-space-nowrap text-overflow-ellipsis">
                {{ tx.description }}
              </div>
            </td>
            <td class="px-4 py-3">
              <span
                class="text-xs px-2 py-1"
                style="
                  border-radius: var(--radius-pill);
                  font-weight: 500;
                "
                :style="{
                  background: tx.type === 'income' ? '#ecfdf5' : '#fef2f2',
                  color: tx.type === 'income' ? '#10b981' : '#ef4444',
                }"
              >
                {{ tx.type === "income" ? "Pemasukan" : "Pengeluaran" }}
              </span>
            </td>
            <td
              class="px-4 py-3 text-sm text-right"
              :style="{
                fontWeight: 600,
                color: tx.type === 'income' ? '#10b981' : '#ef4444',
              }"
            >
              {{ tx.type === "income" ? "+" : "-" }}{{ formatRupiah(tx.amount) }}
            </td>
            <td class="px-4 py-3 text-center">
              <button
                @click="handleDelete(tx.id)"
                class="cursor-pointer"
                style="
                  background: transparent;
                  border: none;
                  color: var(--color-text-secondary);
                  font-size: 14px;
                  padding: 4px 8px;
                "
                title="Hapus"
              >
                <i class="pi pi-trash" />
              </button>
            </td>
          </tr>
          <tr v-if="store.transactions.length === 0 && !store.isLoading">
            <td colspan="5" class="text-center py-6" style="color: var(--color-text-secondary)">
              <p class="text-sm m-0">Tidak ada transaksi</p>
            </td>
          </tr>
        </tbody>
      </table>

      <!-- Pagination -->
      <div
        v-if="store.total > limit"
        class="flex justify-content-between align-items-center px-4 py-3"
        style="border-top: 1px solid var(--color-border-soft)"
      >
        <span class="text-sm" style="color: var(--color-text-secondary)">
          {{ store.total }} transaksi
        </span>
        <div class="flex gap-2">
          <button
            @click="onPageChange(page - 1)"
            :disabled="page <= 1"
            class="px-3 py-1 text-sm cursor-pointer"
            style="
              background: transparent;
              border: 1px solid var(--color-border-soft);
              border-radius: var(--radius-sm);
              color: var(--color-near-black);
            "
          >
            ←
          </button>
          <span class="text-sm px-2 py-1" style="color: var(--color-near-black)">
            {{ page }} / {{ totalPages }}
          </span>
          <button
            @click="onPageChange(page + 1)"
            :disabled="page >= totalPages"
            class="px-3 py-1 text-sm cursor-pointer"
            style="
              background: transparent;
              border: 1px solid var(--color-border-soft);
              border-radius: var(--radius-sm);
              color: var(--color-near-black);
            "
          >
            →
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
