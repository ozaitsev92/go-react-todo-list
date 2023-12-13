
import axios from "axios";

let baseURL = "http://localhost:8080";

if (process.env.API_BASE_URL) {
    baseURL = process.env.API_BASE_URL;
}

const instance = axios.create({baseURL});

export default instance;
