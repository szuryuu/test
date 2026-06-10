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
    class="app-topbar bg-[var(--color-surface)] border-b border-[var(--color-border)] h-[64px] flex items-center justify-between px-[24px] sticky top-0 z-30"
  >
    <!-- Left: hamburger + title -->
    <div class="flex items-center gap-[16px] min-w-0">
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

      <div class="min-w-0">
        <h1 class="m-0 text-[20px] font-bold text-[var(--color-text)] tracking-[-0.02em] truncate">
          {{ pageTitle }}
        </h1>
        <p v-if="!isMobile" class="m-0 text-[13px] text-[var(--color-text-secondary)]">
          {{ pageSubtitle }}
        </p>
      </div>
    </div>

    <!-- Right: user info + logout -->
    <div class="flex items-center gap-[12px]">
      <span
        v-if="!isMobile"
        class="text-[13px] text-[var(--color-text-secondary)]"
      >
        {{ auth.umkm?.name || "" }}
      </span>

      <div
        class="w-[36px] h-[36px] rounded-[8px] bg-[linear-gradient(135deg,#6366f1,#8b5cf6)] flex items-center justify-center text-[12px] text-white font-bold"
        :title="auth.umkm?.name || 'User'"
      >
        {{ initials }}
      </div>

      <button
        @click="handleLogout"
        class="bg-transparent border border-[var(--color-border)] cursor-pointer py-[7px] px-[14px] rounded-[8px] text-[13px] font-medium font-[inherit] text-[var(--color-text-secondary)] transition-all duration-[0.15s] ease flex items-center gap-[6px]"
        @mouseenter="$event.target.style.background = 'var(--color-bg)'"
        @mouseleave="$event.target.style.background = 'transparent'"
      >
        <i class="pi pi-sign-out text-[13px]" />
        <span v-if="!isMobile">Keluar</span>
      </button>
    </div>
  </header>
</template>
