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
  <div
    class="app-layout"
    style="display: flex; min-height: 100vh; background: var(--color-bg)"
  >
    <!-- Sidebar overlay (mobile) -->
    <div
      v-if="isMobile && sidebarOpen"
      @click="toggleSidebar"
      style="
        position: fixed; inset: 0; z-index: 40;
        background: rgba(15, 23, 42, 0.5);
        -webkit-backdrop-filter: blur(2px);
        backdrop-filter: blur(2px);
      "
    />

    <!-- Sidebar -->
    <AppSidebar />

    <!-- Main area -->
    <div
      style="
        flex: 1;
        display: flex;
        flex-direction: column;
        min-width: 0;
        transition: margin-left 0.3s cubic-bezier(0.4, 0, 0.2, 1);
      "
    >
      <AppTopbar />

      <main
        style="
          flex: 1;
          overflow-y: auto;
          padding: 24px;
          max-width: 1280px;
          width: 100%;
          margin: 0 auto;
        "
      >
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
