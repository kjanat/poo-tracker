from __future__ import annotations

from datetime import datetime
from typing import Any

from pydantic import BaseModel, Field


class BristolAnalysis(BaseModel):
    distribution: dict[int, int]
    percentages: dict[int, float]
    most_common: dict[str, Any]
    health_indicator: str | None = None
    trend: str | None = None


class FrequencyStats(BaseModel):
    avg_daily: float
    max_daily: int
    min_daily: int
    total_days: int
    total_entries: int
    consistency_score: float


class TimingPattern(BaseModel):
    hourly_distribution: dict[int, int]
    daily_distribution: dict[str, int]
    peak_hour: int
    most_active_day: str
    regularity_score: float


class Recommendation(BaseModel):
    id: str
    category: str
    title: str
    description: str
    priority: str
    confidence: float
    evidence: list[str]


class RiskFactor(BaseModel):
    factor: str
    severity: str
    description: str
    prevalence: float
    recommendation: str


class HealthScore(BaseModel):
    overall_score: float
    bristol_score: float
    frequency_score: float
    pain_score: float
    satisfaction_score: float
    trend: str


class HealthResponse(BaseModel):
    status: str
    timestamp: str = Field(default_factory=lambda: datetime.now().isoformat())
    redis_connected: bool
    ml_models_loaded: bool
    response_time_ms: float
    version: str


class ErrorResponse(BaseModel):
    error: str
    detail: str | None = None
    timestamp: datetime


class AnalysisResponse(BaseModel):
    patterns: dict[str, Any]
    correlations: dict[str, Any]
    recommendations: list[Recommendation]
    risk_factors: list[RiskFactor]
    bristol_trends: BristolAnalysis
    health_score: HealthScore | None = None
    predictions: Any | None = None
    analysis_metadata: dict[str, Any]
