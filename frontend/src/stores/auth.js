import { defineStore } from "pinia";
import { ref } from "vue";
import { authAPI } from "@/api/auth";

export const useAuthStore = defineStore("auth", () => {
  const token = ref(localStorage.getItem("token") || "");
  const umkm = ref(JSON.parse(localStorage.getItem("umkm") || "null"));

  const isLoggedIn = () => !!token.value;

  async function login(phoneNumber, password) {
    const res = await authAPI.login(phoneNumber, password);
    token.value = res.data.data.token;
    umkm.value = res.data.data.umkm;
    localStorage.setItem("token", token.value);
    localStorage.setItem("umkm", JSON.stringify(umkm.value));
  }

  async function register(data) {
    const res = await authAPI.register(data);
    token.value = res.data.data.token;
    umkm.value = res.data.data.umkm;
    localStorage.setItem("token", token.value);
    localStorage.setItem("umkm", JSON.stringify(umkm.value));
  }

  function logout() {
    token.value = "";
    umkm.value = null;
    localStorage.removeItem("token");
    localStorage.removeItem("umkm");
  }

  return { token, umkm, isLoggedIn, login, register, logout };
});
