// Mirrors the backend's rules; the backend remains the source of truth.

export const MIN_PASSWORD_BYTES = 8;
export const MAX_PASSWORD_BYTES = 72;
export const MAX_TASK_LENGTH = 1000;

const EMAIL_PATTERN = /^[^\s@<>]+@[^\s@<>]+\.[^\s@<>]+$/;

export const isValidEmail = (email) => EMAIL_PATTERN.test(email.trim());

export const byteLength = (value) => new TextEncoder().encode(value).length;

export const isValidPassword = (password) => {
    const n = byteLength(password);
    return n >= MIN_PASSWORD_BYTES && n <= MAX_PASSWORD_BYTES;
};

// safeRedirect only allows same-app paths, never "//host" or absolute URLs.
export const safeRedirect = (path) =>
    typeof path === "string" && path.startsWith("/") && !path.startsWith("//") ? path : "/";
