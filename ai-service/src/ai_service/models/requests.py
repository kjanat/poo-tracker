from __future__ import annotations

from datetime import datetime

from pydantic import BaseModel, Field


class BowelMovementEntry(BaseModel):
    """Request model for a bowel movement entry."""

    id: str
    user_id: str = Field(alias="userId")
    bristol_type: int = Field(alias="bristolType")
    volume: str | None = None
    color: str | None = None
    consistency: str | None = None
    floaters: bool | None = None
    pain: int | None = None
    strain: int | None = None
    satisfaction: int | None = None
    created_at: datetime = Field(alias="createdAt")
    recorded_at: datetime | None = Field(default=None, alias="recordedAt")


class MealEntry(BaseModel):
    """Request model for a meal entry."""

    id: str
    user_id: str = Field(alias="userId")
    name: str | None = None
    meal_time: datetime = Field(alias="mealTime")
    category: str | None = None
    cuisine: str | None = None
    spicy_level: int | None = Field(default=None, alias="spicyLevel")
    fiber_rich: bool | None = Field(default=None, alias="fiberRich")
    dairy: bool | None = None
    gluten: bool | None = None
    created_at: datetime = Field(alias="createdAt")


class SymptomEntry(BaseModel):
    """Request model for a symptom entry."""

    id: str
    user_id: str = Field(alias="userId")
    bowel_movement_id: str | None = Field(default=None, alias="bowelMovementId")
    type: str
    severity: int
    notes: str | None = None
    created_at: datetime = Field(alias="createdAt")
    recorded_at: datetime | None = Field(default=None, alias="recordedAt")


class AnalysisRequest(BaseModel):
    """Request payload for the analysis endpoint."""

    entries: list[BowelMovementEntry]
    meals: list[MealEntry] | None = None
    symptoms: list[SymptomEntry] | None = None
    include_predictions: bool = Field(False, alias="includePredictions")
    include_recommendations: bool = Field(False, alias="includeRecommendations")

    class Config:
        allow_population_by_field_name = True
