import { screen } from "@testing-library/react";
import { delay, http, HttpResponse } from "msw";
import { describe, expect, it } from "vitest";
import { seedUser, server } from "../test/backend";
import { renderApp } from "../test/renderApp";

async function fillAndSubmit(user, email, password) {
    await user.type(await screen.findByLabelText("Email"), email);
    await user.type(screen.getByLabelText("Password"), password);
    await user.click(screen.getByRole("button", { name: "Sign in" }));
}

describe("SignInForm", () => {
    it("signs in and shows the todo list", async () => {
        seedUser();
        const { user } = renderApp("/signin");
        await fillAndSubmit(user, "alice@example.com", "Secret123");
        expect(await screen.findByRole("heading", { name: /plan for today/i })).toBeInTheDocument();
    });

    it("shows the backend message for wrong credentials", async () => {
        seedUser();
        const { user } = renderApp("/signin");
        await fillAndSubmit(user, "alice@example.com", "wrong-password");
        expect(await screen.findByRole("alert")).toHaveTextContent("Invalid email or password.");
    });

    it("shows a network error when the request fails", async () => {
        server.use(http.post("/api/v1/users/login", () => HttpResponse.error()));
        const { user } = renderApp("/signin");
        await fillAndSubmit(user, "alice@example.com", "Secret123");
        expect(await screen.findByRole("alert")).toHaveTextContent(/network error/i);
    });

    it("disables the button while the request is pending", async () => {
        seedUser();
        server.use(
            http.post("/api/v1/users/login", async () => {
                await delay(50);
                return HttpResponse.json({ error: "invalid email or password" }, { status: 401 });
            }),
        );
        const { user } = renderApp("/signin");
        await fillAndSubmit(user, "alice@example.com", "Secret123");
        expect(screen.getByRole("button", { name: "Signing in…" })).toBeDisabled();
        expect(await screen.findByRole("alert")).toBeInTheDocument();
    });

    it("does not send the request when fields are empty", async () => {
        const { user } = renderApp("/signin");
        await user.click(await screen.findByRole("button", { name: "Sign in" }));
        expect(screen.getByRole("form", { name: "Sign in" })).toHaveClass("was-validated");
        expect(screen.queryByRole("alert")).not.toBeInTheDocument();
    });

    it("does not persist anything in localStorage", async () => {
        seedUser();
        const { user } = renderApp("/signin");
        await fillAndSubmit(user, "alice@example.com", "Secret123");
        await screen.findByRole("heading", { name: /plan for today/i });
        expect(window.localStorage.length).toBe(0);
    });
});
