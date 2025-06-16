"""Health metrics calculation utilities."""

import math
from typing import Any

import numpy as np
import pandas as pd

from ..config.logging import get_logger

logger = get_logger("health_metrics")


class HealthMetricsCalculator:
    """Calculator for various health metrics and scores."""

    def __init__(self):
        # Bristol type weights for health scoring
        self.bristol_health_weights = {
            1: 0.1,  # Severe constipation
            2: 0.3,  # Mild constipation
            3: 0.9,  # Normal
            4: 1.0,  # Ideal
            5: 0.7,  # Lacking fiber
            6: 0.3,  # Mild diarrhea
            7: 0.1,  # Severe diarrhea
        }

        # Ideal frequency ranges (bowel movements per day)
        self.ideal_frequency_min = 0.5  # Once every 2 days
        self.ideal_frequency_max = 3.0  # Three times per day

    def calculate_bristol_health_score(self, bristol_types: list[int]) -> float:
        """
        Calculate health score based on Bristol stool types.

        Args:
            bristol_types: List of Bristol stool type values (1-7)

        Returns:
            Health score from 0-100
        """
        if not bristol_types:
            return 0.0

        # Calculate weighted score
        weighted_scores = [
            self.bristol_health_weights.get(bt, 0.5) for bt in bristol_types
        ]

        avg_score = sum(weighted_scores) / len(weighted_scores)

        # Convert to 0-100 scale
        return round(avg_score * 100, 2)

    def calculate_frequency_score(self, daily_frequencies: list[float]) -> float:
        """
        Calculate health score based on bowel movement frequency.

        Args:
            daily_frequencies: List of daily frequencies

        Returns:
            Frequency score from 0-100
        """
        if not daily_frequencies:
            return 0.0

        avg_frequency = sum(daily_frequencies) / len(daily_frequencies)

        # Score based on how close to ideal range
        if self.ideal_frequency_min <= avg_frequency <= self.ideal_frequency_max:
            # Perfect frequency
            base_score = 100
        elif avg_frequency < self.ideal_frequency_min:
            # Too infrequent (constipation)
            base_score = max(0, 100 * (avg_frequency / self.ideal_frequency_min))
        else:
            # Too frequent (possible diarrhea)
            excess_ratio = (
                avg_frequency - self.ideal_frequency_max
            ) / self.ideal_frequency_max
            base_score = max(0, 100 * (1 - excess_ratio))

        # Penalty for inconsistency
        if len(daily_frequencies) > 1:
            consistency_penalty = self._calculate_consistency_penalty(daily_frequencies)
            base_score *= 1 - consistency_penalty

        return round(base_score, 2)

    def calculate_pain_score(self, pain_scores: list[int]) -> float:
        """
        Calculate health score based on pain levels.

        Args:
            pain_scores: List of pain scores (1-10)

        Returns:
            Pain score from 0-100 (higher = less pain = better)
        """
        if not pain_scores:
            return 100.0  # No pain data = assume no pain

        avg_pain = sum(pain_scores) / len(pain_scores)

        # Invert pain score (lower pain = higher health score)
        pain_score = max(0, 100 * (10 - avg_pain) / 9)

        # Penalty for high variance (inconsistent pain)
        if len(pain_scores) > 1:
            pain_std = np.std(pain_scores)
            variance_penalty = min(0.3, pain_std / 10)  # Max 30% penalty
            pain_score *= 1 - variance_penalty

        return round(pain_score, 2)

    def calculate_satisfaction_score(self, satisfaction_scores: list[int]) -> float:
        """
        Calculate score based on satisfaction levels.

        Args:
            satisfaction_scores: List of satisfaction scores (1-10)

        Returns:
            Satisfaction score from 0-100
        """
        if not satisfaction_scores:
            return 50.0  # Neutral if no data

        avg_satisfaction = sum(satisfaction_scores) / len(satisfaction_scores)

        # Convert to 0-100 scale
        satisfaction_score = (avg_satisfaction - 1) * 100 / 9

        return round(max(0, min(100, satisfaction_score)), 2)

    def calculate_overall_health_score(
        self,
        bristol_types: list[int],
        daily_frequencies: list[float],
        pain_scores: list[int] | None = None,
        satisfaction_scores: list[int] | None = None,
        weights: dict[str, float] | None = None,
    ) -> dict[str, float]:
        """
        Calculate overall digestive health score.

        Args:
            bristol_types: Bristol stool types
            daily_frequencies: Daily bowel movement frequencies
            pain_scores: Pain scores (optional)
            satisfaction_scores: Satisfaction scores (optional)
            weights: Custom weights for different components

        Returns:
            Dictionary with component scores and overall score
        """
        # Default weights
        default_weights = {
            "bristol": 0.4,
            "frequency": 0.3,
            "pain": 0.2,
            "satisfaction": 0.1,
        }

        if weights:
            default_weights.update(weights)

        # Calculate component scores
        bristol_score = self.calculate_bristol_health_score(bristol_types)
        frequency_score = self.calculate_frequency_score(daily_frequencies)
        pain_score = self.calculate_pain_score(pain_scores or [])
        satisfaction_score = self.calculate_satisfaction_score(
            satisfaction_scores or []
        )

        # Calculate weighted overall score
        overall_score = (
            bristol_score * default_weights["bristol"]
            + frequency_score * default_weights["frequency"]
            + pain_score * default_weights["pain"]
            + satisfaction_score * default_weights["satisfaction"]
        )

        return {
            "overall_score": round(overall_score, 2),
            "bristol_score": bristol_score,
            "frequency_score": frequency_score,
            "pain_score": pain_score,
            "satisfaction_score": satisfaction_score,
        }

    def calculate_trend_score(
        self, recent_scores: list[float], historical_scores: list[float]
    ) -> dict[str, Any]:
        """
        Calculate trend in health scores over time.

        Args:
            recent_scores: Recent scores (e.g., last 30 days)
            historical_scores: Historical scores for comparison

        Returns:
            Trend analysis with direction and magnitude
        """
        if not recent_scores or not historical_scores:
            return {"trend": "stable", "change": 0.0, "confidence": 0.0}

        recent_avg = sum(recent_scores) / len(recent_scores)
        historical_avg = sum(historical_scores) / len(historical_scores)

        change = recent_avg - historical_avg
        change_percentage = (change / historical_avg) * 100 if historical_avg > 0 else 0

        # Determine trend direction
        if abs(change_percentage) < 5:
            trend = "stable"
        elif change_percentage > 0:
            trend = "improving"
        else:
            trend = "declining"

        # Calculate confidence based on data consistency
        confidence = self._calculate_trend_confidence(recent_scores, historical_scores)

        return {
            "trend": trend,
            "change": round(change, 2),
            "change_percentage": round(change_percentage, 2),
            "confidence": round(confidence, 2),
        }

    def calculate_risk_factors(
        self,
        bristol_types: list[int],
        pain_scores: list[int] | None = None,
        frequencies: list[float] | None = None,
    ) -> list[dict[str, Any]]:
        """
        Identify potential risk factors based on patterns.

        Args:
            bristol_types: Bristol stool types
            pain_scores: Pain scores
            frequencies: Daily frequencies

        Returns:
            List of identified risk factors
        """
        risk_factors = []

        if not bristol_types:
            return risk_factors

        # Check for extreme Bristol types
        extreme_count = sum(1 for bt in bristol_types if bt in [1, 2, 6, 7])
        extreme_ratio = extreme_count / len(bristol_types)

        if extreme_ratio > 0.3:
            risk_factors.append(
                {
                    "factor": "extreme_bristol_types",
                    "severity": "high" if extreme_ratio > 0.5 else "medium",
                    "description": f"{extreme_ratio:.1%} of movements are extreme types (constipation or diarrhea)",
                    "prevalence": extreme_ratio,
                }
            )

        # Check for chronic constipation
        constipation_count = sum(1 for bt in bristol_types if bt <= 2)
        constipation_ratio = constipation_count / len(bristol_types)

        if constipation_ratio > 0.4:
            risk_factors.append(
                {
                    "factor": "chronic_constipation",
                    "severity": "high" if constipation_ratio > 0.6 else "medium",
                    "description": f"{constipation_ratio:.1%} of movements indicate constipation",
                    "prevalence": constipation_ratio,
                }
            )

        # Check for chronic diarrhea
        diarrhea_count = sum(1 for bt in bristol_types if bt >= 6)
        diarrhea_ratio = diarrhea_count / len(bristol_types)

        if diarrhea_ratio > 0.2:
            risk_factors.append(
                {
                    "factor": "chronic_diarrhea",
                    "severity": "high" if diarrhea_ratio > 0.4 else "medium",
                    "description": f"{diarrhea_ratio:.1%} of movements indicate diarrhea",
                    "prevalence": diarrhea_ratio,
                }
            )

        # Check pain levels
        if pain_scores:
            high_pain_count = sum(1 for pain in pain_scores if pain > 7)
            high_pain_ratio = high_pain_count / len(pain_scores)

            if high_pain_ratio > 0.2:
                risk_factors.append(
                    {
                        "factor": "frequent_high_pain",
                        "severity": "high" if high_pain_ratio > 0.4 else "medium",
                        "description": f"{high_pain_ratio:.1%} of movements involve high pain (>7/10)",
                        "prevalence": high_pain_ratio,
                    }
                )

        # Check frequency patterns
        if frequencies:
            avg_frequency = sum(frequencies) / len(frequencies)

            if avg_frequency < 0.3:  # Less than once every 3 days
                risk_factors.append(
                    {
                        "factor": "severe_constipation_frequency",
                        "severity": "high",
                        "description": f"Very low frequency: {avg_frequency:.1f} movements per day",
                        "prevalence": 1.0,
                    }
                )
            elif avg_frequency > 5:  # More than 5 times per day
                risk_factors.append(
                    {
                        "factor": "excessive_frequency",
                        "severity": "high",
                        "description": f"Very high frequency: {avg_frequency:.1f} movements per day",
                        "prevalence": 1.0,
                    }
                )

        return risk_factors

    def calculate_bmi_health_impact(
        self, weight_kg: float, height_cm: float
    ) -> dict[str, Any]:
        """
        Calculate BMI and its potential impact on digestive health.

        Args:
            weight_kg: Weight in kilograms
            height_cm: Height in centimeters

        Returns:
            BMI analysis and health implications
        """
        if weight_kg <= 0 or height_cm <= 0:
            return {"error": "Invalid weight or height values"}

        height_m = height_cm / 100
        bmi = weight_kg / (height_m**2)

        # BMI categories
        if bmi < 18.5:
            category = "underweight"
            digestive_impact = (
                "May increase risk of constipation due to reduced dietary intake"
            )
        elif 18.5 <= bmi < 25:
            category = "normal"
            digestive_impact = "Optimal BMI range for digestive health"
        elif 25 <= bmi < 30:
            category = "overweight"
            digestive_impact = "May increase risk of acid reflux and slower digestion"
        else:
            category = "obese"
            digestive_impact = (
                "Higher risk of digestive issues, acid reflux, and constipation"
            )

        return {
            "bmi": round(bmi, 1),
            "category": category,
            "digestive_impact": digestive_impact,
            "health_risk": "low"
            if 18.5 <= bmi < 25
            else "medium"
            if 25 <= bmi < 30
            else "high",
        }

    def _calculate_consistency_penalty(self, values: list[float]) -> float:
        """Calculate penalty for inconsistent values."""
        if len(values) < 2:
            return 0.0

        std_dev = np.std(values)
        mean_val = np.mean(values)

        if mean_val == 0:
            return 0.3  # High penalty for zero mean

        # Coefficient of variation
        cv = std_dev / mean_val

        # Convert to penalty (0 to 0.3 max)
        penalty = min(0.3, cv * 0.5)

        return penalty

    def _calculate_trend_confidence(
        self, recent_scores: list[float], historical_scores: list[float]
    ) -> float:
        """Calculate confidence in trend analysis."""
        # More data points = higher confidence
        data_confidence = min(1.0, (len(recent_scores) + len(historical_scores)) / 20)

        # Lower variance = higher confidence
        all_scores = recent_scores + historical_scores
        if len(all_scores) > 1:
            variance_confidence = max(0.3, 1 - (np.std(all_scores) / 100))
        else:
            variance_confidence = 0.5

        # Combine confidences
        overall_confidence = (data_confidence + variance_confidence) / 2

        return overall_confidence

    def calculate_digestive_regularity_index(
        self, timing_data: list[datetime], bristol_data: list[int]
    ) -> dict[str, Any]:
        """
        Calculate a comprehensive digestive regularity index.

        Args:
            timing_data: List of bowel movement timestamps
            bristol_data: List of Bristol stool types

        Returns:
            Regularity index and components
        """
        if not timing_data or not bristol_data or len(timing_data) != len(bristol_data):
            return {"regularity_index": 0.0, "components": {}}

        # Time regularity
        time_regularity = self._calculate_time_regularity(timing_data)

        # Bristol consistency
        bristol_consistency = self._calculate_bristol_consistency(bristol_data)

        # Frequency stability
        frequency_stability = self._calculate_frequency_stability(timing_data)

        # Overall index (weighted average)
        regularity_index = (
            time_regularity * 0.4
            + bristol_consistency * 0.3
            + frequency_stability * 0.3
        )

        return {
            "regularity_index": round(regularity_index, 2),
            "components": {
                "time_regularity": round(time_regularity, 2),
                "bristol_consistency": round(bristol_consistency, 2),
                "frequency_stability": round(frequency_stability, 2),
            },
            "interpretation": self._interpret_regularity_index(regularity_index),
        }

    def _calculate_time_regularity(self, timing_data: list[datetime]) -> float:
        """Calculate regularity based on timing patterns."""
        if len(timing_data) < 3:
            return 0.5

        # Calculate hour of day for each movement
        hours = [dt.hour for dt in timing_data]

        # Calculate standard deviation of hours
        hour_std = np.std(hours)

        # Convert to 0-1 score (lower std = higher regularity)
        max_std = 12  # Maximum possible hour standard deviation
        regularity = max(0, 1 - (hour_std / max_std))

        return regularity

    def _calculate_bristol_consistency(self, bristol_data: list[int]) -> float:
        """Calculate consistency of Bristol stool types."""
        if not bristol_data:
            return 0.0

        # Calculate standard deviation
        bristol_std = np.std(bristol_data)

        # Convert to 0-1 score (lower std = higher consistency)
        max_std = 3  # Reasonable max standard deviation for Bristol types
        consistency = max(0, 1 - (bristol_std / max_std))

        return consistency

    def _calculate_frequency_stability(self, timing_data: list[datetime]) -> float:
        """Calculate stability of bowel movement frequency."""
        if len(timing_data) < 7:  # Need at least a week of data
            return 0.5

        # Group by day and count
        daily_counts = {}
        for dt in timing_data:
            date_key = dt.date()
            daily_counts[date_key] = daily_counts.get(date_key, 0) + 1

        # Calculate coefficient of variation
        counts = list(daily_counts.values())
        if len(counts) < 2:
            return 0.5

        mean_count = np.mean(counts)
        std_count = np.std(counts)

        if mean_count == 0:
            return 0.0

        cv = std_count / mean_count

        # Convert to 0-1 score (lower CV = higher stability)
        stability = max(0, 1 - min(cv, 1))

        return stability

    def _interpret_regularity_index(self, index: float) -> str:
        """Interpret regularity index score."""
        if index >= 0.8:
            return "Excellent digestive regularity"
        elif index >= 0.6:
            return "Good digestive regularity"
        elif index >= 0.4:
            return "Moderate digestive regularity"
        elif index >= 0.2:
            return "Poor digestive regularity"
        else:
            return "Very irregular digestive patterns"
