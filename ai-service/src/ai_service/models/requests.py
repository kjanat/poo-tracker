"""Pydantic models for incoming API requests."""

from datetime import datetime
from typing import List

from pydantic import BaseModel, Field


class BowelMovementEntry(BaseModel):
    """Bowel movement entry sent by the client."""

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
    """Meal entry included with the analysis request."""

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
    created_at: datetime | None = Field(default=None, alias="createdAt")


class SymptomEntry(BaseModel):
    """Symptom entry linked to a bowel movement."""

    id: str
    user_id: str = Field(alias="userId")
    bowel_movement_id: str | None = Field(default=None, alias="bowelMovementId")
    type: str
    severity: int
    notes: str | None = None
    created_at: datetime = Field(alias="createdAt")
    recorded_at: datetime | None = Field(default=None, alias="recordedAt")


class AnalysisRequest(BaseModel):
    """Top-level analysis request payload."""

    entries: List[BowelMovementEntry]
    meals: List[MealEntry] | None = None
    symptoms: List[SymptomEntry] | None = None
    include_predictions: bool = False
    include_recommendations: bool = False

    model_config = {
        "populate_by_name": True,
        "extra": "ignore",
    }
