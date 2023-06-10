<template>
  <span class="clock">
    {{ time }}
  </span>
</template>

<script setup>
import { onMounted, onUnmounted, ref } from "vue";
import { useViewStore } from "@/store/viewstore";

const store = useViewStore();

const props = defineProps({
  timezone: {
    type: String,
    default: "UTC",
  },
});

const time = ref("");
let timer;

const updateClock = () => {
  let letter = "L";
  if (props.timezone === "UTC") letter = "Z";
  time.value = `${new Date()
    .toLocaleTimeString("en-UK", {
      timeZone: props.timezone,
      hour12: false,
    })
    .slice(-8)} ${letter}`;

  // Check if seconds are 0, if so, update store.heartbeat.minute...
  // the value doesn't matter.. we just need to change
  if (time.value.slice(-2) === "00") {
    store.heartbeat.minute = time.value;
  }
  store.heartbeat.second = time.value;
};

onMounted(() => {
  timer = setInterval(updateClock, 1000);
  updateClock();
});

onUnmounted(() => {
  clearInterval(timer);
});
</script>

<style scoped>
.clock {
  line-height: 1.5;
  padding: 0.5rem;
  font-size: 1.5rem;
  color: #ffb612;
  background-color: #222;
  border-radius: 15px;
}
</style>
