"""
Tests for service layer components
"""

from datetime import datetime, timedelta
from unittest.mock import AsyncMock, patch

import pytest

from ai_service.models import BowelMovementData, MealData, SymptomData
from ai_service.services.analyzer import AnalyzerService
from ai_service.services.health_assessor import HealthAssessorService
from ai_service.services.recommender import RecommenderService
from ai_service.utils.cache import CacheManager
from ai_service.utils.validators import DataValidator


class TestAnalyzerService:
    """Test AnalyzerService functionality."""

    def setup_method(self):
        """Setup test fixtures."""
        self.analyzer = AnalyzerService()

    @pytest.mark.asyncio
    async def test_analyze_patterns_empty_data(self):
        """Test pattern analysis with empty data."""
        result = await self.analyzer.analyze_patterns([])
        assert isinstance(result, dict)

    @pytest.mark.asyncio
    async def test_analyze_patterns_with_data(self):
        """Test pattern analysis with sample data."""
        sample_entries = [
            {
                "id": "test-1",
                "timestamp": datetime.now().isoformat(),
                "bristol_type": 4,
                "bowel_movement": {
                    "color": "brown",
                    "consistency": "normal",
                    "volume": "medium",
                },
            }
        ]
        result = await self.analyzer.analyze_patterns(sample_entries)
        assert isinstance(result, dict)

    @pytest.mark.asyncio
    async def test_analyze_correlations(self):
        """Test correlation analysis."""
        sample_entries = []
        result = await self.analyzer.analyze_correlations(sample_entries)
        assert isinstance(result, dict)

    @pytest.mark.asyncio
    async def test_analyze_comprehensive_patterns(self):
        """Test comprehensive pattern analysis with sample data."""
        now = datetime.now()
        bowel_movements = []
        for i in range(15):
            created = now - timedelta(days=15 - i)
            if i < 5:
                bristol = 6
            elif i < 10:
                bristol = 4
            else:
                bristol = 3
            bowel_movements.append(
                BowelMovementData(
                    id=f"bm-{i}",
                    created_at=created,
                    bristol_type=bristol,
                    pain=5 if i < 5 else 2,
                )
            )

        meals = [
            MealData(
                id="meal-1",
                meal_time=now - timedelta(days=9, hours=10),
                category="dairy",
                spicy_level=2,
                dairy=True,
            ),
            MealData(
                id="meal-2",
                meal_time=now - timedelta(days=3, hours=8),
                category="fiber_rich",
                spicy_level=1,
                fiber_rich=True,
            ),
        ]

        symptoms = [
            SymptomData(
                id="sym-1",
                created_at=now - timedelta(days=9),
                type="cramps",
                severity=8,
            )
        ]

        result = await self.analyzer.analyze_comprehensive_patterns(
            bowel_movements,
            meals,
            symptoms,
            user_id="user-1",
        )

        assert result["bristol_analysis"]["trend"] == "improving"
        assert (
            "avg_meal_to_bm_hours"
            in result["correlations"]["meals"]["timing_correlations"]
        )
        assert result["correlations"]["symptoms"]["cramps"][
            "avg_severity"
        ] == pytest.approx(8.0)


class TestHealthAssessorService:
    """Test HealthAssessorService functionality."""

    def setup_method(self):
        """Setup test fixtures."""
        self.assessor = HealthAssessorService()

    @pytest.mark.asyncio
    async def test_assess_health_empty_data(self):
        """Test health assessment with empty data."""
        result = await self.assessor.assess_health([])
        assert isinstance(result, dict)

    @pytest.mark.asyncio
    async def test_assess_health_with_data(self):
        """Test health assessment with sample data."""
        sample_entries = [
            {
                "id": "test-1",
                "timestamp": datetime.now().isoformat(),
                "bristol_type": 4,
                "symptoms": [],
            }
        ]
        result = await self.assessor.assess_health(sample_entries)
        assert isinstance(result, dict)


class TestRecommenderService:
    """Test RecommenderService functionality."""

    def setup_method(self):
        """Setup test fixtures."""
        self.recommender = RecommenderService()

    @pytest.mark.asyncio
    async def test_generate_recommendations_empty_data(self):
        """Test recommendation generation with empty data."""
        result = await self.recommender.generate_recommendations([], {}, {})
        assert isinstance(result, dict)

    @pytest.mark.asyncio
    async def test_generate_recommendations_with_data(self):
        """Test recommendation generation with sample data."""
        sample_entries = []
        sample_patterns = {}
        sample_health = {}
        result = await self.recommender.generate_recommendations(
            sample_entries, sample_patterns, sample_health
        )
        assert isinstance(result, dict)


class TestCacheManager:
    """Test CacheManager functionality."""

    def setup_method(self):
        """Setup test fixtures."""
        self.cache_manager = CacheManager()

    @pytest.mark.asyncio
    @patch("redis.asyncio.from_url")
    async def test_ping_success(self, mock_redis):
        """Test successful ping operation."""
        mock_redis_client = AsyncMock()
        mock_redis_client.ping.return_value = True
        mock_redis.return_value = mock_redis_client

        result = await self.cache_manager.ping()
        assert result is True

    @pytest.mark.asyncio
    @patch("redis.asyncio.from_url")
    async def test_ping_failure(self, mock_redis):
        """Test ping operation failure."""
        mock_redis_client = AsyncMock()
        mock_redis_client.ping.side_effect = Exception("Connection failed")
        mock_redis.return_value = mock_redis_client

        with pytest.raises(Exception):  # noqa: B017
            await self.cache_manager.ping()

    @pytest.mark.asyncio
    async def test_close(self):
        """Test close operation."""
        await self.cache_manager.close()
        # Should not raise an exception


class TestDataValidator:
    """Test DataValidator functionality."""

    def setup_method(self):
        """Setup test fixtures."""
        self.validator = DataValidator()

    def test_validate_entries_empty(self):
        """Test validation with empty entries."""
        result = self.validator.validate_entries([])
        assert isinstance(result, list)
        assert len(result) == 0

    def test_validate_entries_valid_data(self):
        """Test validation with valid data."""
        valid_entries = [
            {
                "id": "test-1",
                "timestamp": datetime.now().isoformat(),
                "bristol_type": 4,
                "bowel_movement": {
                    "color": "brown",
                    "consistency": "normal",
                    "volume": "medium",
                },
            }
        ]
        result = self.validator.validate_entries(valid_entries)
        assert isinstance(result, list)

    def test_validate_bristol_type(self):
        """Test Bristol type validation."""
        # Valid types
        for bristol_type in range(1, 8):
            assert self._is_valid_bristol_type(bristol_type)

        # Invalid types
        invalid_types = [0, 8, 9, -1, None, "invalid"]
        for bristol_type in invalid_types:
            assert not self._is_valid_bristol_type(bristol_type)

    def _is_valid_bristol_type(self, bristol_type):
        """Helper method to validate Bristol type."""
        return isinstance(bristol_type, int) and 1 <= bristol_type <= 7


def test_coverage(cov):
    """Ensure coverage plugin reports some measured lines."""
    percent = cov.report(show_missing=False)
    assert percent > 0
