"""Pydantic models for API responses."""

from datetime import datetime
from typing import Any, Dict, List

from pydantic import BaseModel


class BristolAnalysis(BaseModel):
    distribution: Dict[int, int]
    percentages: Dict[int, float]
    most_common: Dict[str, Any]
    health_indicator: str
    trend: str | None = None


class FrequencyStats(BaseModel):
    avg_daily: float
    max_daily: int
    min_daily: int
    total_days: int
    total_entries: int
    consistency_score: float


class TimingPattern(BaseModel):
    hourly_distribution: Dict[int, int]
    daily_distribution: Dict[str, int]
    peak_hour: int
    most_active_day: str
    regularity_score: float


class HealthScore(BaseModel):
    overall_score: float
    bristol_score: float
    frequency_score: float
    pain_score: float
    satisfaction_score: float
    trend: str


class Recommendation(BaseModel):
    id: str
    category: str
    title: str
    description: str
    priority: str
    confidence: float
    evidence: List[str]


class RiskFactor(BaseModel):
    factor: str
    severity: str
    description: str
    prevalence: float
    recommendation: str


class AnalysisResponse(BaseModel):
    patterns: Dict[str, Any]
    correlations: Dict[str, Any]
    bristol_trends: BristolAnalysis
    recommendations: List[Any]
    risk_factors: List[Any]
    health_score: HealthScore | None = None
    predictions: Any | None = None
    analysis_metadata: Dict[str, Any]


class ErrorResponse(BaseModel):
    error: str
    detail: str | None = None
    timestamp: datetime


class HealthResponse(BaseModel):
    status: str
    timestamp: datetime
    redis_connected: bool
    ml_models_loaded: bool
    response_time_ms: float
    version: str
