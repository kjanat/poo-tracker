import { defineConfig } from 'vitest/config';
import { resolve } from 'path';

export default defineConfig({
    test: {
        environment: 'jsdom',
        setupFiles: ['./src/test/setup.ts'],
        coverage: {
            include: [
                'src/**/*.{ts,tsx}'
            ],
            exclude: [
                'src/**/*.d.ts',
                'src/main.tsx',
                'src/vite-env.d.ts'
            ]
        }
  },
    resolve: {
        alias: {
            '@': resolve(__dirname, './src')
        }
    }
});
