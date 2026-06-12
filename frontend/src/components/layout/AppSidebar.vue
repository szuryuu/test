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
    class="sidebar w-[var(--sidebar-w)] bg-[#0f172a] flex flex-col h-screen sticky top-0 shrink-0 z-50 overflow-hidden"
  >
    <!-- Brand -->
    <div
      class="pt-[20px] px-[20px] pb-[16px] border-b border-[rgba(255,255,255,0.06)]"
    >
      <div class="flex items-center gap-[10px]">
        <div
          class="w-[36px] h-[36px] rounded-[10px] bg-[linear-gradient(135deg,#10b981,#059669)] flex items-center justify-center text-[18px] text-white font-bold shrink-0"
        >
          K
        </div>
        <div>
          <h1
            class="m-0 text-[18px] font-bold text-white tracking-[-0.02em] leading-[1.2]"
          >
            Kasir<span class="text-[#34d399]">AI</span>
          </h1>
          <p
            class="m-0 text-[11px] text-[rgba(255,255,255,0.4)] tracking-[0.02em]"
          >
            Keuangan UMKM
          </p>
        </div>
      </div>
    </div>

    <!-- Navigation -->
    <nav class="flex-1 p-[12px] overflow-y-auto">
      <p
        class="m-0 mb-[4px] pt-[8px] pr-[12px] pb-[4px] pl-[12px] text-[11px] font-semibold text-[rgba(255,255,255,0.3)] uppercase tracking-[0.08em]"
      >
        Menu Utama
      </p>

      <ul class="list-none p-0 m-0 flex flex-col gap-[2px]">
        <li v-for="item in menuItems" :key="item.route">
          <button
            @click="navigate(item.route)"
            :title="item.description"
            class="w-full flex items-center gap-[12px] py-[10px] px-[12px] border-0 cursor-pointer text-left text-[14px] font-[inherit] rounded-[8px] transition-all duration-[0.15s] ease"
            :style="{
              background: isActive(item.route)
                ? 'linear-gradient(135deg, rgba(16, 185, 129, 0.15), rgba(5, 150, 105, 0.1))'
                : 'transparent',
              color: isActive(item.route)
                ? '#ffffff'
                : 'rgba(255,255,255,0.55)',
              fontWeight: isActive(item.route) ? 600 : 400,
            }"
            @mouseenter="
              $event.target.style.background = isActive(item.route)
                ? 'linear-gradient(135deg, rgba(16, 185, 129, 0.2), rgba(5, 150, 105, 0.15))'
                : 'rgba(255,255,255,0.06)'
            "
            @mouseleave="
              $event.target.style.background = isActive(item.route)
                ? 'linear-gradient(135deg, rgba(16, 185, 129, 0.15), rgba(5, 150, 105, 0.1))'
                : 'transparent'
            "
          >
            <span
              class="w-[32px] h-[32px] rounded-[8px] flex items-center justify-center shrink-0 text-[15px] transition-all duration-[0.15s] ease"
              :style="{
                background: isActive(item.route)
                  ? 'rgba(16, 185, 129, 0.2)'
                  : 'rgba(255,255,255,0.05)',
                color: isActive(item.route)
                  ? '#34d399'
                  : 'rgba(255,255,255,0.5)',
              }"
            >
              <i :class="item.icon" />
            </span>
            <div>
              <span>{{ item.label }}</span>
              <p
                v-if="!isMobile"
                class="m-0 text-[11px] text-[rgba(255,255,255,0.3)]"
              >
                {{ item.description }}
              </p>
            </div>
          </button>
        </li>
      </ul>
    </nav>

    <!-- WhatsApp Chat -->
    <div class="px-[12px] pb-[12px]">
      <a
        href="https://wa.me/6285924572925?text=bantuan"
        target="_blank"
        rel="noopener noreferrer"
        class="w-full flex items-center justify-center gap-[8px] py-[10px] px-[12px] rounded-[8px] text-[13px] font-semibold text-white no-underline cursor-pointer transition-all duration-[0.2s] ease shadow-[0_4px_12px_rgba(37,211,102,0.2)] hover:shadow-[0_6px_16px_rgba(37,211,102,0.3)] hover:-translate-y-[1px]"
        style="background: linear-gradient(135deg, #25d366, #128c7e)"
      >
        <i class="pi pi-whatsapp text-[16px]" />
        Hubungkan WhatsApp
      </a>
      <p
        class="mt-[6px] mb-0 text-center text-[10px] text-[rgba(255,255,255,0.4)]"
      >
        Catat transaksi via chat
      </p>
    </div>

    <!-- User footer -->
    <div class="p-[12px] border-t border-[rgba(255,255,255,0.06)]">
      <div
        class="flex items-center gap-[10px] py-[10px] px-[12px] rounded-[8px] bg-[rgba(255,255,255,0.04)]"
      >
        <div
          class="w-[34px] h-[34px] rounded-[8px] bg-[linear-gradient(135deg,#6366f1,#8b5cf6)] flex items-center justify-center text-[12px] text-white font-bold shrink-0"
        >
          {{ initials }}
        </div>
        <div class="flex-1 min-w-0">
          <p class="m-0 text-[13px] text-white font-medium truncate">
            {{ auth.umkm?.business_name || "UMKM" }}
          </p>
          <p class="m-0 text-[11px] text-[rgba(255,255,255,0.35)] truncate">
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
