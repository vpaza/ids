import { API } from "@/utils/api";
import fac from "@/facility.json";
import { defineStore } from "pinia";

export const useViewStore = defineStore("view", {
  state: () => ({
    view: fac.views.findIndex((v) => v.name === fac.facility.defaultView),
    sia: {},
    metars: {},
    fetching: [],
    timers: {},
    loggedIn: false,
  }),
  actions: {
    addAirport(airport) {
      if (this.timers[airport] === undefined) {
        this.updateSIA(airport);
        this.timers[airport] = setInterval(() => {
          this.updateSIA(airport);
        }, 15000);
      }
      if (this.sia[airport] !== undefined) return;
      this.sia[airport] = {
        "faa_id": "",
        "icao_id": "",
        "atis": "",
        "atis_time": new Date(),
        "departure_runways": "",
        "arrival_runways": "",
        "metar": "",
      };
    },
    async updateSIA(airport) {
      if (this.fetching.includes(airport)) return;
      this.fetching.push(airport);
      try {
        const response = await API.get("/v1/sia/" + airport);
        this.sia[airport] = response.data;
      } catch (error) {
        console.error(error);
        // Try again in 15 seconds
        setTimeout(() => {
          this.updateSIA(airport);
        }, 15000);
      }
    },
    async patchSIA(airport, field, value) {
      if (field === "ATIS") {
        field = "atis";
      }

      try {
        const data = {}; data[field] = value;
        await API.patch("/v1/sia/" + airport, data);
        this.sia[airport][field] = value;
      } catch (err) {
        console.error(err);
      }
    },
    async getAuthed() {
      try {
        await API.get("/v1/oauth/user");
        this.loggedIn = true;
        return true;
      } catch (err) {
        console.error(err);
        this.loggedIn = false;
        return false;
      }
    }
  }
});
