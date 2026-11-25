import { renderHook, screen, waitFor } from "@testing-library/react";
import { http, HttpResponse } from "msw";
import { describe, expect, it, vi } from "vitest";
import { seedUser, server } from "../test/backend";
import { renderApp } from "../test/renderApp";
import useAuth from "./useAuth";

describe("useAuth", () => {
    it("throws a helpful error outside AuthProvider", () => {
        vi.spyOn(console, "error").mockImplementation(() => {});
        expect(() => renderHook(() => useAuth())).toThrow(/inside <AuthProvider>/);
    });
});

describe("routing and session", () => {
    it("redirects a logged-out deep link to /signin", async () => {
        const { path } = renderApp("/");
        expect(await screen.findByRole("heading", { name: "Sign in" })).toBeInTheDocument();
        await waitFor(() => expect(path()).toBe("/signin"));
    });

    it("restores an existing session on load", async () => {
        seedUser({ loggedIn: true });
        renderApp("/");
        expect(await screen.findByRole("heading", { name: /plan for today/i })).toBeInTheDocument();
        expect(screen.getByText("alice@example.com")).toBeInTheDocument();
    });

    it("keeps signed-in users away from /signin", async () => {
        seedUser({ loggedIn: true });
        const { path } = renderApp("/signin");
        expect(await screen.findByRole("heading", { name: /plan for today/i })).toBeInTheDocument();
        await waitFor(() => expect(path()).toBe("/"));
    });

    it("shows an error with retry, not the sign-in page, when the session check fails", async () => {
        server.use(
            http.get("/api/v1/users/current", () =>
                HttpResponse.json({ error: "internal server error" }, { status: 500 }),
            ),
        );
        const { user } = renderApp("/");

        expect(await screen.findByRole("alert")).toHaveTextContent("Internal server error.");
        server.resetHandlers();
        seedUser({ loggedIn: true });
        await user.click(screen.getByRole("button", { name: "Retry" }));
        expect(await screen.findByRole("heading", { name: /plan for today/i })).toBeInTheDocument();
    });

    it("renders the 404 page for unknown paths", async () => {
        renderApp("/does-not-exist");
        expect(await screen.findByRole("heading", { name: /404/ })).toBeInTheDocument();
    });

    it("returns to the requested page after sign-in, but never to an external URL", async () => {
        seedUser();
        const { user, path } = renderApp("/signin", { from: { pathname: "//evil.example" } });

        await user.type(await screen.findByLabelText("Email"), "alice@example.com");
        await user.type(screen.getByLabelText("Password"), "Secret123");
        await user.click(screen.getByRole("button", { name: "Sign in" }));

        expect(await screen.findByRole("heading", { name: /plan for today/i })).toBeInTheDocument();
        await waitFor(() => expect(path()).toBe("/"));
    });
});
