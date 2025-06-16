const js = require("@eslint/js");
const tseslint = require("@typescript-eslint/eslint-plugin");
const tsparser = require("@typescript-eslint/parser");
const react = require("eslint-plugin-react");
const reactHooks = require("eslint-plugin-react-hooks");
const { resolve } = require("path");
const globals = require("globals");

module.exports = [
  js.configs.recommended,
  {
    files: ["**/*.ts", "**/*.tsx"],
    languageOptions: {
      parser: tsparser,
      parserOptions: { sourceType: "module" }
    },
    plugins: {
      "@typescript-eslint": tseslint,
      react,
      "react-hooks": reactHooks
    },
    settings: { react: { version: "detect" } },
    rules: {
      "no-unused-vars": "off",
      "@typescript-eslint/no-unused-vars": "warn"
    }
  },
  {
    files: ["backend/**/*.ts"],
    languageOptions: {
      parserOptions: { tsconfigRootDir: resolve(__dirname, "backend") },
      globals: globals.node
    }
  },
  {
    files: ["frontend/**/*.{ts,tsx}"],
    languageOptions: {
      parserOptions: { tsconfigRootDir: resolve(__dirname, "frontend") },
      globals: {
        ...globals.browser,
        ...globals.node,
        React: false,
        jest: false,
        __dirname: false
      }
    }
  },
  {
    files: ["backend/src/**/*.test.ts", "backend/src/**/__tests__/*.ts"],
    languageOptions: { globals: globals.jest }
  },
  {
    files: ["frontend/**/*.test.tsx", "frontend/**/*.test.ts"],
    languageOptions: { globals: globals.jest }
  },
  {
    ignores: ["**/dist/**", "**/*.js"]
  }
];
