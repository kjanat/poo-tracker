const {
    defineConfig,
    globalIgnores,
} = require("eslint/config");

const tsParser = require("@typescript-eslint/parser");
const globals = require("globals");
const typescriptEslint = require("@typescript-eslint/eslint-plugin");
const react = require("eslint-plugin-react");
const reactHooks = require("eslint-plugin-react-hooks");
const _import = require("eslint-plugin-import");

const {
    fixupPluginRules,
    fixupConfigRules,
} = require("@eslint/compat");

const js = require("@eslint/js");

const {
    FlatCompat,
} = require("@eslint/eslintrc");

const compat = new FlatCompat({
    baseDirectory: __dirname,
    recommendedConfig: js.configs.recommended,
    allConfig: js.configs.all
});

module.exports = defineConfig([{
    languageOptions: {
        parser: tsParser,
        ecmaVersion: "latest",
        sourceType: "module",

        parserOptions: {
            project: ["./frontend/tsconfig.json", "./backend/tsconfig.json"],
            tsconfigRootDir: __dirname,
        },

        globals: {
            ...globals.browser,
            ...globals.node,
        },
    },

    plugins: {
        "@typescript-eslint": typescriptEslint,
        react,
        "react-hooks": fixupPluginRules(reactHooks),
        import: fixupPluginRules(_import),
    },

    extends: fixupConfigRules(compat.extends(
        "eslint:recommended",
        "plugin:@typescript-eslint/recommended",
        "plugin:react/recommended",
        "plugin:react-hooks/recommended",
        "plugin:import/recommended",
        "plugin:import/typescript",
        "prettier",
    )),

    settings: {
        react: {
            version: "detect",
        },

        "import/resolver": {
            typescript: {},
        },
    },

    rules: {
        "react/react-in-jsx-scope": "off",
        "import/default": "off",
    },
}, globalIgnores([
    "**/dist/",
    "**/build/",
    "**/node_modules/",
    "**/coverage/",
    "**/*.d.ts",
    "**/branding/",
]), globalIgnores([
    "**/node_modules",
    "**/branding",
    "**/dist",
    "**/build",
    "**/coverage",
    "**/out",
    "**/.next",
    "**/*.js",
])]);
