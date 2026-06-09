import { createApp } from "vue";
import { createPinia } from "pinia";
import PrimeVue from "primevue/config";
import App from "./App.vue";
import router from "./router";

/* PrimeIcons */
import "primeicons/primeicons.css";

/* Tailwind CSS 4 + design system */
import "./assets/styles/main.css";

const app = createApp(App);

app.use(createPinia());
app.use(router);
app.use(PrimeVue, {
  ripple: true,
  theme: "none", // Full custom styling via design system
});

app.mount("#app");
