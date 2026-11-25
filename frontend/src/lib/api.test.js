import { http, HttpResponse } from "msw";
import { afterEach, describe, expect, it, vi } from "vitest";
import { server } from "../test/backend";
import { api, errorMessage, setUnauthorizedHandler, tasksApi } from "./api";

describe("api", () => {
    let unregister = () => {};
    afterEach(() => unregister());

    it("uses the same-origin /api base URL", () => {
        expect(api.defaults.baseURL).toBe("/api");
    });

    it("calls the unauthorized handler on 401 from an authenticated endpoint", async () => {
        const onUnauthorized = vi.fn();
        unregister = setUnauthorizedHandler(onUnauthorized);

        await expect(tasksApi.list()).rejects.toMatchObject({ response: { status: 401 } });
        expect(onUnauthorized).toHaveBeenCalledOnce();
    });

    it.each(["/v1/users/login", "/v1/users/current"])(
        "does not treat 401 from %s as an expired session",
        async (path) => {
            const onUnauthorized = vi.fn();
            unregister = setUnauthorizedHandler(onUnauthorized);

            await expect(
                api.request({
                    url: path,
                    method: path.endsWith("login") ? "post" : "get",
                    data: {},
                }),
            ).rejects.toMatchObject({ response: { status: 401 } });
            expect(onUnauthorized).not.toHaveBeenCalled();
        },
    );

    it("surfaces the backend error message", async () => {
        server.use(
            http.get("/api/v1/tasks", () =>
                HttpResponse.json({ error: "task not found" }, { status: 404 }),
            ),
        );
        const err = await tasksApi.list().catch((e) => e);
        expect(errorMessage(err)).toBe("Task not found.");
    });

    it("reports network failures distinctly", async () => {
        server.use(http.get("/api/v1/tasks", () => HttpResponse.error()));
        const err = await tasksApi.list().catch((e) => e);
        expect(errorMessage(err)).toMatch(/network error/i);
    });

    it("falls back to a generic message when the body has no error", async () => {
        server.use(http.get("/api/v1/tasks", () => new HttpResponse("oops", { status: 500 })));
        const err = await tasksApi.list().catch((e) => e);
        expect(errorMessage(err)).toBe("Something went wrong. Please try again.");
    });
});
