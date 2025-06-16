"""Data processing utilities for the AI service."""

from datetime import timedelta
from typing import Any

import numpy as np
import pandas as pd

from ..config.logging import get_logger
from ..models.database import BowelMovementData, MealData, SymptomData

logger = get_logger("data_processor")


class DataProcessor:
    """Utility class for data processing and feature engineering."""

    def __init__(self):
        self.bristol_weights = {1: -2, 2: -1, 3: 0, 4: 1, 5: 0, 6: -1, 7: -2}
        self.meal_categories = {
            "breakfast": 0,
            "lunch": 1,
            "dinner": 2,
            "snack": 3,
            "beverage": 4,
        }
        self.symptom_types = {
            "bloating": 0,
            "cramps": 1,
            "nausea": 2,
            "gas": 3,
            "heartburn": 4,
            "constipation": 5,
            "diarrhea": 6,
        }

    def preprocess_bowel_movements(self, data: list[BowelMovementData]) -> pd.DataFrame:
        """Preprocess bowel movement data for analysis."""
        if not data:
            return pd.DataFrame()

        df = pd.DataFrame([item.to_dict() for item in data])

        # Add time-based features
        df = self._add_time_features(df, "created_at")

        # Add derived features
        df["bristol_health_score"] = df["bristol_type"].map(self.bristol_weights)
        df["is_constipated"] = df["bristol_type"] <= 2
        df["is_diarrhea"] = df["bristol_type"] >= 6
        df["is_ideal"] = df["bristol_type"].isin([3, 4])

        # Handle missing values
        df = self._handle_missing_values(df)

        # Add rolling averages for trends
        df = self._add_rolling_features(df)

        return df

    def preprocess_meals(self, data: list[MealData]) -> pd.DataFrame:
        """Preprocess meal data for analysis."""
        if not data:
            return pd.DataFrame()

        df = pd.DataFrame([item.to_dict() for item in data])

        # Add time-based features
        df = self._add_time_features(df, "meal_time")

        # Encode categorical features
        df["category_encoded"] = df["category"].map(self.meal_categories).fillna(-1)

        # Create meal characteristic features
        df["is_problematic"] = (
            (df["spicy_level"].fillna(0) > 7)
            | (df["dairy"] == True)
            | (df["gluten"] == True)
        )
        df["is_healthy"] = (df["fiber_rich"] == True) & (
            df["spicy_level"].fillna(0) <= 3
        )

        return df

    def preprocess_symptoms(self, data: list[SymptomData]) -> pd.DataFrame:
        """Preprocess symptom data for analysis."""
        if not data:
            return pd.DataFrame()

        df = pd.DataFrame([item.to_dict() for item in data])

        # Add time-based features
        df = self._add_time_features(df, "created_at")

        # Encode symptom types
        df["type_encoded"] = df["type"].map(self.symptom_types).fillna(-1)

        # Categorize severity
        df["severity_category"] = pd.cut(
            df["severity"], bins=[0, 3, 6, 10], labels=["mild", "moderate", "severe"]
        )

        return df

    def create_time_windows(
        self, df: pd.DataFrame, window_size: str = "1D", time_col: str = "created_at"
    ) -> pd.DataFrame:
        """Create time-based windows for analysis."""
        if df.empty or time_col not in df.columns:
            return df

        df = df.copy()
        df[time_col] = pd.to_datetime(df[time_col])
        df = df.set_index(time_col)

        # Resample to create time windows
        agg_dict = {}
        if "bristol_type" in df.columns:
            agg_dict["bristol_type"] = ["mean", "std", "count"]
        if "pain" in df.columns:
            agg_dict["pain"] = ["mean", "max"]
        if "satisfaction" in df.columns:
            agg_dict["satisfaction"] = ["mean", "min"]

        if not agg_dict:
            return df.reset_index()

        windowed = df.resample(window_size).agg(agg_dict).reset_index()

        # Flatten column names
        windowed.columns = [
            f"{col[0]}_{col[1]}" if col[1] else col[0] for col in windowed.columns
        ]

        return windowed

    def detect_patterns(
        self, df: pd.DataFrame, pattern_type: str = "daily"
    ) -> dict[str, Any]:
        """Detect patterns in time series data."""
        if df.empty:
            return {}

        patterns = {}

        if pattern_type == "daily":
            patterns = self._detect_daily_patterns(df)
        elif pattern_type == "weekly":
            patterns = self._detect_weekly_patterns(df)
        elif pattern_type == "cyclical":
            patterns = self._detect_cyclical_patterns(df)

        return patterns

    def calculate_correlations(
        self,
        bm_df: pd.DataFrame,
        meal_df: pd.DataFrame,
        lag_hours: list[int] = [6, 12, 24, 48],
    ) -> dict[str, float]:
        """Calculate time-lagged correlations between meals and bowel movements."""
        correlations = {}

        if bm_df.empty or meal_df.empty:
            return correlations

        # Ensure datetime index
        bm_df = bm_df.copy()
        meal_df = meal_df.copy()

        bm_df["timestamp"] = pd.to_datetime(bm_df["created_at"])
        meal_df["timestamp"] = pd.to_datetime(meal_df["meal_time"])

        for lag in lag_hours:
            # Shift meal times by lag
            meal_shifted = meal_df.copy()
            meal_shifted["timestamp"] = meal_shifted["timestamp"] + timedelta(hours=lag)

            # Find correlations within time windows
            correlation = self._calculate_windowed_correlation(
                bm_df, meal_shifted, window_hours=2
            )

            if correlation is not None:
                correlations[f"lag_{lag}h"] = correlation

        return correlations

    def normalize_features(self, df: pd.DataFrame, features: list[str]) -> pd.DataFrame:
        """Normalize specified features using z-score normalization."""
        df = df.copy()

        for feature in features:
            if feature in df.columns:
                mean_val = df[feature].mean()
                std_val = df[feature].std()

                if std_val > 0:
                    df[f"{feature}_normalized"] = (df[feature] - mean_val) / std_val
                else:
                    df[f"{feature}_normalized"] = 0

        return df

    def create_feature_matrix(
        self,
        bm_df: pd.DataFrame,
        meal_df: pd.DataFrame | None = None,
        symptom_df: pd.DataFrame | None = None,
    ) -> pd.DataFrame:
        """Create a comprehensive feature matrix for ML models."""
        if bm_df.empty:
            return pd.DataFrame()

        features = bm_df[["user_id", "created_at"]].copy()

        # Bowel movement features
        bm_features = self._extract_bowel_movement_features(bm_df)
        features = features.merge(bm_features, on=["user_id", "created_at"], how="left")

        # Meal features (if available)
        if meal_df is not None and not meal_df.empty:
            meal_features = self._extract_meal_features(meal_df, bm_df)
            features = features.merge(meal_features, on="user_id", how="left")

        # Symptom features (if available)
        if symptom_df is not None and not symptom_df.empty:
            symptom_features = self._extract_symptom_features(symptom_df, bm_df)
            features = features.merge(symptom_features, on="user_id", how="left")

        # Fill missing values
        features = features.fillna(0)

        return features

    def _add_time_features(self, df: pd.DataFrame, time_col: str) -> pd.DataFrame:
        """Add time-based features to dataframe."""
        df = df.copy()
        df[time_col] = pd.to_datetime(df[time_col])

        df["hour"] = df[time_col].dt.hour
        df["day_of_week"] = df[time_col].dt.dayofweek
        df["day_of_month"] = df[time_col].dt.day
        df["month"] = df[time_col].dt.month
        df["quarter"] = df[time_col].dt.quarter
        df["is_weekend"] = df["day_of_week"].isin([5, 6])
        df["is_morning"] = df["hour"].between(6, 11)
        df["is_afternoon"] = df["hour"].between(12, 17)
        df["is_evening"] = df["hour"].between(18, 23)
        df["is_night"] = df["hour"].isin([0, 1, 2, 3, 4, 5])

        return df

    def _handle_missing_values(self, df: pd.DataFrame) -> pd.DataFrame:
        """Handle missing values in the dataframe."""
        df = df.copy()

        # Fill numeric columns with median
        numeric_cols = df.select_dtypes(include=[np.number]).columns
        for col in numeric_cols:
            if col in ["pain", "strain", "satisfaction"]:
                # Use conservative defaults for health metrics
                defaults = {"pain": 1, "strain": 1, "satisfaction": 5}
                df[col] = df[col].fillna(defaults.get(col, df[col].median()))

        # Fill categorical columns with mode or default
        categorical_cols = df.select_dtypes(include=["object"]).columns
        for col in categorical_cols:
            if col in ["volume", "color", "consistency"]:
                df[col] = df[col].fillna("unknown")

        return df

    def _add_rolling_features(self, df: pd.DataFrame) -> pd.DataFrame:
        """Add rolling window features for trend analysis."""
        df = df.copy()
        df = df.sort_values("created_at")

        # Rolling averages for key metrics
        for window in [3, 7, 14]:  # 3, 7, and 14-day windows
            if "bristol_type" in df.columns:
                df[f"bristol_rolling_{window}d"] = (
                    df["bristol_type"].rolling(window=window, min_periods=1).mean()
                )

            if "pain" in df.columns:
                df[f"pain_rolling_{window}d"] = (
                    df["pain"].rolling(window=window, min_periods=1).mean()
                )

            if "satisfaction" in df.columns:
                df[f"satisfaction_rolling_{window}d"] = (
                    df["satisfaction"].rolling(window=window, min_periods=1).mean()
                )

        return df

    def _detect_daily_patterns(self, df: pd.DataFrame) -> dict[str, Any]:
        """Detect daily patterns in the data."""
        if "hour" not in df.columns:
            return {}

        hourly_freq = df["hour"].value_counts().sort_index()

        return {
            "peak_hours": hourly_freq.nlargest(3).index.tolist(),
            "low_hours": hourly_freq.nsmallest(3).index.tolist(),
            "hourly_distribution": hourly_freq.to_dict(),
            "morning_frequency": len(df[df["is_morning"] == True]),
            "evening_frequency": len(df[df["is_evening"] == True]),
        }

    def _detect_weekly_patterns(self, df: pd.DataFrame) -> dict[str, Any]:
        """Detect weekly patterns in the data."""
        if "day_of_week" not in df.columns:
            return {}

        daily_freq = df["day_of_week"].value_counts().sort_index()
        day_names = [
            "Monday",
            "Tuesday",
            "Wednesday",
            "Thursday",
            "Friday",
            "Saturday",
            "Sunday",
        ]

        return {
            "busiest_days": [day_names[i] for i in daily_freq.nlargest(3).index],
            "quietest_days": [day_names[i] for i in daily_freq.nsmallest(3).index],
            "weekend_vs_weekday": {
                "weekend": len(df[df["is_weekend"] == True]),
                "weekday": len(df[df["is_weekend"] == False]),
            },
        }

    def _detect_cyclical_patterns(self, df: pd.DataFrame) -> dict[str, Any]:
        """Detect cyclical patterns using FFT analysis."""
        if "bristol_type" not in df.columns or len(df) < 14:
            return {}

        # Simple cyclical pattern detection
        # Group by day and calculate daily averages
        df_daily = df.groupby(df["created_at"].dt.date)["bristol_type"].mean()

        if len(df_daily) < 7:
            return {}

        # Calculate autocorrelation for different lags
        autocorr = {}
        for lag in [1, 3, 7, 14]:  # 1, 3, 7, 14 day cycles
            if len(df_daily) > lag:
                corr = df_daily.autocorr(lag=lag)
                if not np.isnan(corr):
                    autocorr[f"{lag}_day_cycle"] = corr

        return {
            "autocorrelations": autocorr,
            "potential_cycles": [k for k, v in autocorr.items() if v > 0.3],
        }

    def _calculate_windowed_correlation(
        self, bm_df: pd.DataFrame, meal_df: pd.DataFrame, window_hours: int = 2
    ) -> float | None:
        """Calculate correlation within time windows."""
        correlations = []

        for _, meal in meal_df.iterrows():
            meal_time = meal["timestamp"]
            window_start = meal_time - timedelta(hours=window_hours / 2)
            window_end = meal_time + timedelta(hours=window_hours / 2)

            nearby_bms = bm_df[
                (bm_df["timestamp"] >= window_start)
                & (bm_df["timestamp"] <= window_end)
            ]

            if len(nearby_bms) > 0:
                # Simple correlation based on presence/absence
                correlations.append(1.0)
            else:
                correlations.append(0.0)

        return np.mean(correlations) if correlations else None

    def _extract_bowel_movement_features(self, df: pd.DataFrame) -> pd.DataFrame:
        """Extract features specific to bowel movements."""
        features = df[["user_id", "created_at"]].copy()

        # Basic features
        if "bristol_type" in df.columns:
            features["bristol_mean"] = df.groupby("user_id")["bristol_type"].transform(
                "mean"
            )
            features["bristol_std"] = df.groupby("user_id")["bristol_type"].transform(
                "std"
            )

        if "pain" in df.columns:
            features["pain_mean"] = df.groupby("user_id")["pain"].transform("mean")
            features["pain_max"] = df.groupby("user_id")["pain"].transform("max")

        if "satisfaction" in df.columns:
            features["satisfaction_mean"] = df.groupby("user_id")[
                "satisfaction"
            ].transform("mean")
            features["satisfaction_min"] = df.groupby("user_id")[
                "satisfaction"
            ].transform("min")

        # Frequency features
        features["daily_frequency"] = (
            df.groupby(["user_id", df["created_at"].dt.date])
            .size()
            .groupby("user_id")
            .mean()
        )

        return features

    def _extract_meal_features(
        self, meal_df: pd.DataFrame, bm_df: pd.DataFrame
    ) -> pd.DataFrame:
        """Extract meal-related features."""
        features = meal_df[["user_id"]].drop_duplicates()

        # Meal characteristics
        features["spicy_meal_ratio"] = (
            meal_df.groupby("user_id")["spicy_level"]
            .apply(lambda x: (x > 5).mean() if not x.isna().all() else 0)
            .reindex(features["user_id"])
            .fillna(0)
            .values
        )

        features["dairy_meal_ratio"] = (
            meal_df.groupby("user_id")["dairy"]
            .mean()
            .reindex(features["user_id"])
            .fillna(0)
            .values
        )
        features["fiber_meal_ratio"] = (
            meal_df.groupby("user_id")["fiber_rich"]
            .mean()
            .reindex(features["user_id"])
            .fillna(0)
            .values
        )

        return features

    def _extract_symptom_features(
        self, symptom_df: pd.DataFrame, bm_df: pd.DataFrame
    ) -> pd.DataFrame:
        """Extract symptom-related features."""
        features = symptom_df[["user_id"]].drop_duplicates()

        # Symptom characteristics
        features["avg_symptom_severity"] = (
            symptom_df.groupby("user_id")["severity"]
            .mean()
            .reindex(features["user_id"])
            .fillna(0)
            .values
        )
        features["symptom_frequency"] = (
            symptom_df.groupby("user_id")
            .size()
            .reindex(features["user_id"])
            .fillna(0)
            .values
        )

        # Most common symptom type
        most_common_symptoms = symptom_df.groupby("user_id")["type"].apply(
            lambda x: x.mode().iloc[0] if len(x.mode()) > 0 else "none"
        )
        features["most_common_symptom"] = (
            most_common_symptoms.reindex(features["user_id"]).fillna("none").values
        )

        return features
