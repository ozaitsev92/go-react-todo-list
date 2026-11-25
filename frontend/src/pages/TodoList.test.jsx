import { screen, waitFor, within } from "@testing-library/react";
import { http, HttpResponse } from "msw";
import { describe, expect, it } from "vitest";
import { expireSession, seedTask, seedUser, server } from "../test/backend";
import { renderApp } from "../test/renderApp";

async function renderList() {
    seedUser({ loggedIn: true });
    const utils = renderApp("/");
    await screen.findByRole("heading", { name: /plan for today/i });
    return utils;
}

const taskList = () => screen.findByRole("list", { name: "Tasks" });

describe("TodoList page", () => {
    it("renders the tasks from the API, oldest first", async () => {
        seedUser({ loggedIn: true });
        seedTask("first");
        seedTask("second", true);
        renderApp("/");

        const items = within(await taskList()).getAllByRole("listitem");
        expect(items.map((li) => li.textContent)).toEqual(["first", "second"]);
        expect(screen.getByRole("checkbox", { name: "second" })).toBeChecked();
    });

    it("shows an empty state", async () => {
        await renderList();
        expect(await screen.findByText(/no tasks yet/i)).toBeInTheDocument();
    });

    it("shows an error and retries when loading fails", async () => {
        seedUser({ loggedIn: true });
        server.use(http.get("/api/v1/tasks", () => HttpResponse.error(), { once: true }));
        seedTask("after retry");
        const { user } = renderApp("/");

        expect(await screen.findByRole("alert")).toHaveTextContent(/network error/i);
        await user.click(screen.getByRole("button", { name: "Try again" }));
        expect(await screen.findByRole("checkbox", { name: "after retry" })).toBeInTheDocument();
    });

    it("adds a task and clears the input", async () => {
        const { user } = await renderList();
        const input = screen.getByLabelText("New task");

        await user.type(input, "  buy milk  {Enter}");

        expect(await screen.findByRole("checkbox", { name: "buy milk" })).toBeInTheDocument();
        expect(input).toHaveValue("");
    });

    it("keeps the typed text and shows the error when adding fails", async () => {
        const { user } = await renderList();
        server.use(http.post("/api/v1/tasks", () => HttpResponse.error()));
        const input = screen.getByLabelText("New task");

        await user.type(input, "lost?{Enter}");

        expect(await screen.findByRole("alert")).toHaveTextContent(/network error/i);
        expect(input).toHaveValue("lost?");
    });

    it("toggles completion from the keyboard", async () => {
        seedUser({ loggedIn: true });
        seedTask("walk dog");
        const { user } = renderApp("/");
        const checkbox = await screen.findByRole("checkbox", { name: "walk dog" });

        checkbox.focus();
        await user.keyboard(" ");

        await waitFor(() => expect(checkbox).toBeChecked());
        expect(screen.getByText("walk dog")).toHaveClass("completed");
    });

    it("edits a task with Enter and returns focus to the Edit button", async () => {
        seedUser({ loggedIn: true });
        seedTask("old text");
        const { user } = renderApp("/");

        await user.click(await screen.findByRole("button", { name: 'Edit "old text"' }));
        const field = screen.getByLabelText("Task text");
        expect(field).toHaveFocus();
        await user.clear(field);
        await user.type(field, "new text{Enter}");

        expect(await screen.findByRole("checkbox", { name: "new text" })).toBeInTheDocument();
        expect(screen.getByRole("button", { name: 'Edit "new text"' })).toHaveFocus();
    });

    it("cancels editing with Escape without saving", async () => {
        seedUser({ loggedIn: true });
        seedTask("keep me");
        const { user } = renderApp("/");

        await user.click(await screen.findByRole("button", { name: 'Edit "keep me"' }));
        await user.type(screen.getByLabelText("Task text"), " changed{Escape}");

        expect(await screen.findByRole("checkbox", { name: "keep me" })).toBeInTheDocument();
    });

    it("deletes a task", async () => {
        seedUser({ loggedIn: true });
        seedTask("remove me");
        const { user } = renderApp("/");

        await user.click(await screen.findByRole("button", { name: 'Delete "remove me"' }));

        await waitFor(() => expect(screen.queryByText("remove me")).not.toBeInTheDocument());
        expect(screen.getByText(/no tasks yet/i)).toBeInTheDocument();
    });

    it("sends the user to /signin when the session expires mid-use", async () => {
        const { user, path } = await renderList();
        expireSession();

        await user.type(screen.getByLabelText("New task"), "too late{Enter}");

        expect(await screen.findByRole("heading", { name: "Sign in" })).toBeInTheDocument();
        await waitFor(() => expect(path()).toBe("/signin"));
    });

    it("logs out and no longer shows the list", async () => {
        const { user, path } = await renderList();

        await user.click(screen.getByRole("button", { name: "Log out" }));

        expect(await screen.findByRole("heading", { name: "Sign in" })).toBeInTheDocument();
        await waitFor(() => expect(path()).toBe("/signin"));
        expect(screen.queryByText("alice@example.com")).not.toBeInTheDocument();
    });
});
