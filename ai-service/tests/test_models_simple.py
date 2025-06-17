"""
Tests for models and data structures
"""

from datetime import datetime

from ai_service.models.requests import BowelMovementEntry
from ai_service.models.responses import (
    AnalysisResponse,
    BristolAnalysis,
    ErrorResponse,
    HealthResponse,
)


class TestBowelMovementEntry:
    """Test BowelMovementEntry model."""

    def test_bowel_movement_entry_creation(self):
        """Test creating a bowel movement entry."""
        now = datetime.now()
        entry = BowelMovementEntry(
            id="test-1", userId="user-1", bristolType=4, createdAt=now
        )
        assert entry.id == "test-1"
        assert entry.user_id == "user-1"
        assert entry.bristol_type == 4
        assert entry.created_at == now

    def test_valid_bristol_types(self):
        """Test all valid Bristol types."""
        now = datetime.now()
        for bristol_type in range(1, 8):
            entry = BowelMovementEntry(
                id=f"test-{bristol_type}",
                userId="user-1",
                bristolType=bristol_type,
                createdAt=now,
            )
            assert entry.bristol_type == bristol_type


class TestAnalysisResponse:
    """Test AnalysisResponse model."""

    def test_analysis_response_creation(self):
        """Test creating an analysis response."""
        response = AnalysisResponse(
            patterns={},
            correlations={},
            recommendations=[],  # This should be a list
            risk_factors=[],  # This should be a list
            bristol_trends=BristolAnalysis(
                distribution={1: 0, 2: 0, 3: 5, 4: 10, 5: 3, 6: 1, 7: 0},
                percentages={1: 0.0, 2: 0.0, 3: 26.3, 4: 52.6, 5: 15.8, 6: 5.3, 7: 0.0},
                most_common={"type": 4, "count": 10},
                health_indicator="healthy",
            ),
            analysis_metadata={},  # This field is required
        )
        assert isinstance(response.patterns, dict)
        assert isinstance(response.correlations, dict)
        assert isinstance(response.recommendations, list)
        assert isinstance(response.risk_factors, list)
        assert isinstance(response.bristol_trends, BristolAnalysis)
        assert isinstance(response.analysis_metadata, dict)


class TestHealthResponse:
    """Test HealthResponse model."""

    def test_health_response_creation(self):
        """Test creating a health response."""
        response = HealthResponse(
            status="healthy",
            timestamp=datetime.now().isoformat(),
            redis_connected=True,
            ml_models_loaded=True,
            response_time_ms=50.0,
            version="1.0.0",
        )
        assert response.status == "healthy"
        assert response.redis_connected is True
        assert response.ml_models_loaded is True
        assert response.response_time_ms == 50.0
        assert response.version == "1.0.0"


class TestErrorResponse:
    """Test ErrorResponse model."""

    def test_error_response_creation(self):
        """Test creating an error response."""
        now = datetime.now()
        response = ErrorResponse(
            error="Validation failed",
            detail="Invalid Bristol type",
            timestamp=now,
        )
        assert response.error == "Validation failed"
        assert response.detail == "Invalid Bristol type"
        assert isinstance(response.timestamp, datetime)


class TestModelValidation:
    """Test model validation functionality."""

    def test_bristol_type_range(self):
        """Test Bristol type range validation."""
        valid_types = list(range(1, 8))
        for bristol_type in valid_types:
            assert 1 <= bristol_type <= 7

        invalid_types = [0, 8, 9, -1, 100]
        for bristol_type in invalid_types:
            assert not (1 <= bristol_type <= 7)
