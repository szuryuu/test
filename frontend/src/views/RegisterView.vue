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
    class="flex align-items-center justify-content-center min-h-screen py-4"
    style="background: var(--color-pale-gray)"
  >
    <div
      class="w-full px-4"
      style="max-width: 440px"
    >
      <!-- Brand -->
      <div class="text-center mb-5">
        <h1
          class="text-3xl m-0 mb-2"
          style="font-weight: 700; color: var(--color-near-black); letter-spacing: -0.03em"
        >
          Daftar Kasir<span style="color: var(--color-apple-blue)">AI</span>
        </h1>
        <p class="m-0" style="color: var(--color-text-secondary); font-size: 15px">
          Mulai catat keuangan UMKM Anda
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
        <div
          v-if="errorMsg"
          class="px-4 py-3 mb-4"
          style="
            background: #fef2f2;
            color: #dc2626;
            border-radius: var(--radius-sm);
            font-size: 13px;
          "
        >
          {{ errorMsg }}
        </div>

        <div class="flex flex-column gap-3">
          <div class="flex flex-column gap-2">
            <label class="text-sm" style="font-weight: 500; color: var(--color-near-black)">
              Nama Lengkap
            </label>
            <input
              v-model="form.name"
              type="text"
              placeholder="Budi Santoso"
              class="w-full px-4 py-3 text-sm"
              style="
                background: var(--color-pale-gray);
                border: 1px solid var(--color-border-soft);
                border-radius: var(--radius-sm);
                color: var(--color-near-black);
              "
            />
          </div>

          <div class="flex flex-column gap-2">
            <label class="text-sm" style="font-weight: 500; color: var(--color-near-black)">
              Nama Usaha
            </label>
            <input
              v-model="form.business_name"
              type="text"
              placeholder="Warung Nasi Budi"
              class="w-full px-4 py-3 text-sm"
              style="
                background: var(--color-pale-gray);
                border: 1px solid var(--color-border-soft);
                border-radius: var(--radius-sm);
                color: var(--color-near-black);
              "
            />
          </div>

          <div class="flex flex-column gap-2">
            <label class="text-sm" style="font-weight: 500; color: var(--color-near-black)">
              Nomor WhatsApp
            </label>
            <input
              v-model="form.phone_number"
              type="text"
              placeholder="6281234567890"
              class="w-full px-4 py-3 text-sm"
              style="
                background: var(--color-pale-gray);
                border: 1px solid var(--color-border-soft);
                border-radius: var(--radius-sm);
                color: var(--color-near-black);
              "
            />
          </div>

          <div class="flex flex-column gap-2">
            <label class="text-sm" style="font-weight: 500; color: var(--color-near-black)">
              Password
            </label>
            <input
              v-model="form.password"
              type="password"
              placeholder="Minimal 6 karakter"
              class="w-full px-4 py-3 text-sm"
              style="
                background: var(--color-pale-gray);
                border: 1px solid var(--color-border-soft);
                border-radius: var(--radius-sm);
                color: var(--color-near-black);
              "
            />
          </div>

          <div class="flex flex-column gap-2">
            <label class="text-sm" style="font-weight: 500; color: var(--color-near-black)">
              Jenis Usaha <span style="color: var(--color-text-secondary); font-weight: 400">(opsional)</span>
            </label>
            <select
              v-model="form.business_type"
              class="w-full px-4 py-3 text-sm"
              style="
                background: var(--color-pale-gray);
                border: 1px solid var(--color-border-soft);
                border-radius: var(--radius-sm);
                color: var(--color-near-black);
              "
            >
              <option value="">-- Pilih jenis usaha --</option>
              <option v-for="bt in businessTypes" :key="bt.value" :value="bt.value">
                {{ bt.label }}
              </option>
            </select>
          </div>

          <button
            @click="handleRegister"
            :disabled="isLoading"
            class="w-full py-3 mt-2 cursor-pointer text-sm transition-all transition-duration-200"
            style="
              background: var(--color-apple-blue);
              color: var(--color-pure-white);
              border: none;
              border-radius: var(--radius-sm);
              font-weight: 600;
              opacity: isLoading ? 0.7 : 1;
            "
          >
            {{ isLoading ? "Mendaftarkan..." : "Daftar" }}
          </button>
        </div>

        <p class="text-center mt-4 mb-0 text-sm" style="color: var(--color-text-secondary)">
          Sudah punya akun?
          <router-link
            to="/login"
            style="color: var(--color-apple-blue); font-weight: 500"
          >
            Masuk
          </router-link>
        </p>
      </div>
    </div>
  </div>
</template>
