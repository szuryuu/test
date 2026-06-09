<script setup>
import AppSidebar from "./AppSidebar.vue";
import AppTopbar from "./AppTopbar.vue";
</script>

<template>
  <div class="flex h-screen overflow-hidden">
    <AppSidebar />
    <div class="flex flex-column flex-1 overflow-hidden">
      <AppTopbar />
      <main
        class="flex-1 overflow-auto p-6"
        style="background: var(--color-pale-gray)"
      >
        <slot />
      </main>
    </div>
  </div>
</template>
