type Constructor<T = Record<string, unknown>> = new (...args: unknown[]) => T
type ServiceIdentifier<T> = Constructor<T> | string | symbol

interface ServiceBinding<T> {
  factory: () => T
  singleton: boolean
  instance?: T
}

export class DependencyContainer {
  private readonly bindings = new Map<ServiceIdentifier<unknown>, ServiceBinding<unknown>>()

  bind<T>(identifier: ServiceIdentifier<T>): ServiceBuilder<T> {
    return new ServiceBuilder<T>(this, identifier)
  }

  get<T>(identifier: ServiceIdentifier<T>): T {
    const binding = this.bindings.get(identifier)
    if (binding == null) {
      throw new Error(`Service ${String(identifier)} not registered`)
    }

    if (binding.singleton) {
      if (binding.instance == null) {
        binding.instance = binding.factory()
      }
      return binding.instance as T
    }

    return binding.factory() as T
  }

  internal_register<T>(identifier: ServiceIdentifier<T>, binding: ServiceBinding<T>): void {
    this.bindings.set(identifier, binding)
  }
}

class ServiceBuilder<T> {
  constructor(
    private readonly container: DependencyContainer,
    private readonly identifier: ServiceIdentifier<T>
  ) {}

  toFactory(factory: () => T): ServiceLifecycleBuilder<T> {
    return new ServiceLifecycleBuilder<T>(this.container, this.identifier, factory)
  }

  toClass<U extends T>(constructor: Constructor<U>): ServiceLifecycleBuilder<T> {
    const factory = (): T => new constructor()
    return new ServiceLifecycleBuilder<T>(this.container, this.identifier, factory)
  }

  toInstance(instance: T): void {
    this.container.internal_register(this.identifier, {
      factory: () => instance,
      singleton: true,
      instance
    })
  }
}

class ServiceLifecycleBuilder<T> {
  constructor(
    private readonly container: DependencyContainer,
    private readonly identifier: ServiceIdentifier<T>,
    private readonly factory: () => T
  ) {}

  asSingleton(): void {
    this.container.internal_register(this.identifier, {
      factory: this.factory,
      singleton: true
    })
  }

  asTransient(): void {
    this.container.internal_register(this.identifier, {
      factory: this.factory,
      singleton: false
    })
  }
}

// Global container instance
export const container = new DependencyContainer()
