"""Pattern detection service for identifying digestive health patterns."""

from datetime import timedelta
from typing import Any

import numpy as np

from ..config.logging import get_logger
from ..models.database import BowelMovementData, MealData, SymptomData
from ..utils.data_processing import DataProcessor

logger = get_logger("pattern_detector")


class PatternDetectorService:
    """Service for detecting patterns in digestive health data."""

    def __init__(self):
        self.data_processor = DataProcessor()

    async def detect_all_patterns(
        self,
        bowel_movements: list[BowelMovementData],
        meals: list[MealData] | None = None,
        symptoms: list[SymptomData] | None = None,
        confidence_threshold: float = 0.7,
    ) -> dict[str, Any]:
        """
        Detect all types of patterns in the data.

        Args:
            bowel_movements: Bowel movement data
            meals: Optional meal data
            symptoms: Optional symptom data
            confidence_threshold: Minimum confidence for pattern detection

        Returns:
            Dictionary containing all detected patterns
        """
        logger.info("Starting comprehensive pattern detection")

        patterns = {}

        # Timing patterns
        timing_patterns = await self._detect_timing_patterns(bowel_movements)
        if timing_patterns:
            patterns["timing"] = timing_patterns

        # Frequency patterns
        frequency_patterns = await self._detect_frequency_patterns(bowel_movements)
        if frequency_patterns:
            patterns["frequency"] = frequency_patterns

        # Bristol type patterns
        bristol_patterns = await self._detect_bristol_patterns(bowel_movements)
        if bristol_patterns:
            patterns["bristol"] = bristol_patterns

        # Meal correlation patterns
        if meals:
            meal_patterns = await self._detect_meal_patterns(bowel_movements, meals)
            if meal_patterns:
                patterns["meal_correlations"] = meal_patterns

        # Symptom correlation patterns
        if symptoms:
            symptom_patterns = await self._detect_symptom_patterns(
                bowel_movements, symptoms
            )
            if symptom_patterns:
                patterns["symptom_correlations"] = symptom_patterns

        # Weekly/cyclical patterns
        cyclical_patterns = await self._detect_cyclical_patterns(bowel_movements)
        if cyclical_patterns:
            patterns["cyclical"] = cyclical_patterns

        return patterns

    async def _detect_timing_patterns(
        self, bowel_movements: list[BowelMovementData]
    ) -> dict[str, Any] | None:
        """Detect timing-based patterns."""
        if len(bowel_movements) < 7:  # Need at least a week of data
            return None

        timing_data = [bm.created_at for bm in bowel_movements]
        hours = [dt.hour for dt in timing_data]

        # Find peak hours
        hour_counts = {}
        for hour in hours:
            hour_counts[hour] = hour_counts.get(hour, 0) + 1

        total_movements = len(hours)
        peak_hours = []

        for hour, count in hour_counts.items():
            frequency = count / total_movements
            if frequency > 0.2:  # At least 20% of movements
                peak_hours.append(
                    {
                        "hour": hour,
                        "frequency": frequency,
                        "confidence": min(
                            0.9, frequency * 2
                        ),  # Higher frequency = higher confidence
                    }
                )

        # Calculate regularity
        hour_std = np.std(hours)
        regularity_score = max(
            0, 1 - (hour_std / 12)
        )  # 12 is max possible std for hours

        return {
            "peak_hours": peak_hours,
            "regularity_score": regularity_score,
            "hour_standard_deviation": hour_std,
            "pattern_strength": regularity_score,
        }

    async def _detect_frequency_patterns(
        self, bowel_movements: list[BowelMovementData]
    ) -> dict[str, Any] | None:
        """Detect frequency-based patterns."""
        if len(bowel_movements) < 14:  # Need at least 2 weeks
            return None

        # Group by date
        daily_counts = {}
        for bm in bowel_movements:
            date_key = bm.created_at.date()
            daily_counts[date_key] = daily_counts.get(date_key, 0) + 1

        counts = list(daily_counts.values())
        avg_frequency = np.mean(counts)
        frequency_std = np.std(counts)

        # Calculate consistency (lower coefficient of variation = more consistent)
        cv = frequency_std / avg_frequency if avg_frequency > 0 else 1
        consistency_score = max(0, 1 - min(cv, 1))

        # Identify patterns
        patterns = []

        if avg_frequency < 0.5:
            patterns.append(
                {
                    "type": "low_frequency",
                    "description": f"Infrequent bowel movements (avg {avg_frequency:.1f} per day)",
                    "confidence": 0.8,
                }
            )
        elif avg_frequency > 3:
            patterns.append(
                {
                    "type": "high_frequency",
                    "description": f"Frequent bowel movements (avg {avg_frequency:.1f} per day)",
                    "confidence": 0.8,
                }
            )

        if consistency_score > 0.7:
            patterns.append(
                {
                    "type": "consistent_frequency",
                    "description": "Highly consistent bowel movement frequency",
                    "confidence": consistency_score,
                }
            )

        return {
            "average_daily_frequency": avg_frequency,
            "consistency_score": consistency_score,
            "frequency_standard_deviation": frequency_std,
            "detected_patterns": patterns,
        }

    async def _detect_bristol_patterns(
        self, bowel_movements: list[BowelMovementData]
    ) -> dict[str, Any] | None:
        """Detect Bristol stool type patterns."""
        if not bowel_movements:
            return None

        bristol_types = [bm.bristol_type for bm in bowel_movements]
        bristol_counts = {}
        for bt in bristol_types:
            bristol_counts[bt] = bristol_counts.get(bt, 0) + 1

        total = len(bristol_types)
        patterns = []

        # Check for dominant patterns
        for bristol_type, count in bristol_counts.items():
            prevalence = count / total
            if prevalence > 0.4:  # More than 40% of movements
                if bristol_type <= 2:
                    patterns.append(
                        {
                            "type": "constipation_pattern",
                            "bristol_type": bristol_type,
                            "prevalence": prevalence,
                            "confidence": min(0.9, prevalence * 1.5),
                            "description": f"Constipation pattern (Type {bristol_type})",
                        }
                    )
                elif bristol_type >= 6:
                    patterns.append(
                        {
                            "type": "diarrhea_pattern",
                            "bristol_type": bristol_type,
                            "prevalence": prevalence,
                            "confidence": min(0.9, prevalence * 1.5),
                            "description": f"Loose stool pattern (Type {bristol_type})",
                        }
                    )
                elif bristol_type in [3, 4]:
                    patterns.append(
                        {
                            "type": "healthy_pattern",
                            "bristol_type": bristol_type,
                            "prevalence": prevalence,
                            "confidence": min(0.9, prevalence * 1.2),
                            "description": f"Healthy pattern (Type {bristol_type})",
                        }
                    )

        # Check for variability
        bristol_std = np.std(bristol_types)
        if bristol_std > 2:
            patterns.append(
                {
                    "type": "high_variability",
                    "description": "High variability in stool types",
                    "confidence": min(0.8, bristol_std / 3),
                    "standard_deviation": bristol_std,
                }
            )

        return {
            "bristol_distribution": bristol_counts,
            "variability": bristol_std,
            "detected_patterns": patterns,
        }

    async def _detect_meal_patterns(
        self, bowel_movements: list[BowelMovementData], meals: list[MealData]
    ) -> dict[str, Any] | None:
        """Detect meal correlation patterns."""
        if not meals or len(bowel_movements) < 5:
            return None

        patterns = []
        correlations = {}

        # Group meals by characteristics
        meal_groups = {
            "spicy": [m for m in meals if m.spicy_level and m.spicy_level > 6],
            "dairy": [m for m in meals if m.dairy],
            "gluten": [m for m in meals if m.gluten],
            "fiber_rich": [m for m in meals if m.fiber_rich],
        }

        for food_type, type_meals in meal_groups.items():
            if not type_meals:
                continue

            # Find bowel movements 6-48 hours after these meals
            related_movements = []
            for meal in type_meals:
                meal_time = meal.meal_time
                window_start = meal_time + timedelta(hours=6)
                window_end = meal_time + timedelta(hours=48)

                for bm in bowel_movements:
                    if window_start <= bm.created_at <= window_end:
                        related_movements.append(bm)

            if len(related_movements) >= 3:  # Need at least 3 related movements
                # Calculate correlation metrics
                bristol_types = [bm.bristol_type for bm in related_movements]
                avg_bristol = sum(bristol_types) / len(bristol_types)

                pain_scores = [
                    bm.pain for bm in related_movements if bm.pain is not None
                ]
                avg_pain = sum(pain_scores) / len(pain_scores) if pain_scores else None

                # Determine if this is a trigger or beneficial food
                is_trigger = (
                    avg_bristol <= 2 or avg_bristol >= 6 or (avg_pain and avg_pain > 6)
                )
                is_beneficial = 3 <= avg_bristol <= 4 and (
                    not avg_pain or avg_pain <= 3
                )

                if is_trigger:
                    patterns.append(
                        {
                            "type": "trigger_food",
                            "food_type": food_type,
                            "average_bristol": avg_bristol,
                            "average_pain": avg_pain,
                            "sample_size": len(related_movements),
                            "confidence": min(0.8, len(related_movements) / 10),
                            "description": f"{food_type.title()} foods may trigger digestive issues",
                        }
                    )
                elif is_beneficial:
                    patterns.append(
                        {
                            "type": "beneficial_food",
                            "food_type": food_type,
                            "average_bristol": avg_bristol,
                            "average_pain": avg_pain,
                            "sample_size": len(related_movements),
                            "confidence": min(0.7, len(related_movements) / 10),
                            "description": f"{food_type.title()} foods may be beneficial for digestion",
                        }
                    )

                correlations[food_type] = {
                    "average_bristol": avg_bristol,
                    "average_pain": avg_pain,
                    "sample_size": len(related_movements),
                }

        return {
            "detected_patterns": patterns,
            "food_correlations": correlations,
        }

    async def _detect_symptom_patterns(
        self, bowel_movements: list[BowelMovementData], symptoms: list[SymptomData]
    ) -> dict[str, Any] | None:
        """Detect symptom correlation patterns."""
        if not symptoms:
            return None

        patterns = []

        # Group symptoms by type
        symptom_groups = {}
        for symptom in symptoms:
            if symptom.type not in symptom_groups:
                symptom_groups[symptom.type] = []
            symptom_groups[symptom.type].append(symptom)

        for symptom_type, type_symptoms in symptom_groups.items():
            if len(type_symptoms) < 3:  # Need at least 3 occurrences
                continue

            # Find bowel movements around symptom times
            related_movements = []
            for symptom in type_symptoms:
                symptom_time = symptom.created_at
                window_start = symptom_time - timedelta(hours=12)
                window_end = symptom_time + timedelta(hours=12)

                for bm in bowel_movements:
                    if window_start <= bm.created_at <= window_end:
                        related_movements.append(bm)

            if related_movements:
                bristol_types = [bm.bristol_type for bm in related_movements]
                avg_bristol = sum(bristol_types) / len(bristol_types)

                avg_severity = sum(s.severity for s in type_symptoms) / len(
                    type_symptoms
                )

                patterns.append(
                    {
                        "type": "symptom_correlation",
                        "symptom_type": symptom_type,
                        "average_severity": avg_severity,
                        "associated_bristol": avg_bristol,
                        "occurrence_count": len(type_symptoms),
                        "confidence": min(0.8, len(type_symptoms) / 5),
                        "description": f"{symptom_type.title()} symptoms correlate with digestive patterns",
                    }
                )

        return {
            "detected_patterns": patterns,
        }

    async def _detect_cyclical_patterns(
        self, bowel_movements: list[BowelMovementData]
    ) -> dict[str, Any] | None:
        """Detect weekly/cyclical patterns."""
        if len(bowel_movements) < 21:  # Need at least 3 weeks
            return None

        patterns = []

        # Day of week patterns
        day_counts = {}
        for bm in bowel_movements:
            day = bm.created_at.weekday()  # 0 = Monday, 6 = Sunday
            day_counts[day] = day_counts.get(day, 0) + 1

        total = len(bowel_movements)
        day_names = [
            "Monday",
            "Tuesday",
            "Wednesday",
            "Thursday",
            "Friday",
            "Saturday",
            "Sunday",
        ]

        # Find significantly high/low days
        for day, count in day_counts.items():
            frequency = count / total
            expected = 1 / 7  # Expected frequency for uniform distribution

            if frequency > expected * 1.5:  # 50% above expected
                patterns.append(
                    {
                        "type": "high_frequency_day",
                        "day": day_names[day],
                        "frequency": frequency,
                        "confidence": min(0.8, frequency * 3),
                        "description": f"Higher frequency on {day_names[day]}",
                    }
                )
            elif frequency < expected * 0.5:  # 50% below expected
                patterns.append(
                    {
                        "type": "low_frequency_day",
                        "day": day_names[day],
                        "frequency": frequency,
                        "confidence": min(0.8, (expected - frequency) * 3),
                        "description": f"Lower frequency on {day_names[day]}",
                    }
                )

        # Weekend vs weekday pattern
        weekday_count = sum(day_counts.get(i, 0) for i in range(5))  # Mon-Fri
        weekend_count = sum(day_counts.get(i, 0) for i in [5, 6])  # Sat-Sun

        weekday_freq = weekday_count / total if total > 0 else 0
        weekend_freq = weekend_count / total if total > 0 else 0

        expected_weekday = 5 / 7

        if abs(weekday_freq - expected_weekday) > 0.15:  # Significant difference
            pattern_type = (
                "weekday_dominant"
                if weekday_freq > expected_weekday
                else "weekend_dominant"
            )
            patterns.append(
                {
                    "type": pattern_type,
                    "weekday_frequency": weekday_freq,
                    "weekend_frequency": weekend_freq,
                    "confidence": 0.7,
                    "description": "Different patterns between weekdays and weekends",
                }
            )

        return {
            "detected_patterns": patterns,
            "day_distribution": {day_names[k]: v for k, v in day_counts.items()},
            "weekday_vs_weekend": {
                "weekday_frequency": weekday_freq,
                "weekend_frequency": weekend_freq,
            },
        }
