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
        width: 100%; max-width: 460px;
        background: var(--color-surface);
        border-radius: 16px;
        box-shadow: 0 20px 60px rgba(15, 23, 42, 0.08), 0 8px 24px rgba(15, 23, 42, 0.04);
        overflow: hidden;
      "
    >
      <!-- Brand header -->
      <div
        style="
          padding: 28px 32px 0;
          text-align: center;
        "
      >
        <h2
          style="
            margin: 0 0 4px; font-size: 22px; font-weight: 700;
            color: var(--color-text); letter-spacing: -0.03em;
          "
        >
          Daftar Kasir<span style="color: #10b981">AI</span>
        </h2>
        <p
          style="
            margin: 0 0 20px; font-size: 14px;
            color: var(--color-text-secondary);
          "
        >
          Mulai catat keuangan UMKM Anda
        </p>
      </div>

      <!-- Form -->
      <div style="padding: 0 32px 28px">
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

        <div
          style="
            display: grid;
            grid-template-columns: 1fr 1fr;
            gap: 12px;
          "
        >
          <div
            style="display: flex; flex-direction: column; gap: 6px; grid-column: span 2"
          >
            <label style="font-size: 13px; font-weight: 500; color: var(--color-text)">
              Nama Lengkap
            </label>
            <input
              v-model="form.name"
              type="text"
              placeholder="Budi Santoso"
              style="
                width: 100%; padding: 11px 14px; font-size: 14px; font-family: inherit;
                background: var(--color-bg);
                border: 1px solid var(--color-border);
                border-radius: 10px;
                color: var(--color-text);
                transition: all 0.15s ease;
                outline: none;
              "
              @focus="$event.target.style.borderColor = 'var(--color-brand-500)'; $event.target.style.boxShadow = '0 0 0 3px rgba(16, 185, 129, 0.1)'"
              @blur="$event.target.style.borderColor = 'var(--color-border)'; $event.target.style.boxShadow = 'none'"
            />
          </div>

          <div
            style="display: flex; flex-direction: column; gap: 6px; grid-column: span 2"
          >
            <label style="font-size: 13px; font-weight: 500; color: var(--color-text)">
              Nama Usaha
            </label>
            <input
              v-model="form.business_name"
              type="text"
              placeholder="Warung Nasi Budi"
              style="
                width: 100%; padding: 11px 14px; font-size: 14px; font-family: inherit;
                background: var(--color-bg);
                border: 1px solid var(--color-border);
                border-radius: 10px;
                color: var(--color-text);
                transition: all 0.15s ease;
                outline: none;
              "
              @focus="$event.target.style.borderColor = 'var(--color-brand-500)'; $event.target.style.boxShadow = '0 0 0 3px rgba(16, 185, 129, 0.1)'"
              @blur="$event.target.style.borderColor = 'var(--color-border)'; $event.target.style.boxShadow = 'none'"
            />
          </div>

          <div
            style="display: flex; flex-direction: column; gap: 6px; grid-column: span 2"
          >
            <label style="font-size: 13px; font-weight: 500; color: var(--color-text)">
              Nomor WhatsApp
            </label>
            <input
              v-model="form.phone_number"
              type="text"
              placeholder="6281234567890"
              style="
                width: 100%; padding: 11px 14px; font-size: 14px; font-family: inherit;
                background: var(--color-bg);
                border: 1px solid var(--color-border);
                border-radius: 10px;
                color: var(--color-text);
                transition: all 0.15s ease;
                outline: none;
              "
              @focus="$event.target.style.borderColor = 'var(--color-brand-500)'; $event.target.style.boxShadow = '0 0 0 3px rgba(16, 185, 129, 0.1)'"
              @blur="$event.target.style.borderColor = 'var(--color-border)'; $event.target.style.boxShadow = 'none'"
            />
          </div>

          <div
            style="display: flex; flex-direction: column; gap: 6px; grid-column: span 2"
          >
            <label style="font-size: 13px; font-weight: 500; color: var(--color-text)">
              Password
            </label>
            <input
              v-model="form.password"
              type="password"
              placeholder="Minimal 6 karakter"
              style="
                width: 100%; padding: 11px 14px; font-size: 14px; font-family: inherit;
                background: var(--color-bg);
                border: 1px solid var(--color-border);
                border-radius: 10px;
                color: var(--color-text);
                transition: all 0.15s ease;
                outline: none;
              "
              @focus="$event.target.style.borderColor = 'var(--color-brand-500)'; $event.target.style.boxShadow = '0 0 0 3px rgba(16, 185, 129, 0.1)'"
              @blur="$event.target.style.borderColor = 'var(--color-border)'; $event.target.style.boxShadow = 'none'"
            />
          </div>

          <div
            style="display: flex; flex-direction: column; gap: 6px; grid-column: span 2"
          >
            <label style="font-size: 13px; font-weight: 500; color: var(--color-text)">
              Jenis Usaha
              <span style="font-weight: 400; color: var(--color-text-tertiary)">(opsional)</span>
            </label>
            <select
              v-model="form.business_type"
              style="
                width: 100%; padding: 11px 14px; font-size: 14px; font-family: inherit;
                background: var(--color-bg);
                border: 1px solid var(--color-border);
                border-radius: 10px;
                color: var(--color-text);
                cursor: pointer; transition: all 0.15s ease;
                outline: none;
              "
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
          style="
            width: 100%; padding: 12px; margin-top: 20px;
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
          {{ isLoading ? "Mendaftarkan..." : "Daftar" }}
        </button>

        <div
          style="
            text-align: center; margin-top: 20px; padding-top: 16px;
            border-top: 1px solid var(--color-border-light);
          "
        >
          <p style="margin: 0; font-size: 14px; color: var(--color-text-secondary)">
            Sudah punya akun?
            <router-link
              to="/login"
              style="
                color: #10b981; font-weight: 600; text-decoration: none;
              "
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
  div[style*="grid-template-columns: 1fr 1fr"] {
    grid-template-columns: 1fr !important;
  }
}
</style>
