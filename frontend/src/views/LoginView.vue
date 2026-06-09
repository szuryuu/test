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
    style="
      min-height: 100vh;
      display: flex;
      align-items: center;
      justify-content: center;
      background: linear-gradient(135deg, #f6f8fa 0%, #e9edf2 100%);
      padding: 24px 16px;
    "
  >
    <!-- Auth card -->
    <div
      style="
        width: 100%; max-width: 420px;
        background: var(--color-surface);
        border-radius: 16px;
        box-shadow: 0 20px 60px rgba(15, 23, 42, 0.08), 0 8px 24px rgba(15, 23, 42, 0.04);
        overflow: hidden;
      "
    >
      <!-- Brand header -->
      <div
        style="
          padding: 32px 32px 0;
          text-align: center;
        "
      >
        <div
          style="
            width: 56px; height: 56px; border-radius: 14px;
            background: linear-gradient(135deg, #10b981, #059669);
            display: flex; align-items: center; justify-content: center;
            margin: 0 auto 16px;
            font-size: 28px; color: white; font-weight: 800;
          "
        >
          K
        </div>
        <h1
          style="
            margin: 0 0 4px; font-size: 24px; font-weight: 700;
            color: var(--color-text); letter-spacing: -0.03em;
          "
        >
          Kasir<span style="color: #10b981">AI</span>
        </h1>
        <p
          style="
            margin: 0 0 24px; font-size: 14px;
            color: var(--color-text-secondary);
          "
        >
          Catat keuangan UMKM via WhatsApp
        </p>
      </div>

      <!-- Form -->
      <div style="padding: 0 32px 32px">
        <!-- Error -->
        <div
          v-if="errorMsg"
          style="
            padding: 12px 16px; margin-bottom: 20px;
            background: var(--color-expense-bg);
            border: 1px solid rgba(239, 68, 68, 0.2);
            border-radius: 10px;
            display: flex; align-items: flex-start; gap: 10px;
          "
        >
          <i class="pi pi-exclamation-circle" style="color: var(--color-expense); font-size: 16px; margin-top: 1px; flex-shrink: 0" />
          <p style="margin: 0; font-size: 13px; color: #dc2626; line-height: 1.4">
            {{ errorMsg }}
          </p>
        </div>

        <div style="display: flex; flex-direction: column; gap: 16px">
          <div style="display: flex; flex-direction: column; gap: 6px">
            <label
              for="phone"
              style="font-size: 13px; font-weight: 500; color: var(--color-text)"
            >
              Nomor WhatsApp
            </label>
            <div style="position: relative">
              <i
                class="pi pi-phone"
                style="
                  position: absolute; left: 14px; top: 50%;
                  transform: translateY(-50%); font-size: 14px;
                  color: var(--color-text-tertiary);
                "
              />
              <input
                id="phone"
                v-model="phoneNumber"
                type="text"
                placeholder="6281234567890"
                style="
                  width: 100%; padding: 12px 12px 12px 44px;
                  font-size: 14px; font-family: inherit;
                  background: var(--color-bg);
                  border: 1px solid var(--color-border);
                  border-radius: 10px;
                  color: var(--color-text);
                  transition: all 0.15s ease;
                  outline: none;
                "
                @focus="$event.target.style.borderColor = 'var(--color-brand-500)'; $event.target.style.boxShadow = '0 0 0 3px rgba(16, 185, 129, 0.1)'"
                @blur="$event.target.style.borderColor = 'var(--color-border)'; $event.target.style.boxShadow = 'none'"
                @keyup.enter="handleLogin"
              />
            </div>
          </div>

          <div style="display: flex; flex-direction: column; gap: 6px">
            <label
              for="password"
              style="font-size: 13px; font-weight: 500; color: var(--color-text)"
            >
              Password
            </label>
            <div style="position: relative">
              <i
                class="pi pi-lock"
                style="
                  position: absolute; left: 14px; top: 50%;
                  transform: translateY(-50%); font-size: 14px;
                  color: var(--color-text-tertiary);
                "
              />
              <input
                id="password"
                v-model="password"
                type="password"
                placeholder="Masukkan password"
                style="
                  width: 100%; padding: 12px 12px 12px 44px;
                  font-size: 14px; font-family: inherit;
                  background: var(--color-bg);
                  border: 1px solid var(--color-border);
                  border-radius: 10px;
                  color: var(--color-text);
                  transition: all 0.15s ease;
                  outline: none;
                "
                @focus="$event.target.style.borderColor = 'var(--color-brand-500)'; $event.target.style.boxShadow = '0 0 0 3px rgba(16, 185, 129, 0.1)'"
                @blur="$event.target.style.borderColor = 'var(--color-border)'; $event.target.style.boxShadow = 'none'"
                @keyup.enter="handleLogin"
              />
            </div>
          </div>

          <button
            @click="handleLogin"
            :disabled="isLoading"
            style="
              width: 100%; padding: 12px; margin-top: 4px;
              font-size: 15px; font-weight: 600; font-family: inherit;
              border: none; border-radius: 10px; cursor: pointer;
              transition: all 0.2s ease;
              background: linear-gradient(135deg, #10b981, #059669);
              color: white;
            "
            :style="{ opacity: isLoading ? 0.7 : 1 }"
            @mouseenter="!isLoading && ($event.target.style.boxShadow = '0 4px 16px rgba(16, 185, 129, 0.35)')"
            @mouseleave="!isLoading && ($event.target.style.boxShadow = 'none')"
          >
            <i v-if="isLoading" class="pi pi-spin pi-spinner" style="margin-right: 8px" />
            {{ isLoading ? "Memproses..." : "Masuk" }}
          </button>
        </div>

        <div
          style="
            text-align: center; margin-top: 24px; padding-top: 20px;
            border-top: 1px solid var(--color-border-light);
          "
        >
          <p style="margin: 0; font-size: 14px; color: var(--color-text-secondary)">
            Belum punya akun?
            <router-link
              to="/register"
              style="
                color: #10b981; font-weight: 600; text-decoration: none;
              "
              @mouseenter="$event.target.style.textDecoration = 'underline'"
              @mouseleave="$event.target.style.textDecoration = 'none'"
            >
              Daftar di sini
            </router-link>
          </p>
        </div>
      </div>
    </div>
  </div>
</template>
