import json
import os
from datetime import datetime, timedelta
from typing import Any, Dict, List, Optional

import numpy as np
import pandas as pd
import redis
from fastapi import FastAPI, HTTPException
from pydantic import BaseModel

app = FastAPI(
    title="Poo Tracker AI Service",
    description="AI-powered analysis for bowel movement patterns and correlations",
    version="1.0.0",
)

# Redis connection
redis_client = redis.from_url(os.getenv("REDIS_URL", "redis://localhost:6379"))


class EntryData(BaseModel):
    id: str
    userId: str
    bristolType: int
    volume: Optional[str] = None
    color: Optional[str] = None
    consistency: Optional[str] = None
    floaters: bool = False
    pain: Optional[int] = None
    strain: Optional[int] = None
    satisfaction: Optional[int] = None
    createdAt: str


class MealData(BaseModel):
    id: str
    userId: str
    name: str
    mealTime: str
    category: Optional[str] = None
    spicyLevel: Optional[int] = None
    fiberRich: bool = False
    dairy: bool = False
    gluten: bool = False


class AnalysisRequest(BaseModel):
    entries: List[EntryData]
    meals: List[MealData]


class AnalysisResponse(BaseModel):
    patterns: Dict[str, Any]
    correlations: Dict[str, Any]
    recommendations: List[str]
    risk_factors: List[str]
    bristol_trends: Dict[str, Any]


@app.get("/health")
async def health_check():
    """Health check endpoint"""
    # Check Redis connection without failing if Redis is unavailable
    redis_connected = False
    try:
        redis_connected = redis_client.ping()
    except Exception:
        # Redis connection failed, but service can still function
        pass

    return {
        "status": "healthy",
        "timestamp": datetime.now().isoformat(),
        "redis_connected": redis_connected,
    }


@app.post("/analyze", response_model=AnalysisResponse)
async def analyze_patterns(request: AnalysisRequest):
    """
    Analyze bowel movement patterns and correlations with meals.
    This is where the AI magic happens, you magnificent bastard.
    """
    try:
        if not request.entries:
            raise HTTPException(
                status_code=400, detail="No entries provided for analysis"
            )

        # Convert to DataFrames for analysis
        entries_df = pd.DataFrame([entry.dict() for entry in request.entries])
        entries_df["createdAt"] = pd.to_datetime(entries_df["createdAt"])

        meals_df = pd.DataFrame([meal.dict() for meal in request.meals])
        if not meals_df.empty:
            meals_df["mealTime"] = pd.to_datetime(meals_df["mealTime"])

        # Perform analysis
        analysis_result = perform_comprehensive_analysis(entries_df, meals_df)

        # Cache results
        cache_key = (
            f"analysis:{request.entries[0].userId}:{datetime.now().strftime('%Y%m%d')}"
        )
        redis_client.setex(
            cache_key, 3600, json.dumps(analysis_result)
        )  # Cache for 1 hour

        return AnalysisResponse(**analysis_result)

    except Exception as e:
        raise HTTPException(status_code=500, detail=f"Analysis failed: {str(e)}")


@app.get("/summary/{user_id}", response_model=AnalysisResponse)
async def get_cached_summary(user_id: str) -> AnalysisResponse:
    """Return the most recent cached analysis for ``user_id``."""
    try:
        pattern = f"analysis:{user_id}:*"
        keys = sorted(redis_client.keys(pattern))
        if not keys:
            raise HTTPException(status_code=404, detail="No cached analysis found")

        latest_key = keys[-1]
        cached = redis_client.get(latest_key)
        if cached is None:
            raise HTTPException(status_code=404, detail="Cached analysis missing")

        data = json.loads(cached)
        return AnalysisResponse(**data)
    except HTTPException:
        raise
    except Exception as exc:  # pragma: no cover - unexpected errors
        raise HTTPException(status_code=500, detail=f"Summary retrieval failed: {exc}")


def perform_comprehensive_analysis(
    entries_df: pd.DataFrame, meals_df: pd.DataFrame
) -> Dict[str, Any]:
    """
    Perform comprehensive analysis of bowel movements and correlations.
    """

    # Bristol Stool Chart analysis
    bristol_analysis = analyze_bristol_patterns(entries_df)

    # Timing patterns
    timing_patterns = analyze_timing_patterns(entries_df)

    # Meal correlations (if meal data available)
    meal_correlations = {}
    if not meals_df.empty:
        meal_correlations = analyze_meal_correlations(entries_df, meals_df)

    # Generate recommendations
    recommendations = generate_recommendations(entries_df, meals_df)

    # Identify risk factors
    risk_factors = identify_risk_factors(entries_df)

    return {
        "patterns": {
            "timing": timing_patterns,
            "frequency": calculate_frequency_stats(entries_df),
            "consistency_trends": analyze_consistency_trends(entries_df),
        },
        "correlations": meal_correlations,
        "recommendations": recommendations,
        "risk_factors": risk_factors,
        "bristol_trends": bristol_analysis,
    }


def analyze_bristol_patterns(df: pd.DataFrame) -> Dict[str, Any]:
    """Analyze Bristol Stool Chart patterns"""
    bristol_counts = df["bristolType"].value_counts().sort_index()

    # Bristol type descriptions
    bristol_descriptions = {
        1: "Severe constipation (hard lumps)",
        2: "Mild constipation (lumpy sausage)",
        3: "Normal (cracked sausage)",
        4: "Ideal (smooth sausage)",
        5: "Lacking fiber (soft blobs)",
        6: "Mild diarrhea (fluffy pieces)",
        7: "Severe diarrhea (watery)",
    }

    # Calculate percentages
    total_entries = len(df)
    bristol_percentages = (bristol_counts / total_entries * 100).round(2)

    # Determine dominant pattern
    most_common_type = bristol_counts.idxmax()

    return {
        "distribution": bristol_counts.to_dict(),
        "percentages": bristol_percentages.to_dict(),
        "most_common": {
            "type": int(most_common_type),
            "description": bristol_descriptions.get(
                int(most_common_type), "Unknown type"
            ),
            "percentage": float(bristol_percentages[most_common_type]),
        },
        "health_indicator": assess_bristol_health(bristol_percentages),
    }


def analyze_timing_patterns(df: pd.DataFrame) -> Dict[str, Any]:
    """Analyze timing patterns of bowel movements"""
    df["hour"] = df["createdAt"].dt.hour
    df["day_of_week"] = df["createdAt"].dt.day_name()

    hourly_distribution = df["hour"].value_counts().sort_index()
    daily_distribution = df["day_of_week"].value_counts()

    return {
        "hourly_distribution": hourly_distribution.to_dict(),
        "daily_distribution": daily_distribution.to_dict(),
        "peak_hour": int(hourly_distribution.idxmax()),
        "most_active_day": daily_distribution.idxmax(),
    }


def analyze_meal_correlations(
    entries_df: pd.DataFrame, meals_df: pd.DataFrame
) -> Dict[str, Any]:
    """Analyze correlations between meals and bowel movements"""
    correlations: Dict[str, Any] = {}

    # Look for patterns within 24 hours of meals
    for _, meal in meals_df.iterrows():
        meal_time = meal["mealTime"]

        # Find entries within 6-48 hours after meal (digestion time)
        time_window_start = meal_time + timedelta(hours=6)
        time_window_end = meal_time + timedelta(hours=48)

        related_entries = entries_df[
            (entries_df["createdAt"] >= time_window_start)
            & (entries_df["createdAt"] <= time_window_end)
        ]

        if not related_entries.empty:
            # Analyze bristol types following this meal type
            if meal["category"] not in correlations:
                correlations[meal["category"]] = {
                    "bristol_avg": [],
                    "satisfaction_avg": [],
                    "pain_avg": [],
                }

            correlations[meal["category"]]["bristol_avg"].extend(
                related_entries["bristolType"].tolist()
            )
            if "satisfaction" in related_entries.columns:
                correlations[meal["category"]]["satisfaction_avg"].extend(
                    related_entries["satisfaction"].dropna().tolist()
                )
            if "pain" in related_entries.columns:
                correlations[meal["category"]]["pain_avg"].extend(
                    related_entries["pain"].dropna().tolist()
                )

    # Calculate averages
    for category in correlations:
        if correlations[category]["bristol_avg"]:
            correlations[category]["avg_bristol"] = np.mean(
                correlations[category]["bristol_avg"]
            )
        if correlations[category]["satisfaction_avg"]:
            correlations[category]["avg_satisfaction"] = np.mean(
                correlations[category]["satisfaction_avg"]
            )
        if correlations[category]["pain_avg"]:
            correlations[category]["avg_pain"] = np.mean(
                correlations[category]["pain_avg"]
            )

    return correlations


def calculate_frequency_stats(df: pd.DataFrame) -> Dict[str, Any]:
    """Calculate frequency statistics"""
    df["date"] = df["createdAt"].dt.date
    daily_counts = df.groupby("date").size()

    return {
        "avg_daily": float(daily_counts.mean()),
        "max_daily": int(daily_counts.max()),
        "min_daily": int(daily_counts.min()),
        "total_days": len(daily_counts),
        "total_entries": len(df),
    }


def analyze_consistency_trends(df: pd.DataFrame) -> Dict[str, Any]:
    """Analyze consistency trends over time"""
    if "consistency" not in df.columns:
        return {}

    consistency_counts = df["consistency"].value_counts()
    return {
        "distribution": consistency_counts.to_dict(),
        "most_common": (
            consistency_counts.idxmax() if not consistency_counts.empty else None
        ),
    }


def generate_recommendations(
    entries_df: pd.DataFrame, meals_df: pd.DataFrame
) -> List[str]:
    """Generate personalized recommendations based on patterns"""
    recommendations = []

    # Bristol type recommendations
    avg_bristol = entries_df["bristolType"].mean()

    if avg_bristol < 3:
        recommendations.extend(
            [
                "Consider increasing fiber intake - you might be constipated",
                "Try drinking more water throughout the day",
                "Consider adding prunes or other natural laxatives to your diet",
            ]
        )
    elif avg_bristol > 5:
        recommendations.extend(
            [
                "Your stool is quite loose - consider reducing dairy or high-fat foods",
                "You might want to increase binding foods like bananas and rice",
                "Consider keeping a detailed food diary to identify triggers",
            ]
        )
    else:
        recommendations.append(
            "Your Bristol scores look healthy overall! Keep doing what you're doing."
        )

    # Frequency recommendations
    daily_avg = len(entries_df) / max(
        1, (entries_df["createdAt"].max() - entries_df["createdAt"].min()).days + 1
    )

    if daily_avg < 0.5:
        recommendations.append(
            "You're not pooping very often - increase fiber and water intake"
        )
    elif daily_avg > 3:
        recommendations.append(
            "You're pooping quite frequently - monitor for any digestive issues"
        )

    # Pain recommendations
    if "pain" in entries_df.columns and entries_df["pain"].notna().any():
        avg_pain = entries_df["pain"].dropna().mean()
        if avg_pain > 5:
            recommendations.append(
                "You're experiencing significant pain - consider consulting a healthcare provider"
            )

    return recommendations


def identify_risk_factors(df: pd.DataFrame) -> List[str]:
    """Identify potential risk factors based on patterns"""
    risk_factors = []

    # Check for extreme Bristol types
    extreme_bristol = df[df["bristolType"].isin([1, 2, 6, 7])]
    if len(extreme_bristol) / len(df) > 0.3:
        risk_factors.append(
            "High frequency of extreme Bristol types (chronic constipation or diarrhea)"
        )

    # Check for high pain scores
    if "pain" in df.columns and df["pain"].notna().any():
        high_pain_entries = df[df["pain"] > 7]
        if len(high_pain_entries) / len(df.dropna(subset=["pain"])) > 0.2:
            risk_factors.append("Frequent high pain scores during bowel movements")

    # Check for very low satisfaction
    if "satisfaction" in df.columns and df["satisfaction"].notna().any():
        low_satisfaction = df[df["satisfaction"] < 3]
        if len(low_satisfaction) / len(df.dropna(subset=["satisfaction"])) > 0.4:
            risk_factors.append("Consistently low satisfaction with bowel movements")

    return risk_factors


def assess_bristol_health(percentages: pd.Series) -> str:
    """Assess overall digestive health based on Bristol distribution"""
    healthy_range = percentages.get(3, 0) + percentages.get(4, 0)  # Types 3-4 are ideal

    if healthy_range > 70:
        return "Excellent - Your bowel movements are consistently healthy"
    elif healthy_range > 50:
        return "Good - Mostly healthy with some room for improvement"
    elif healthy_range > 30:
        return "Fair - Consider dietary adjustments for better digestive health"
    else:
        return "Poor - Significant digestive issues detected, consider medical consultation"


if __name__ == "__main__":
    import uvicorn

    uvicorn.run(app, host="0.0.0.0", port=8000)
