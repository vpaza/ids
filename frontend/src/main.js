import { createApp } from "vue";
import { createPinia } from "pinia";
import App from "./App.vue";
import router from "./router";

import "./assets/css/style.css";

const pinia = createPinia();

const app = createApp(App);
app.config.devtools = true;

app.use(router).use(pinia).mount("#app");
