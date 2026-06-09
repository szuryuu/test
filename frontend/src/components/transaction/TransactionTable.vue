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
      style="
        display: flex; flex-wrap: wrap; align-items: center; gap: 8px;
        margin-bottom: 16px; padding: 16px;
        background: var(--color-surface);
        border: 1px solid var(--color-border);
        border-radius: var(--radius-card);
        box-shadow: var(--shadow-card);
      "
    >
      <select
        v-model="filters.type"
        style="
          padding: 8px 12px; font-size: 13px; font-family: inherit;
          border: 1px solid var(--color-border); border-radius: 8px;
          background: var(--color-surface); color: var(--color-text);
          cursor: pointer; min-width: 130px;
        "
      >
        <option v-for="opt in typeOptions" :key="opt.value" :value="opt.value">
          {{ opt.label }}
        </option>
      </select>

      <input
        v-model="filters.start_date"
        type="date"
        style="
          padding: 8px 12px; font-size: 13px; font-family: inherit;
          border: 1px solid var(--color-border); border-radius: 8px;
          background: var(--color-surface); color: var(--color-text);
        "
        placeholder="Dari"
      />

      <input
        v-model="filters.end_date"
        type="date"
        style="
          padding: 8px 12px; font-size: 13px; font-family: inherit;
          border: 1px solid var(--color-border); border-radius: 8px;
          background: var(--color-surface); color: var(--color-text);
        "
        placeholder="Sampai"
      />

      <button
        @click="clearFilters"
        style="
          padding: 8px 14px; font-size: 13px; font-family: inherit;
          border: 1px solid var(--color-border); border-radius: 8px;
          background: transparent; color: var(--color-text-secondary);
          cursor: pointer; transition: all 0.15s ease;
        "
        @mouseenter="$event.target.style.background = 'var(--color-bg)'"
        @mouseleave="$event.target.style.background = 'transparent'"
      >
        Reset
      </button>
    </div>

    <!-- Loading state -->
    <div
      v-if="store.isLoading && !store.transactions.length"
      style="
        background: var(--color-surface);
        border: 1px solid var(--color-border);
        border-radius: var(--radius-card);
        box-shadow: var(--shadow-card);
        overflow: hidden;
      "
    >
      <div style="padding: 16px; display: flex; flex-direction: column; gap: 12px">
        <div v-for="i in 5" :key="i" class="skeleton" style="height: 48px; border-radius: 8px" />
      </div>
    </div>

    <!-- Table -->
    <div
      v-else
      style="
        background: var(--color-surface);
        border: 1px solid var(--color-border);
        border-radius: var(--radius-card);
        box-shadow: var(--shadow-card);
        overflow: hidden;
      "
    >
      <!-- Desktop table -->
      <table
        class="responsive-table"
        style="width: 100%; border-collapse: collapse"
      >
        <thead>
          <tr style="border-bottom: 1px solid var(--color-border); background: var(--color-bg)">
            <th
              style="
                text-align: left; padding: 12px 16px; font-size: 11px;
                font-weight: 600; color: var(--color-text-secondary);
                text-transform: uppercase; letter-spacing: 0.05em;
              "
            >
              Tanggal
            </th>
            <th
              style="
                text-align: left; padding: 12px 16px; font-size: 11px;
                font-weight: 600; color: var(--color-text-secondary);
                text-transform: uppercase; letter-spacing: 0.05em;
              "
            >
              Deskripsi
            </th>
            <th
              style="
                text-align: left; padding: 12px 16px; font-size: 11px;
                font-weight: 600; color: var(--color-text-secondary);
                text-transform: uppercase; letter-spacing: 0.05em;
              "
            >
              Tipe
            </th>
            <th
              style="
                text-align: right; padding: 12px 16px; font-size: 11px;
                font-weight: 600; color: var(--color-text-secondary);
                text-transform: uppercase; letter-spacing: 0.05em;
              "
            >
              Jumlah
            </th>
            <th
              style="
                text-align: center; padding: 12px 16px; font-size: 11px;
                font-weight: 600; color: var(--color-text-secondary);
                text-transform: uppercase; letter-spacing: 0.05em;
                width: 60px;
              "
            >
              Aksi
            </th>
          </tr>
        </thead>
        <tbody>
          <tr
            v-for="(tx, idx) in store.transactions"
            :key="tx.id"
            style="transition: background 0.15s ease;"
            :style="{
              borderBottom: '1px solid var(--color-border-light)',
              background: idx % 2 === 0 ? 'transparent' : 'rgba(248, 250, 252, 0.5)',
            }"
            @mouseenter="$event.target.style.background = 'var(--color-bg)'"
            @mouseleave="$event.target.style.background = idx % 2 === 0 ? 'transparent' : 'rgba(248, 250, 252, 0.5)'"
          >
            <td
              data-label="Tanggal"
              style="padding: 12px 16px; font-size: 13px; color: var(--color-text);"
            >
              {{ formatDate(tx.transaction_date) }}
            </td>
            <td
              data-label="Deskripsi"
              style="padding: 12px 16px; font-size: 13px; color: var(--color-text); max-width: 240px;"
            >
              <span class="truncate" style="display: block">{{ tx.description }}</span>
            </td>
            <td data-label="Tipe" style="padding: 12px 16px">
              <span
                style="
                  display: inline-block; font-size: 11px; font-weight: 500;
                  padding: 3px 10px; border-radius: 999px;
                "
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
              style="padding: 12px 16px; font-size: 13px; font-weight: 600; text-align: right;"
              :style="{
                color: tx.type === 'income' ? 'var(--color-income)' : 'var(--color-expense)',
              }"
            >
              {{ tx.type === "income" ? "+" : "−" }}{{ formatRupiah(tx.amount) }}
            </td>
            <td data-label="Aksi" style="padding: 12px 16px; text-align: center">
              <button
                @click="handleDelete(tx.id)"
                style="
                  background: transparent; border: none; cursor: pointer;
                  padding: 6px 8px; border-radius: 6px;
                  color: var(--color-text-tertiary); font-size: 15px;
                  transition: all 0.15s ease;
                "
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
            <td colspan="5" style="padding: 48px 16px; text-align: center">
              <div
                style="
                  width: 56px; height: 56px; border-radius: 12px;
                  background: var(--color-bg); margin: 0 auto 12px;
                  display: flex; align-items: center; justify-content: center;
                  font-size: 24px; color: var(--color-text-tertiary);
                "
              >
                <i class="pi pi-inbox" />
              </div>
              <p style="font-size: 14px; color: var(--color-text-secondary); margin: 0">
                Tidak ada transaksi
              </p>
              <p style="font-size: 12px; color: var(--color-text-tertiary); margin: 4px 0 0">
                {{ filters.type || filters.start_date || filters.end_date ? 'Coba ubah filter pencarian' : 'Tambahkan transaksi baru untuk memulai' }}
              </p>
            </td>
          </tr>
        </tbody>
      </table>

      <!-- Pagination -->
      <div
        v-if="store.total > limit"
        style="
          display: flex; justify-content: space-between; align-items: center;
          padding: 12px 16px; border-top: 1px solid var(--color-border);
          background: var(--color-bg);
        "
      >
        <span style="font-size: 12px; color: var(--color-text-secondary)">
          {{ store.total }} transaksi
        </span>
        <div style="display: flex; align-items: center; gap: 4px">
          <button
            @click="onPageChange(page - 1)"
            :disabled="page <= 1"
            style="
              padding: 6px 10px; font-size: 13px; font-family: inherit;
              border: 1px solid var(--color-border); border-radius: 6px;
              background: var(--color-surface); color: var(--color-text);
              cursor: pointer; transition: all 0.15s ease;
            "
            :style="{ opacity: page <= 1 ? 0.4 : 1, cursor: page <= 1 ? 'not-allowed' : 'pointer' }"
          >
            <i class="pi pi-chevron-left" />
          </button>
          <span style="font-size: 12px; color: var(--color-text-secondary); padding: 0 8px">
            {{ page }} / {{ totalPages }}
          </span>
          <button
            @click="onPageChange(page + 1)"
            :disabled="page >= totalPages"
            style="
              padding: 6px 10px; font-size: 13px; font-family: inherit;
              border: 1px solid var(--color-border); border-radius: 6px;
              background: var(--color-surface); color: var(--color-text);
              cursor: pointer; transition: all 0.15s ease;
            "
            :style="{ opacity: page >= totalPages ? 0.4 : 1, cursor: page >= totalPages ? 'not-allowed' : 'pointer' }"
          >
            <i class="pi pi-chevron-right" />
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
