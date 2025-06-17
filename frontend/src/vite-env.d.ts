// This file provides type definitions for Vite's special environment variables

interface ImportMetaEnv {
  readonly VITE_API_URL?: string
  // Add other env variables as needed
}

interface ImportMeta {
  readonly env: ImportMetaEnv
}
