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
    <td rowspan="2" class="text-5xl w-[3em] border border-gray-500 text-center">
      {{ props.airport }}
    </td>
    <td
      v-if="!isClosed && hasDualATIS"
      ref="arratisbox"
      rowspan="2"
      class="text-5xl text-center w-[1em] border border-gray-500"
      @click.right.stop.prevent="checkHideArr()"
      @click.middle.stop.prevent="clearField('arrival_atis')"
      @click="editArrATIS()"
    >
      {{ sia[props.airport].arrival_atis != "" ? sia[props.airport].arrival_atis : "-" }}
    </td>
    <td
      v-if="!isClosed && hasDualATIS"
      ref="atisbox"
      rowspan="2"
      class="text-5xl text-center w-[1em] border border-gray-500"
      @click.middle.stop.prevent="clearField('atis')"
      @click="editATIS()"
    >
      {{ sia[props.airport].atis != "" ? sia[props.airport].atis : "-" }}
    </td>
    <td
      v-if="!isClosed && !hasDualATIS"
      ref="atisbox"
      colspan="2"
      rowspan="2"
      class="text-5xl text-center w-[2em] border border-gray-500"
      @click.right.stop.prevent="overrideArrival = true"
      @click.middle.stop.prevent="clearField('atis')"
      @click="editATIS()"
    >
      {{ sia[props.airport].atis != "" ? sia[props.airport].atis : "-" }}
    </td>
    <td
      v-if="!isClosed"
      ref="arrrwybox"
      class="border border-gray-500 w-[15em]"
      @click="editArrRwy()"
      @click.middle.stop.prevent="clearField('arrival_runways')"
    >
      <span class="pr-2">ARR:</span>
      {{ sia[props.airport].arrival_runways != "" ? sia[props.airport].arrival_runways : "______" }}
    </td>
    <td v-if="isClosed" rowspan="2" colspan="3" class="border border-gray-500 light-red text-center">CLOSED</td>
    <td class="w-[10em] border border-gray-500 text-center text-yellow-400 font-bold">{{ wind }}</td>
    <td ref="metarbox" rowspan="2" class="border border-gray-500 align-top">{{ sia[props.airport].metar }}</td>
  </tr>
  <tr v-if="sia[props.airport] !== undefined">
    <td
      v-if="!isClosed"
      ref="deprwybox"
      class="border border-gray-500"
      @click="editDepRwy()"
      @click.middle.stop.prevent="clearField('departure_runways')"
    >
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
          xmlns="http://www.w3.org/2000/svg"
          class="w-8 h-8 text-red-900 cursor-pointer"
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
          @click="showModal = false"
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
        <div class="mb-4 text-sm">
          <div class="flex items-center">
            <label class="block text-gray-100 font-bold pr-4 capitalize">{{ editing.replace("_", " ") }}:</label>
            <input
              ref="modaleditbox"
              v-model="modalText"
              type="text"
              class="w-full px-4 py-2 text-gray-100 bg-gray-600 rounded focus:bg-gray-600 focus:outline-none uppercase"
              @keyup.esc="showModal = false"
              @keyup.enter="saveModal"
            />
          </div>
        </div>
        <button class="px-6 py-2 text-gray-100 border border-red-600 light-red rounded" @click="showModal = false">
          Cancel
        </button>
        <button class="px-6 py-2 ml-2 text-blue-100 bg-blue-600 rounded" @click="saveModal">Save</button>
      </div>
    </div>
  </div>
</template>

<script setup>
// @TODO Refactor this... a lot....

import { DateTime } from "luxon";
import { storeToRefs } from "pinia";
import { computed, onBeforeUnmount, onMounted, ref, watch } from "vue";
import { useViewStore } from "@/store/viewstore";
import fac from "@/facility.json";
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
const editing = ref("");
const modaleditbox = ref(null);
const flashes = {};
const metarbox = ref(null);
const atisbox = ref(null);
const arratisbox = ref(null);
const deprwybox = ref(null);
const arrrwybox = ref(null);
const isClosed = ref(false);
const overrideArrival = ref(false);
let initialized = false;
const hasDualATIS = computed(
  () =>
    fac.airports.filter((a) => a.name === props.airport)[0]["dual-atis"] === true &&
    (overrideArrival.value || sia.value[props.airport].arrival_atis !== "")
);
const airport = fac.airports.filter((a) => a.name === props.airport)[0];

const clearField = (field) => {
  store.patchSIA(props.airport, field, "");
};

const saveModal = () => {
  store.patchSIA(props.airport, editing.value, modalText.value.toUpperCase());

  modalText.value = "";
  editing.value = "";
  showModal.value = false;
};

const fields = {
  metar: metarbox,
  arratis: arratisbox,
  atis: atisbox,
  deprwy: deprwybox,
  arrrwy: arrrwybox,
};

onMounted(() => {
  isClosed.value = closed();
})

watch(
  () => store.sia[props.airport].metar,
  () => {
    if (store.sia[props.airport].first) return;
    if (typeof metarbox.value === "undefined") return;
    flashes.metar = fac.facility.update_flash_duration;
  }
);

watch(
  () => store.sia[props.airport].atis,
  () => {
    if (store.sia[props.airport].first) return;
    if (typeof atisbox.value === "undefined") return;
    flashes.atis = fac.facility.update_flash_duration;
  }
);

watch(
  () => store.sia[props.airport].arrival_atis,
  () => {
    if (store.sia[props.airport].first) return;
    if (typeof arratisbox.value === "undefined") return;
    flashes.arratis = fac.facility.update_flash_duration;
  }
);

watch(
  () => store.sia[props.airport].departure_runways,
  () => {
    if (store.sia[props.airport].first) return;
    if (typeof deprwybox.value === "undefined") return;
    flashes.deprwy = fac.facility.update_flash_duration;
  }
);

watch(
  () => store.sia[props.airport].arrival_runways,
  () => {
    if (store.sia[props.airport].first) return;
    if (typeof arrrwybox.value === "undefined") return;
    flashes.arrrwy = fac.facility.update_flash_duration;
  }
);

watch(
  () => store.heartbeat.second,
  () => {
    Object.keys(flashes).forEach((field) => {
      if (flashes[field] > 0) {
        flashes[field]--;
        if (fields[field] === undefined) {
          flashes[field] = undefined;
          return;
        }
        const f = fields[field].value;
        if (f.classList.contains("bg-blue-900")) {
          f.classList.remove("bg-blue-900");
        } else {
          f.classList.add("bg-blue-900");
        }
      } else {
        fields[field].value.classList.remove("bg-blue-900");
        flashes[field] = undefined;
      }
    });
  }
);

watch(
  () => store.heartbeat.minute,
  () => {
    if (airport.hours.continuous) return;
    isClosed.value = closed();
  }
);

const checkHideArr = () => {
  overrideArrival.value = sia.value[props.airport].arrival_atis !== "" && overrideArrival.value;
};

const openModal = () => {
  showModal.value = true;
  setTimeout(() => {
    modaleditbox.value.select();
  }, 0);
};

const editATIS = () => {
  editing.value = "ATIS";
  modalText.value = sia.value[props.airport].atis;
  openModal();
};

const editArrATIS = () => {
  editing.value = "arrival_ATIS";
  modalText.value = sia.value[props.airport].arrival_atis;
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

  const m = parseMetar(sia.value[props.airport].metar);
  if (m.wind === undefined) {
    return "??";
  }

  if (m.wind.speed_kts < 3) {
    return "CALM";
  }

  let wind = `${calcWindDir(m.wind.degrees, sia.value[props.airport].mag_var).toString().padStart(3, "0")} @ ${
    m.wind.speed_kts
  }`;
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
  const new_dir = dir + magvar;
  if (new_dir < 0) {
    new_dir = 360 + new_dir;
  }
  if (magvar > 360) {
    new_dir = 360 - new_dir;
  }

  return dir + magvar;
};

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
  if (dtend < dtstart) {
    dtend = dtend.plus({ days: 1 });
  }

  return now >= dtstart && now < dtend;
};

const inDST = () => {
  const offset = getTimezoneOffset() / 60;
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
