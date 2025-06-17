"""Data models for the AI service."""

from .database import (
    AnalysisMetadata,
    BowelMovementData,
    MealData,
    SymptomData,
)
from .requests import (
    AnalysisRequest,
    BowelMovementEntry,
    MealEntry,
    SymptomEntry,
)
from .responses import (
    AnalysisResponse,
    BristolAnalysis,
    Recommendation,
    RiskFactor,
    ErrorResponse,
    FrequencyStats,
    HealthResponse,
    HealthScore,
    TimingPattern,
)

__all__ = [
    "AnalysisMetadata",
    "BowelMovementData",
    "MealData",
    "SymptomData",
    "AnalysisRequest",
    "BowelMovementEntry",
    "MealEntry",
    "SymptomEntry",
    "AnalysisResponse",
    "BristolAnalysis",
    "Recommendation",
    "RiskFactor",
    "ErrorResponse",
    "FrequencyStats",
    "HealthResponse",
    "HealthScore",
    "TimingPattern",
]
