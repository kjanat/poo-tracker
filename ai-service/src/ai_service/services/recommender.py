"""Recommendation engine for generating personalized health advice."""

import uuid
from typing import Any

from ..config.logging import get_logger
from ..models.database import BowelMovementData, MealData
from ..models.responses import Recommendation, RiskFactor

logger = get_logger("recommender")


class RecommenderService:
    """Service for generating personalized recommendations and identifying risk factors."""

    def __init__(self):
        self.bristol_recommendations = {
            1: {
                "dietary": [
                    "Increase fiber intake with fruits and vegetables",
                    "Drink more water throughout the day",
                    "Add prunes or other natural laxatives",
                ],
                "lifestyle": [
                    "Increase physical activity",
                    "Establish regular bathroom routine",
                    "Reduce stress levels",
                ],
                "medical": [
                    "Consider consulting a healthcare provider for chronic constipation"
                ],
            },
            2: {
                "dietary": [
                    "Gradually increase fiber intake",
                    "Ensure adequate hydration",
                    "Include whole grains in diet",
                ],
                "lifestyle": [
                    "Regular exercise can help",
                    "Don't delay when you feel the urge",
                ],
                "medical": ["Monitor symptoms and consult doctor if persistent"],
            },
            6: {
                "dietary": [
                    "Reduce high-fat foods",
                    "Limit dairy if lactose intolerant",
                    "Avoid spicy foods temporarily",
                ],
                "lifestyle": ["Stay hydrated with electrolytes", "Rest when possible"],
                "medical": ["Track symptoms and see doctor if persists >3 days"],
            },
            7: {
                "dietary": [
                    "Follow BRAT diet (bananas, rice, applesauce, toast)",
                    "Avoid dairy and fatty foods",
                    "Small frequent meals",
                ],
                "lifestyle": [
                    "Stay very hydrated",
                    "Get plenty of rest",
                    "Avoid strenuous activity",
                ],
                "medical": ["Seek medical attention if severe or persistent"],
            },
        }

    async def generate_recommendations(
        self,
        analysis_result: dict[str, Any],
        bowel_movements: list[BowelMovementData],
        meals: list[MealData] | None = None,
    ) -> list[Recommendation]:
        """
        Generate personalized recommendations based on analysis.

        Args:
            analysis_result: Analysis results from analyzer service
            bowel_movements: Bowel movement data
            meals: Optional meal data

        Returns:
            List of personalized recommendations
        """
        logger.info("Generating personalized recommendations")

        recommendations = []

        # Bristol-based recommendations
        bristol_recs = await self._generate_bristol_recommendations(bowel_movements)
        recommendations.extend(bristol_recs)

        # Frequency-based recommendations
        frequency_recs = await self._generate_frequency_recommendations(analysis_result)
        recommendations.extend(frequency_recs)

        # Pain-based recommendations
        pain_recs = await self._generate_pain_recommendations(bowel_movements)
        recommendations.extend(pain_recs)

        # Meal correlation recommendations
        if meals:
            meal_recs = await self._generate_meal_recommendations(
                analysis_result, meals
            )
            recommendations.extend(meal_recs)

        # Timing recommendations
        timing_recs = await self._generate_timing_recommendations(analysis_result)
        recommendations.extend(timing_recs)

        # Sort by priority and limit to top recommendations
        recommendations.sort(
            key=lambda x: self._get_priority_weight(x.priority), reverse=True
        )
        return recommendations[:10]  # Return top 10 recommendations

    async def identify_risk_factors(
        self,
        bowel_movements: list[BowelMovementData],
        analysis_result: dict[str, Any],
    ) -> list[RiskFactor]:
        """
        Identify potential health risk factors.

        Args:
            bowel_movements: Bowel movement data
            analysis_result: Analysis results

        Returns:
            List of identified risk factors
        """
        logger.info("Identifying health risk factors")

        risk_factors = []

        # Bristol type risk factors
        bristol_risks = await self._identify_bristol_risks(bowel_movements)
        risk_factors.extend(bristol_risks)

        # Frequency risk factors
        frequency_risks = await self._identify_frequency_risks(analysis_result)
        risk_factors.extend(frequency_risks)

        # Pain risk factors
        pain_risks = await self._identify_pain_risks(bowel_movements)
        risk_factors.extend(pain_risks)

        # Pattern risk factors
        pattern_risks = await self._identify_pattern_risks(analysis_result)
        risk_factors.extend(pattern_risks)

        return risk_factors

    async def _generate_bristol_recommendations(
        self, bowel_movements: list[BowelMovementData]
    ) -> list[Recommendation]:
        """Generate recommendations based on Bristol stool types."""
        recommendations = []

        if not bowel_movements:
            return recommendations

        bristol_types = [bm.bristol_type for bm in bowel_movements]
        avg_bristol = sum(bristol_types) / len(bristol_types)

        # Recommendations for constipation (types 1-2)
        constipation_ratio = sum(1 for bt in bristol_types if bt <= 2) / len(
            bristol_types
        )
        if constipation_ratio > 0.3:
            recommendations.append(
                Recommendation(
                    id=str(uuid.uuid4()),
                    category="diet",
                    title="Increase Fiber Intake",
                    description="Add more fruits, vegetables, and whole grains to help with constipation",
                    priority="high" if constipation_ratio > 0.6 else "medium",
                    confidence=0.9,
                    evidence=[
                        f"{constipation_ratio:.1%} of movements indicate constipation"
                    ],
                )
            )

        # Recommendations for diarrhea (types 6-7)
        diarrhea_ratio = sum(1 for bt in bristol_types if bt >= 6) / len(bristol_types)
        if diarrhea_ratio > 0.2:
            recommendations.append(
                Recommendation(
                    id=str(uuid.uuid4()),
                    category="diet",
                    title="Follow BRAT Diet",
                    description="Eat bananas, rice, applesauce, and toast to firm up stools",
                    priority="high" if diarrhea_ratio > 0.5 else "medium",
                    confidence=0.8,
                    evidence=[
                        f"{diarrhea_ratio:.1%} of movements indicate loose stools"
                    ],
                )
            )

        return recommendations

    async def _generate_frequency_recommendations(
        self, analysis_result: dict[str, Any]
    ) -> list[Recommendation]:
        """Generate recommendations based on frequency patterns."""
        recommendations = []

        patterns = analysis_result.get("patterns", {})
        frequency = patterns.get("frequency", {})

        avg_daily = frequency.get("avg_daily", 1.0)

        if avg_daily < 0.5:  # Less than once every 2 days
            recommendations.append(
                Recommendation(
                    id=str(uuid.uuid4()),
                    category="lifestyle",
                    title="Establish Regular Bathroom Routine",
                    description="Try to use the bathroom at the same time each day, especially after meals",
                    priority="medium",
                    confidence=0.7,
                    evidence=[f"Average frequency: {avg_daily:.1f} times per day"],
                )
            )

        elif avg_daily > 3:  # More than 3 times per day
            recommendations.append(
                Recommendation(
                    id=str(uuid.uuid4()),
                    category="medical",
                    title="Monitor Frequent Bowel Movements",
                    description="Consider tracking triggers and consulting healthcare provider if persistent",
                    priority="medium",
                    confidence=0.6,
                    evidence=[f"High frequency: {avg_daily:.1f} times per day"],
                )
            )

        return recommendations

    async def _generate_pain_recommendations(
        self, bowel_movements: list[BowelMovementData]
    ) -> list[Recommendation]:
        """Generate recommendations based on pain patterns."""
        recommendations = []

        pain_scores = [bm.pain for bm in bowel_movements if bm.pain is not None]

        if not pain_scores:
            return recommendations

        avg_pain = sum(pain_scores) / len(pain_scores)
        high_pain_ratio = sum(1 for pain in pain_scores if pain > 7) / len(pain_scores)

        if avg_pain > 5:
            recommendations.append(
                Recommendation(
                    id=str(uuid.uuid4()),
                    category="medical",
                    title="Address Bowel Movement Pain",
                    description="Consider consulting healthcare provider about persistent pain during bowel movements",
                    priority="high" if avg_pain > 7 else "medium",
                    confidence=0.8,
                    evidence=[f"Average pain level: {avg_pain:.1f}/10"],
                )
            )

        if high_pain_ratio > 0.3:
            recommendations.append(
                Recommendation(
                    id=str(uuid.uuid4()),
                    category="lifestyle",
                    title="Stress Reduction Techniques",
                    description="Practice relaxation techniques as stress can contribute to digestive discomfort",
                    priority="medium",
                    confidence=0.6,
                    evidence=[f"{high_pain_ratio:.1%} of movements involve high pain"],
                )
            )

        return recommendations

    async def _generate_meal_recommendations(
        self, analysis_result: dict[str, Any], meals: list[MealData]
    ) -> list[Recommendation]:
        """Generate recommendations based on meal correlations."""
        recommendations = []

        correlations = analysis_result.get("correlations", {}).get("meals", {})
        trigger_foods = correlations.get("trigger_foods", [])
        beneficial_foods = correlations.get("beneficial_foods", [])

        # Trigger food recommendations
        for trigger in trigger_foods:
            if trigger.get("severity") == "high":
                recommendations.append(
                    Recommendation(
                        id=str(uuid.uuid4()),
                        category="diet",
                        title=f"Limit {trigger['type'].title()} Foods",
                        description=f"Consider reducing {trigger['type']} foods as they may trigger digestive issues",
                        priority="medium",
                        confidence=0.7,
                        evidence=[
                            f"Associated with Bristol type {trigger.get('avg_bristol', 'unknown')}"
                        ],
                    )
                )

        # Beneficial food recommendations
        for beneficial in beneficial_foods:
            recommendations.append(
                Recommendation(
                    id=str(uuid.uuid4()),
                    category="diet",
                    title=f"Include More {beneficial['type'].title()} Foods",
                    description=f"Consider adding more {beneficial['type']} foods to your diet",
                    priority="low",
                    confidence=0.6,
                    evidence=["Associated with healthy Bristol types"],
                )
            )

        return recommendations

    async def _generate_timing_recommendations(
        self, analysis_result: dict[str, Any]
    ) -> list[Recommendation]:
        """Generate recommendations based on timing patterns."""
        recommendations = []

        patterns = analysis_result.get("patterns", {})
        timing = patterns.get("timing", {})

        regularity_score = timing.get("regularity_score", 0.5)

        if regularity_score < 0.3:
            recommendations.append(
                Recommendation(
                    id=str(uuid.uuid4()),
                    category="lifestyle",
                    title="Improve Schedule Regularity",
                    description="Try to maintain consistent meal times and sleep schedule to improve digestive regularity",
                    priority="medium",
                    confidence=0.6,
                    evidence=[f"Regularity score: {regularity_score:.1f}"],
                )
            )

        return recommendations

    async def _identify_bristol_risks(
        self, bowel_movements: list[BowelMovementData]
    ) -> list[RiskFactor]:
        """Identify risk factors based on Bristol patterns."""
        risk_factors = []

        if not bowel_movements:
            return risk_factors

        bristol_types = [bm.bristol_type for bm in bowel_movements]

        # Chronic constipation risk
        constipation_ratio = sum(1 for bt in bristol_types if bt <= 2) / len(
            bristol_types
        )
        if constipation_ratio > 0.5:
            risk_factors.append(
                RiskFactor(
                    factor="chronic_constipation",
                    severity="high" if constipation_ratio > 0.7 else "medium",
                    description=f"Chronic constipation pattern detected ({constipation_ratio:.1%} of movements)",
                    prevalence=constipation_ratio,
                    recommendation="Increase fiber, water intake, and consider medical consultation",
                )
            )

        # Chronic diarrhea risk
        diarrhea_ratio = sum(1 for bt in bristol_types if bt >= 6) / len(bristol_types)
        if diarrhea_ratio > 0.3:
            risk_factors.append(
                RiskFactor(
                    factor="chronic_diarrhea",
                    severity="high" if diarrhea_ratio > 0.5 else "medium",
                    description=f"Chronic loose stool pattern detected ({diarrhea_ratio:.1%} of movements)",
                    prevalence=diarrhea_ratio,
                    recommendation="Monitor hydration and consider medical evaluation",
                )
            )

        return risk_factors

    async def _identify_frequency_risks(
        self, analysis_result: dict[str, Any]
    ) -> list[RiskFactor]:
        """Identify risk factors based on frequency patterns."""
        risk_factors = []

        patterns = analysis_result.get("patterns", {})
        frequency = patterns.get("frequency", {})

        avg_daily = frequency.get("avg_daily", 1.0)

        if avg_daily < 0.3:
            risk_factors.append(
                RiskFactor(
                    factor="severe_constipation",
                    severity="high",
                    description=f"Very infrequent bowel movements ({avg_daily:.1f} per day)",
                    prevalence=1.0,
                    recommendation="Immediate dietary and lifestyle changes, consider medical consultation",
                )
            )

        elif avg_daily > 5:
            risk_factors.append(
                RiskFactor(
                    factor="excessive_frequency",
                    severity="medium",
                    description=f"Very frequent bowel movements ({avg_daily:.1f} per day)",
                    prevalence=1.0,
                    recommendation="Monitor for dehydration and consider medical evaluation",
                )
            )

        return risk_factors

    async def _identify_pain_risks(
        self, bowel_movements: list[BowelMovementData]
    ) -> list[RiskFactor]:
        """Identify risk factors based on pain patterns."""
        risk_factors = []

        pain_scores = [bm.pain for bm in bowel_movements if bm.pain is not None]

        if not pain_scores:
            return risk_factors

        high_pain_ratio = sum(1 for pain in pain_scores if pain > 7) / len(pain_scores)

        if high_pain_ratio > 0.3:
            risk_factors.append(
                RiskFactor(
                    factor="frequent_severe_pain",
                    severity="high",
                    description=f"Frequent severe pain during bowel movements ({high_pain_ratio:.1%})",
                    prevalence=high_pain_ratio,
                    recommendation="Medical consultation recommended to rule out underlying conditions",
                )
            )

        return risk_factors

    async def _identify_pattern_risks(
        self, analysis_result: dict[str, Any]
    ) -> list[RiskFactor]:
        """Identify risk factors based on overall patterns."""
        risk_factors = []

        # Check for highly irregular patterns
        patterns = analysis_result.get("patterns", {})
        timing = patterns.get("timing", {})

        regularity_score = timing.get("regularity_score", 0.5)

        if regularity_score < 0.2:
            risk_factors.append(
                RiskFactor(
                    factor="highly_irregular_pattern",
                    severity="medium",
                    description="Highly irregular bowel movement patterns detected",
                    prevalence=1.0,
                    recommendation="Focus on establishing regular routine and reducing stress",
                )
            )

        return risk_factors

    def _get_priority_weight(self, priority: str) -> int:
        """Get numeric weight for priority sorting."""
        weights = {"urgent": 4, "high": 3, "medium": 2, "low": 1}
        return weights.get(priority, 1)
