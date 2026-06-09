<script setup>
import { ref } from "vue";
import { useRouter } from "vue-router";
import { useAuthStore } from "@/stores/auth";

const router = useRouter();
const auth = useAuthStore();

const phoneNumber = ref("");
const password = ref("");
const isLoading = ref(false);
const errorMsg = ref("");

async function handleLogin() {
  if (!phoneNumber.value || !password.value) {
    errorMsg.value = "Nomor WhatsApp dan password wajib diisi";
    return;
  }
  isLoading.value = true;
  errorMsg.value = "";
  try {
    await auth.login(phoneNumber.value, password.value);
    router.push("/dashboard");
  } catch (err) {
    errorMsg.value =
      err.response?.data?.message || "Login gagal. Periksa kembali nomor dan password.";
  } finally {
    isLoading.value = false;
  }
}
</script>

<template>
  <div
    class="flex align-items-center justify-content-center min-h-screen"
    style="background: var(--color-pale-gray)"
  >
    <div
      class="w-full px-4"
      style="max-width: 400px"
    >
      <!-- Brand -->
      <div class="text-center mb-6">
        <h1
          class="text-4xl m-0 mb-2"
          style="font-weight: 700; color: var(--color-near-black); letter-spacing: -0.03em"
        >
          Kasir<span style="color: var(--color-apple-blue)">AI</span>
        </h1>
        <p class="m-0" style="color: var(--color-text-secondary); font-size: 15px">
          Catat keuangan UMKM via WhatsApp
        </p>
      </div>

      <!-- Card -->
      <div
        class="p-5"
        style="
          background: var(--color-pure-white);
          border-radius: var(--radius-md);
          box-shadow: 0 1px 3px rgba(0,0,0,0.04);
        "
      >
        <h2 class="text-xl m-0 mb-4" style="font-weight: 600; color: var(--color-near-black)">
          Masuk
        </h2>

        <!-- Error -->
        <div
          v-if="errorMsg"
          class="px-4 py-3 mb-4 text-sm"
          style="
            background: #fef2f2;
            color: #dc2626;
            border-radius: var(--radius-sm);
            font-size: 13px;
          "
        >
          {{ errorMsg }}
        </div>

        <!-- Form -->
        <div class="flex flex-column gap-4">
          <div class="flex flex-column gap-2">
            <label
              for="phone"
              class="text-sm"
              style="font-weight: 500; color: var(--color-near-black)"
            >
              Nomor WhatsApp
            </label>
            <input
              id="phone"
              v-model="phoneNumber"
              type="text"
              placeholder="6281234567890"
              class="w-full px-4 py-3 text-sm border-none outline-none transition-all transition-duration-200"
              style="
                background: var(--color-pale-gray);
                color: var(--color-near-black);
                border-radius: var(--radius-sm);
                border: 1px solid var(--color-border-soft);
              "
              @keyup.enter="handleLogin"
            />
          </div>

          <div class="flex flex-column gap-2">
            <label
              for="password"
              class="text-sm"
              style="font-weight: 500; color: var(--color-near-black)"
            >
              Password
            </label>
            <input
              id="password"
              v-model="password"
              type="password"
              placeholder="••••••"
              class="w-full px-4 py-3 text-sm border-none outline-none transition-all transition-duration-200"
              style="
                background: var(--color-pale-gray);
                color: var(--color-near-black);
                border-radius: var(--radius-sm);
                border: 1px solid var(--color-border-soft);
              "
              @keyup.enter="handleLogin"
            />
          </div>

          <button
            @click="handleLogin"
            :disabled="isLoading"
            class="w-full py-3 cursor-pointer text-sm transition-all transition-duration-200"
            style="
              background: var(--color-apple-blue);
              color: var(--color-pure-white);
              border: none;
              border-radius: var(--radius-sm);
              font-weight: 600;
              opacity: isLoading ? 0.7 : 1;
            "
          >
            {{ isLoading ? "Memproses..." : "Masuk" }}
          </button>
        </div>

        <p class="text-center mt-4 mb-0 text-sm" style="color: var(--color-text-secondary)">
          Belum punya akun?
          <router-link
            to="/register"
            style="color: var(--color-apple-blue); font-weight: 500"
          >
            Daftar di sini
          </router-link>
        </p>
      </div>
    </div>
  </div>
</template>
