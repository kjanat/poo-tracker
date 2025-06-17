"""Internal data models used by the service logic."""

from datetime import datetime
from typing import Any

from pydantic import BaseModel


class BowelMovementData(BaseModel):
    """Normalized bowel movement record."""

    id: str
    user_id: str
    bristol_type: int
    volume: str | None = None
    color: str | None = None
    consistency: str | None = None
    floaters: bool | None = None
    pain: int | None = None
    strain: int | None = None
    satisfaction: int | None = None
    created_at: datetime
    recorded_at: datetime | None = None

    def to_dict(self) -> dict[str, Any]:
        """Return a plain dictionary representation."""
        return self.model_dump()


class MealData(BaseModel):
    """Meal record used for correlation analysis."""

    id: str
    user_id: str
    name: str | None = None
    meal_time: datetime
    category: str | None = None
    cuisine: str | None = None
    spicy_level: int | None = None
    fiber_rich: bool | None = None
    dairy: bool | None = None
    gluten: bool | None = None
    created_at: datetime | None = None

    def to_dict(self) -> dict[str, Any]:
        return self.model_dump()


class SymptomData(BaseModel):
    """Symptom record linked to a bowel movement."""

    id: str
    user_id: str
    bowel_movement_id: str | None = None
    type: str
    severity: int
    notes: str | None = None
    created_at: datetime
    recorded_at: datetime | None = None

    def to_dict(self) -> dict[str, Any]:
        return self.model_dump()


class AnalysisMetadata(BaseModel):
    """Metadata describing an analysis run."""

    analysis_id: str
    user_id: str
    analysis_type: str
    data_period_start: datetime
    data_period_end: datetime
    total_entries: int
    total_meals: int
    total_symptoms: int
    ml_models_used: list[str]
    processing_time_seconds: float
    cache_hit: bool
    confidence_score: float
    data_quality_score: float

    def to_dict(self) -> dict[str, Any]:
        return self.model_dump()
