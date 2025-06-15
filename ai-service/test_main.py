"""
Basic tests for the Poo Tracker AI Service
"""

import pytest
from fastapi.testclient import TestClient
from main import app

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
    from fastapi import FastAPI

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
    import os

    from main import redis_client

    # Test that the redis client is properly initialized
    assert redis_client is not None

    # Test health endpoint with Redis (should work in CI with Redis service)
    response = client.get("/health")
    assert response.status_code == 200
    data = response.json()

    # If REDIS_URL is set, Redis should be connected
    if os.getenv("REDIS_URL"):
        assert data["redis_connected"] is True
    # If no Redis URL, connection might be False but endpoint should still work
    assert isinstance(data["redis_connected"], bool)


def test_analyze_and_fetch_cached_result():
    """Ensure analyze endpoint caches results and retrieval works."""
    from datetime import datetime, timedelta

    user_id = "test-user"
    now = datetime.now()

    entries = [
        {
            "id": "e1",
            "userId": user_id,
            "bristolType": 4,
            "createdAt": now.isoformat(),
        }
    ]

    meals = [
        {
            "id": "m1",
            "userId": user_id,
            "name": "Toast",
            "category": "breakfast",
            "mealTime": (now - timedelta(hours=8)).isoformat(),
        }
    ]

    analyze_resp = client.post("/analyze", json={"entries": entries, "meals": meals})
    assert analyze_resp.status_code == 200
    analysis_data = analyze_resp.json()

    cache_resp = client.get(f"/analysis/{user_id}")

    from main import redis_client

    try:
        redis_up = redis_client.ping()
    except Exception:
        redis_up = False

    if redis_up:
        assert cache_resp.status_code == 200
        cached_data = cache_resp.json()
        assert cached_data == analysis_data
    else:
        assert cache_resp.status_code == 404
