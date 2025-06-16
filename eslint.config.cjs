const tsParser = require("@typescript-eslint/parser");
const globals = require("globals");
const typescriptEslint = require("@typescript-eslint/eslint-plugin");
const react = require("eslint-plugin-react");
const reactHooks = require("eslint-plugin-react-hooks");
const _import = require("eslint-plugin-import");
const eslintConfigPrettier = require("eslint-config-prettier");

const {
    fixupPluginRules,
} = require("@eslint/compat");

const js = require("@eslint/js");

module.exports = [
    // Global ignores
    {
        ignores: [
            "**/dist/",
            "**/build/",
            "**/node_modules/",
            "**/coverage/",
            "**/*.d.ts",
            "**/branding/",
            "**/out/",
            "**/.next/",
        ],
    },

    // Base config for all files
    js.configs.recommended,

    // TypeScript files
    {
        files: [ "**/*.ts", "**/*.tsx" ],
        languageOptions: {
            parser: tsParser,
            ecmaVersion: "latest",
            sourceType: "module",
            parserOptions: {
                projectService: true,
                tsconfigRootDir: __dirname,
                ecmaFeatures: {
                    jsx: true,
                },
            },
            globals: {
                ...globals.browser,
                ...globals.node,
                React: "readonly",
                JSX: "readonly",
            },
        },
        plugins: {
            "@typescript-eslint": typescriptEslint,
            react,
            "react-hooks": fixupPluginRules(reactHooks),
            import: fixupPluginRules(_import),
        },
        rules: {
            // TypeScript recommended rules
            ...typescriptEslint.configs.recommended.rules,

            // React rules
            ...react.configs.recommended.rules,
            "react/react-in-jsx-scope": "off",

            // React hooks rules
            ...reactHooks.configs.recommended.rules,

            // Import rules
            "import/default": "off",
            "import/no-unresolved": "error",
            "import/named": "error",
            "import/namespace": "error",
            "import/no-absolute-path": "error",
            "import/no-dynamic-require": "error",
            "import/no-self-import": "error",
            "import/no-cycle": "error",
            "import/no-useless-path-segments": "error",
        },
        settings: {
            react: {
                version: "detect",
            },
            "import/resolver": {
                typescript: {
                    project: [ "./frontend/tsconfig.json", "./backend/tsconfig.json" ],
                },
            },
        },
    },

    // Test files configuration
    {
        files: [ "**/*.test.ts", "**/*.test.tsx", "**/*.spec.ts", "**/*.spec.tsx" ],
        languageOptions: {
            globals: {
                ...globals.node,
                describe: "readonly",
                it: "readonly",
                expect: "readonly",
                beforeEach: "readonly",
                afterEach: "readonly",
                beforeAll: "readonly",
                afterAll: "readonly",
                jest: "readonly",
                test: "readonly",
            },
        },
    },

    // Node.js configuration files
    {
        files: [ "**/*.js", "**/*.cjs", "**/*.mjs" ],
        languageOptions: {
            globals: {
                ...globals.node,
                require: "readonly",
                module: "readonly",
                exports: "readonly",
                console: "readonly",
                process: "readonly",
                __dirname: "readonly",
                __filename: "readonly",
            },
        },
    },

    // Prettier config (must be last to override other configs)
    eslintConfigPrettier,
];
