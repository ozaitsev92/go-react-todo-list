import { describe, expect, it } from "vitest";
import { isValidEmail, isValidPassword, safeRedirect } from "./validation";

describe("validation", () => {
    it.each(["a@example.com", "a+tag@example.museum", "first.last@sub.example.technology"])(
        "accepts %s",
        (email) => expect(isValidEmail(email)).toBe(true),
    );

    it.each(["", "plain", "a@b", "Alice <a@example.com>", "a b@example.com"])(
        "rejects %j",
        (email) => expect(isValidEmail(email)).toBe(false),
    );

    it("measures password length in bytes like the backend (8..72)", () => {
        expect(isValidPassword("a".repeat(7))).toBe(false);
        expect(isValidPassword("a".repeat(8))).toBe(true);
        expect(isValidPassword("a".repeat(72))).toBe(true);
        expect(isValidPassword("a".repeat(73))).toBe(false);
        expect(isValidPassword("ж".repeat(37))).toBe(false); // 74 bytes
    });

    it.each([
        ["/", "/"],
        ["/tasks", "/tasks"],
        ["//evil.example", "/"],
        ["https://evil.example", "/"],
        [undefined, "/"],
    ])("safeRedirect(%j) = %j", (input, want) => expect(safeRedirect(input)).toBe(want));
});
