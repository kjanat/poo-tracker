"""Core analysis service for bowel movement pattern analysis."""

import uuid
from datetime import datetime, timedelta
from typing import Any

import pandas as pd

from ..config.logging import get_logger
from ..models.database import AnalysisMetadata, BowelMovementData, MealData, SymptomData
from ..models.responses import (
    BristolAnalysis,
    FrequencyStats,
    TimingPattern,
)
from ..utils.data_processing import DataProcessor
from ..utils.health_metrics import HealthMetricsCalculator

logger = get_logger("analyzer")


class AnalyzerService:
    """Core service for analyzing bowel movement patterns and correlations."""

    def __init__(self):
        self.data_processor = DataProcessor()
        self.health_calculator = HealthMetricsCalculator()

    async def analyze_comprehensive_patterns(
        self,
        bowel_movements: list[BowelMovementData],
        meals: list[MealData] | None = None,
        symptoms: list[SymptomData] | None = None,
        user_id: str | None = None,
    ) -> dict[str, Any]:
        """
        Perform comprehensive analysis of bowel movement patterns.

        Args:
            bowel_movements: List of bowel movement data
            meals: Optional list of meal data for correlation analysis
            symptoms: Optional list of symptom data
            user_id: User ID for personalized analysis

        Returns:
            Comprehensive analysis results
        """
        logger.info(
            f"Starting comprehensive analysis for {len(bowel_movements)} entries"
        )

        if not bowel_movements:
            raise ValueError("No bowel movement data provided")

        # Convert to DataFrames for analysis
        bm_df = self._create_bowel_movement_dataframe(bowel_movements)
        meal_df = self._create_meal_dataframe(meals) if meals else pd.DataFrame()
        symptom_df = (
            self._create_symptom_dataframe(symptoms) if symptoms else pd.DataFrame()
        )

        # Perform individual analyses
        bristol_analysis = await self._analyze_bristol_patterns(bm_df)
        timing_patterns = await self._analyze_timing_patterns(bm_df)
        frequency_stats = await self._calculate_frequency_stats(bm_df)

        # Correlation analysis if meal data available
        meal_correlations = {}
        if not meal_df.empty:
            meal_correlations = await self._analyze_meal_correlations(bm_df, meal_df)

        # Symptom correlations if symptom data available
        symptom_correlations = {}
        if not symptom_df.empty:
            symptom_correlations = await self._analyze_symptom_correlations(
                bm_df, symptom_df
            )

        # Create metadata
        metadata = self._create_analysis_metadata(
            user_id=user_id or "unknown",
            bowel_movements=bowel_movements,
            meals=meals or [],
            symptoms=symptoms or [],
        )

        return {
            "patterns": {
                "timing": timing_patterns.dict(),
                "frequency": frequency_stats.dict(),
                "consistency_trends": self._analyze_consistency_trends(bm_df),
            },
            "bristol_analysis": bristol_analysis.dict(),
            "correlations": {
                "meals": meal_correlations,
                "symptoms": symptom_correlations,
            },
            "metadata": metadata.dict(),
        }

    async def _analyze_bristol_patterns(self, df: pd.DataFrame) -> BristolAnalysis:
        """Analyze Bristol Stool Chart patterns."""
        bristol_counts = df["bristol_type"].value_counts().sort_index()
        total_entries = len(df)

        # Bristol type descriptions
        bristol_descriptions = {
            1: "Severe constipation (hard lumps)",
            2: "Mild constipation (lumpy sausage)",
            3: "Normal (cracked sausage)",
            4: "Ideal (smooth sausage)",
            5: "Lacking fiber (soft blobs)",
            6: "Mild diarrhea (fluffy pieces)",
            7: "Severe diarrhea (watery)",
        }

        # Calculate percentages
        bristol_percentages = (bristol_counts / total_entries * 100).round(2)
        most_common_type = bristol_counts.idxmax()

        # Assess health indicator
        health_indicator = self._assess_bristol_health(bristol_percentages)

        # Determine trend if we have time data
        trend = None
        if len(df) > 10:  # Need sufficient data for trend
            recent_data = df.tail(len(df) // 3)  # Last third of data
            early_data = df.head(len(df) // 3)  # First third of data

            recent_avg = recent_data["bristol_type"].mean()
            early_avg = early_data["bristol_type"].mean()

            # Ideal range is 3-4
            recent_distance = min(abs(recent_avg - 3), abs(recent_avg - 4))
            early_distance = min(abs(early_avg - 3), abs(early_avg - 4))

            if recent_distance < early_distance - 0.3:
                trend = "improving"
            elif recent_distance > early_distance + 0.3:
                trend = "declining"
            else:
                trend = "stable"

        return BristolAnalysis(
            distribution=bristol_counts.to_dict(),
            percentages=bristol_percentages.to_dict(),
            most_common={
                "type": int(most_common_type),
                "description": bristol_descriptions.get(
                    int(most_common_type), "Unknown"
                ),
                "percentage": float(bristol_percentages[most_common_type]),
            },
            health_indicator=health_indicator,
            trend=trend,
        )

    async def _analyze_timing_patterns(self, df: pd.DataFrame) -> TimingPattern:
        """Analyze timing patterns of bowel movements."""
        # Extract time features
        df["hour"] = df["created_at"].dt.hour
        df["day_of_week_name"] = df["created_at"].dt.day_name()
        df["day_of_week"] = df["created_at"].dt.dayofweek

        # Calculate distributions
        hourly_dist = df["hour"].value_counts().sort_index()
        daily_dist = df["day_of_week_name"].value_counts()

        # Calculate regularity score
        regularity_score = self._calculate_regularity_score(df)

        return TimingPattern(
            hourly_distribution=hourly_dist.to_dict(),
            daily_distribution=daily_dist.to_dict(),
            peak_hour=int(hourly_dist.idxmax()),
            most_active_day=daily_dist.idxmax(),
            regularity_score=regularity_score,
        )

    async def _calculate_frequency_stats(self, df: pd.DataFrame) -> FrequencyStats:
        """Calculate frequency statistics."""
        df["date"] = df["created_at"].dt.date
        daily_counts = df.groupby("date").size()

        # Calculate consistency score based on variance
        consistency_score = self._calculate_consistency_score(daily_counts)

        return FrequencyStats(
            avg_daily=float(daily_counts.mean()),
            max_daily=int(daily_counts.max()),
            min_daily=int(daily_counts.min()),
            total_days=len(daily_counts),
            total_entries=len(df),
            consistency_score=consistency_score,
        )

    async def _analyze_meal_correlations(
        self, bm_df: pd.DataFrame, meal_df: pd.DataFrame
    ) -> dict[str, Any]:
        """Analyze correlations between meals and bowel movements."""
        category_correlations = {}
        trigger_foods = []
        beneficial_foods = []

        # Group meals by category
        if "category" in meal_df.columns:
            for category in meal_df["category"].dropna().unique():
                category_meals = meal_df[meal_df["category"] == category]
                correlation_data = self._calculate_meal_category_correlation(
                    bm_df, category_meals, category
                )
                if correlation_data:
                    category_correlations[category] = correlation_data

        # Identify trigger and beneficial foods
        trigger_foods, beneficial_foods = self._identify_trigger_beneficial_foods(
            bm_df, meal_df
        )

        # Calculate timing correlations
        timing_correlations = self._calculate_meal_timing_correlations(bm_df, meal_df)

        return {
            "category_correlations": category_correlations,
            "trigger_foods": trigger_foods,
            "beneficial_foods": beneficial_foods,
            "timing_correlations": timing_correlations,
        }

    async def _analyze_symptom_correlations(
        self, bm_df: pd.DataFrame, symptom_df: pd.DataFrame
    ) -> dict[str, Any]:
        """Analyze correlations between symptoms and bowel movements."""
        correlations = {}

        # Group symptoms by type
        for symptom_type in symptom_df["type"].unique():
            type_symptoms = symptom_df[symptom_df["type"] == symptom_type]
            correlation = self._calculate_symptom_correlation(bm_df, type_symptoms)
            if correlation:
                correlations[symptom_type] = correlation

        return correlations

    def _create_bowel_movement_dataframe(
        self, data: list[BowelMovementData]
    ) -> pd.DataFrame:
        """Convert bowel movement data to DataFrame."""
        return pd.DataFrame([item.to_dict() for item in data])

    def _create_meal_dataframe(self, data: list[MealData]) -> pd.DataFrame:
        """Convert meal data to DataFrame."""
        return pd.DataFrame([item.to_dict() for item in data])

    def _create_symptom_dataframe(self, data: list[SymptomData]) -> pd.DataFrame:
        """Convert symptom data to DataFrame."""
        return pd.DataFrame([item.to_dict() for item in data])

    def _assess_bristol_health(self, percentages: pd.Series) -> str:
        """Assess overall digestive health based on Bristol distribution."""
        healthy_range = percentages.get(3, 0) + percentages.get(4, 0)

        if healthy_range > 70:
            return "Excellent - Your bowel movements are consistently healthy"
        elif healthy_range > 50:
            return "Good - Mostly healthy with some room for improvement"
        elif healthy_range > 30:
            return "Fair - Consider dietary adjustments for better digestive health"
        else:
            return "Poor - Significant digestive issues detected, consider medical consultation"

    def _calculate_regularity_score(self, df: pd.DataFrame) -> float:
        """Calculate regularity score based on timing consistency."""
        if len(df) < 7:  # Need at least a week of data
            return 0.5

        # Calculate hour variance
        hour_std = df["hour"].std()
        max_hour_std = 12  # Maximum possible standard deviation for hours
        hour_score = max(0, 1 - (hour_std / max_hour_std))

        # Calculate day-to-day consistency
        daily_counts = df.groupby(df["created_at"].dt.date).size()
        frequency_std = daily_counts.std()
        max_frequency_std = daily_counts.mean()  # Normalize by mean
        frequency_score = (
            max(0, 1 - (frequency_std / max_frequency_std))
            if max_frequency_std > 0
            else 0
        )

        # Combined score
        return (hour_score + frequency_score) / 2

    def _calculate_consistency_score(self, daily_counts: pd.Series) -> float:
        """Calculate consistency score based on daily frequency variance."""
        if len(daily_counts) < 3:
            return 0.5

        mean_frequency = daily_counts.mean()
        std_frequency = daily_counts.std()

        if mean_frequency == 0:
            return 0.0

        # Coefficient of variation (lower is more consistent)
        cv = std_frequency / mean_frequency

        # Convert to 0-1 score (lower CV = higher score)
        # Assume CV > 1.0 is very inconsistent
        consistency_score = max(0, 1 - min(cv, 1.0))

        return consistency_score

    def _calculate_meal_category_correlation(
        self, bm_df: pd.DataFrame, meal_df: pd.DataFrame, category: str
    ) -> dict[str, float] | None:
        """Calculate correlation between a meal category and bowel movements."""
        if meal_df.empty:
            return None

        correlations = {}

        # Look for bowel movements 6-48 hours after meals
        for _, meal in meal_df.iterrows():
            meal_time = meal["meal_time"]
            window_start = meal_time + timedelta(hours=6)
            window_end = meal_time + timedelta(hours=48)

            related_bms = bm_df[
                (bm_df["created_at"] >= window_start)
                & (bm_df["created_at"] <= window_end)
            ]

            if not related_bms.empty:
                # Calculate average metrics
                avg_bristol = related_bms["bristol_type"].mean()
                avg_pain = (
                    related_bms["pain"].dropna().mean()
                    if "pain" in related_bms.columns
                    else None
                )
                avg_satisfaction = (
                    related_bms["satisfaction"].dropna().mean()
                    if "satisfaction" in related_bms.columns
                    else None
                )

                correlations.setdefault("bristol_scores", []).append(avg_bristol)
                if avg_pain is not None:
                    correlations.setdefault("pain_scores", []).append(avg_pain)
                if avg_satisfaction is not None:
                    correlations.setdefault("satisfaction_scores", []).append(
                        avg_satisfaction
                    )

        # Calculate final averages
        final_correlations = {}
        for metric, scores in correlations.items():
            if scores:
                final_correlations[f"avg_{metric}"] = sum(scores) / len(scores)

        return final_correlations if final_correlations else None

    def _identify_trigger_beneficial_foods(
        self, bm_df: pd.DataFrame, meal_df: pd.DataFrame
    ) -> tuple[list[dict], list[dict]]:
        """Identify foods that trigger problems or provide benefits."""
        trigger_foods = []
        beneficial_foods = []

        # Group by food characteristics
        food_groups = {
            "spicy": meal_df[meal_df["spicy_level"].fillna(0) > 5],
            "dairy": meal_df[meal_df["dairy"]],
            "gluten": meal_df[meal_df["gluten"]],
            "fiber_rich": meal_df[meal_df["fiber_rich"]],
        }

        for food_type, food_meals in food_groups.items():
            if food_meals.empty:
                continue

            correlation = self._calculate_meal_category_correlation(
                bm_df, food_meals, food_type
            )
            if correlation:
                avg_bristol = correlation.get("avg_bristol_scores", 4)
                avg_pain = correlation.get("avg_pain_scores", 1)

                # Identify triggers (extreme bristol types or high pain)
                if avg_bristol <= 2 or avg_bristol >= 6 or (avg_pain and avg_pain > 5):
                    trigger_foods.append(
                        {
                            "type": food_type,
                            "avg_bristol": avg_bristol,
                            "avg_pain": avg_pain,
                            "severity": (
                                "high" if avg_pain and avg_pain > 7 else "medium"
                            ),
                        }
                    )

                # Identify beneficial foods (ideal bristol range and low pain)
                elif 3 <= avg_bristol <= 4 and (not avg_pain or avg_pain <= 2):
                    beneficial_foods.append(
                        {
                            "type": food_type,
                            "avg_bristol": avg_bristol,
                            "avg_pain": avg_pain,
                            "benefit_level": (
                                "high" if 3.5 <= avg_bristol <= 4 else "medium"
                            ),
                        }
                    )

        return trigger_foods, beneficial_foods

    def _calculate_meal_timing_correlations(
        self, bm_df: pd.DataFrame, meal_df: pd.DataFrame
    ) -> dict[str, float]:
        """Calculate timing correlations between meals and bowel movements."""
        if meal_df.empty:
            return {}

        timing_correlations = {}

        # Calculate average time between meals and subsequent bowel movements
        time_deltas = []
        for _, meal in meal_df.iterrows():
            meal_time = meal["meal_time"]

            # Find next bowel movement after meal
            future_bms = bm_df[bm_df["created_at"] > meal_time]
            if not future_bms.empty:
                next_bm = future_bms.loc[future_bms["created_at"].idxmin()]
                time_delta = (
                    next_bm["created_at"] - meal_time
                ).total_seconds() / 3600  # hours
                if time_delta <= 72:  # Within 3 days
                    time_deltas.append(time_delta)

        if time_deltas:
            timing_correlations["avg_meal_to_bm_hours"] = sum(time_deltas) / len(
                time_deltas
            )
            timing_correlations["median_meal_to_bm_hours"] = sorted(time_deltas)[
                len(time_deltas) // 2
            ]

        return timing_correlations

    def _calculate_symptom_correlation(
        self, bm_df: pd.DataFrame, symptom_df: pd.DataFrame
    ) -> dict[str, Any] | None:
        """Calculate correlation between symptoms and bowel movements."""
        if symptom_df.empty:
            return None

        correlations = {}

        # Calculate average severity and timing
        avg_severity = symptom_df["severity"].mean()
        correlations["avg_severity"] = avg_severity

        # Find bowel movements around the time of symptoms
        related_bms = []
        for _, symptom in symptom_df.iterrows():
            symptom_time = symptom["created_at"]

            # Look for BMs within 24 hours of symptom
            window_start = symptom_time - timedelta(hours=12)
            window_end = symptom_time + timedelta(hours=12)

            nearby_bms = bm_df[
                (bm_df["created_at"] >= window_start)
                & (bm_df["created_at"] <= window_end)
            ]

            related_bms.extend(nearby_bms.to_dict("records"))

        if related_bms:
            related_df = pd.DataFrame(related_bms)
            correlations["avg_bristol_with_symptoms"] = related_df[
                "bristol_type"
            ].mean()
            if "pain" in related_df.columns:
                correlations["avg_pain_with_symptoms"] = (
                    related_df["pain"].dropna().mean()
                )

        return correlations

    def _analyze_consistency_trends(self, df: pd.DataFrame) -> dict[str, Any]:
        """Analyze consistency trends over time."""
        if "consistency" not in df.columns or df["consistency"].isna().all():
            return {"message": "No consistency data available"}

        consistency_counts = df["consistency"].value_counts()
        return {
            "distribution": consistency_counts.to_dict(),
            "most_common": (
                consistency_counts.idxmax() if not consistency_counts.empty else None
            ),
        }

    def _create_analysis_metadata(
        self,
        user_id: str,
        bowel_movements: list[BowelMovementData],
        meals: list[MealData],
        symptoms: list[SymptomData],
    ) -> AnalysisMetadata:
        """Create analysis metadata."""
        start_time = datetime.now()

        # Calculate data period
        dates = [bm.created_at for bm in bowel_movements]
        data_period_start = min(dates) if dates else start_time
        data_period_end = max(dates) if dates else start_time

        # Processing time (placeholder - would be calculated at the end)
        processing_time = (
            0.1  # This would be calculated properly in real implementation
        )

        return AnalysisMetadata(
            analysis_id=str(uuid.uuid4()),
            user_id=user_id,
            analysis_type="comprehensive_pattern_analysis",
            data_period_start=data_period_start,
            data_period_end=data_period_end,
            total_entries=len(bowel_movements),
            total_meals=len(meals),
            total_symptoms=len(symptoms),
            ml_models_used=[],  # Would be populated if ML models were used
            processing_time_seconds=processing_time,
            cache_hit=False,
            confidence_score=0.8,  # Would be calculated based on data quality
            data_quality_score=self._calculate_data_quality_score(bowel_movements),
        )

    def _calculate_data_quality_score(
        self, bowel_movements: list[BowelMovementData]
    ) -> float:
        """Calculate data quality score based on completeness and consistency."""
        if not bowel_movements:
            return 0.0

        total_fields = len(bowel_movements) * 5  # 5 key fields per entry
        filled_fields = 0

        for bm in bowel_movements:
            # Count non-null key fields
            if bm.bristol_type is not None:
                filled_fields += 1
            if bm.pain is not None:
                filled_fields += 1
            if bm.strain is not None:
                filled_fields += 1
            if bm.satisfaction is not None:
                filled_fields += 1
            if bm.volume is not None:
                filled_fields += 1

        completeness_score = filled_fields / total_fields if total_fields > 0 else 0

        # Adjust for data span (more data over longer period = higher quality)
        if len(bowel_movements) > 1:
            dates = [bm.created_at for bm in bowel_movements]
            date_span_days = (max(dates) - min(dates)).days
            span_bonus = min(
                0.2, date_span_days / 30 * 0.2
            )  # Bonus up to 0.2 for 30+ days
            completeness_score += span_bonus

        return min(1.0, completeness_score)
