import js from '@eslint/js'
import tseslint from '@typescript-eslint/eslint-plugin'
import tsparser from '@typescript-eslint/parser'
import reactHooks from 'eslint-plugin-react-hooks'
import reactRefresh from 'eslint-plugin-react-refresh'

export default [
  js.configs.recommended,
  // Backend TypeScript configuration
  {
    files: ['backend/src/**/*.{ts,tsx}'],
    ignores: ['backend/src/**/__tests__/**/*'],
    languageOptions: {
      parser: tsparser,
      parserOptions: {
        ecmaVersion: 'latest',
        sourceType: 'module',
        project: './backend/tsconfig.json'
      },
      globals: {
        console: 'readonly',
        process: 'readonly',
        Buffer: 'readonly',
        __dirname: 'readonly',
        __filename: 'readonly',
        global: 'readonly',
        module: 'readonly',
        require: 'readonly',
        exports: 'readonly'
      }
    },
    plugins: {
      '@typescript-eslint': tseslint
    },
    rules: {
      // StandardJS style rules (compatible with Prettier)
      'indent': 'off', // Let Prettier handle indentation
      'quotes': ['error', 'single'],
      'semi': ['error', 'never'],
      'no-unused-vars': 'off', // Use TypeScript version instead
      'no-undef': 'error',
      'comma-dangle': 'off', // Let Prettier handle trailing commas
      'space-before-function-paren': 'off', // Let Prettier handle spacing
      
      // TypeScript specific rules
      '@typescript-eslint/no-unused-vars': 'error',
      '@typescript-eslint/explicit-function-return-type': 'warn',
      '@typescript-eslint/no-explicit-any': 'warn',
      'no-console': 'warn'
    }
  },
  // Frontend React TypeScript configuration
  {
    files: ['frontend/src/**/*.{ts,tsx}'],
    languageOptions: {
      parser: tsparser,
      parserOptions: {
        ecmaVersion: 'latest',
        sourceType: 'module',
        project: './frontend/tsconfig.json',
        ecmaFeatures: {
          jsx: true
        }
      },
      globals: {
        console: 'readonly',
        window: 'readonly',
        document: 'readonly',
        localStorage: 'readonly',
        sessionStorage: 'readonly',
        fetch: 'readonly',
        confirm: 'readonly',
        alert: 'readonly',
        setTimeout: 'readonly',
        clearTimeout: 'readonly',
        setInterval: 'readonly',
        clearInterval: 'readonly'
      }
    },
    plugins: {
      '@typescript-eslint': tseslint,
      'react-hooks': reactHooks,
      'react-refresh': reactRefresh
    },
    rules: {
      // StandardJS style rules (compatible with Prettier)
      'indent': 'off', // Let Prettier handle indentation
      'quotes': ['error', 'single'],
      'semi': ['error', 'never'],
      'no-unused-vars': 'off', // Use TypeScript version instead
      'no-undef': 'error',
      'comma-dangle': 'off', // Let Prettier handle trailing commas
      'space-before-function-paren': 'off', // Let Prettier handle spacing
      
      // TypeScript specific rules
      '@typescript-eslint/no-unused-vars': 'error',
      '@typescript-eslint/explicit-function-return-type': 'off', // Too strict for React components
      '@typescript-eslint/no-explicit-any': 'warn',
      'no-console': 'warn',
      
      // React specific rules
      ...reactHooks.configs.recommended.rules,
      'react-refresh/only-export-components': ['warn', { allowConstantExport: true }]
    }
  },
  // Test files configuration
  {
    files: ['**/__tests__/**/*.{ts,tsx}', '**/*.test.{ts,tsx}', '**/*.spec.{ts,tsx}', '**/test/**/*.{ts,tsx}'],
    languageOptions: {
      parser: tsparser,
      parserOptions: {
        ecmaVersion: 'latest',
        sourceType: 'module'
      },
      globals: {
        console: 'readonly',
        process: 'readonly',
        Buffer: 'readonly',
        __dirname: 'readonly',
        __filename: 'readonly',
        global: 'readonly',
        module: 'readonly',
        require: 'readonly',
        exports: 'readonly',
        // Test globals
        jest: 'readonly',
        describe: 'readonly',
        it: 'readonly',
        test: 'readonly',
        expect: 'readonly',
        beforeEach: 'readonly',
        afterEach: 'readonly',
        beforeAll: 'readonly',
        afterAll: 'readonly',
        vi: 'readonly',
        vitest: 'readonly',
        // Browser globals for frontend tests
        window: 'readonly',
        document: 'readonly',
        localStorage: 'readonly',
        sessionStorage: 'readonly',
        fetch: 'readonly'
      }
    },
    plugins: {
      '@typescript-eslint': tseslint
    },
    rules: {
      // Relaxed rules for test files
      'no-console': 'off',
      '@typescript-eslint/no-explicit-any': 'off',
      '@typescript-eslint/explicit-function-return-type': 'off',
      '@typescript-eslint/no-unused-vars': 'warn'
    }
  },
  // Global ignores
  {
    ignores: [
      '**/dist/',
      '**/build/',
      '**/node_modules/',
      '**/coverage/',
      '**/.next/',
      '**/public/',
      '**/*.js',
      '**/prisma/migrations/',
      'ai-service/**'
    ]
  }
]
