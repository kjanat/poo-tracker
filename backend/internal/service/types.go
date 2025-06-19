package service

// AnalyticsResponse represents a structured analytics response
type AnalyticsResponse struct {
	Total               int         `json:"total"`
	AverageBristol      float64     `json:"avgBristol,omitempty"`
	AveragePerDay       float64     `json:"averagePerDay,omitempty"`
	MostCommonType      int         `json:"mostCommonBristolType,omitempty"`
	BristolDistribution map[int]int `json:"bristolTypeDistribution,omitempty"`
	// Add more fields as analytics features expand
}
