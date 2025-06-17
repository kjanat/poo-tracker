"""Response models returned by the API."""

from __future__ import annotations

from typing import Any

from pydantic import BaseModel


class BristolAnalysis(BaseModel):
    """Analysis details for Bristol stool chart."""

    distribution: dict[int, int]
    percentages: dict[int, float]
    most_common: dict[str, Any]
    health_indicator: str
    trend: str | None = None


class TimingPattern(BaseModel):
    """Patterns related to timing of bowel movements."""

    hourly_distribution: dict[int, int]
    daily_distribution: dict[str, int]
    peak_hour: int
    most_active_day: str
    regularity_score: float


class FrequencyStats(BaseModel):
    """Statistics on bowel movement frequency."""

    avg_daily: float
    max_daily: int
    min_daily: int
    total_days: int
    total_entries: int
    consistency_score: float


class AnalysisResponse(BaseModel):
    """Full analysis response payload."""

    patterns: dict[str, Any]
    correlations: dict[str, Any]
    recommendations: list[Any]
    risk_factors: list[Any]
    bristol_trends: BristolAnalysis
    analysis_metadata: dict[str, Any]


class HealthScore(BaseModel):
    """Aggregated health score details."""

    overall_score: float
    bristol_score: float
    frequency_score: float
    pain_score: float
    satisfaction_score: float
    trend: str


class HealthResponse(BaseModel):
    """Health check response."""

    status: str
    timestamp: str
    redis_connected: bool
    ml_models_loaded: bool
    response_time_ms: float
    version: str


class ErrorResponse(BaseModel):
    """Standard error response."""

    error: str
    detail: str
    timestamp: Any


class Recommendation(BaseModel):
    """A single personalized recommendation."""

    id: str
    category: str
    title: str
    description: str
    priority: str
    confidence: float
    evidence: list[str]


class RiskFactor(BaseModel):
    """Identified risk factor derived from analysis."""

    factor: str
    severity: str
    description: str
    prevalence: float
    recommendation: str
