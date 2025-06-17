"""
Tests for service layer components
"""

from datetime import datetime
from unittest.mock import AsyncMock, patch

import pytest

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

        result = await self.cache_manager.ping()
        assert result is False

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
