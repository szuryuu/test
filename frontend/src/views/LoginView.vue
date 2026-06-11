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
  <div class="min-h-screen flex items-center justify-center bg-[linear-gradient(135deg,#f6f8fa_0%,#e9edf2_100%)] p-[24px_16px]">
    <!-- Auth card -->
    <div
      class="w-full max-w-[420px] bg-[var(--color-surface)] rounded-[16px] shadow-[0_20px_60px_rgba(15,23,42,0.08),_0_8px_24px_rgba(15,23,42,0.04)] overflow-hidden"
    >
      <!-- Brand header -->
      <div class="pt-[32px] px-[32px] pb-0 text-center">
        <div
          class="w-[56px] h-[56px] rounded-[14px] bg-[linear-gradient(135deg,#10b981,#059669)] flex items-center justify-center mx-auto mb-[16px] text-[28px] text-white font-extrabold"
        >
          K
        </div>
        <h1 class="m-0 mb-[4px] text-[24px] font-bold text-[var(--color-text)] tracking-[-0.03em]">
          Kasir<span class="text-[#10b981]">AI</span>
        </h1>
        <p class="m-0 mb-[24px] text-[14px] text-[var(--color-text-secondary)]">
          Catat keuangan UMKM via WhatsApp
        </p>
      </div>

      <!-- Form -->
      <div class="px-[32px] pb-[32px]">
        <!-- Error -->
        <div
          v-if="errorMsg"
          class="p-[12px_16px] mb-[20px] bg-[var(--color-expense-bg)] border border-[rgba(239,68,68,0.2)] rounded-[10px] flex items-start gap-[10px]"
        >
          <i class="pi pi-exclamation-circle text-[var(--color-expense)] text-[16px] mt-[1px] shrink-0" />
          <p class="m-0 text-[13px] text-[#dc2626] leading-[1.4]">
            {{ errorMsg }}
          </p>
        </div>

        <div class="flex flex-col gap-[16px]">
          <div class="flex flex-col gap-[6px]">
            <label for="phone" class="text-[13px] font-medium text-[var(--color-text)]">
              Nomor WhatsApp
            </label>
            <div class="relative">
              <i
                class="pi pi-phone absolute left-[14px] top-[50%] -translate-y-1/2 text-[14px] text-[var(--color-text-tertiary)]"
              />
              <input
                id="phone"
                v-model="phoneNumber"
                type="tel"
                inputmode="numeric"
                maxlength="13"
                pattern="628[0-9]{8,10}"
                placeholder="6281234567890"
                class="w-full py-[12px] pr-[12px] pl-[44px] text-[14px] font-[inherit] bg-[var(--color-bg)] border border-[var(--color-border)] rounded-[10px] text-[var(--color-text)] transition-all duration-[0.15s] ease outline-0"
                @focus="$event.target.style.borderColor = 'var(--color-brand-500)'; $event.target.style.boxShadow = '0 0 0 3px rgba(16, 185, 129, 0.1)'"
                @blur="$event.target.style.borderColor = 'var(--color-border)'; $event.target.style.boxShadow = 'none'"
                @keyup.enter="handleLogin"
              />
            </div>
          </div>

          <div class="flex flex-col gap-[6px]">
            <label for="password" class="text-[13px] font-medium text-[var(--color-text)]">
              Password
            </label>
            <div class="relative">
              <i
                class="pi pi-lock absolute left-[14px] top-[50%] -translate-y-1/2 text-[14px] text-[var(--color-text-tertiary)]"
              />
              <input
                id="password"
                v-model="password"
                type="password"
                placeholder="Masukkan password"
                class="w-full py-[12px] pr-[12px] pl-[44px] text-[14px] font-[inherit] bg-[var(--color-bg)] border border-[var(--color-border)] rounded-[10px] text-[var(--color-text)] transition-all duration-[0.15s] ease outline-0"
                @focus="$event.target.style.borderColor = 'var(--color-brand-500)'; $event.target.style.boxShadow = '0 0 0 3px rgba(16, 185, 129, 0.1)'"
                @blur="$event.target.style.borderColor = 'var(--color-border)'; $event.target.style.boxShadow = 'none'"
                @keyup.enter="handleLogin"
              />
            </div>
          </div>

          <button
            @click="handleLogin"
            :disabled="isLoading"
            class="w-full p-[12px] mt-[4px] text-[15px] font-semibold font-[inherit] border-0 rounded-[10px] cursor-pointer transition-all duration-[0.2s] ease bg-[linear-gradient(135deg,#10b981,#059669)] text-white"
            :style="{ opacity: isLoading ? 0.7 : 1 }"
            @mouseenter="!isLoading && ($event.target.style.boxShadow = '0 4px 16px rgba(16, 185, 129, 0.35)')"
            @mouseleave="!isLoading && ($event.target.style.boxShadow = 'none')"
          >
            <i v-if="isLoading" class="pi pi-spin pi-spinner mr-[8px]" />
            {{ isLoading ? "Memproses..." : "Masuk" }}
          </button>
        </div>

        <div class="text-center mt-[24px] pt-[20px] border-t border-[var(--color-border-light)]">
          <p class="m-0 text-[14px] text-[var(--color-text-secondary)]">
            Belum punya akun?
            <router-link
              to="/register"
              class="text-[#10b981] font-semibold no-underline"
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
