"""Internal data models used by service layer."""

from __future__ import annotations

from dataclasses import asdict, dataclass
from datetime import datetime
from typing import Any


@dataclass
class BowelMovementData:
    """Representation of a bowel movement entry."""

    id: str
    created_at: datetime
    bristol_type: int
    pain: int | None = None
    strain: int | None = None
    satisfaction: int | None = None
    volume: str | None = None
    consistency: str | None = None
    color: str | None = None

    def to_dict(self) -> dict[str, Any]:
        return asdict(self)


@dataclass
class MealData:
    """Representation of a meal entry."""

    id: str
    meal_time: datetime
    name: str | None = None
    category: str | None = None
    spicy_level: int | None = None
    dairy: bool = False
    gluten: bool = False
    fiber_rich: bool = False

    def to_dict(self) -> dict[str, Any]:
        return asdict(self)


@dataclass
class SymptomData:
    """Representation of a symptom entry."""

    id: str
    created_at: datetime
    type: str
    severity: int

    def to_dict(self) -> dict[str, Any]:
        return asdict(self)


@dataclass
class AnalysisMetadata:
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
        return asdict(self)
