"""Health assessment service for calculating health scores and trends."""

from typing import Any

from ..config.logging import get_logger
from ..models.database import BowelMovementData, SymptomData
from ..models.responses import HealthScore
from ..utils.health_metrics import HealthMetricsCalculator

logger = get_logger("health_assessor")


class HealthAssessorService:
    """Service for calculating comprehensive health scores and assessments."""

    def __init__(self):
        self.health_calculator = HealthMetricsCalculator()

    async def assess_health(self, entries: list[dict[str, Any]]) -> dict[str, Any]:
        """Placeholder for tests."""
        return {}

    async def calculate_health_score(
        self,
        bowel_movements: list[BowelMovementData],
        symptoms: list[SymptomData] | None = None,
    ) -> HealthScore:
        """
        Calculate comprehensive health score.

        Args:
            bowel_movements: List of bowel movement data
            symptoms: Optional list of symptom data

        Returns:
            Comprehensive health score with component breakdowns
        """
        logger.info(f"Calculating health score for {len(bowel_movements)} entries")

        if not bowel_movements:
            return HealthScore(
                overall_score=0.0,
                bristol_score=0.0,
                frequency_score=0.0,
                pain_score=100.0,
                satisfaction_score=50.0,
                trend="stable",
            )

        # Extract data for calculations
        bristol_types = [bm.bristol_type for bm in bowel_movements]
        pain_scores = [bm.pain for bm in bowel_movements if bm.pain is not None]
        satisfaction_scores = [
            bm.satisfaction for bm in bowel_movements if bm.satisfaction is not None
        ]

        # Calculate daily frequencies
        daily_frequencies = self._calculate_daily_frequencies(bowel_movements)

        # Calculate component scores
        scores = self.health_calculator.calculate_overall_health_score(
            bristol_types=bristol_types,
            daily_frequencies=daily_frequencies,
            pain_scores=pain_scores if pain_scores else None,
            satisfaction_scores=satisfaction_scores if satisfaction_scores else None,
        )

        # Determine trend
        trend = await self._calculate_trend(bowel_movements)

        return HealthScore(
            overall_score=scores["overall_score"],
            bristol_score=scores["bristol_score"],
            frequency_score=scores["frequency_score"],
            pain_score=scores["pain_score"],
            satisfaction_score=scores["satisfaction_score"],
            trend=trend,
        )

    async def assess_digestive_health(
        self,
        bowel_movements: list[BowelMovementData],
        symptoms: list[SymptomData] | None = None,
    ) -> dict[str, Any]:
        """
        Perform comprehensive digestive health assessment.

        Args:
            bowel_movements: Bowel movement data
            symptoms: Optional symptom data

        Returns:
            Detailed health assessment
        """
        assessment = {}

        # Basic health score
        health_score = await self.calculate_health_score(bowel_movements, symptoms)
        assessment["health_score"] = health_score.dict()

        # Risk factors
        bristol_types = [bm.bristol_type for bm in bowel_movements]
        pain_scores = [bm.pain for bm in bowel_movements if bm.pain is not None]
        frequencies = self._calculate_daily_frequencies(bowel_movements)

        risk_factors = self.health_calculator.calculate_risk_factors(
            bristol_types=bristol_types,
            pain_scores=pain_scores if pain_scores else None,
            frequencies=frequencies,
        )
        assessment["risk_factors"] = risk_factors

        # Regularity assessment
        timing_data = [bm.created_at for bm in bowel_movements]
        regularity = self.health_calculator.calculate_digestive_regularity_index(
            timing_data=timing_data,
            bristol_data=bristol_types,
        )
        assessment["regularity"] = regularity

        return assessment

    def _calculate_daily_frequencies(
        self, bowel_movements: list[BowelMovementData]
    ) -> list[float]:
        """Calculate daily frequencies from bowel movement data."""
        if not bowel_movements:
            return []

        # Group by date
        daily_counts = {}
        for bm in bowel_movements:
            date_key = bm.created_at.date()
            daily_counts[date_key] = daily_counts.get(date_key, 0) + 1

        return list(daily_counts.values())

    async def _calculate_trend(self, bowel_movements: list[BowelMovementData]) -> str:
        """Calculate health trend from recent vs historical data."""
        if len(bowel_movements) < 14:  # Need at least 2 weeks of data
            return "stable"

        # Sort by date
        sorted_movements = sorted(bowel_movements, key=lambda x: x.created_at)

        # Split into recent and historical (50-50 split)
        midpoint = len(sorted_movements) // 2
        historical = sorted_movements[:midpoint]
        recent = sorted_movements[midpoint:]

        # Calculate average bristol scores for both periods
        historical_bristol = sum(bm.bristol_type for bm in historical) / len(historical)
        recent_bristol = sum(bm.bristol_type for bm in recent) / len(recent)

        # Calculate average pain scores (if available)
        historical_pain = [bm.pain for bm in historical if bm.pain is not None]
        recent_pain = [bm.pain for bm in recent if bm.pain is not None]

        bristol_trend = self._assess_bristol_trend(historical_bristol, recent_bristol)
        pain_trend = self._assess_pain_trend(historical_pain, recent_pain)

        # Combine trends (bristol is more important)
        if bristol_trend == "improving" and pain_trend != "declining":
            return "improving"
        elif bristol_trend == "declining" or pain_trend == "declining":
            return "declining"
        else:
            return "stable"

    def _assess_bristol_trend(self, historical_avg: float, recent_avg: float) -> str:
        """Assess trend based on Bristol type averages."""
        ideal_range = (3, 4)  # Ideal Bristol types

        # Calculate distance from ideal for both periods
        historical_distance = min(
            abs(historical_avg - ideal_range[0]), abs(historical_avg - ideal_range[1])
        )
        recent_distance = min(
            abs(recent_avg - ideal_range[0]), abs(recent_avg - ideal_range[1])
        )

        # Determine trend
        improvement_threshold = 0.3
        if recent_distance < historical_distance - improvement_threshold:
            return "improving"
        elif recent_distance > historical_distance + improvement_threshold:
            return "declining"
        else:
            return "stable"

    def _assess_pain_trend(
        self, historical_pain: list[int], recent_pain: list[int]
    ) -> str:
        """Assess trend based on pain scores."""
        if not historical_pain or not recent_pain:
            return "stable"

        historical_avg = sum(historical_pain) / len(historical_pain)
        recent_avg = sum(recent_pain) / len(recent_pain)

        # Lower pain is better
        improvement_threshold = 1.0
        if recent_avg < historical_avg - improvement_threshold:
            return "improving"
        elif recent_avg > historical_avg + improvement_threshold:
            return "declining"
        else:
            return "stable"
