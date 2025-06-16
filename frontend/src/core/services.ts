import { container } from './di/Container'
import { ApiClient, type AuthProvider, type ApiConfig } from './api/ApiClient'
import { EntryService } from '../domains/entries/EntryService'
import { MealService } from '../domains/meals/MealService'
import { AnalyticsService } from '../domains/analytics/AnalyticsService'
import { useAuthStore } from '../stores/authStore'

// Service identifiers
export const SERVICE_IDENTIFIERS = {
  API_CONFIG: Symbol('ApiConfig'),
  AUTH_PROVIDER: Symbol('AuthProvider'),
  API_CLIENT: Symbol('ApiClient'),
  ENTRY_SERVICE: Symbol('EntryService'),
  MEAL_SERVICE: Symbol('MealService'),
  ANALYTICS_SERVICE: Symbol('AnalyticsService')
} as const

// Auth provider implementation
class ZustandAuthProvider implements AuthProvider {
  getToken(): string | null {
    return useAuthStore.getState().token
  }
}

// Initialize services
export function initializeServices(): void {
  // Configuration
  container.bind(SERVICE_IDENTIFIERS.API_CONFIG).toInstance({
    baseUrl: import.meta.env.VITE_API_URL ?? 'http://localhost:3002',
    timeout: 10000
  } satisfies ApiConfig)

  // Auth provider
  container.bind(SERVICE_IDENTIFIERS.AUTH_PROVIDER).toClass(ZustandAuthProvider).asSingleton()

  // API client
  container.bind(SERVICE_IDENTIFIERS.API_CLIENT).toFactory(() => {
    const config = container.get<ApiConfig>(SERVICE_IDENTIFIERS.API_CONFIG)
    const authProvider = container.get<AuthProvider>(SERVICE_IDENTIFIERS.AUTH_PROVIDER)
    return new ApiClient(config, authProvider)
  }).asSingleton()

  // Domain services
  container.bind(SERVICE_IDENTIFIERS.ENTRY_SERVICE).toFactory(() => {
    const apiClient = container.get<ApiClient>(SERVICE_IDENTIFIERS.API_CLIENT)
    return new EntryService(apiClient)
  }).asSingleton()

  container.bind(SERVICE_IDENTIFIERS.MEAL_SERVICE).toFactory(() => {
    const apiClient = container.get<ApiClient>(SERVICE_IDENTIFIERS.API_CLIENT)
    return new MealService(apiClient)
  }).asSingleton()

  container.bind(SERVICE_IDENTIFIERS.ANALYTICS_SERVICE).toFactory(() => {
    const apiClient = container.get<ApiClient>(SERVICE_IDENTIFIERS.API_CLIENT)
    return new AnalyticsService(apiClient)
  }).asSingleton()
}

// Convenience getters
export const getEntryService = (): EntryService => 
  container.get<EntryService>(SERVICE_IDENTIFIERS.ENTRY_SERVICE)

export const getMealService = (): MealService => 
  container.get<MealService>(SERVICE_IDENTIFIERS.MEAL_SERVICE)

export const getAnalyticsService = (): AnalyticsService => 
  container.get<AnalyticsService>(SERVICE_IDENTIFIERS.ANALYTICS_SERVICE)

// Export the container for direct access
export { container }
