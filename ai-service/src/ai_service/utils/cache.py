"""Cache management utilities using Redis."""

import hashlib
import json
from typing import Any

import redis.asyncio as redis

from ..config.logging import get_logger
from ..config.settings import get_settings

logger = get_logger("cache")
settings = get_settings()


class CacheManager:
    """Redis cache manager for analysis results and temporary data."""

    def __init__(self):
        self.redis_client = None
        self._connect()

    def _connect(self):
        """Initialize Redis connection."""
        try:
            self.redis_client = redis.from_url(**settings.redis_config)
        except Exception as e:
            logger.warning(f"Failed to connect to Redis: {e}")
            self.redis_client = None

    async def ping(self) -> bool:
        """Test Redis connection."""
        if not self.redis_client:
            return False

        try:
            await self.redis_client.ping()
            return True
        except Exception as e:
            logger.warning(f"Redis ping failed: {e}")
            raise

    async def get(self, key: str) -> Any | None:
        """Get value from cache."""
        if not self.redis_client:
            return None

        try:
            value = await self.redis_client.get(key)
            if value:
                return json.loads(value)
            return None
        except Exception as e:
            logger.warning(f"Cache get failed for key {key}: {e}")
            return None

    async def set(self, key: str, value: Any, ttl: int = 3600) -> bool:
        """Set value in cache with TTL."""
        if not self.redis_client:
            return False

        try:
            serialized_value = json.dumps(value, default=str)
            await self.redis_client.setex(key, ttl, serialized_value)
            return True
        except Exception as e:
            logger.warning(f"Cache set failed for key {key}: {e}")
            return False

    async def delete(self, key: str) -> bool:
        """Delete key from cache."""
        if not self.redis_client:
            return False

        try:
            await self.redis_client.delete(key)
            return True
        except Exception as e:
            logger.warning(f"Cache delete failed for key {key}: {e}")
            return False

    async def generate_analysis_cache_key(self, request: Any) -> str:
        """Generate cache key for analysis request."""
        # Create a hash of the request data
        request_str = json.dumps(request.dict(), sort_keys=True, default=str)
        hash_digest = hashlib.md5(request_str.encode()).hexdigest()

        return f"{settings.cache_prefix}:analysis:{hash_digest}"

    async def get_analysis_result(self, cache_key: str) -> dict[str, Any] | None:
        """Get cached analysis result."""
        return await self.get(cache_key)

    async def cache_analysis_result(
        self, cache_key: str, result: dict[str, Any], ttl: int = None
    ) -> bool:
        """Cache analysis result."""
        cache_ttl = ttl or settings.cache_ttl
        return await self.set(cache_key, result, cache_ttl)

    async def get_cache_stats(self) -> dict[str, Any]:
        """Get cache statistics."""
        if not self.redis_client:
            return {"connected": False}

        try:
            info = await self.redis_client.info()

            return {
                "connected": True,
                "used_memory": info.get("used_memory_human", "unknown"),
                "connected_clients": info.get("connected_clients", 0),
                "total_commands_processed": info.get("total_commands_processed", 0),
                "keyspace_hits": info.get("keyspace_hits", 0),
                "keyspace_misses": info.get("keyspace_misses", 0),
                "hit_rate": self._calculate_hit_rate(
                    info.get("keyspace_hits", 0), info.get("keyspace_misses", 0)
                ),
            }
        except Exception as e:
            logger.warning(f"Failed to get cache stats: {e}")
            return {"connected": False, "error": str(e)}

    def _calculate_hit_rate(self, hits: int, misses: int) -> float:
        """Calculate cache hit rate."""
        total = hits + misses
        if total == 0:
            return 0.0
        return round((hits / total) * 100, 2)

    async def close(self):
        """Close Redis connection."""
        if self.redis_client:
            await self.redis_client.aclose()
