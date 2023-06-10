import { defineStore } from "pinia";
import { API } from "@/utils/api";
import fac from "@/facility.json";

export const useViewStore = defineStore("view", {
  state: () => ({
    view: fac.views.findIndex((v) => v.name === fac.facility.defaultView),
    sia: {},
    metars: {},
    fetching: [],
    metarFetching: false,
    timers: {},
    loggedIn: false,
    initialized: false,
    heartbeat: {
      second: 0,
      minute: 0,
    },
    metarTimer: null,
  }),
  actions: {
    async addAirport(airport) {
      if (this.sia[airport] !== undefined) return;
      this.sia[airport] = {
        faa_id: "",
        icao_id: "",
        atis: "",
        arrival_atis: "",
        atis_time: new Date(),
        arrival_atis_time: new Date(),
        departure_runways: "",
        arrival_runways: "",
        metar: "",
        first: true,
      };
      if (this.timers[airport] === undefined) {
        await this.updateSIA(airport);
        this.sia[airport].first = true;
        setTimeout(() => {
          this.sia[airport].first = false;
        }, 1000);
        this.timers[airport] = setInterval(() => {
          this.updateSIA(airport);
        }, 15000);
      }
    },
    async updateSIA(airport) {
      if (this.fetching.includes(airport)) return;
      this.fetching.push(airport);
      try {
        const response = await API.get(`/v1/sia/${airport}`);
        this.sia[airport] = response.data;
      } catch (error) {
        console.error(error);
        // Try again in 15 seconds
        setTimeout(() => {
          this.updateSIA(airport);
        }, 15000);
      } finally {
        this.fetching.splice(this.fetching.indexOf(airport), 1);
      }
    },
    async updateMetars() {
      try {
        const response = await API.get("/v1/weather/metar/all");
        Object.keys(response.data).forEach((key) => {
          if (response.data[key] !== "") {
            this.metars[key] = response.data[key];
          }
        });
      } catch (error) {
        console.error(error);
      }
      if (this.metarTimer === null) {
        this.metarTimer = setInterval(() => {
          this.updateMetars();
        }, 60000);
      }
    },
    async patchSIA(airport, field, value) {
      if (field === "ATIS") {
        field = "atis";
      }
      if (field === "arrival_ATIS") {
        field = "arrival_atis";
      }

      try {
        const data = {};
        data[field] = value;
        await API.patch(`/v1/sia/${airport}`, data);
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
    },
  },
});
