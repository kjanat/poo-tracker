"""Application settings and configuration."""

import os
from functools import lru_cache
from typing import Any

from pydantic import BaseSettings, Field


class Settings(BaseSettings):
    """Application settings loaded from environment variables."""

    # Basic app settings
    app_name: str = "Poo Tracker AI Service"
    app_version: str = "1.0.0"
    debug: bool = Field(default=False, env="DEBUG")
    environment: str = Field(default="development", env="ENVIRONMENT")

    # Server settings
    host: str = Field(default="0.0.0.0", env="HOST")
    port: int = Field(default=8000, env="PORT")
    workers: int = Field(default=1, env="WORKERS")

    # Redis settings
    redis_url: str = Field(default="redis://localhost:6379", env="REDIS_URL")
    redis_timeout: int = Field(default=5, env="REDIS_TIMEOUT")
    redis_retry_on_timeout: bool = Field(default=True, env="REDIS_RETRY_ON_TIMEOUT")

    # Cache settings
    cache_ttl: int = Field(default=3600, env="CACHE_TTL")  # 1 hour default
    cache_prefix: str = Field(default="poo_tracker", env="CACHE_PREFIX")

    # ML settings
    ml_model_path: str = Field(default="./data/models", env="ML_MODEL_PATH")
    enable_ml_features: bool = Field(default=True, env="ENABLE_ML_FEATURES")

    # Analysis settings
    max_analysis_days: int = Field(default=365, env="MAX_ANALYSIS_DAYS")
    min_entries_for_ml: int = Field(default=10, env="MIN_ENTRIES_FOR_ML")

    # Security settings
    api_key: str | None = Field(default=None, env="API_KEY")
    enable_auth: bool = Field(default=False, env="ENABLE_AUTH")

    # Rate limiting
    rate_limit_requests: int = Field(default=100, env="RATE_LIMIT_REQUESTS")
    rate_limit_window: int = Field(default=60, env="RATE_LIMIT_WINDOW")  # seconds

    # Logging
    log_level: str = Field(default="INFO", env="LOG_LEVEL")
    log_format: str = Field(default="json", env="LOG_FORMAT")  # json or text

    # Backend integration
    backend_url: str = Field(default="http://localhost:3002", env="BACKEND_URL")
    backend_timeout: int = Field(default=30, env="BACKEND_TIMEOUT")

    class Config:
        """Pydantic config."""

        env_file = ".env"
        env_file_encoding = "utf-8"

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
