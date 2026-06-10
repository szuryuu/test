<script setup>
import { ref, computed, watch, onMounted } from "vue";
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
    ...(filters.value.type ? { type: filters.value.type } : {}),
    ...(filters.value.start_date ? { start_date: filters.value.start_date } : {}),
    ...(filters.value.end_date ? { end_date: filters.value.end_date } : {}),
  });
}

function onPageChange(p) {
  page.value = p;
  loadData();
}

function handleDelete(id) {
  emit("delete", id);
}

onMounted(loadData);

const typeOptions = [
  { label: "Semua Tipe", value: "" },
  { label: "Pemasukan", value: "income" },
  { label: "Pengeluaran", value: "expense" },
];

const totalPages = computed(() => Math.ceil(store.total / limit));

function clearFilters() {
  filters.value = { type: "", start_date: "", end_date: "" };
}
</script>

<template>
  <div>
    <!-- Filters -->
    <div
      class="flex flex-wrap items-center gap-[8px] mb-[16px] p-[16px] bg-[var(--color-surface)] border border-[var(--color-border)] rounded-[var(--radius-card)] shadow-[var(--shadow-card)]"
    >
      <select
        v-model="filters.type"
        class="py-[8px] px-[12px] text-[13px] font-[inherit] border border-[var(--color-border)] rounded-[8px] bg-[var(--color-surface)] text-[var(--color-text)] cursor-pointer min-w-[130px]"
      >
        <option v-for="opt in typeOptions" :key="opt.value" :value="opt.value">
          {{ opt.label }}
        </option>
      </select>

      <input
        v-model="filters.start_date"
        type="date"
        class="py-[8px] px-[12px] text-[13px] font-[inherit] border border-[var(--color-border)] rounded-[8px] bg-[var(--color-surface)] text-[var(--color-text)]"
        placeholder="Dari"
      />

      <input
        v-model="filters.end_date"
        type="date"
        class="py-[8px] px-[12px] text-[13px] font-[inherit] border border-[var(--color-border)] rounded-[8px] bg-[var(--color-surface)] text-[var(--color-text)]"
        placeholder="Sampai"
      />

      <button
        @click="clearFilters"
        class="py-[8px] px-[14px] text-[13px] font-[inherit] border border-[var(--color-border)] rounded-[8px] bg-transparent text-[var(--color-text-secondary)] cursor-pointer transition-all duration-[0.15s] ease"
        @mouseenter="$event.target.style.background = 'var(--color-bg)'"
        @mouseleave="$event.target.style.background = 'transparent'"
      >
        Reset
      </button>
    </div>

    <!-- Loading state -->
    <div
      v-if="store.isLoading && !store.transactions.length"
      class="bg-[var(--color-surface)] border border-[var(--color-border)] rounded-[var(--radius-card)] shadow-[var(--shadow-card)] overflow-hidden"
    >
      <div class="p-[16px] flex flex-col gap-[12px]">
        <div v-for="i in 5" :key="i" class="skeleton h-[48px] rounded-[8px]" />
      </div>
    </div>

    <!-- Table -->
    <div
      v-else
      class="bg-[var(--color-surface)] border border-[var(--color-border)] rounded-[var(--radius-card)] shadow-[var(--shadow-card)] overflow-hidden"
    >
      <!-- Desktop table -->
      <table class="responsive-table w-full border-collapse">
        <thead>
          <tr class="border-b border-[var(--color-border)] bg-[var(--color-bg)]">
            <th
              class="text-left p-[12px_16px] text-[11px] font-semibold text-[var(--color-text-secondary)] uppercase tracking-[0.05em]"
            >
              Tanggal
            </th>
            <th
              class="text-left p-[12px_16px] text-[11px] font-semibold text-[var(--color-text-secondary)] uppercase tracking-[0.05em]"
            >
              Deskripsi
            </th>
            <th
              class="text-left p-[12px_16px] text-[11px] font-semibold text-[var(--color-text-secondary)] uppercase tracking-[0.05em]"
            >
              Tipe
            </th>
            <th
              class="text-right p-[12px_16px] text-[11px] font-semibold text-[var(--color-text-secondary)] uppercase tracking-[0.05em]"
            >
              Jumlah
            </th>
            <th
              class="text-center p-[12px_16px] text-[11px] font-semibold text-[var(--color-text-secondary)] uppercase tracking-[0.05em] w-[60px]"
            >
              Aksi
            </th>
          </tr>
        </thead>
        <tbody>
          <tr
            v-for="(tx, idx) in store.transactions"
            :key="tx.id"
            class="transition-[background] duration-[0.15s] ease"
            :style="{
              borderBottom: '1px solid var(--color-border-light)',
              background: idx % 2 === 0 ? 'transparent' : 'rgba(248, 250, 252, 0.5)',
            }"
            @mouseenter="$event.target.style.background = 'var(--color-bg)'"
            @mouseleave="$event.target.style.background = idx % 2 === 0 ? 'transparent' : 'rgba(248, 250, 252, 0.5)'"
          >
            <td
              data-label="Tanggal"
              class="p-[12px_16px] text-[13px] text-[var(--color-text)]"
            >
              {{ formatDate(tx.transaction_date) }}
            </td>
            <td
              data-label="Deskripsi"
              class="p-[12px_16px] text-[13px] text-[var(--color-text)] max-w-[240px]"
            >
              <span class="block truncate">{{ tx.description }}</span>
            </td>
            <td data-label="Tipe" class="p-[12px_16px]">
              <span
                class="inline-block text-[11px] font-medium py-[3px] px-[10px] rounded-[999px]"
                :style="{
                  background: tx.type === 'income' ? 'var(--color-income-bg)' : 'var(--color-expense-bg)',
                  color: tx.type === 'income' ? 'var(--color-income)' : 'var(--color-expense)',
                }"
              >
                {{ tx.type === "income" ? "Pemasukan" : "Pengeluaran" }}
              </span>
            </td>
            <td
              data-label="Jumlah"
              class="p-[12px_16px] text-[13px] font-semibold text-right"
              :style="{
                color: tx.type === 'income' ? 'var(--color-income)' : 'var(--color-expense)',
              }"
            >
              {{ tx.type === "income" ? "+" : "−" }}{{ formatRupiah(tx.amount) }}
            </td>
            <td data-label="Aksi" class="p-[12px_16px] text-center">
              <button
                @click="handleDelete(tx.id)"
                class="bg-transparent border-0 cursor-pointer p-[6px_8px] rounded-[6px] text-[var(--color-text-tertiary)] text-[15px] transition-all duration-[0.15s] ease"
                @mouseenter="
                  $event.target.style.background = 'var(--color-expense-bg)';
                  $event.target.style.color = 'var(--color-expense)';
                "
                @mouseleave="
                  $event.target.style.background = 'transparent';
                  $event.target.style.color = 'var(--color-text-tertiary)';
                "
                title="Hapus transaksi"
              >
                <i class="pi pi-trash" />
              </button>
            </td>
          </tr>

          <!-- Empty state -->
          <tr v-if="store.transactions.length === 0 && !store.isLoading">
            <td colspan="5" class="p-[48px_16px] text-center">
              <div
                class="w-[56px] h-[56px] rounded-[12px] bg-[var(--color-bg)] mx-auto mb-[12px] flex items-center justify-center text-[24px] text-[var(--color-text-tertiary)]"
              >
                <i class="pi pi-inbox" />
              </div>
              <p class="text-[14px] text-[var(--color-text-secondary)] m-0">
                Tidak ada transaksi
              </p>
              <p class="text-[12px] text-[var(--color-text-tertiary)] mt-[4px]">
                {{ filters.type || filters.start_date || filters.end_date ? 'Coba ubah filter pencarian' : 'Tambahkan transaksi baru untuk memulai' }}
              </p>
            </td>
          </tr>
        </tbody>
      </table>

      <!-- Pagination -->
      <div
        v-if="store.total > limit"
        class="flex justify-between items-center p-[12px_16px] border-t border-[var(--color-border)] bg-[var(--color-bg)]"
      >
        <span class="text-[12px] text-[var(--color-text-secondary)]">
          {{ store.total }} transaksi
        </span>
        <div class="flex items-center gap-[4px]">
          <button
            @click="onPageChange(page - 1)"
            :disabled="page <= 1"
            class="py-[6px] px-[10px] text-[13px] font-[inherit] border border-[var(--color-border)] rounded-[6px] bg-[var(--color-surface)] text-[var(--color-text)] cursor-pointer transition-all duration-[0.15s] ease"
            :style="{ opacity: page <= 1 ? 0.4 : 1, cursor: page <= 1 ? 'not-allowed' : 'pointer' }"
          >
            <i class="pi pi-chevron-left" />
          </button>
          <span class="text-[12px] text-[var(--color-text-secondary)] px-[8px]">
            {{ page }} / {{ totalPages }}
          </span>
          <button
            @click="onPageChange(page + 1)"
            :disabled="page >= totalPages"
            class="py-[6px] px-[10px] text-[13px] font-[inherit] border border-[var(--color-border)] rounded-[6px] bg-[var(--color-surface)] text-[var(--color-text)] cursor-pointer transition-all duration-[0.15s] ease"
            :style="{ opacity: page >= totalPages ? 0.4 : 1, cursor: page >= totalPages ? 'not-allowed' : 'pointer' }"
          >
            <i class="pi pi-chevron-right" />
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
