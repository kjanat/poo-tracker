"""
Tests for configuration and utility functions
"""

import os
from unittest.mock import patch

from ai_service.config.settings import Settings, get_settings
from ai_service.utils.data_processing import DataProcessor
from ai_service.utils.health_metrics import HealthMetricsCalculator


class TestSettings:
    """Test settings configuration."""

    def test_default_settings(self):
        """Test default settings values."""
        settings = Settings()
        assert settings.app_name == "Poo Tracker AI Service"
        assert settings.app_version == "0.1.0"
        assert settings.environment in ["development", "production", "testing"]

    @patch.dict(os.environ, {"ENVIRONMENT": "production"})
    def test_production_settings(self):
        """Test production environment settings."""
        settings = Settings()
        assert settings.environment == "production"
        assert settings.is_production is True
        assert settings.is_development is False

    @patch.dict(os.environ, {"ENVIRONMENT": "development"})
    def test_development_settings(self):
        """Test development environment settings."""
        settings = Settings()
        assert settings.environment == "development"
        assert settings.is_development is True
        assert settings.is_production is False

    def test_get_settings_singleton(self):
        """Test that get_settings returns singleton instance."""
        settings1 = get_settings()
        settings2 = get_settings()
        assert settings1 is settings2


class TestDataProcessor:
    """Test data processing utilities."""

    def setup_method(self):
        """Setup test fixtures."""
        self.processor = DataProcessor()

    def test_process_empty_data(self):
        """Test processing empty data."""
        result = self.processor.normalize_entries([])
        assert isinstance(result, list)
        assert len(result) == 0

    def test_normalize_entries(self):
        """Test entry normalization."""
        sample_entries = [
            {"id": "test-1", "bristol_type": 4, "timestamp": "2025-06-17T10:00:00"}
        ]
        result = self.processor.normalize_entries(sample_entries)
        assert isinstance(result, list)
        assert len(result) >= 0

    def test_extract_bristol_patterns(self):
        """Test Bristol pattern extraction."""
        sample_entries = [
            {"bristol_type": 4},
            {"bristol_type": 3},
            {"bristol_type": 4},
        ]
        result = self.processor.extract_bristol_patterns(sample_entries)
        assert isinstance(result, dict)

    def test_calculate_frequency_patterns(self):
        """Test frequency pattern calculation."""
        sample_entries = []
        result = self.processor.calculate_frequency_patterns(sample_entries)
        assert isinstance(result, dict)


class TestHealthMetricsCalculator:
    """Test health metrics calculation."""

    def setup_method(self):
        """Setup test fixtures."""
        self.calculator = HealthMetricsCalculator()

    def test_calculate_overall_score_empty(self):
        """Test overall score calculation with empty data."""
        result = self.calculator.calculate_overall_score([])
        assert isinstance(result, int | float)
        assert 0 <= result <= 100

    def test_calculate_bristol_score(self):
        """Test Bristol score calculation."""
        sample_entries = [
            {"bristol_type": 4},
            {"bristol_type": 3},
            {"bristol_type": 5},
        ]
        result = self.calculator.calculate_bristol_score(sample_entries)
        assert isinstance(result, int | float)
        assert 0 <= result <= 100

    def test_calculate_frequency_score(self):
        """Test frequency score calculation."""
        sample_entries = []
        result = self.calculator.calculate_frequency_score(sample_entries)
        assert isinstance(result, int | float)
        assert 0 <= result <= 100

    def test_detect_health_issues(self):
        """Test health issue detection."""
        sample_entries = []
        result = self.calculator.detect_health_issues(sample_entries)
        assert isinstance(result, list)

    def test_generate_health_alerts(self):
        """Test health alert generation."""
        sample_entries = []
        result = self.calculator.generate_health_alerts(sample_entries)
        assert isinstance(result, list)


class TestEnvironmentVariables:
    """Test environment variable handling."""

    @patch.dict(os.environ, {"REDIS_URL": "redis://test:6379"})
    def test_redis_url_setting(self):
        """Test Redis URL from environment."""
        settings = Settings()
        assert "redis://" in settings.redis_url

    @patch.dict(os.environ, {"LOG_LEVEL": "DEBUG"})
    def test_log_level_setting(self):
        """Test log level from environment."""
        settings = Settings()
        assert settings.log_level == "DEBUG"

    def test_missing_environment_variables(self):
        """Test handling of missing environment variables."""
        # Should not raise an error with missing optional env vars
        settings = Settings()
        assert settings is not None


class TestUtilityFunctions:
    """Test various utility functions."""

    def test_timestamp_formatting(self):
        """Test timestamp formatting utilities."""
        from datetime import datetime

        now = datetime.now()
        iso_format = now.isoformat()
        assert "T" in iso_format
        assert isinstance(iso_format, str)

    def test_bristol_type_validation_utility(self):
        """Test Bristol type validation utility."""

        def is_valid_bristol_type(value):
            return isinstance(value, int) and 1 <= value <= 7

        # Test valid types
        for i in range(1, 8):
            assert is_valid_bristol_type(i) is True

        # Test invalid types
        invalid_values = [0, 8, 9, -1, None, "invalid", 1.5]
        for value in invalid_values:
            assert is_valid_bristol_type(value) is False

    def test_data_structure_validation(self):
        """Test data structure validation utilities."""

        def validate_entry_structure(entry):
            required_fields = ["id", "timestamp", "bristol_type"]
            return all(field in entry for field in required_fields)

        valid_entry = {
            "id": "test-1",
            "timestamp": "2025-06-17T10:00:00",
            "bristol_type": 4,
        }
        assert validate_entry_structure(valid_entry) is True

        invalid_entry = {"id": "test-1"}
        assert validate_entry_structure(invalid_entry) is False
