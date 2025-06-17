"""Application settings and configuration."""

from functools import lru_cache
from typing import Any

from pydantic_settings import BaseSettings, SettingsConfigDict


class Settings(BaseSettings):
    """Application settings loaded from environment variables."""

    model_config = SettingsConfigDict(
        env_file=".env",
        env_file_encoding="utf-8",
        extra="ignore",  # Ignore extra fields from environment
    )

    # Basic app settings
    app_name: str = "Poo Tracker AI Service"
    app_version: str = "1.0.0"
    debug: bool = False
    environment: str = "development"

    # Server settings
    host: str = "0.0.0.0"
    port: int = 8000
    workers: int = 1

    # Redis settings
    redis_url: str = "redis://localhost:6379"
    redis_timeout: int = 5
    redis_retry_on_timeout: bool = True

    # Cache settings
    cache_ttl: int = 3600  # 1 hour default
    cache_prefix: str = "poo_tracker"

    # ML settings
    ml_model_path: str = "./data/models"
    enable_ml_features: bool = True

    # Analysis settings
    max_analysis_days: int = 365
    min_entries_for_ml: int = 10

    # Security settings
    api_key: str | None = None
    enable_auth: bool = False

    # Rate limiting
    rate_limit_requests: int = 100
    rate_limit_window: int = 60  # seconds

    # Logging
    log_level: str = "INFO"
    log_format: str = "json"  # json or text

    # Backend integration
    backend_url: str = "http://localhost:3002"
    backend_timeout: int = 30

    @property
    def redis_config(self) -> dict[str, Any]:
        """Get Redis configuration dictionary."""
        return {
            "url": self.redis_url,
            "socket_connect_timeout": self.redis_timeout,
            "socket_timeout": self.redis_timeout,
            "retry_on_timeout": self.redis_retry_on_timeout,
            "health_check_interval": 30,
        }

    @property
    def is_production(self) -> bool:
        """Check if running in production environment."""
        return self.environment.lower() == "production"

    @property
    def is_development(self) -> bool:
        """Check if running in development environment."""
        return self.environment.lower() == "development"


@lru_cache
def get_settings() -> Settings:
    """Get cached application settings."""
    return Settings()
