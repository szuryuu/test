<script setup>
import { ref, provide, onMounted, onUnmounted } from "vue";
import AppSidebar from "./AppSidebar.vue";
import AppTopbar from "./AppTopbar.vue";

const sidebarOpen = ref(window.innerWidth >= 1024);
const isMobile = ref(window.innerWidth < 1024);

function handleResize() {
  const mobile = window.innerWidth < 1024;
  isMobile.value = mobile;
  if (!mobile) sidebarOpen.value = true;
  else sidebarOpen.value = false;
}

onMounted(() => window.addEventListener("resize", handleResize));
onUnmounted(() => window.removeEventListener("resize", handleResize));

function toggleSidebar() {
  sidebarOpen.value = !sidebarOpen.value;
}

provide("sidebarOpen", sidebarOpen);
provide("toggleSidebar", toggleSidebar);
provide("isMobile", isMobile);
</script>

<template>
  <div class="app-layout flex min-h-screen bg-[var(--color-bg)]">
    <!-- Sidebar overlay (mobile) -->
    <div
      v-if="isMobile && sidebarOpen"
      @click="toggleSidebar"
      class="fixed inset-0 z-40 bg-[rgba(15,23,42,0.5)] backdrop-blur-[2px]"
    />

    <!-- Sidebar -->
    <AppSidebar />

    <!-- Main area -->
    <div class="flex-1 flex flex-col min-w-0 transition-[margin-left] duration-[0.3s] ease-[cubic-bezier(0.4,0,0.2,1)]">
      <AppTopbar />

      <main class="flex-1 overflow-y-auto p-[24px] max-w-[1280px] w-full mx-auto">
        <router-view v-slot="{ Component }">
          <transition name="page" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </main>
    </div>
  </div>
</template>

<style scoped>
@media (max-width: 640px) {
  main {
    padding: 16px !important;
  }
}
</style>
