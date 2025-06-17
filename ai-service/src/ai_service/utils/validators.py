"""Data validation utilities."""

from datetime import datetime, timedelta
from typing import Any

from pydantic import BaseModel

from ..config.logging import get_logger
from ..models.requests import AnalysisRequest

logger = get_logger("validator")


class ValidationResult(BaseModel):
    """Result of data validation."""

    is_valid: bool
    errors: list[str] = []
    warnings: list[str] = []


class DataValidator:
    """Data validation utilities for AI service."""

    def __init__(self):
        self.max_entries_per_request = 1000
        self.max_time_range_days = 365
        self.min_bristol_type = 1
        self.max_bristol_type = 7

    def validate_entries(self, entries: list[dict[str, Any]]) -> list[dict[str, Any]]:
        """Validate basic entry structure for tests."""
        if not isinstance(entries, list):
            return []

        required = {"id", "timestamp", "bristol_type"}
        valid: list[dict[str, Any]] = []
        for entry in entries:
            if isinstance(entry, dict) and required.issubset(entry):
                valid.append(entry)

        return valid

    def validate_analysis_request(self, request: AnalysisRequest) -> ValidationResult:
        """
        Validate analysis request data.

        Args:
            request: Analysis request to validate

        Returns:
            Validation result with errors and warnings
        """
        errors = []
        warnings = []

        # Check entry count
        if len(request.entries) == 0:
            errors.append("No bowel movement entries provided")
        elif len(request.entries) > self.max_entries_per_request:
            errors.append(
                f"Too many entries: {len(request.entries)}. Maximum: {self.max_entries_per_request}"
            )

        # Validate time range
        if request.entries:
            time_validation = self._validate_time_range(request.entries)
            errors.extend(time_validation.errors)
            warnings.extend(time_validation.warnings)

        # Validate Bristol types
        bristol_validation = self._validate_bristol_types(request.entries)
        errors.extend(bristol_validation.errors)
        warnings.extend(bristol_validation.warnings)

        # Validate pain/satisfaction scores
        score_validation = self._validate_scores(request.entries)
        errors.extend(score_validation.errors)
        warnings.extend(score_validation.warnings)

        # Validate meals if present
        if request.meals:
            meal_validation = self._validate_meals(request.meals)
            errors.extend(meal_validation.errors)
            warnings.extend(meal_validation.warnings)

        # Validate symptoms if present
        if request.symptoms:
            symptom_validation = self._validate_symptoms(request.symptoms)
            errors.extend(symptom_validation.errors)
            warnings.extend(symptom_validation.warnings)

        # Data quality warnings
        quality_warnings = self._check_data_quality(request)
        warnings.extend(quality_warnings)

        return ValidationResult(
            is_valid=len(errors) == 0,
            errors=errors,
            warnings=warnings,
        )

    def _validate_time_range(self, entries: list[Any]) -> ValidationResult:
        """Validate time range of entries."""
        errors = []
        warnings = []

        if not entries:
            return ValidationResult(is_valid=True)

        dates = [entry.created_at for entry in entries]
        min_date = min(dates)
        max_date = max(dates)

        # Check for future dates
        now = datetime.now()
        future_entries = [d for d in dates if d > now]
        if future_entries:
            errors.append(f"Found {len(future_entries)} entries with future timestamps")

        # Check time range
        time_range = (max_date - min_date).days
        if time_range > self.max_time_range_days:
            warnings.append(
                f"Large time range: {time_range} days. Consider splitting into smaller batches."
            )

        # Check for very old data
        very_old_threshold = now - timedelta(days=730)  # 2 years
        old_entries = [d for d in dates if d < very_old_threshold]
        if old_entries:
            warnings.append(f"Found {len(old_entries)} entries older than 2 years")

        return ValidationResult(
            is_valid=len(errors) == 0,
            errors=errors,
            warnings=warnings,
        )

    def _validate_bristol_types(self, entries: list[Any]) -> ValidationResult:
        """Validate Bristol stool types."""
        errors = []
        warnings = []

        invalid_bristol = [
            entry
            for entry in entries
            if entry.bristol_type < self.min_bristol_type
            or entry.bristol_type > self.max_bristol_type
        ]

        if invalid_bristol:
            errors.append(
                f"Found {len(invalid_bristol)} entries with invalid Bristol types (must be 1-7)"
            )

        # Check for unusual patterns
        bristol_types = [entry.bristol_type for entry in entries]
        if len(set(bristol_types)) == 1 and len(bristol_types) > 10:
            warnings.append(
                "All entries have the same Bristol type - unusual pattern detected"
            )

        return ValidationResult(
            is_valid=len(errors) == 0,
            errors=errors,
            warnings=warnings,
        )

    def _validate_scores(self, entries: list[Any]) -> ValidationResult:
        """Validate pain, strain, and satisfaction scores."""
        errors = []
        warnings = []

        for entry in entries:
            # Validate pain scores
            if entry.pain is not None and (entry.pain < 1 or entry.pain > 10):
                errors.append(f"Invalid pain score: {entry.pain} (must be 1-10)")

            # Validate strain scores
            if entry.strain is not None and (entry.strain < 1 or entry.strain > 10):
                errors.append(f"Invalid strain score: {entry.strain} (must be 1-10)")

            # Validate satisfaction scores
            if entry.satisfaction is not None and (
                entry.satisfaction < 1 or entry.satisfaction > 10
            ):
                errors.append(
                    f"Invalid satisfaction score: {entry.satisfaction} (must be 1-10)"
                )

        # Check for missing critical data
        pain_missing = sum(1 for entry in entries if entry.pain is None)
        if pain_missing > len(entries) * 0.8:
            warnings.append(
                f"Pain data missing for {pain_missing}/{len(entries)} entries"
            )

        return ValidationResult(
            is_valid=len(errors) == 0,
            errors=errors,
            warnings=warnings,
        )

    def _validate_meals(self, meals: list[Any]) -> ValidationResult:
        """Validate meal data."""
        errors = []
        warnings = []

        # Check spicy levels
        invalid_spicy = [
            meal
            for meal in meals
            if meal.spicy_level is not None
            and (meal.spicy_level < 1 or meal.spicy_level > 10)
        ]

        if invalid_spicy:
            errors.append(
                f"Found {len(invalid_spicy)} meals with invalid spicy levels (must be 1-10)"
            )

        # Check for missing meal times
        missing_times = [meal for meal in meals if not meal.meal_time]
        if missing_times:
            errors.append(f"Found {len(missing_times)} meals without meal times")

        # Warn about missing categories
        missing_categories = [meal for meal in meals if not meal.category]
        if len(missing_categories) > len(meals) * 0.5:
            warnings.append(
                f"Meal category missing for {len(missing_categories)}/{len(meals)} meals"
            )

        return ValidationResult(
            is_valid=len(errors) == 0,
            errors=errors,
            warnings=warnings,
        )

    def _validate_symptoms(self, symptoms: list[Any]) -> ValidationResult:
        """Validate symptom data."""
        errors = []
        warnings = []

        # Check severity scores
        invalid_severity = [
            symptom
            for symptom in symptoms
            if symptom.severity < 1 or symptom.severity > 10
        ]

        if invalid_severity:
            errors.append(
                f"Found {len(invalid_severity)} symptoms with invalid severity (must be 1-10)"
            )

        # Check for valid symptom types
        valid_types = {
            "bloating",
            "cramps",
            "nausea",
            "gas",
            "heartburn",
            "constipation",
            "diarrhea",
        }
        invalid_types = [
            symptom for symptom in symptoms if symptom.type not in valid_types
        ]

        if invalid_types:
            warnings.append(f"Found {len(invalid_types)} symptoms with unknown types")

        return ValidationResult(
            is_valid=len(errors) == 0,
            errors=errors,
            warnings=warnings,
        )

    def _check_data_quality(self, request: AnalysisRequest) -> list[str]:
        """Check overall data quality and provide warnings."""
        warnings = []

        # Check data completeness
        total_entries = len(request.entries)

        # Bristol type completeness (always required)
        if total_entries == 0:
            return warnings

        # Pain data completeness
        pain_entries = sum(1 for entry in request.entries if entry.pain is not None)
        pain_completeness = pain_entries / total_entries
        if pain_completeness < 0.5:
            warnings.append(f"Low pain data completeness: {pain_completeness:.1%}")

        # Satisfaction data completeness
        satisfaction_entries = sum(
            1 for entry in request.entries if entry.satisfaction is not None
        )
        satisfaction_completeness = satisfaction_entries / total_entries
        if satisfaction_completeness < 0.5:
            warnings.append(
                f"Low satisfaction data completeness: {satisfaction_completeness:.1%}"
            )

        # Volume/color/consistency data
        detail_entries = sum(
            1
            for entry in request.entries
            if entry.volume or entry.color or entry.consistency
        )
        detail_completeness = detail_entries / total_entries
        if detail_completeness < 0.3:
            warnings.append(f"Low detail data completeness: {detail_completeness:.1%}")

        # Check for sufficient data for analysis
        if total_entries < 7:
            warnings.append("Limited data (<7 entries) may reduce analysis accuracy")
        elif total_entries < 30:
            warnings.append("Moderate data (<30 entries) may limit trend analysis")

        # Check meal-to-entry ratio
        if request.meals:
            meal_ratio = len(request.meals) / total_entries
            if meal_ratio < 0.3:
                warnings.append(
                    "Low meal data relative to bowel movements may limit correlation analysis"
                )

        return warnings

    def validate_bristol_type(self, bristol_type: int) -> bool:
        """Validate individual Bristol type."""
        return self.min_bristol_type <= bristol_type <= self.max_bristol_type

    def validate_score(self, score: int | None) -> bool:
        """Validate individual score (1-10 scale)."""
        if score is None:
            return True
        return 1 <= score <= 10

    def validate_datetime_range(self, dt: datetime, max_future_hours: int = 24) -> bool:
        """Validate datetime is within reasonable range."""
        now = datetime.now()
        max_future = now + timedelta(hours=max_future_hours)
        min_past = now - timedelta(days=self.max_time_range_days)

        return min_past <= dt <= max_future
