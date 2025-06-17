from __future__ import annotations

from datetime import datetime

from pydantic import BaseModel


class BowelMovementData(BaseModel):
    """Internal representation of a bowel movement entry."""

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


class MealData(BaseModel):
    """Internal representation of a meal entry."""

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
    created_at: datetime


class SymptomData(BaseModel):
    """Internal representation of a symptom entry."""

    id: str
    user_id: str
    bowel_movement_id: str | None = None
    type: str
    severity: int
    notes: str | None = None
    created_at: datetime
    recorded_at: datetime | None = None


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
    ml_models_used: list[str] = []
    processing_time_seconds: float
    cache_hit: bool = False
    confidence_score: float = 0.0
    data_quality_score: float = 0.0
