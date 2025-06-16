import type { TestingLibraryMatchers } from '@testing-library/jest-dom/matchers'
import type 'vitest'
import type 'vite/client'

declare module 'vitest' {
  interface Assertion<T = unknown>
    extends jest.Matchers<void>,
    TestingLibraryMatchers<T, void> {}
  interface AsymmetricMatchersContaining
    extends TestingLibraryMatchers<unknown, void> {}
}
