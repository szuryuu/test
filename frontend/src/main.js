import { createApp } from "vue";
import { createPinia } from "pinia";
import PrimeVue from "primevue/config";
import App from "./App.vue";
import router from "./router";

/* PrimeVue core + theme */
import "primeicons/primeicons.css";

/* Tailwind CSS 4 + design tokens */
import "./assets/styles/main.css";

const app = createApp(App);

app.use(createPinia());
app.use(router);
app.use(PrimeVue, {
  ripple: true,
  theme: "none", // Tailwind handles all styling
});

app.mount("#app");
