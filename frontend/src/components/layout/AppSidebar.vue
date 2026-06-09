<script setup>
import { ref } from "vue";
import { useRouter, useRoute } from "vue-router";
import { useAuthStore } from "@/stores/auth";

const router = useRouter();
const route = useRoute();
const auth = useAuthStore();

const menuItems = ref([
  {
    label: "Dashboard",
    icon: "pi pi-chart-bar",
    route: "/dashboard",
  },
  {
    label: "Transaksi",
    icon: "pi pi-list",
    route: "/transactions",
  },
  {
    label: "Skor KUR",
    icon: "pi pi-star",
    route: "/kur-score",
  },
  {
    label: "Laporan",
    icon: "pi pi-file",
    route: "/reports",
  },
]);

function isActive(itemRoute) {
  return route.path === itemRoute;
}

function navigate(itemRoute) {
  router.push(itemRoute);
}
</script>

<template>
  <aside
    class="flex flex-column h-screen"
    style="background: var(--color-near-black); width: 256px; flex-shrink: 0"
  >
    <!-- Brand -->
    <div class="px-5 pt-6 pb-4">
      <h1
        class="text-2xl m-0"
        style="color: var(--color-pure-white); font-weight: 600; letter-spacing: -0.02em"
      >
        Kasir<span style="color: var(--color-apple-blue)">AI</span>
      </h1>
      <p class="text-sm mt-1 mb-0" style="color: var(--color-text-secondary)">
        Keuangan UMKM
      </p>
    </div>

    <!-- Navigation -->
    <nav class="flex-1 px-3 pt-2">
      <ul class="list-none p-0 m-0 flex flex-column gap-1">
        <li v-for="item in menuItems" :key="item.route">
          <button
            @click="navigate(item.route)"
            class="w-full flex align-items-center gap-3 px-4 py-3 border-none cursor-pointer text-left transition-all transition-duration-200"
            :style="{
              background: isActive(item.route) ? 'rgba(255,255,255,0.1)' : 'transparent',
              color: isActive(item.route) ? 'var(--color-pure-white)' : 'var(--color-text-secondary)',
              borderRadius: 'var(--radius-sm)',
              fontSize: '14px',
              fontWeight: isActive(item.route) ? 600 : 400,
            }"
          >
            <i :class="item.icon" style="font-size: 18px" />
            <span>{{ item.label }}</span>
          </button>
        </li>
      </ul>
    </nav>

    <!-- Footer -->
    <div class="px-5 pb-5">
      <div
        class="px-4 py-3 border-round-md"
        style="background: rgba(255,255,255,0.05)"
      >
        <p class="text-xs m-0 mb-1" style="color: var(--color-text-secondary)">
          {{ auth.umkm?.business_name || "UMKM" }}
        </p>
        <p class="text-xs m-0" style="color: var(--color-border-mid)">
          {{ auth.umkm?.phone_number || "" }}
        </p>
      </div>
    </div>
  </aside>
</template>
