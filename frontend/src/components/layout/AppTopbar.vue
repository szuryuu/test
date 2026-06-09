<script setup>
import { useRouter } from "vue-router";
import { useAuthStore } from "@/stores/auth";

const router = useRouter();
const auth = useAuthStore();

function logout() {
  auth.logout();
  router.push("/login");
}
</script>

<template>
  <header
    class="flex justify-content-between align-items-center px-5 py-3"
    style="
      background: var(--color-pure-white);
      border-bottom: 1px solid var(--color-border-soft);
      height: 64px;
    "
  >
    <div>
      <h2 class="text-lg m-0" style="font-weight: 600; color: var(--color-near-black)">
        {{ auth.umkm?.business_name || "KasirAI" }}
      </h2>
    </div>

    <div class="flex align-items-center gap-3">
      <span class="text-sm" style="color: var(--color-text-secondary)">
        {{ auth.umkm?.name || "" }}
      </span>
      <button
        @click="logout"
        class="flex align-items-center gap-2 px-4 py-2 cursor-pointer transition-all transition-duration-200"
        style="
          background: transparent;
          color: var(--color-text-secondary);
          border: 1px solid var(--color-border-soft);
          border-radius: var(--radius-sm);
          font-size: 13px;
          font-weight: 500;
        "
      >
        <i class="pi pi-sign-out" style="font-size: 14px" />
        Keluar
      </button>
    </div>
  </header>
</template>
