import { createRouter, createWebHistory } from "vue-router";
import SIA from "@/views/infoarea.vue";
import { nextTick } from "vue";

const routes = [
  {
    path: "/",
    name: "Home",
    component: SIA,
  }
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

export default router;
