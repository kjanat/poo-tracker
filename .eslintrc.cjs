module.exports = {
  root: true,
  parser: '@typescript-eslint/parser',
  parserOptions: {
    ecmaVersion: 'latest',
    sourceType: 'module',
    project: ['./frontend/tsconfig.json', './backend/tsconfig.json'],
    tsconfigRootDir: __dirname
  },
  env: {
    browser: true,
    node: true,
    es2022: true
  },
  plugins: ['@typescript-eslint', 'react', 'react-hooks', 'import'],
  extends: [
    'eslint:recommended',
    'plugin:@typescript-eslint/recommended',
    'plugin:react/recommended',
    'plugin:react-hooks/recommended',
    'plugin:import/recommended',
    'plugin:import/typescript',
    'prettier'
  ],
  settings: {
    react: { version: 'detect' },
    'import/resolver': {
      typescript: {}
    }
  },
  rules: {
    'react/react-in-jsx-scope': 'off',
    'import/default': 'off'
  },
  ignorePatterns: ['dist/', 'build/', 'node_modules/', 'coverage/', '**/*.d.ts', 'branding/']
}
