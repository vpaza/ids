<template>
  <div class="flex flex-col min-h-screen">
    <header class="fixed bg-neutral-700 text-white top-0 p-4 w-full shadow-md mb-[10px] h-[83px]">
      <div class="py-1">
        <div class="grid grid-cols-4 items-center justify-between">
          <div class="col-span-1 title">
            <span class="facility">{{ fac.facility.id }}</span> <span class="idscolor">IDS</span>
          </div>
          <div class="col-span-2">
            <div v-if="loggedIn" class="flex items-center">
              <label class="block text-gray-100 font-bold pr-4"> Views: </label>
              <select
                id="view"
                v-model="view"
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
    <main class="pt-[95px] flex-1 overflow-y-auto p-5 pb-[44px] flex flex-col">
      <router-view v-if="loggedIn" />
    </main>
    <footer class="fixed z-50 bg-neutral-700 text-white bottom-0 p-0 w-full">
      <div v-if="loggedIn">
        <button
          class="border-2 border-gray-500 bg-slate-800 hover:bg-gray-700 text-white font-bold py-2 px-4 rounded h-full w-[10rem] mr-1"
          @click="router.push('/')"
        >
          SIA
        </button>
        <button
          class="border-2 border-gray-500 bg-blue-800 hover:bg-gray-700 text-white font-bold py-2 px-4 rounded h-full w-[10rem] mr-1"
          @click="router.push('/weather')"
        >
          WX
        </button>
        <button
          class="border-2 border-gray-500 bg-green-800 hover:bg-gray-700 text-white font-bold py-2 px-4 rounded h-full w-[10rem] mr-1"
          @click="router.push('/sops')"
        >
          SOP
        </button>
        <button
          class="border-2 border-gray-500 bg-yellow-800 hover:bg-gray-700 text-white font-bold py-2 px-4 rounded h-full w-[10rem] mr-1"
          @click="router.push('/briefing')"
        >
          BRIEF
        </button>
      </div>
      <div v-else>
        <a :href="`${fac.facility.api}/v1/oauth/login?redirect=${location}`">Login</a>
      </div>
    </footer>
  </div>
</template>

<script setup>
import { onMounted, ref } from "vue";
import { useRouter } from "vue-router";
import { storeToRefs } from "pinia";
import { useViewStore } from "@/store/viewstore";
import Clock from "./components/Clock.vue";
import fac from "./facility.json";

const store = useViewStore();
const router = useRouter();
const { view, loggedIn } = storeToRefs(store);
const location = ref(window.location.href);

const changeView = () => {
  store.view = view.value;
};

store.getAuthed();

onMounted(() => {
  fac.views.forEach((v) => {
    v.facilities.forEach(async (f) => {
      if (store.sia[f] === undefined) {
        await store.addAirport(f);
      }
    });
  });

  store.updateMetars();
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
