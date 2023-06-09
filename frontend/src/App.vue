<template>
  <div class="flex flex-col min-h-screen">
    <header class="z-50 bg-neutral-700 text-white top-0 p-4 w-full shadow-md mb-[10px]">
      <div class="py-1">
        <div class="grid grid-cols-4 items-center justify-between">
          <div class="col-span-1 title">
            <span class="facility">{{ fac.facility.id }}</span> <span class="idscolor">IDS</span>
          </div>
          <div class="col-span-2">
            <div class="flex items-center" v-if="loggedIn">
              <label class="block text-gray-100 font-bold pr-4"> Views: </label>
              <select
                v-model="view"
                id="view"
                class="block w-full bg-neutral-800 text-white py-3 px-4 pr-8 rounded leading-tight focus:outline-none focus:border-gray-500"
                @change="changeView()"
              >
                <option v-for="(v, i) in fac.views" :key="i" :value="i">{{ v.name }}</option>
              </select>
            </div>
          </div>
          <div class="col-span-1 text-right ml-auto">
            <Clock :timezone="fac.timezone.name" />
          </div>
        </div>
      </div>
    </header>
      <main class="flex-1 overflow-y-auto p-5">
        <router-view v-if="loggedIn" />
      </main>
    <footer class="z-50 bg-neutral-700 text-white bottom-0 p-4">
      <div v-if="loggedIn">
        footer
      </div>
      <div v-else>
        <a :href="`${fac.facility.api}/v1/oauth/login?redirect=${location}`">Login</a>
      </div>
    </footer>
  </div>
</template>

<script setup>
import { ref, onMounted } from "vue";
import { storeToRefs } from "pinia";
import { useViewStore } from "@/store/viewstore";
import Clock from "./components/Clock.vue";
import fac from "./facility.json";

const store = useViewStore();
const {view, loggedIn} = storeToRefs(store);
const location = ref(window.location.href)

const changeView = () => {
  store.view = view.value;
};

store.getAuthed();

onMounted(() => {
  fac.views.forEach((v) => {
    v.facilities.forEach((f) => {
      if (store.sia[f] === undefined) {
        store.addAirport(f);
      }
    });
  })
});
</script>

<style scoped>
.title {
  font-size: 1.5rem;
  font-weight: 700;
  color: #ffb612;
}
.facility {
  color: v-bind("fac.facility.color");
}
.idscolor {
  color: v-bind("fac.facility.idscolor");
}
</style>
