import axios from "axios";
import fac from "../facility.json";

export const API = axios.create({
  baseURL: fac.facility.api,
  withCredentials: true,
});
