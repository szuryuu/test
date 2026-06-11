<script setup>
import { ref } from "vue";
import { useRouter } from "vue-router";
import { useAuthStore } from "@/stores/auth";

const router = useRouter();
const auth = useAuthStore();

const form = ref({
  name: "",
  business_name: "",
  phone_number: "",
  password: "",
  business_type: "",
});
const isLoading = ref(false);
const errorMsg = ref("");

const businessTypes = [
  { label: "Kuliner", value: "kuliner" },
  { label: "Fashion", value: "fashion" },
  { label: "Jasa", value: "jasa" },
  { label: "Pertanian", value: "pertanian" },
  { label: "Lainnya", value: "lainnya" },
];

async function handleRegister() {
  if (!form.value.name || !form.value.business_name || !form.value.phone_number || !form.value.password) {
    errorMsg.value = "Semua field wajib diisi kecuali jenis usaha";
    return;
  }
  if (form.value.password.length < 6) {
    errorMsg.value = "Password minimal 6 karakter";
    return;
  }
  isLoading.value = true;
  errorMsg.value = "";
  try {
    await auth.register(form.value);
    router.push("/dashboard");
  } catch (err) {
    errorMsg.value =
      err.response?.data?.message || "Pendaftaran gagal. Coba lagi.";
  } finally {
    isLoading.value = false;
  }
}
</script>

<template>
  <div class="min-h-screen flex items-center justify-center bg-[linear-gradient(135deg,#f6f8fa_0%,#e9edf2_100%)] p-[24px_16px]">
    <!-- Auth card -->
    <div
      class="w-full max-w-[460px] bg-[var(--color-surface)] rounded-[16px] shadow-[0_20px_60px_rgba(15,23,42,0.08),_0_8px_24px_rgba(15,23,42,0.04)] overflow-hidden"
    >
      <!-- Brand header -->
      <div class="pt-[28px] px-[32px] pb-0 text-center">
        <h2 class="m-0 mb-[4px] text-[22px] font-bold text-[var(--color-text)] tracking-[-0.03em]">
          Daftar Kasir<span class="text-[#10b981]">AI</span>
        </h2>
        <p class="m-0 mb-[20px] text-[14px] text-[var(--color-text-secondary)]">
          Mulai catat keuangan UMKM Anda
        </p>
      </div>

      <!-- Form -->
      <div class="px-[32px] pb-[28px]">
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

        <div class="register-form-grid grid grid-cols-2 gap-[12px]">
          <div class="flex flex-col gap-[6px] col-span-2">
            <label class="text-[13px] font-medium text-[var(--color-text)]">
              Nama Lengkap
            </label>
            <input
              v-model="form.name"
              type="text"
              placeholder="Budi Santoso"
              class="w-full py-[11px] px-[14px] text-[14px] font-[inherit] bg-[var(--color-bg)] border border-[var(--color-border)] rounded-[10px] text-[var(--color-text)] transition-all duration-[0.15s] ease outline-0"
              @focus="$event.target.style.borderColor = 'var(--color-brand-500)'; $event.target.style.boxShadow = '0 0 0 3px rgba(16, 185, 129, 0.1)'"
              @blur="$event.target.style.borderColor = 'var(--color-border)'; $event.target.style.boxShadow = 'none'"
            />
          </div>

          <div class="flex flex-col gap-[6px] col-span-2">
            <label class="text-[13px] font-medium text-[var(--color-text)]">
              Nama Usaha
            </label>
            <input
              v-model="form.business_name"
              type="text"
              placeholder="Warung Nasi Budi"
              class="w-full py-[11px] px-[14px] text-[14px] font-[inherit] bg-[var(--color-bg)] border border-[var(--color-border)] rounded-[10px] text-[var(--color-text)] transition-all duration-[0.15s] ease outline-0"
              @focus="$event.target.style.borderColor = 'var(--color-brand-500)'; $event.target.style.boxShadow = '0 0 0 3px rgba(16, 185, 129, 0.1)'"
              @blur="$event.target.style.borderColor = 'var(--color-border)'; $event.target.style.boxShadow = 'none'"
            />
          </div>

          <div class="flex flex-col gap-[6px] col-span-2">
            <label class="text-[13px] font-medium text-[var(--color-text)]">
              Nomor WhatsApp
            </label>
            <input
              v-model="form.phone_number"
              type="tel"
              inputmode="numeric"
              maxlength="13"
              pattern="628[0-9]{8,10}"
              placeholder="6281234567890"
              class="w-full py-[11px] px-[14px] text-[14px] font-[inherit] bg-[var(--color-bg)] border border-[var(--color-border)] rounded-[10px] text-[var(--color-text)] transition-all duration-[0.15s] ease outline-0"
              @focus="$event.target.style.borderColor = 'var(--color-brand-500)'; $event.target.style.boxShadow = '0 0 0 3px rgba(16, 185, 129, 0.1)'"
              @blur="$event.target.style.borderColor = 'var(--color-border)'; $event.target.style.boxShadow = 'none'"
            />
          </div>

          <div class="flex flex-col gap-[6px] col-span-2">
            <label class="text-[13px] font-medium text-[var(--color-text)]">
              Password
            </label>
            <input
              v-model="form.password"
              type="password"
              placeholder="Minimal 6 karakter"
              class="w-full py-[11px] px-[14px] text-[14px] font-[inherit] bg-[var(--color-bg)] border border-[var(--color-border)] rounded-[10px] text-[var(--color-text)] transition-all duration-[0.15s] ease outline-0"
              @focus="$event.target.style.borderColor = 'var(--color-brand-500)'; $event.target.style.boxShadow = '0 0 0 3px rgba(16, 185, 129, 0.1)'"
              @blur="$event.target.style.borderColor = 'var(--color-border)'; $event.target.style.boxShadow = 'none'"
            />
          </div>

          <div class="flex flex-col gap-[6px] col-span-2">
            <label class="text-[13px] font-medium text-[var(--color-text)]">
              Jenis Usaha
              <span class="font-normal text-[var(--color-text-tertiary)]">(opsional)</span>
            </label>
            <select
              v-model="form.business_type"
              class="w-full py-[11px] px-[14px] text-[14px] font-[inherit] bg-[var(--color-bg)] border border-[var(--color-border)] rounded-[10px] text-[var(--color-text)] cursor-pointer transition-all duration-[0.15s] ease outline-0"
              @focus="$event.target.style.borderColor = 'var(--color-brand-500)'"
              @blur="$event.target.style.borderColor = 'var(--color-border)'"
            >
              <option value="">-- Pilih jenis usaha --</option>
              <option v-for="bt in businessTypes" :key="bt.value" :value="bt.value">
                {{ bt.label }}
              </option>
            </select>
          </div>
        </div>

        <button
          @click="handleRegister"
          :disabled="isLoading"
          class="w-full p-[12px] mt-[20px] text-[15px] font-semibold font-[inherit] border-0 rounded-[10px] cursor-pointer transition-all duration-[0.2s] ease bg-[linear-gradient(135deg,#10b981,#059669)] text-white"
          :style="{ opacity: isLoading ? 0.7 : 1 }"
          @mouseenter="!isLoading && ($event.target.style.boxShadow = '0 4px 16px rgba(16, 185, 129, 0.35)')"
          @mouseleave="!isLoading && ($event.target.style.boxShadow = 'none')"
        >
          <i v-if="isLoading" class="pi pi-spin pi-spinner mr-[8px]" />
          {{ isLoading ? "Mendaftarkan..." : "Daftar" }}
        </button>

        <div class="text-center mt-[20px] pt-[16px] border-t border-[var(--color-border-light)]">
          <p class="m-0 text-[14px] text-[var(--color-text-secondary)]">
            Sudah punya akun?
            <router-link
              to="/login"
              class="text-[#10b981] font-semibold no-underline"
              @mouseenter="$event.target.style.textDecoration = 'underline'"
              @mouseleave="$event.target.style.textDecoration = 'none'"
            >
              Masuk
            </router-link>
          </p>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
@media (max-width: 480px) {
  .register-form-grid {
    grid-template-columns: 1fr !important;
  }
}
</style>
