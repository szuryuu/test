<script setup>
import { ref, inject, computed } from "vue";
import { useRouter, useRoute } from "vue-router";
import { useAuthStore } from "@/stores/auth";

const router = useRouter();
const route = useRoute();
const auth = useAuthStore();
const sidebarOpen = inject("sidebarOpen");
const toggleSidebar = inject("toggleSidebar");
const isMobile = inject("isMobile");

const menuItems = ref([
  {
    label: "Dashboard",
    icon: "pi pi-chart-bar",
    route: "/dashboard",
    description: "Ringkasan keuangan",
  },
  {
    label: "Transaksi",
    icon: "pi pi-list",
    route: "/transactions",
    description: "Catat & kelola",
  },
  {
    label: "Skor KUR",
    icon: "pi pi-star",
    route: "/kur-score",
    description: "Kesiapan kredit",
  },
  {
    label: "Laporan",
    icon: "pi pi-file",
    route: "/reports",
    description: "Generate bulanan",
  },
]);

function isActive(itemRoute) {
  return route.path === itemRoute;
}

function navigate(itemRoute) {
  if (isMobile.value) toggleSidebar();
  router.push(itemRoute);
}

const initials = computed(() => {
  const name = auth.umkm?.name || "U";
  return name
    .split(" ")
    .map((w) => w[0])
    .join("")
    .slice(0, 2)
    .toUpperCase();
});
</script>

<template>
  <aside
    v-show="sidebarOpen"
    class="sidebar"
    style="
      width: var(--sidebar-w);
      background: #0f172a;
      display: flex;
      flex-direction: column;
      height: 100vh;
      position: sticky;
      top: 0;
      flex-shrink: 0;
      z-index: 50;
      overflow: hidden;
    "
  >
    <!-- Brand -->
    <div
      style="
        padding: 20px 20px 16px;
        border-bottom: 1px solid rgba(255, 255, 255, 0.06);
      "
    >
      <div style="display: flex; align-items: center; gap: 10px">
        <div
          style="
            width: 36px; height: 36px; border-radius: 10px;
            background: linear-gradient(135deg, #10b981, #059669);
            display: flex; align-items: center; justify-content: center;
            font-size: 18px; color: white; font-weight: 700;
            flex-shrink: 0;
          "
        >
          K
        </div>
        <div>
          <h1
            style="
              margin: 0; font-size: 18px; font-weight: 700;
              color: white; letter-spacing: -0.02em; line-height: 1.2;
            "
          >
            Kasir<span style="color: #34d399">AI</span>
          </h1>
          <p
            style="
              margin: 0; font-size: 11px; color: rgba(255,255,255,0.4);
              letter-spacing: 0.02em;
            "
          >
            Keuangan UMKM
          </p>
        </div>
      </div>
    </div>

    <!-- Navigation -->
    <nav
      style="
        flex: 1; padding: 12px 12px; overflow-y: auto;
      "
    >
      <p
        style="
          margin: 0 0 4px; padding: 8px 12px 4px;
          font-size: 11px; font-weight: 600; color: rgba(255,255,255,0.3);
          text-transform: uppercase; letter-spacing: 0.08em;
        "
      >
        Menu Utama
      </p>

      <ul
        style="
          list-style: none; padding: 0; margin: 0;
          display: flex; flex-direction: column; gap: 2px;
        "
      >
        <li v-for="item in menuItems" :key="item.route">
          <button
            @click="navigate(item.route)"
            :title="item.description"
            style="
              width: 100%; display: flex; align-items: center; gap: 12px;
              padding: 10px 12px; border: none; cursor: pointer;
              text-align: left; font-size: 14px; font-family: inherit;
              border-radius: 8px; transition: all 0.15s ease;
            "
            :style="{
              background: isActive(item.route)
                ? 'linear-gradient(135deg, rgba(16, 185, 129, 0.15), rgba(5, 150, 105, 0.1))'
                : 'transparent',
              color: isActive(item.route) ? '#ffffff' : 'rgba(255,255,255,0.55)',
              fontWeight: isActive(item.route) ? 600 : 400,
            }"
            @mouseenter="$event.target.style.background = isActive(item.route)
              ? 'linear-gradient(135deg, rgba(16, 185, 129, 0.2), rgba(5, 150, 105, 0.15))'
              : 'rgba(255,255,255,0.06)'"
            @mouseleave="$event.target.style.background = isActive(item.route)
              ? 'linear-gradient(135deg, rgba(16, 185, 129, 0.15), rgba(5, 150, 105, 0.1))'
              : 'transparent'"
          >
            <span
              style="
                width: 32px; height: 32px; border-radius: 8px;
                display: flex; align-items: center; justify-content: center;
                flex-shrink: 0; font-size: 15px;
                transition: all 0.15s ease;
              "
              :style="{
                background: isActive(item.route)
                  ? 'rgba(16, 185, 129, 0.2)'
                  : 'rgba(255,255,255,0.05)',
                color: isActive(item.route) ? '#34d399' : 'rgba(255,255,255,0.5)',
              }"
            >
              <i :class="item.icon" />
            </span>
            <div>
              <span>{{ item.label }}</span>
              <p
                v-if="!isMobile"
                style="
                  margin: 0; font-size: 11px;
                  color: rgba(255,255,255,0.3);
                "
              >
                {{ item.description }}
              </p>
            </div>
          </button>
        </li>
      </ul>
    </nav>

    <!-- User footer -->
    <div
      style="
        padding: 12px; border-top: 1px solid rgba(255, 255, 255, 0.06);
      "
    >
      <div
        style="
          display: flex; align-items: center; gap: 10px;
          padding: 10px 12px; border-radius: 8px;
          background: rgba(255, 255, 255, 0.04);
        "
      >
        <div
          style="
            width: 34px; height: 34px; border-radius: 8px;
            background: linear-gradient(135deg, #6366f1, #8b5cf6);
            display: flex; align-items: center; justify-content: center;
            font-size: 12px; color: white; font-weight: 700;
            flex-shrink: 0;
          "
        >
          {{ initials }}
        </div>
        <div style="flex: 1; min-width: 0">
          <p
            style="
              margin: 0; font-size: 13px; color: white; font-weight: 500;
              overflow: hidden; text-overflow: ellipsis; white-space: nowrap;
            "
          >
            {{ auth.umkm?.business_name || "UMKM" }}
          </p>
          <p
            style="
              margin: 0; font-size: 11px; color: rgba(255,255,255,0.35);
              overflow: hidden; text-overflow: ellipsis; white-space: nowrap;
            "
          >
            {{ auth.umkm?.phone_number || "" }}
          </p>
        </div>
      </div>
    </div>
  </aside>
</template>

<style scoped>
.sidebar {
  transition: transform 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

@media (max-width: 1023px) {
  .sidebar {
    position: fixed !important;
    left: 0;
    top: 0;
    z-index: 50;
    box-shadow: 4px 0 24px rgba(0, 0, 0, 0.2);
  }
}
</style>
