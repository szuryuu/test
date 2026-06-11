<script setup>
import { ref, reactive } from "vue";
import { useRouter } from "vue-router";
import { useAuthStore } from "@/stores/auth";

const router = useRouter();
const auth = useAuthStore();

const phoneNumber = ref("");
const password = ref("");
const isLoading = ref(false);
const errorMsg = ref("");

const errors = reactive({
  phone: "",
  password: "",
});

function validatePhone() {
  const phone = phoneNumber.value.trim();
  if (!phone) {
    errors.phone = "Wajib diisi";
  } else if (!/^628[0-9]{8,10}$/.test(phone)) {
    errors.phone = "Format: 628 (11-13 digit)";
  } else {
    errors.phone = "";
  }
}

function validatePassword() {
  if (!password.value) {
    errors.password = "Wajib diisi";
  } else if (password.value.length < 6) {
    errors.password = "Minimal 6 karakter";
  } else {
    errors.password = "";
  }
}

function validateAll() {
  validatePhone();
  validatePassword();
  return !errors.phone && !errors.password;
}

async function handleLogin() {
  if (!validateAll()) {
    errorMsg.value = "Mohon periksa kembali data yang Anda masukkan";
    return;
  }

  isLoading.value = true;
  errorMsg.value = "";
  try {
    await auth.login(phoneNumber.value, password.value);
    router.push("/dashboard");
  } catch (err) {
    errorMsg.value =
      err.response?.data?.message ||
      "Login gagal. Periksa kembali nomor dan password.";
  } finally {
    isLoading.value = false;
  }
}
</script>

<template>
  <div
    class="min-h-screen flex items-center justify-center bg-[linear-gradient(135deg,#f6f8fa_0%,#e9edf2_100%)] px-4 py-6"
  >
    <div
      class="w-full max-w-[420px] bg-[var(--color-surface)] rounded-2xl shadow-[0_20px_60px_rgba(15,23,42,0.08),_0_8px_24px_rgba(15,23,42,0.04)] overflow-hidden"
    >
      <div class="pt-8 px-8 pb-0 text-center">
        <div
          class="w-14 h-14 rounded-xl bg-[linear-gradient(135deg,#10b981,#059669)] flex items-center justify-center mx-auto mb-4 text-[28px] text-white font-extrabold"
        >
          K
        </div>
        <h1
          class="m-0 mb-1 text-2xl font-bold text-[var(--color-text)] tracking-[-0.03em]"
        >
          Kasir<span class="text-brand-500">AI</span>
        </h1>
        <p class="m-0 mb-6 text-sm text-[var(--color-text-secondary)]">
          Catat keuangan UMKM via WhatsApp
        </p>
      </div>

      <div class="px-8 pb-8">
        <div
          v-if="errorMsg"
          class="px-4 py-3 mb-5 bg-[var(--color-expense-bg)] border border-red-500/20 rounded-[10px] flex items-start gap-2.5"
        >
          <i
            class="pi pi-exclamation-circle text-[var(--color-expense)] text-base mt-0.5 shrink-0"
          />
          <p class="m-0 text-[13px] text-red-500 leading-relaxed">
            {{ errorMsg }}
          </p>
        </div>

        <div class="flex flex-col gap-4">
          <div class="flex flex-col gap-1.5">
            <div class="flex justify-between items-center">
              <label
                for="phone"
                class="text-[13px] font-medium text-[var(--color-text)]"
                >Nomor WhatsApp</label
              >
              <span
                v-if="errors.phone"
                class="text-xs font-medium text-red-500"
                >{{ errors.phone }}</span
              >
            </div>
            <div class="relative">
              <i
                class="pi pi-phone absolute left-3.5 top-1/2 -translate-y-1/2 text-sm text-[var(--color-text-tertiary)] pointer-events-none"
              />
              <input
                id="phone"
                v-model="phoneNumber"
                type="tel"
                inputmode="numeric"
                maxlength="13"
                placeholder="6281234567890"
                class="w-full py-3 pr-10 pl-11 text-sm font-[inherit] bg-[var(--color-bg)] border rounded-[10px] text-[var(--color-text)] transition-colors duration-150 ease-out outline-none focus:outline-none"
                :class="
                  errors.phone
                    ? 'border-red-500 focus:border-red-600'
                    : 'border-[var(--color-border)] focus:border-brand-500'
                "
                @input="
                  phoneNumber = $event.target.value.replace(/\D/g, '');
                  validatePhone();
                "
                @keyup.enter="handleLogin"
                @blur="validatePhone"
              />
              <i
                v-if="errors.phone"
                class="pi pi-exclamation-circle absolute right-3.5 top-1/2 -translate-y-1/2 text-sm text-red-500 pointer-events-none"
              />
            </div>
          </div>

          <div class="flex flex-col gap-1.5">
            <div class="flex justify-between items-center">
              <label
                for="password"
                class="text-[13px] font-medium text-[var(--color-text)]"
                >Password</label
              >
              <span
                v-if="errors.password"
                class="text-xs font-medium text-red-500"
                >{{ errors.password }}</span
              >
            </div>
            <div class="relative">
              <i
                class="pi pi-lock absolute left-3.5 top-1/2 -translate-y-1/2 text-sm text-[var(--color-text-tertiary)] pointer-events-none"
              />
              <input
                id="password"
                v-model="password"
                type="password"
                placeholder="Masukkan password"
                class="w-full py-3 pr-10 pl-11 text-sm font-[inherit] bg-[var(--color-bg)] border rounded-[10px] text-[var(--color-text)] transition-colors duration-150 ease-out outline-none focus:outline-none"
                :class="
                  errors.password
                    ? 'border-red-500 focus:border-red-600'
                    : 'border-[var(--color-border)] focus:border-brand-500'
                "
                @input="validatePassword"
                @keyup.enter="handleLogin"
                @blur="validatePassword"
              />
              <i
                v-if="errors.password"
                class="pi pi-exclamation-circle absolute right-3.5 top-1/2 -translate-y-1/2 text-sm text-red-500 pointer-events-none"
              />
            </div>
          </div>

          <button
            @click="handleLogin"
            :disabled="isLoading"
            class="w-full p-3 mt-2 text-[15px] font-semibold font-[inherit] border-0 rounded-[10px] cursor-pointer transition-all duration-200 ease-out bg-[linear-gradient(135deg,#10b981,#059669)] text-white disabled:opacity-70 hover:shadow-[0_4px_16px_rgba(16,185,129,0.35)] outline-none focus:outline-none"
          >
            <i v-if="isLoading" class="pi pi-spin pi-spinner mr-2" />
            {{ isLoading ? "Memproses..." : "Masuk" }}
          </button>
        </div>

        <div
          class="text-center mt-6 pt-5 border-t border-[var(--color-border-light)]"
        >
          <p class="m-0 text-sm text-[var(--color-text-secondary)]">
            Belum punya akun?
            <router-link
              to="/register"
              class="text-brand-500 font-semibold no-underline hover:underline"
            >
              Daftar di sini
            </router-link>
          </p>
        </div>
      </div>
    </div>
  </div>
</template>
