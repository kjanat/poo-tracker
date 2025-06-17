import { defineConfig } from 'vitest/config'
import { resolve } from 'path'

export default defineConfig({
  test: {
    environment: 'node',
    include: ['src/**/__tests__/**/*.ts', 'src/**/?(*.)+(spec|test).ts'],
    reporters: [
      'default',
      ['junit', { outputFile: './../backend-junit.xml' }]
    ],
    coverage: {
      provider: 'v8',
      reporter: ['text', 'lcov', 'html'],
      include: ['src/**/*.ts'],
      exclude: ['src/**/*.d.ts', 'src/index.ts', 'src/**/__tests__/**']
    },
    globals: true
  },
  resolve: {
    alias: {
      '@': resolve(__dirname, './src'),
    }
  }
})
