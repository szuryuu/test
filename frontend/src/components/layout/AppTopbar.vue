<script setup>
import { computed, inject } from "vue";
import { useRouter, useRoute } from "vue-router";
import { useAuthStore } from "@/stores/auth";

const router = useRouter();
const route = useRoute();
const auth = useAuthStore();
const toggleSidebar = inject("toggleSidebar");
const isMobile = inject("isMobile");

const pageTitle = computed(() => {
  const titles = {
    "/dashboard": "Dashboard",
    "/transactions": "Transaksi",
    "/kur-score": "Skor KUR",
    "/reports": "Laporan",
  };
  return titles[route.path] || "KasirAI";
});

const pageSubtitle = computed(() => {
  const subtitles = {
    "/dashboard": "Ringkasan keuangan usaha Anda",
    "/transactions": "Catat dan kelola semua transaksi",
    "/kur-score": "Kesiapan kredit UMKM Anda",
    "/reports": "Generate laporan keuangan bulanan",
  };
  return subtitles[route.path] || "";
});

function handleLogout() {
  auth.logout();
  router.push("/login");
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
  <header
    style="
      background: var(--color-surface);
      border-bottom: 1px solid var(--color-border);
      height: 64px;
      display: flex;
      align-items: center;
      justify-content: space-between;
      padding: 0 24px;
      position: sticky;
      top: 0;
      z-index: 30;
    "
  >
    <!-- Left: hamburger + title -->
    <div style="display: flex; align-items: center; gap: 16px; min-width: 0">
      <!-- Hamburger (always visible on mobile) -->
      <button
        @click="toggleSidebar"
        :style="{
          background: 'transparent',
          border: 'none',
          cursor: 'pointer',
          width: 36,
          height: 36,
          borderRadius: 8,
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          color: 'var(--color-text-secondary)',
          fontSize: 18,
          transition: 'all 0.15s ease',
        }"
        @mouseenter="$event.target.style.background = 'var(--color-bg)'"
        @mouseleave="$event.target.style.background = 'transparent'"
        aria-label="Toggle sidebar"
      >
        <i class="pi pi-bars" />
      </button>

      <div style="min-width: 0">
        <h1
          style="
            margin: 0; font-size: 20px; font-weight: 700;
            color: var(--color-text); letter-spacing: -0.02em;
            overflow: hidden; text-overflow: ellipsis; white-space: nowrap;
          "
        >
          {{ pageTitle }}
        </h1>
        <p
          v-if="!isMobile"
          style="
            margin: 0; font-size: 13px; color: var(--color-text-secondary);
          "
        >
          {{ pageSubtitle }}
        </p>
      </div>
    </div>

    <!-- Right: user info + logout -->
    <div style="display: flex; align-items: center; gap: 12px">
      <span
        v-if="!isMobile"
        style="font-size: 13px; color: var(--color-text-secondary)"
      >
        {{ auth.umkm?.name || "" }}
      </span>

      <div
        style="
          width: 36px; height: 36px; border-radius: 8px;
          background: linear-gradient(135deg, #6366f1, #8b5cf6);
          display: flex; align-items: center; justify-content: center;
          font-size: 12px; color: white; font-weight: 700;
        "
        :title="auth.umkm?.name || 'User'"
      >
        {{ initials }}
      </div>

      <button
        @click="handleLogout"
        style="
          background: transparent;
          border: 1px solid var(--color-border);
          cursor: pointer;
          padding: 7px 14px;
          border-radius: 8px;
          font-size: 13px;
          font-weight: 500;
          font-family: inherit;
          color: var(--color-text-secondary);
          transition: all 0.15s ease;
          display: flex;
          align-items: center;
          gap: 6px;
        "
        @mouseenter="$event.target.style.background = 'var(--color-bg)'"
        @mouseleave="$event.target.style.background = 'transparent'"
      >
        <i class="pi pi-sign-out" style="font-size: 13px" />
        <span v-if="!isMobile">Keluar</span>
      </button>
    </div>
  </header>
</template>
