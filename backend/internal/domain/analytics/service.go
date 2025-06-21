// Package analytics provides the analytics service interface for health data analysis.
// This file contains the Service interface that defines the contract for analytics operations,
// including health overviews, correlation analysis, trend analysis, pattern recognition,
// and health insights.
//
// The service interface provides comprehensive analytics capabilities:
// - Cross-domain health overviews combining bowel, meal, symptom, and medication data
// - Correlation analysis to identify relationships between different health factors
// - Trend analysis for temporal patterns and changes over time
// - Behavior pattern recognition for lifestyle insights
// - Health scoring and personalized recommendations
package analytics

import (
	"context"
	"time"
)

// Service defines the interface for analytics operations
type Service interface {
	// Cross-domain analytics
	GetUserHealthOverview(ctx context.Context, userID string, start, end time.Time) (*HealthOverview, error)
	GetCorrelationAnalysis(ctx context.Context, userID string, start, end time.Time) (*CorrelationAnalysis, error)
	GetTrendAnalysis(ctx context.Context, userID string, start, end time.Time) (*TrendAnalysis, error)

	// Pattern recognition
	GetBehaviorPatterns(ctx context.Context, userID string, start, end time.Time) (*BehaviorPatterns, error)
	GetHealthInsights(ctx context.Context, userID string, start, end time.Time) (*HealthInsights, error)

	// Predictive analytics
	GetHealthScore(ctx context.Context, userID string) (*HealthScore, error)
	GetRecommendations(ctx context.Context, userID string) ([]*Recommendation, error)
}
