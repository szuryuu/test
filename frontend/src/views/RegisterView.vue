<script setup>
import { ref, reactive } from "vue";
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

const errors = reactive({
  name: "",
  business_name: "",
  phone_number: "",
  password: "",
});

const businessTypes = [
  { label: "Kuliner", value: "kuliner" },
  { label: "Fashion", value: "fashion" },
  { label: "Jasa", value: "jasa" },
  { label: "Pertanian", value: "pertanian" },
  { label: "Lainnya", value: "lainnya" },
];

function validateName() {
  const v = form.value.name.trim();
  if (!v) {
    errors.name = "Wajib diisi";
  } else if (v.length < 2) {
    errors.name = "Minimal 2 karakter";
  } else {
    errors.name = "";
  }
}

function validateBusinessName() {
  if (!form.value.business_name.trim()) {
    errors.business_name = "Wajib diisi";
  } else {
    errors.business_name = "";
  }
}

function validatePhone() {
  const phone = form.value.phone_number.trim();
  if (!phone) {
    errors.phone_number = "Wajib diisi";
  } else if (!/^628[0-9]{8,10}$/.test(phone)) {
    errors.phone_number = "Format: 628 (11-13 digit)";
  } else {
    errors.phone_number = "";
  }
}

function validatePassword() {
  if (!form.value.password) {
    errors.password = "Wajib diisi";
  } else if (form.value.password.length < 6) {
    errors.password = "Minimal 6 karakter";
  } else {
    errors.password = "";
  }
}

function validateAll() {
  validateName();
  validateBusinessName();
  validatePhone();
  validatePassword();
  return !Object.values(errors).some(Boolean);
}

async function handleRegister() {
  if (!validateAll()) {
    errorMsg.value = "Mohon periksa kembali data yang Anda masukkan";
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
  <div
    class="min-h-screen flex items-center justify-center bg-[linear-gradient(135deg,#f6f8fa_0%,#e9edf2_100%)] px-4 py-6"
  >
    <div
      class="w-full max-w-[460px] bg-[var(--color-surface)] rounded-2xl shadow-[0_20px_60px_rgba(15,23,42,0.08),_0_8px_24px_rgba(15,23,42,0.04)] overflow-hidden"
    >
      <div class="pt-7 px-8 pb-0 text-center">
        <h2
          class="m-0 mb-1 text-[22px] font-bold text-[var(--color-text)] tracking-[-0.03em]"
        >
          Daftar Kasir<span class="text-brand-500">AI</span>
        </h2>
        <p class="m-0 mb-5 text-sm text-[var(--color-text-secondary)]">
          Mulai catat keuangan UMKM Anda
        </p>
      </div>

      <div class="px-8 pb-7">
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
              <label class="text-[13px] font-medium text-[var(--color-text)]"
                >Nama Lengkap</label
              >
              <span
                v-if="errors.name"
                class="text-xs font-medium text-red-500"
                >{{ errors.name }}</span
              >
            </div>
            <div class="relative">
              <input
                v-model="form.name"
                type="text"
                placeholder="Budi Santoso"
                class="w-full py-2.5 pr-10 pl-3.5 text-sm font-[inherit] bg-[var(--color-bg)] border rounded-[10px] text-[var(--color-text)] transition-colors duration-150 ease-out outline-none focus:outline-none"
                :class="
                  errors.name
                    ? 'border-red-500 focus:border-red-600'
                    : 'border-[var(--color-border)] focus:border-brand-500'
                "
                @input="validateName"
                @blur="validateName"
              />
              <i
                v-if="errors.name"
                class="pi pi-exclamation-circle absolute right-3.5 top-1/2 -translate-y-1/2 text-sm text-red-500 pointer-events-none"
              />
            </div>
          </div>

          <div class="flex flex-col gap-1.5">
            <div class="flex justify-between items-center">
              <label class="text-[13px] font-medium text-[var(--color-text)]"
                >Nama Usaha</label
              >
              <span
                v-if="errors.business_name"
                class="text-xs font-medium text-red-500"
                >{{ errors.business_name }}</span
              >
            </div>
            <div class="relative">
              <input
                v-model="form.business_name"
                type="text"
                placeholder="Warung Nasi Budi"
                class="w-full py-2.5 pr-10 pl-3.5 text-sm font-[inherit] bg-[var(--color-bg)] border rounded-[10px] text-[var(--color-text)] transition-colors duration-150 ease-out outline-none focus:outline-none"
                :class="
                  errors.business_name
                    ? 'border-red-500 focus:border-red-600'
                    : 'border-[var(--color-border)] focus:border-brand-500'
                "
                @input="validateBusinessName"
                @blur="validateBusinessName"
              />
              <i
                v-if="errors.business_name"
                class="pi pi-exclamation-circle absolute right-3.5 top-1/2 -translate-y-1/2 text-sm text-red-500 pointer-events-none"
              />
            </div>
          </div>

          <div class="flex flex-col gap-1.5">
            <div class="flex justify-between items-center">
              <label class="text-[13px] font-medium text-[var(--color-text)]"
                >Nomor WhatsApp</label
              >
              <span
                v-if="errors.phone_number"
                class="text-xs font-medium text-red-500"
                >{{ errors.phone_number }}</span
              >
            </div>
            <div class="relative">
              <input
                v-model="form.phone_number"
                type="tel"
                inputmode="numeric"
                maxlength="13"
                placeholder="6281234567890"
                class="w-full py-2.5 pr-10 pl-3.5 text-sm font-[inherit] bg-[var(--color-bg)] border rounded-[10px] text-[var(--color-text)] transition-colors duration-150 ease-out outline-none focus:outline-none"
                :class="
                  errors.phone_number
                    ? 'border-red-500 focus:border-red-600'
                    : 'border-[var(--color-border)] focus:border-brand-500'
                "
                @input="
                  form.phone_number = $event.target.value.replace(/\D/g, '');
                  validatePhone();
                "
                @blur="validatePhone"
              />
              <i
                v-if="errors.phone_number"
                class="pi pi-exclamation-circle absolute right-3.5 top-1/2 -translate-y-1/2 text-sm text-red-500 pointer-events-none"
              />
            </div>
          </div>

          <div class="flex flex-col gap-1.5">
            <div class="flex justify-between items-center">
              <label class="text-[13px] font-medium text-[var(--color-text)]"
                >Password</label
              >
              <span
                v-if="errors.password"
                class="text-xs font-medium text-red-500"
                >{{ errors.password }}</span
              >
            </div>
            <div class="relative">
              <input
                v-model="form.password"
                type="password"
                placeholder="Minimal 6 karakter"
                class="w-full py-2.5 pr-10 pl-3.5 text-sm font-[inherit] bg-[var(--color-bg)] border rounded-[10px] text-[var(--color-text)] transition-colors duration-150 ease-out outline-none focus:outline-none"
                :class="
                  errors.password
                    ? 'border-red-500 focus:border-red-600'
                    : 'border-[var(--color-border)] focus:border-brand-500'
                "
                @input="validatePassword"
                @blur="validatePassword"
              />
              <i
                v-if="errors.password"
                class="pi pi-exclamation-circle absolute right-3.5 top-1/2 -translate-y-1/2 text-sm text-red-500 pointer-events-none"
              />
            </div>
          </div>

          <div class="flex flex-col gap-1.5">
            <div class="flex justify-between items-center">
              <label class="text-[13px] font-medium text-[var(--color-text)]">
                Jenis Usaha
                <span class="font-normal text-[var(--color-text-tertiary)]"
                  >(opsional)</span
                >
              </label>
            </div>
            <div class="relative">
              <select
                v-model="form.business_type"
                class="w-full py-2.5 pr-10 pl-3.5 text-sm font-[inherit] bg-[var(--color-bg)] border border-[var(--color-border)] rounded-[10px] text-[var(--color-text)] cursor-pointer transition-colors duration-150 ease-out outline-none focus:outline-none focus:border-brand-500 appearance-none"
              >
                <option value="">-- Pilih jenis usaha --</option>
                <option
                  v-for="bt in businessTypes"
                  :key="bt.value"
                  :value="bt.value"
                >
                  {{ bt.label }}
                </option>
              </select>
              <i
                class="pi pi-chevron-down absolute right-3.5 top-1/2 -translate-y-1/2 text-xs text-[var(--color-text-tertiary)] pointer-events-none"
              />
            </div>
          </div>
        </div>

        <button
          @click="handleRegister"
          :disabled="isLoading"
          class="w-full p-3 mt-5 text-[15px] font-semibold font-[inherit] border-0 rounded-[10px] cursor-pointer transition-all duration-200 ease-out bg-[linear-gradient(135deg,#10b981,#059669)] text-white disabled:opacity-70 hover:shadow-[0_4px_16px_rgba(16,185,129,0.35)] outline-none focus:outline-none"
        >
          <i v-if="isLoading" class="pi pi-spin pi-spinner mr-2" />
          {{ isLoading ? "Mendaftarkan..." : "Daftar" }}
        </button>

        <div
          class="text-center mt-5 pt-4 border-t border-[var(--color-border-light)]"
        >
          <p class="m-0 text-sm text-[var(--color-text-secondary)]">
            Sudah punya akun?
            <router-link
              to="/login"
              class="text-brand-500 font-semibold no-underline hover:underline"
            >
              Masuk
            </router-link>
          </p>
        </div>
      </div>
    </div>
  </div>
</template>
