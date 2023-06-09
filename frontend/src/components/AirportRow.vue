<template>
  <tr v-if="sia[props.airport] === undefined">
    <td rowspan="2" class="text-5xl w-[10%] border border-gray-500 text-center">
      {{ props.airport }}
    </td>
    <td rowspan="2" class="text-5xl text-center w-[5%] border border-gray-500">??</td>
    <td class="w-[10%] border border-gray-500"><span class="pr-2">ARR:</span> ??</td>
    <td class="w-[10%] border border-gray-500 text-center">??</td>
    <td rowspan="2" class="border border-gray-500 align-top">??</td>
  </tr>
  <tr v-if="sia[props.airport] === undefined">
    <td class="border border-gray-500"><span class="pr-2">DEP:</span> ??</td>
    <td class="border border-gray-500 text-center">??</td>
  </tr>
  <tr v-if="sia[props.airport] !== undefined">
    <td rowspan="2" class="text-5xl w-[10%] border border-gray-500 text-center">
      {{ props.airport }}
    </td>
    <td ref="atisbox" rowspan="2" class="text-5xl text-center w-[5%] border border-gray-500" @click="editATIS()" v-if="!isClosed">
      {{ sia[props.airport].atis != "" ? sia[props.airport].atis : "-" }}
    </td>
    <td ref="arrrwybox" class="w-[10%] border border-gray-500" v-if="!isClosed" @click="editArrRwy()">
      <span class="pr-2">ARR:</span>
      {{ sia[props.airport].arrival_runways != "" ? sia[props.airport].arrival_runways : "______" }}
    </td>
    <td rowspan="2" colspan="2" class="w-[10%] border border-gray-500 light-red text-center" v-if="isClosed">CLOSED</td>
    <td class="w-[10%] border border-gray-500 text-center text-yellow-400 font-bold">{{ wind }}</td>
    <td ref="metarbox" rowspan="2" class="border border-gray-500 align-top">{{ sia[props.airport].metar }}</td>
  </tr>
  <tr v-if="sia[props.airport] !== undefined">
    <td ref="deprwybox" class="border border-gray-500" v-if="!isClosed" @click="editDepRwy()">
      <span class="pr-2">DEP:</span>
      {{ sia[props.airport].departure_runways != "" ? sia[props.airport].departure_runways : "______" }}
    </td>
    <td class="border border-gray-500 text-center text-blue-400 font-bold">{{ altimeter }}</td>
  </tr>
  <tr>
    <td colspan="5" class="h-[1rem]"></td>
  </tr>

  <div v-show="showModal" class="absolute inset-0 flex items-center justify-center bg-neutral-700 bg-opacity-50">
    <div class="max-w-2xl p-6 mx-4 bg-neutral-800 rounded-md shadow-xl">
      <div class="flex items-center justify-between">
        <h3 class="text-2xl">Edit</h3>
        <svg
          @click="showModal = false"
          xmlns="http://www.w3.org/2000/svg"
          class="w-8 h-8 text-red-900 cursor-pointer"
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z"
          />
        </svg>
      </div>
      <div class="mt-4">
        <p class="mb-4 text-sm">
          <div class="flex items-center">
            <label class="block text-gray-100 font-bold pr-4 capitalize">{{ editing.replace("_", " ") }}:</label>
            <input ref="modaleditbox" type="text" class="w-full px-4 py-2 text-gray-100 bg-gray-600 rounded focus:bg-gray-600 focus:outline-none uppercase" v-model="modalText" @keyup.enter="saveModal" />
          </div>
        </p>
        <button @click="showModal = false" class="px-6 py-2 text-gray-100 border border-red-600 light-red rounded">Cancel</button>
        <button @click="saveModal" class="px-6 py-2 ml-2 text-blue-100 bg-blue-600 rounded">Save</button>
      </div>
    </div>
  </div>
</template>

<script setup>
// @TODO Refactor this... a lot....

import { DateTime } from "luxon";
import { storeToRefs } from "pinia";
import { useViewStore } from "@/store/viewstore";
import fac from "@/facility.json";
import { onBeforeUnmount, onMounted, computed, ref, watch } from "vue";
import { parseMetar } from "@/utils/metar.js";

const props = defineProps({
  airport: {
    type: String,
    required: true,
  },
});
const store = useViewStore();
const { sia } = storeToRefs(store);
const showModal = ref(false);
const modalText = ref("");
let editing = ref("");
const modaleditbox = ref(null);
const flashes = ref({});
const timers = {};
const metarbox = ref(null);
const atisbox = ref(null);
const deprwybox = ref(null);
const arrrwybox = ref(null);
const isClosed = ref(false);
let closeTimer = undefined;
let syncTimer = undefined;

onMounted(() => {
  const airport = fac.airports.filter((a) => a.name === props.airport)[0];
  if (airport === undefined) return false;

  if (airport.hours.continuous) return;
  isClosed.value = closed();
  syncTimer = setInterval(() => {
    // Check if we are 0 seconds into the next minute and then schedule
    // a check for the next minute
    if (DateTime.local().second === 0) {
      closeTimer = setInterval(() => {
        isClosed.value = closed();
      }, 60000);
      clearInterval(syncTimer);
    }
  }, 1000);
});

const saveModal = () => {
  store.patchSIA(props.airport, editing.value, modalText.value.toUpperCase());

  modalText.value = "";
  editing.value = "";
  showModal.value = false;
};

watch(() => store.sia[props.airport].metar, () => {
  flashes.metar = fac.facility.update_flash_duration;
  if (timers.metar) clearTimeout(timers.metar);
  timers.metar = setInterval(() => {
    flashField("metar");
  }, 1000);
});

watch(() => store.sia[props.airport].atis, () => {
  flashes.atis = fac.facility.update_flash_duration;
  if (timers.atis) clearTimeout(timers.atis);
  timers.atis = setInterval(() => {
    flashField("atis");
  }, 1000);
});

watch(() => store.sia[props.airport].departure_runways, () => {
  flashes.deprwy = fac.facility.update_flash_duration;
  if (timers.deprwy) clearTimeout(timers.deprwy);
  timers.deprwy = setInterval(() => {
    flashField("deprwy");
  }, 1000);
});

watch(() => store.sia[props.airport].arrival_runways, () => {
  flashes.arrrwy = fac.facility.update_flash_duration;
  if (timers.arrrwy) clearTimeout(timers.arrrwy);
  timers.arrrwy = setInterval(() => {
    flashField("arrrwy");
  }, 1000);
});

onBeforeUnmount(() => {
  clearInterval(timers.metar);
  clearInterval(timers.atis);
  clearInterval(timers.deprwy);
  clearInterval(timers.arrrwy);
  clearInterval(syncTimer);
  clearInterval(closeTimer);
});

const flashField = (field) => {
  if (flashes[field] > 0) {
    flashes[field]--;
    // Check if field has class bg-blue-500
    let f = null;
    if (field === "metar") {
      f = metarbox.value;
    } else if (field === "atis") {
      f = atisbox.value;
    } else if (field === "deprwy") {
      f = deprwybox.value;
    } else if (field === "arrrwy") {
      f = arrrwybox.value;
    }

    if (f.classList.contains("bg-blue-900")) {
      f.classList.remove("bg-blue-900");
    } else {
      f.classList.add("bg-blue-900");
    }
  } else {
    clearTimeout(timers[field]);
    timers[field] = undefined;
  }
}

const openModal = () => {
  showModal.value = true;
  setTimeout(() => {
    modaleditbox.value.select();
  }, 0);
}

const editATIS = () => {
  editing.value = "ATIS";
  modalText.value = sia.value[props.airport].atis;
  openModal();
};

const editArrRwy = () => {
  editing.value = "arrival_runways";
  modalText.value = sia.value[props.airport].arrival_runways;
  openModal();
};

const editDepRwy = () => {
  editing.value = "departure_runways";
  modalText.value = sia.value[props.airport].departure_runways;
  openModal();
};

const wind = computed(() => {
  if (
    sia.value[props.airport] === undefined ||
    sia.value[props.airport].metar === undefined ||
    sia.value[props.airport].metar === ""
  ) {
    return "??";
  }

  let m = parseMetar(sia.value[props.airport].metar);
  if (m.wind === undefined) {
    return "??";
  }

  if (m.wind.speed_kts < 3) {
    return "CALM";
  }

  let wind = `${calcWindDir(m.wind.degrees, sia.value[props.airport].mag_var).toString().padStart(3, "0")} @ ${m.wind.speed_kts}`;
  if (m.wind.gust_kts > m.wind.speed_kts + 6) {
    wind += `G${m.wind.gust_kts}`;
  }

  return wind;
});

const calcWindDir = (dir, magvar) => {
  // Round magvar to nearest 10
  if (magvar < 0) {
    magvar = Math.round((-1 * magvar) / 10) * -10;
  } else {
    magvar = Math.round(magvar / 10) * 10;
  }

  return dir + magvar;
}

const altimeter = computed(() => {
  if (
    sia.value[props.airport] === undefined ||
    sia.value[props.airport].metar === undefined ||
    sia.value[props.airport].metar === ""
  ) {
    return "??";
  }

  return parseMetar(sia.value[props.airport].metar).barometer.hg.toFixed(2) || "??";
});

const closed = () => {
  if (sia.value[props.airport] === undefined) return false;

  const airport = fac.airports.filter((a) => a.name === props.airport)[0];
  if (airport === undefined) return false;

  if (airport.hours.continuous) return false;

  for (let i = 0; i < airport.hours.schedule.length; i++) {
    const schedule = airport.hours.schedule[i];
    if (schedule.whenDST !== undefined) {
      if (inDST() && schedule.whenDST) {
        if (!betweenTimes(schedule.open, schedule.close, schedule.local, schedule.days)) {
          return true;
        }
      } else if (!inDST() && !schedule.whenDST) {
        if (!betweenTimes(schedule.open, schedule.close, schedule.local, schedule.days)) {
          return true;
        }
      }
    } else {
      // Get current month number and day of month in local time (fac.timezone.name)
      const currentDate = new Date();
      const month = currentDate.toLocaleString("en-US", { timeZone: fac.timezone.name, month: "numeric" });
      const day = currentDate.toLocaleString("en-US", { timeZone: fac.timezone.name, day: "numeric" });

      // Check if inbetween schedule.start.month and schedule.start.day and schedule.end.month and schedule.end.day
      if (
        month >= schedule.start.month &&
        day >= schedule.start.day &&
        month <= schedule.end.month &&
        day <= schedule.end.day
      ) {
        if (!betweenTimes(schedule.open, schedule.close, schedule.local, schedule.days)) {
          return true;
        }
      }
    }
  }

  return false;
};

const dowlist = ["Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"];

const betweenTimes = (start, end, local, days) => {
  if (days !== undefined) {
    // Get day of week in local time (fac.timezone.name)
    const day = new Date().toLocaleString("en-US", { timeZone: fac.timezone.name, weekday: "long" });
    if (!days.includes(dowlist.indexOf(day))) {
      return false;
    }
  }

  let now = DateTime.fromObject({}, { zone: fac.timezone.name });
  let dtstart = DateTime.fromObject({}, { zone: fac.timezone.name });
  let dtend = DateTime.fromObject({}, { zone: fac.timezone.name });
  if (local) {
    now = DateTime.fromObject({}, { zone: fac.timezone.name });
    dtstart = DateTime.fromObject({}, { zone: fac.timezone.name });
    dtend = DateTime.fromObject({}, { zone: fac.timezone.name });
  } else {
    now = DateTime.fromObject({}, { zone: "UTC" });
    dtstart = DateTime.fromObject({}, { zone: "UTC" });
    dtend = DateTime.fromObject({}, { zone: "UTC" });
  }

  const start_hr = parseInt(start.split(":")[0]);
  const start_min = parseInt(start.split(":")[1]);
  const end_hr = parseInt(end.split(":")[0]);
  const end_min = parseInt(end.split(":")[1]);

  dtstart = dtstart.set({ hour: start_hr, minute: start_min });
  dtend = dtend.set({ hour: end_hr, minute: end_min });

  return now >= dtstart && now < dtend;
};

const inDST = () => {
  let offset = getTimezoneOffset() / 60;
  return fac.timezone.dst === offset;
};

const getTimezoneOffset = () => {
  const now = new Date();
  const localizedTime = new Date(now.toLocaleString("en-US", { timeZone: fac.timezone.name }));
  const utcTime = new Date(now.toLocaleString("en-US", { timeZone: "UTC" }));
  return Math.round((localizedTime.getTime() - utcTime.getTime()) / (60 * 1000));
};
</script>

<style scoped>
.light-red {
  background: rgba(255, 0, 0, 0.1);
}
td {
  padding-left: 5px;
  padding-right: 5px;
}
</style>
