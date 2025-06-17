"""Logging configuration for the AI service."""

import logging
import sys
from typing import Any

from .settings import get_settings


def setup_logging() -> None:
    """Configure application logging."""
    settings = get_settings()

    # Set log level
    log_level = getattr(logging, settings.log_level.upper(), logging.INFO)

    # Configure root logger
    logging.basicConfig(
        level=log_level,
        format="%(asctime)s - %(name)s - %(levelname)s - %(message)s",
        handlers=[
            logging.StreamHandler(sys.stdout),
        ],
    )

    # Configure specific loggers
    logger = logging.getLogger("ai_service")
    logger.setLevel(log_level)

    # Suppress verbose logs from external libraries in production
    if settings.is_production:
        logging.getLogger("uvicorn").setLevel(logging.WARNING)
        logging.getLogger("fastapi").setLevel(logging.WARNING)
        logging.getLogger("redis").setLevel(logging.WARNING)


def get_logger(name: str) -> logging.Logger:
    """Get a logger instance for the given name."""
    return logging.getLogger(f"ai_service.{name}")


class ContextLogger:
    """Logger with context information."""

    def __init__(self, logger: logging.Logger, context: dict[str, Any] | None = None):
        self.logger = logger
        self.context = context or {}

    def _format_message(self, message: str) -> str:
        """Add context to log message."""
        if self.context:
            context_str = " | ".join(f"{k}={v}" for k, v in self.context.items())
            return f"{message} | {context_str}"
        return message

    def info(self, message: str, **kwargs: Any) -> None:
        """Log info message with context."""
        self.logger.info(self._format_message(message), **kwargs)

    def warning(self, message: str, **kwargs: Any) -> None:
        """Log warning message with context."""
        self.logger.warning(self._format_message(message), **kwargs)

    def error(self, message: str, **kwargs: Any) -> None:
        """Log error message with context."""
        self.logger.error(self._format_message(message), **kwargs)

    def debug(self, message: str, **kwargs: Any) -> None:
        """Log debug message with context."""
        self.logger.debug(self._format_message(message), **kwargs)
