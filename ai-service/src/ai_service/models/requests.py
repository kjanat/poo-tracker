"""Request models used for API input."""

from __future__ import annotations

from datetime import datetime

from pydantic import BaseModel, Field


class BowelMovementEntry(BaseModel):
    """Bowel movement entry submitted by a user."""

    id: str = Field(..., alias="id")
    user_id: str = Field(..., alias="userId")
    bristol_type: int = Field(..., alias="bristolType")
    created_at: datetime = Field(..., alias="createdAt")

    model_config = {"populate_by_name": True}


class AnalysisRequest(BaseModel):
    """Top-level analysis request."""

    entries: list[BowelMovementEntry]
    meals: list | None = None
    symptoms: list | None = None

    model_config = {"populate_by_name": True}
