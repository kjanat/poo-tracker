import { defineConfig } from 'vitest/config'
import path from 'path'

export default defineConfig({
  test: {
    environment: 'node',
    include: ['src/**/__tests__/**/*.ts', 'src/**/?(*.)+(spec|test).ts'],
    reporters: [
      'default',
      ['html', { outputFile: './test-report.html' }],
      ['text', { outputFile: './test-report.txt' }],
      ['junit', { outputFile: './junit.xml' }],
      ['json', { outputFile: './test-report.json' }]
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
      '@': path.resolve(__dirname, './src')
    }
  }
})
