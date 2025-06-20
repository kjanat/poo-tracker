package analytics

import (
	"errors"
	"fmt"
)

// Domain errors
var (
	ErrInsufficientData     = errors.New("insufficient data for analysis")
	ErrInvalidDateRange     = errors.New("invalid date range for analysis")
	ErrAnalysisNotAvailable = errors.New("analysis not available for this user")
	ErrCorrelationFailed    = errors.New("correlation analysis failed")
	ErrTrendAnalysisFailed  = errors.New("trend analysis failed")
	ErrPredictionFailed     = errors.New("prediction analysis failed")
	ErrUserNotAuthorized    = errors.New("user not authorized to access analytics")
	ErrInvalidAnalysisType  = errors.New("invalid analysis type requested")
)

// AnalysisError represents an analysis error with additional context
type AnalysisError struct {
	Type    string `json:"type"`
	Message string `json:"message"`
	Context string `json:"context"`
}

func (e AnalysisError) Error() string {
	return fmt.Sprintf("analysis error '%s': %s (context: %s)", e.Type, e.Message, e.Context)
}

// NewAnalysisError creates a new analysis error
func NewAnalysisError(analysisType, message, context string) AnalysisError {
	return AnalysisError{
		Type:    analysisType,
		Message: message,
		Context: context,
	}
}
