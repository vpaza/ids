<template>
  <div class="w-full">
    <ul class="flex mb-0 list-none flex-wrap pt-3 pb-4 flex-row">
      <li class="-mb-px mr-2 last:mr-0 flex-auto text-center">
        <a
          class="cursor-pointer text-xs font-bold uppercase px-5 py-3 shadow-lg rounded block leading-normal"
          :class="{
            'text-white bg-gray-900': openTab !== 1,
            'text-white bg-blue-900': openTab === 1,
          }"
          @click="toggleTabs(1)"
        >
          METAR
        </a>
      </li>
      <li class="-mb-px mr-2 last:mr-0 flex-auto text-center">
        <a
          class="cursor-pointer text-xs font-bold uppercase px-5 py-3 shadow-lg rounded block leading-normal"
          :class="{
            'text-white bg-gray-900': openTab !== 2,
            'text-white bg-blue-900': openTab === 2,
          }"
          @click="toggleTabs(2)"
        >
          PRE-DUTY BRIEF
        </a>
      </li>
      <li class="-mb-px mr-2 last:mr-0 flex-auto text-center">
        <a
          class="cursor-pointer text-xs font-bold uppercase px-5 py-3 shadow-lg rounded block leading-normal"
          :class="{
            'text-white bg-gray-900': openTab !== 3,
            'text-white bg-blue-900': openTab === 3,
          }"
          @click="toggleTabs(3)"
        >
          SATELITE
        </a>
      </li>
      <li class="-mb-px mr-2 last:mr-0 flex-auto text-center">
        <a
          class="cursor-pointer text-xs font-bold uppercase px-5 py-3 shadow-lg rounded block leading-normal"
          :class="{
            'text-white bg-gray-900': openTab !== 4,
            'text-white bg-blue-900': openTab === 4,
          }"
          @click="toggleTabs(4)"
        >
          REQ PIREP
        </a>
      </li>
      <li class="-mb-px mr-2 last:mr-0 flex-auto text-center">
        <a
          class="cursor-pointer text-xs font-bold uppercase px-5 py-3 shadow-lg rounded block leading-normal"
          :class="{
            'text-white bg-gray-900': openTab !== 5,
            'text-white bg-blue-900': openTab === 5,
          }"
          @click="toggleTabs(5)"
        >
          ICING+CONV
        </a>
      </li>
      <li class="-mb-px mr-2 last:mr-0 flex-auto text-center">
        <a
          class="cursor-pointer text-xs font-bold uppercase px-5 py-3 shadow-lg rounded block leading-normal"
          :class="{
            'text-white bg-gray-900': openTab !== 6,
            'text-white bg-blue-900': openTab === 6,
          }"
          @click="toggleTabs(6)"
        >
          PROG + SIG WX CHART
        </a>
      </li>
    </ul>
    <div class="relative flex flex-col min-w-0 break-words w-full mb-6 shadow-lg rounded">
      <div class="px-4 py-5 flex-auto">
        <div class="tab-content tab-space">
          <div :class="{ hidden: openTab !== 1, block: openTab === 1 }">
            <div class="flex w-full justify-left mb-2">
              <div class="flex items-center">
                <label class="block text-gray-100 font-bold pr-4"> Filter: </label>
                <input
                  v-model="filter"
                  type="text"
                  class="w-[10em] px-4 py-2 text-gray-100 bg-gray-600 rounded focus:bg-gray-600 focus:outline-none uppercase"
                />
              </div>
            </div>
            <table class="table-fixed w-full text-white border-collapse">
              <MetarRow
                v-for="metar in metarFiltered()"
                :key="metar"
                :mag_var="metars[metar].mag_var"
                :airport="metar"
                :metar="metars[metar].metar"
              ></MetarRow>
            </table>
          </div>
          <div :class="{ hidden: openTab !== 2, block: openTab === 2 }" class="align-center">
            <center>
              <video controls="" width="956" height="717">
                <source :src="fac.weather.preduty_briefing" type="video/mp4" />
              </video>
            </center>
          </div>
          <div :class="{ hidden: openTab !== 3, block: openTab === 3 }" class="text-white grid grid-cols-3 gap-4">
            <div>
              <p>IR View:</p>
              <img :src="fac.weather.satelite_ir" alt="Satelite IR" />
            </div>
            <div>
              <p>Visible View:</p>
              <img :src="fac.weather.satelite_vis" alt="Satelite VIS" />
            </div>
            <div>
              <p>Water Vapor View:</p>
              <img :src="fac.weather.satelite_water" alt="Satelite Water Vapor" />
            </div>
          </div>
          <div :class="{ hidden: openTab !== 4, block: openTab === 4 }">
            <center>
              <img :src="fac.weather.required_pirep" alt="Required PIREP" />
            </center>
          </div>
          <div :class="{ hidden: openTab !== 5, block: openTab === 5 }" class="grid grid-cols-4 gap-4">
            <div>
              <img :src="fac.weather.icing_1" alt="Icing" />
            </div>
            <div>
              <img :src="fac.weather.icing_2" alt="Icing" />
            </div>
            <div>
              <img :src="fac.weather.icing_3" alt="Icing" />
            </div>
            <div>
              <img :src="fac.weather.icing_4" alt="Icing" />
            </div>
            <div>
              <img :src="fac.weather.convective_1" alt="Convective" />
            </div>
            <div>
              <img :src="fac.weather.convective_2" alt="Convective" />
            </div>
            <div>
              <img :src="fac.weather.convective_3" alt="Convective" />
            </div>
            <div>
              <img :src="fac.weather.convective_4" alt="Convective" />
            </div>
          </div>
          <div :class="{ hidden: openTab !== 6, block: openTab === 6 }" class="grid grid-cols-2 gap-4">
            <div>
              <img :src="fac.weather.prog_1" alt="Prog" />
            </div>
            <div>
              <img :src="fac.weather.prog_2" alt="Prog" />
            </div>
            <div>
              <img :src="fac.weather.sigwx_24" alt="Sigwx" />
            </div>
            <div>
              <img :src="fac.weather.sigwx_36" alt="Sigwx" />
            </div>
            <div>
              <img :src="fac.weather.sigwx_48" alt="Sigwx" />
            </div>
            <div>
              <img :src="fac.weather.sigwx_60" alt="Sigwx" />
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, ref } from "vue";
import { storeToRefs } from "pinia";
import { stringify } from "postcss";
import { useViewStore } from "../store/viewstore";
import fac from "@/facility.json";
import MetarRow from "@/components/MetarRow.vue";

const store = useViewStore();
const { metars } = storeToRefs(store);
const openTab = ref(1);
const filter = ref("");

const metarFiltered = () => {
  // Return keys of metars filtered by filter
  return Object.keys(metars.value)
    .filter((key) => {
      return key.toUpperCase().includes(filter.value.toUpperCase());
    })
    .sort();
};

const toggleTabs = (tab) => {
  openTab.value = tab;
};
</script>

<style lang="scss" scoped></style>
