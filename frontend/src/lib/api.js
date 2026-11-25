import axios from "axios";

// Same-origin /api: Vite, nginx or the ingress route it to the backend.
export const api = axios.create({
    baseURL: "/api",
    withCredentials: true,
    timeout: 10_000,
    headers: { Accept: "application/json" },
});

let unauthorizedHandler = null;

// setUnauthorizedHandler sets the callback for 401s on authenticated calls and
// returns an unregister function.
export function setUnauthorizedHandler(handler) {
    unauthorizedHandler = handler;
    return () => {
        if (unauthorizedHandler === handler) {
            unauthorizedHandler = null;
        }
    };
}

// Here 401 means bad credentials or no session yet, not an expired one.
const SESSION_ENDPOINTS = ["/v1/users/login", "/v1/users/current"];

api.interceptors.response.use(
    (response) => response,
    (error) => {
        const url = error.config?.url ?? "";
        if (error.response?.status === 401 && !SESSION_ENDPOINTS.includes(url)) {
            unauthorizedHandler?.();
        }
        return Promise.reject(error);
    },
);

export const isCancel = (error) => axios.isCancel(error);

// errorMessage returns user-facing text, preferring the backend's {"error"}.
export function errorMessage(error, fallback = "Something went wrong. Please try again.") {
    if (!error?.response) {
        return "Network error. Check your connection and try again.";
    }
    const message = error.response.data?.error;
    if (typeof message !== "string" || message === "") {
        return fallback;
    }
    return message.charAt(0).toUpperCase() + message.slice(1) + ".";
}

const data = (response) => response.data;

export const authApi = {
    register: (email, password) => api.post("/v1/users", { email, password }).then(data),
    login: (email, password) => api.post("/v1/users/login", { email, password }).then(data),
    logout: () => api.post("/v1/users/logout"),
    current: (signal) => api.get("/v1/users/current", { signal }).then(data),
};

const taskPath = (id) => `/v1/tasks/${encodeURIComponent(id)}`;

export const tasksApi = {
    list: (signal) => api.get("/v1/tasks", { signal }).then(data),
    create: (text) => api.post("/v1/tasks", { text }).then(data),
    update: (id, text) => api.put(taskPath(id), { text }).then(data),
    setCompleted: (id, completed) =>
        api
            .put(`${taskPath(id)}/${completed ? "mark-completed" : "mark-not-completed"}`)
            .then(data),
    remove: (id) => api.delete(taskPath(id)),
};
