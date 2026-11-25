import js from "@eslint/js";
import react from "eslint-plugin-react";
import reactHooks from "eslint-plugin-react-hooks";
import jsxA11y from "eslint-plugin-jsx-a11y";
import globals from "globals";

export default [
    { ignores: ["dist", "coverage"] },
    js.configs.recommended,
    react.configs.flat.recommended,
    react.configs.flat["jsx-runtime"],
    jsxA11y.flatConfigs.recommended,
    {
        files: ["**/*.{js,jsx}"],
        plugins: { "react-hooks": reactHooks },
        languageOptions: {
            ecmaVersion: "latest",
            sourceType: "module",
            globals: { ...globals.browser },
        },
        settings: { react: { version: "detect" } },
        rules: {
            ...reactHooks.configs.recommended.rules,
            "no-unused-vars": ["error", { argsIgnorePattern: "^_" }],
            eqeqeq: ["error", "always"],
            "no-console": "warn",
        },
    },
    {
        files: ["**/*.test.{js,jsx}", "src/test/**", "*.config.js"],
        languageOptions: { globals: { ...globals.node } },
    },
];
