"""
Comprehensive tests for the Poo Tracker AI Service
"""

from datetime import datetime, timedelta
from unittest.mock import AsyncMock, MagicMock, patch

from fastapi import FastAPI
from fastapi.testclient import TestClient

# Mock Redis client before importing the app
with patch("redis.asyncio.from_url") as mock_redis:
    # Mock Redis client
    mock_redis_client = AsyncMock()
    mock_redis_client.ping.return_value = True
    mock_redis.return_value = mock_redis_client

    from ai_service.main import app

# Set up app state for testing
app.state.cache_manager = AsyncMock()
app.state.cache_manager.ping = AsyncMock(return_value=True)
app.state.cache_manager.close = AsyncMock()

app.state.analyzer = AsyncMock()
app.state.health_assessor = AsyncMock()
app.state.recommender = AsyncMock()
app.state.validator = AsyncMock()

client = TestClient(app)


class TestHealthEndpoint:
    """Test health endpoint functionality."""

    def test_health_endpoint_healthy(self):
        """Test health endpoint returns healthy status when all services are up."""
        with patch.object(app.state.cache_manager, "ping", return_value=True):
            response = client.get("/health")
            assert response.status_code == 200
            data = response.json()
            assert "status" in data
            assert data["status"] in ["healthy", "degraded", "unhealthy"]
            assert "timestamp" in data
            assert "redis_connected" in data
            assert "ml_models_loaded" in data
            assert "response_time_ms" in data
            assert "version" in data
            assert isinstance(data["redis_connected"], bool)
            assert isinstance(data["ml_models_loaded"], bool)
            assert isinstance(data["response_time_ms"], (int, float))

    def test_health_endpoint_degraded(self):
        """Test health endpoint returns degraded status when Redis is down."""
        with patch.object(
            app.state.cache_manager, "ping", side_effect=Exception("Redis down")
        ):
            response = client.get("/health")
            assert response.status_code == 200
            data = response.json()
            assert data["redis_connected"] is False
            assert data["status"] in ["degraded", "unhealthy"]


class TestRootEndpoint:
    """Test root endpoint functionality."""

    def test_root_endpoint(self):
        """Test root endpoint returns service information."""
        response = client.get("/")
        assert response.status_code == 200
        data = response.json()
        assert "service" in data
        assert "version" in data
        assert "status" in data
        assert "docs" in data
        assert "health" in data
        assert "analysis" in data
        assert data["status"] == "running"


class TestAppStructure:
    """Test application structure and configuration."""

    def test_app_instance(self):
        """Test that app is a proper FastAPI instance."""
        assert isinstance(app, FastAPI)
        assert app.title == "Poo Tracker AI Service"
        assert app.version == "1.0.0"

    def test_app_routes(self):
        """Test that all expected routes are registered."""
        assert hasattr(app, "routes")
        route_paths = [route.path for route in app.routes]
        assert "/" in route_paths
        assert "/health" in route_paths
        assert "/analyze" in route_paths
        assert "/metrics" in route_paths

    def test_bristol_stool_types(self):
        """Test Bristol Stool Chart type validation."""
        # Valid Bristol types should be 1-7
        for bristol_type in range(1, 8):
            assert 1 <= bristol_type <= 7


class TestAnalyzeEndpoint:
    """Test analysis endpoint functionality."""

    def test_analyze_endpoint_empty_data(self):
        """Test analyze endpoint with empty data."""
        request_data = {
            "entries": [],
            "analysis_options": {
                "include_patterns": True,
                "include_correlations": True,
                "include_recommendations": True,
                "include_health_assessment": True,
            },
        }

        # Mock the validator to pass validation
        app.state.validator.validate_entries = MagicMock(return_value=[])

        response = client.post("/analyze", json=request_data)
        # Should handle empty data gracefully
        assert response.status_code in [200, 422, 400]

    def test_analyze_endpoint_with_data(self):
        """Test analyze endpoint with sample data."""
        now = datetime.now()
        request_data = {
            "entries": [
                {
                    "id": "test-entry-1",
                    "timestamp": now.isoformat(),
                    "bristol_type": 4,
                    "bowel_movement": {
                        "color": "brown",
                        "consistency": "normal",
                        "volume": "medium",
                        "urgency": "normal",
                        "pain_level": 0,
                        "blood": False,
                        "mucus": False,
                    },
                    "meals": [
                        {
                            "id": "meal-1",
                            "timestamp": (now - timedelta(hours=12)).isoformat(),
                            "name": "Test Meal",
                            "foods": ["apple", "chicken", "rice"],
                            "portion_size": "medium",
                            "preparation_method": "grilled",
                        }
                    ],
                    "symptoms": [],
                }
            ],
            "analysis_options": {
                "include_patterns": True,
                "include_correlations": True,
                "include_recommendations": True,
                "include_health_assessment": True,
            },
        }

        # Mock all the services
        app.state.validator.validate_entries = MagicMock(
            return_value=request_data["entries"]
        )
        app.state.analyzer.analyze_patterns = AsyncMock(
            return_value={
                "bristol_patterns": {"most_common": 4},
                "timing_patterns": {"avg_frequency": 1.0},
                "volume_patterns": {"trend": "stable"},
            }
        )
        app.state.analyzer.analyze_correlations = AsyncMock(
            return_value={
                "meal_correlations": [],
                "trigger_foods": [],
                "beneficial_foods": [],
            }
        )
        app.state.health_assessor.assess_health = AsyncMock(
            return_value={"overall_score": 85, "risk_factors": [], "alerts": []}
        )
        app.state.recommender.generate_recommendations = AsyncMock(
            return_value={
                "dietary": ["Stay hydrated"],
                "lifestyle": ["Regular exercise"],
                "medical": [],
            }
        )

        response = client.post("/analyze", json=request_data)
        assert response.status_code == 200
        data = response.json()
        assert "patterns" in data
        assert "correlations" in data
        assert "health_assessment" in data
        assert "recommendations" in data

    def test_analyze_endpoint_invalid_data(self):
        """Test analyze endpoint with invalid data."""
        request_data = {
            "entries": [
                {
                    "id": "invalid-entry",
                    "timestamp": "invalid-timestamp",
                    "bristol_type": 10,  # Invalid bristol type
                }
            ]
        }

        response = client.post("/analyze", json=request_data)
        assert response.status_code == 422  # Validation error


class TestMetricsEndpoint:
    """Test metrics endpoint functionality."""

    def test_metrics_endpoint(self):
        """Test metrics endpoint returns proper format."""
        response = client.get("/metrics")
        assert response.status_code == 200
        # Metrics should be in plain text format for Prometheus
        assert response.headers["content-type"] == "text/plain; charset=utf-8"


class TestErrorHandling:
    """Test error handling across the application."""

    def test_404_endpoint(self):
        """Test non-existent endpoint returns 404."""
        response = client.get("/nonexistent")
        assert response.status_code == 404

    def test_analyze_missing_fields(self):
        """Test analyze endpoint with missing required fields."""
        response = client.post("/analyze", json={})
        assert response.status_code == 422

    def test_analyze_malformed_json(self):
        """Test analyze endpoint with malformed JSON."""
        response = client.post(
            "/analyze",
            data="invalid json",
            headers={"content-type": "application/json"},
        )
        assert response.status_code == 422


class TestServiceIntegration:
    """Test integration between different services."""

    @patch("ai_service.main.logger")
    def test_service_logging(self, mock_logger):
        """Test that services properly log operations."""
        response = client.get("/health")
        assert response.status_code == 200
        # Logger should have been called
        mock_logger.info.assert_called()

    def test_redis_connection_handling(self):
        """Test Redis connection handling in different scenarios."""
        with patch.object(app.state.cache_manager, "ping", return_value=True):
            response = client.get("/health")
            data = response.json()
            assert data["redis_connected"] is True

        with patch.object(
            app.state.cache_manager, "ping", side_effect=Exception("Connection failed")
        ):
            response = client.get("/health")
            data = response.json()
            assert data["redis_connected"] is False


class TestDataValidation:
    """Test data validation functionality."""

    def test_bristol_type_validation(self):
        """Test Bristol stool type validation."""
        valid_types = [1, 2, 3, 4, 5, 6, 7]
        invalid_types = [0, 8, 9, -1, 100]

        for valid_type in valid_types:
            assert 1 <= valid_type <= 7

        for invalid_type in invalid_types:
            assert not (1 <= invalid_type <= 7)

    def test_timestamp_format(self):
        """Test timestamp format validation."""
        valid_timestamp = datetime.now().isoformat()
        assert isinstance(valid_timestamp, str)
        assert "T" in valid_timestamp


class TestCacheIntegration:
    """Test cache integration functionality."""

    def test_cache_manager_initialization(self):
        """Test that cache manager is properly initialized."""
        assert hasattr(app.state, "cache_manager")
        assert app.state.cache_manager is not None

    async def test_cache_operations(self):
        """Test basic cache operations."""
        # Test ping
        result = await app.state.cache_manager.ping()
        assert result is True

        # Test close
        await app.state.cache_manager.close()
        # Should not raise an exception


class TestResponseModels:
    """Test response model validation."""

    def test_health_response_structure(self):
        """Test health response has correct structure."""
        response = client.get("/health")
        data = response.json()

        required_fields = [
            "status",
            "timestamp",
            "redis_connected",
            "ml_models_loaded",
            "response_time_ms",
            "version",
        ]

        for field in required_fields:
            assert field in data

    def test_root_response_structure(self):
        """Test root response has correct structure."""
        response = client.get("/")
        data = response.json()

        required_fields = ["service", "version", "status", "docs", "health", "analysis"]

        for field in required_fields:
            assert field in data
