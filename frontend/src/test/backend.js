import { http, HttpResponse } from "msw";
import { setupServer } from "msw/node";

// In-memory fake of the Go API contract (cookie session, UUIDs, status codes,
// {"error"} bodies), exercised through the real axios instance.

const state = { users: new Map(), session: null, tasks: [] };

const fail = (status, error) => HttpResponse.json({ error }, { status });
const publicUser = ({ id, email }) => ({ id, email });

function requireSession() {
    return state.session ? null : fail(401, "unauthorized");
}

export const handlers = [
    http.post("/api/v1/users", async ({ request }) => {
        const { email, password } = await request.json();
        const normalized = String(email).trim().toLowerCase();
        if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(normalized)) {
            return fail(422, "email is invalid");
        }
        const bytes = new TextEncoder().encode(password).length;
        if (bytes < 8 || bytes > 72) {
            return fail(422, "password must be between 8 and 72 bytes");
        }
        if (state.users.has(normalized)) {
            return fail(409, "email is already registered");
        }
        const user = { id: crypto.randomUUID(), email: normalized, password };
        state.users.set(normalized, user);
        return HttpResponse.json(publicUser(user), { status: 201 });
    }),

    http.post("/api/v1/users/login", async ({ request }) => {
        const { email, password } = await request.json();
        const user = state.users.get(String(email).trim().toLowerCase());
        if (!user || user.password !== password) {
            return fail(401, "invalid email or password");
        }
        state.session = user;
        return HttpResponse.json(publicUser(user));
    }),

    http.post("/api/v1/users/logout", () => {
        state.session = null;
        return new HttpResponse(null, { status: 204 });
    }),

    http.get(
        "/api/v1/users/current",
        () => requireSession() ?? HttpResponse.json(publicUser(state.session)),
    ),

    http.get("/api/v1/tasks", () => requireSession() ?? HttpResponse.json(state.tasks)),

    http.post("/api/v1/tasks", async ({ request }) => {
        const denied = requireSession();
        if (denied) {
            return denied;
        }
        const text = String((await request.json()).text ?? "").trim();
        if (text === "" || text.length > 1000) {
            return fail(422, "text must be 1 to 1000 characters");
        }
        const now = new Date().toISOString();
        const task = {
            id: crypto.randomUUID(),
            text,
            completed: false,
            user_id: state.session.id,
            created_at: now,
            updated_at: now,
        };
        state.tasks.push(task);
        return HttpResponse.json(task, { status: 201 });
    }),

    http.put(
        "/api/v1/tasks/:id",
        async ({ params, request }) =>
            requireSession() ??
            withTask(params.id, async (task) => {
                task.text = String((await request.json()).text).trim();
                return HttpResponse.json(task);
            }),
    ),

    http.put(
        "/api/v1/tasks/:id/:action",
        ({ params }) =>
            requireSession() ??
            withTask(params.id, (task) => {
                task.completed = params.action === "mark-completed";
                return HttpResponse.json(task);
            }),
    ),

    http.delete(
        "/api/v1/tasks/:id",
        ({ params }) =>
            requireSession() ??
            withTask(params.id, () => {
                state.tasks = state.tasks.filter((t) => t.id !== params.id);
                return new HttpResponse(null, { status: 204 });
            }),
    ),
];

function withTask(id, fn) {
    const task = state.tasks.find((t) => t.id === id);
    return task ? fn(task) : fail(404, "task not found");
}

export const server = setupServer(...handlers);

export function resetBackend() {
    state.users.clear();
    state.session = null;
    state.tasks = [];
}

// seedUser creates an account and, if loggedIn, an active session for it.
export function seedUser({
    email = "alice@example.com",
    password = "Secret123",
    loggedIn = false,
} = {}) {
    const user = { id: crypto.randomUUID(), email, password };
    state.users.set(email, user);
    if (loggedIn) {
        state.session = user;
    }
    return user;
}

export function seedTask(text, completed = false) {
    const now = new Date().toISOString();
    const task = {
        id: crypto.randomUUID(),
        text,
        completed,
        user_id: state.session?.id ?? "",
        created_at: now,
        updated_at: now,
    };
    state.tasks.push(task);
    return task;
}

// expireSession simulates the cookie expiring server-side.
export function expireSession() {
    state.session = null;
}
