import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";

const backend = process.env.API_UPSTREAM ?? "http://localhost:8080";

export default defineConfig({
    plugins: [react()],
    server: {
        port: 8081,
        proxy: {
            "/api": {
                target: backend,
                rewrite: (path) => path.replace(/^\/api/, ""),
            },
        },
    },
    test: {
        environment: "jsdom",
        setupFiles: ["./src/test/setup.js"],
        restoreMocks: true,
        coverage: {
            include: ["src/**/*.{js,jsx}"],
            exclude: ["src/test/**", "src/main.jsx", "**/*.test.{js,jsx}"],
        },
    },
});
