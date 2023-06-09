<template>
  <span class="clock">
    {{ time }}
  </span>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'

const props = defineProps({
  timezone: {
    type: String,
    default: 'UTC',
  },
});

const time = ref("");
let timer = undefined;

const updateClock = () => {
  let letter = "L";
  if (props.timezone === 'UTC') letter="Z";
  time.value = new Date().toLocaleTimeString('en-UK', {
    timeZone: props.timezone,
    hour12: false,
  }).slice(-8) + " " + letter;
}

onMounted(() => {
  timer = setInterval(updateClock, 1000)
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
