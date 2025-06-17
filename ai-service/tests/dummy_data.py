from __future__ import annotations

from dataclasses import dataclass
from datetime import datetime


@dataclass
class DummyBM:
    """Simple stand-in for bowel movement model."""

    id: str
    user_id: str
    bristol_type: int
    created_at: datetime
    pain: int | None = None
    strain: int | None = None
    satisfaction: int | None = None
    volume: str | None = None

    def to_dict(self) -> dict:
        return self.__dict__


@dataclass
class DummyMeal:
    """Simple stand-in for meal model."""

    id: str
    user_id: str
    meal_time: datetime
    category: str | None = None
    name: str | None = None
    spicy_level: int | None = None
    fiber_rich: bool | None = None
    dairy: bool | None = None
    gluten: bool | None = None
    created_at: datetime | None = None

    def __post_init__(self) -> None:
        if self.created_at is None:
            self.created_at = self.meal_time

    def to_dict(self) -> dict:
        return self.__dict__


@dataclass
class DummySymptom:
    """Simple stand-in for symptom model."""

    id: str
    user_id: str
    type: str
    severity: int
    created_at: datetime
    bowel_movement_id: str | None = None

    def to_dict(self) -> dict:
        return self.__dict__
