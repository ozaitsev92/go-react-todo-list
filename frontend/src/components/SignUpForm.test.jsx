import { screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";
import { seedUser } from "../test/backend";
import { renderApp } from "../test/renderApp";

async function signUp(user, email, password, confirm = password) {
    await user.type(await screen.findByLabelText("Email"), email);
    await user.type(screen.getByLabelText("Password"), password);
    await user.type(screen.getByLabelText("Confirm password"), confirm);
    await user.click(screen.getByRole("button", { name: "Sign up" }));
}

describe("SignUpForm", () => {
    it("creates an account with a plus-addressed email and a long TLD", async () => {
        const { user } = renderApp("/signup");
        await signUp(user, "a+tag@example.museum", "Secret123");
        expect(await screen.findByRole("heading", { name: "Account created" })).toBeInTheDocument();
        expect(screen.getByRole("link", { name: "sign in" })).toHaveAttribute("href", "/signin");
    });

    it("shows the backend message when the email is taken", async () => {
        seedUser({ email: "bob@example.com" });
        const { user } = renderApp("/signup");
        await signUp(user, "bob@example.com", "Secret123");
        expect(await screen.findByRole("alert")).toHaveTextContent("Email is already registered.");
    });

    it("rejects passwords over 72 bytes before calling the API", async () => {
        const { user } = renderApp("/signup");
        await signUp(user, "c@example.com", "a".repeat(73));
        expect(screen.getByLabelText("Password")).toHaveClass("is-invalid");
        expect(screen.queryByRole("heading", { name: "Account created" })).not.toBeInTheDocument();
    });

    it("marks a mismatched confirmation as invalid", async () => {
        const { user } = renderApp("/signup");
        await signUp(user, "d@example.com", "Secret123", "Secret124");
        expect(screen.getByLabelText("Confirm password")).toHaveClass("is-invalid");
        expect(screen.getByText("Passwords must match.")).toBeInTheDocument();
    });

    it("sets autocomplete hints for password managers", async () => {
        renderApp("/signup");
        expect(await screen.findByLabelText("Email")).toHaveAttribute("autocomplete", "username");
        expect(screen.getByLabelText("Password")).toHaveAttribute("autocomplete", "new-password");
    });
});
