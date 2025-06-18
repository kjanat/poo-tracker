const tsParser = require('@typescript-eslint/parser')
const globals = require('globals')
const typescriptEslint = require('@typescript-eslint/eslint-plugin')
const react = require('eslint-plugin-react')
const reactHooks = require('eslint-plugin-react-hooks')
const _import = require('eslint-plugin-import')
const eslintConfigPrettier = require('eslint-config-prettier')

const { fixupPluginRules } = require('@eslint/compat')

const js = require('@eslint/js')

module.exports = [
  // Global ignores
  {
    ignores: [
      '**/dist/',
      '**/build/',
      '**/node_modules/',
      '**/coverage/',
      '**/*.d.ts',
      '**/branding/',
      '**/out/',
      '**/.next/',
      'ai-service/'
    ]
  },

  // Base config for all files
  js.configs.recommended,

  // TypeScript files
  {
    files: ['**/*.ts', '**/*.tsx'],
    languageOptions: {
      parser: tsParser,
      ecmaVersion: 2024, // Updated for ES2024
      sourceType: 'module',
      parserOptions: {
        projectService: true,
        tsconfigRootDir: __dirname,
        ecmaFeatures: {
          jsx: true
        }
      },
      globals: {
        ...globals.browser,
        ...globals.node,
        ...globals.es2024, // Add ES2024 globals
        React: 'readonly',
        JSX: 'readonly'
      }
    },
    plugins: {
      '@typescript-eslint': typescriptEslint,
      react,
      'react-hooks': fixupPluginRules(reactHooks),
      import: fixupPluginRules(_import)
    },
    rules: {
      // TypeScript recommended rules
      ...typescriptEslint.configs.recommended.rules,

      // React rules
      ...react.configs.recommended.rules,
      'react/react-in-jsx-scope': 'off',

      // React hooks rules
      ...reactHooks.configs.recommended.rules,

      // Import rules
      'import/default': 'off',
      'import/no-unresolved': 'error',
      'import/named': 'error',
      'import/namespace': 'error',
      'import/no-absolute-path': 'error',
      'import/no-dynamic-require': 'error',
      'import/no-self-import': 'error',
      'import/no-cycle': 'error',
      'import/no-useless-path-segments': 'error'
    },
    settings: {
      react: {
        version: 'detect'
      },
      'import/resolver': {
        typescript: {
          project: ['./frontend/tsconfig.json', './backend/tsconfig.json']
        }
      }
    }
  },

  // Backend-specific configuration (no React)
  {
    files: ['backend/**/*.ts', 'backend/**/*.tsx'],
    languageOptions: {
      globals: {
        ...globals.node,
        ...globals.es2024 // Backend also gets ES2024 globals
      }
    },
    settings: {
      react: {
        version: '999.999.999' // Disable React detection for backend
      }
    },
    rules: {
      // Backend can be stricter about React
      'react/jsx-uses-react': 'off',
      'react/jsx-uses-vars': 'off'
    }
  },

  // Frontend-specific configuration
  {
    files: ['frontend/**/*.ts', 'frontend/**/*.tsx'],
    languageOptions: {
      globals: {
        ...globals.browser,
        ...globals.es2024 // Frontend gets ES2024 globals
      }
    }
  },

  // Test files configuration (Vitest + Jest compatibility)
  {
    files: [
      '**/*.test.ts',
      '**/*.test.tsx',
      '**/*.spec.ts',
      '**/*.spec.tsx',
      '**/test/**/*.ts',
      '**/test/**/*.tsx'
    ],
    languageOptions: {
      globals: {
        ...globals.node,
        ...globals.es2024,

        // Vitest globals
        vi: 'readonly',
        describe: 'readonly',
        it: 'readonly',
        test: 'readonly',
        expect: 'readonly',
        beforeEach: 'readonly',
        afterEach: 'readonly',
        beforeAll: 'readonly',
        afterAll: 'readonly',

        // Jest compatibility (in case some tests still use Jest)
        jest: 'readonly'
      }
    },
    rules: {
      // Allow any types in test files for mocking
      '@typescript-eslint/no-explicit-any': 'off',
      // Allow non-null assertions in tests
      '@typescript-eslint/no-non-null-assertion': 'off',
      // Allow empty functions in mocks
      '@typescript-eslint/no-empty-function': 'off'
    }
  },

  // Node.js configuration files
  {
    files: ['**/*.js', '**/*.cjs', '**/*.mjs'],
    languageOptions: {
      ecmaVersion: 2024, // Updated for ES2024
      globals: {
        ...globals.node,
        ...globals.es2024, // Config files can use ES2024 too
        require: 'readonly',
        module: 'readonly',
        exports: 'readonly',
        console: 'readonly',
        process: 'readonly',
        __dirname: 'readonly',
        __filename: 'readonly'
      }
    }
  },

  // Prettier config (must be last to override other configs)
  eslintConfigPrettier
]
