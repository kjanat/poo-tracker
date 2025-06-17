"""
Basic tests for the Poo Tracker AI Service
"""

from unittest.mock import AsyncMock, patch

import pytest
from fastapi import FastAPI
from fastapi.testclient import TestClient

# Mock Redis client before importing the app
with patch("redis.asyncio.from_url") as mock_redis:
    # Mock Redis client
    mock_redis_client = AsyncMock()
    mock_redis_client.ping.return_value = True
    mock_redis.return_value = mock_redis_client

    from src.ai_service.main import app

# Set up app state for testing
app.state.cache_manager = AsyncMock()
app.state.cache_manager.ping = AsyncMock(return_value=True)

client = TestClient(app)


def test_health_endpoint():
    """Test that the health endpoint returns 200"""
    response = client.get("/health")
    assert response.status_code == 200
    data = response.json()
    assert "status" in data
    assert data["status"] == "healthy"
    assert "timestamp" in data
    assert "redis_connected" in data
    # redis_connected should be a boolean
    assert isinstance(data["redis_connected"], bool)


def test_app_info():
    """Test that the app has correct metadata"""
    assert app.title == "Poo Tracker AI Service"
    assert app.version == "1.0.0"


def test_root_endpoint():
    """Test the root endpoint if it exists"""
    response = client.get("/")
    # This might return 404 if no root endpoint exists, which is fine
    assert response.status_code in [200, 404]


@pytest.mark.asyncio
async def test_api_structure():
    """Test that the API is properly structured"""
    # Check that the app is a FastAPI instance
    assert isinstance(app, FastAPI)

    # Check that routes exist
    assert hasattr(app, "routes")
    assert len(app.routes) > 0


def test_bristol_types_validation():
    """Test Bristol stool chart type validation (basic)"""
    # This is a placeholder - you can expand based on your actual endpoints
    valid_bristol_types = [1, 2, 3, 4, 5, 6, 7]
    for bristol_type in valid_bristol_types:
        assert 1 <= bristol_type <= 7


def test_redis_connection_handling():
    """Test that Redis connection is handled gracefully"""

    # Test health endpoint (Redis should be mocked)
    response = client.get("/health")
    assert response.status_code == 200
    data = response.json()

    # With our mock, redis_connected should be True
    assert data["redis_connected"] is True
    assert isinstance(data["redis_connected"], bool)
