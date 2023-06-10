<template>
  <table class="table-fixed w-full text-white border-collapse">
    <AirportRow v-for="airport in viewFacility.facilities" :key="airport" :airport="airport" />
  </table>
</template>

<script setup>
import { ref, watch } from "vue";
import { storeToRefs } from "pinia";
import fac from "@/facility.json";
import AirportRow from "@/components/AirportRow.vue";
import { useViewStore } from "../store/viewstore";

const store = useViewStore();
const { view } = storeToRefs(store);

watch(view, () => {
  viewFacility.value = fac.views[view.value];
});
const viewFacility = ref(fac.views[view.value]);
</script>

<style lang="scss" scoped></style>
