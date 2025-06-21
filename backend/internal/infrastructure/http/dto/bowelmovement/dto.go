package bowelmovement

import (
	"time"

	"github.com/kjanat/poo-tracker/backend/internal/domain/bowelmovement"
	"github.com/kjanat/poo-tracker/backend/internal/domain/shared"
)

// CreateBowelMovementRequest represents the HTTP request for creating a bowel movement
type CreateBowelMovementRequest struct {
	BristolType  int       `json:"bristol_type" binding:"required,min=1,max=7"`
	RecordedAt   time.Time `json:"recorded_at" binding:"required"`
	Volume       *string   `json:"volume,omitempty"`
	Color        *string   `json:"color,omitempty"`
	Consistency  *string   `json:"consistency,omitempty"`
	Floaters     bool      `json:"floaters"`
	Pain         int       `json:"pain" binding:"required,min=1,max=10"`
	Strain       int       `json:"strain" binding:"required,min=1,max=10"`
	Satisfaction int       `json:"satisfaction" binding:"required,min=1,max=10"`
	PhotoURL     string    `json:"photo_url,omitempty"`
	SmellLevel   *string   `json:"smell_level,omitempty"`
}

// UpdateBowelMovementRequest represents the HTTP request for updating a bowel movement
type UpdateBowelMovementRequest struct {
	BristolType  *int       `json:"bristol_type,omitempty" binding:"omitempty,min=1,max=7"`
	RecordedAt   *time.Time `json:"recorded_at,omitempty"`
	Volume       *string    `json:"volume,omitempty"`
	Color        *string    `json:"color,omitempty"`
	Consistency  *string    `json:"consistency,omitempty"`
	Floaters     *bool      `json:"floaters,omitempty"`
	Pain         *int       `json:"pain,omitempty" binding:"omitempty,min=1,max=10"`
	Strain       *int       `json:"strain,omitempty" binding:"omitempty,min=1,max=10"`
	Satisfaction *int       `json:"satisfaction,omitempty" binding:"omitempty,min=1,max=10"`
	PhotoURL     *string    `json:"photo_url,omitempty"`
	SmellLevel   *string    `json:"smell_level,omitempty"`
}

// CreateBowelMovementDetailsRequest represents the HTTP request for creating bowel movement details
type CreateBowelMovementDetailsRequest struct {
	Notes             string   `json:"notes,omitempty"`
	DetailedNotes     string   `json:"detailed_notes,omitempty"`
	Environment       string   `json:"environment,omitempty"`
	PreConditions     string   `json:"pre_conditions,omitempty"`
	PostConditions    string   `json:"post_conditions,omitempty"`
	Tags              []string `json:"tags,omitempty"`
	WeatherCondition  string   `json:"weather_condition,omitempty"`
	StressLevel       *int     `json:"stress_level,omitempty" binding:"omitempty,min=1,max=10"`
	SleepQuality      *int     `json:"sleep_quality,omitempty" binding:"omitempty,min=1,max=10"`
	ExerciseIntensity *int     `json:"exercise_intensity,omitempty" binding:"omitempty,min=1,max=10"`
}

// UpdateBowelMovementDetailsRequest represents the HTTP request for updating bowel movement details
type UpdateBowelMovementDetailsRequest struct {
	Notes             *string  `json:"notes,omitempty"`
	DetailedNotes     *string  `json:"detailed_notes,omitempty"`
	Environment       *string  `json:"environment,omitempty"`
	PreConditions     *string  `json:"pre_conditions,omitempty"`
	PostConditions    *string  `json:"post_conditions,omitempty"`
	Tags              []string `json:"tags,omitempty"`
	WeatherCondition  *string  `json:"weather_condition,omitempty"`
	StressLevel       *int     `json:"stress_level,omitempty" binding:"omitempty,min=1,max=10"`
	SleepQuality      *int     `json:"sleep_quality,omitempty" binding:"omitempty,min=1,max=10"`
	ExerciseIntensity *int     `json:"exercise_intensity,omitempty" binding:"omitempty,min=1,max=10"`
}

// BowelMovementResponse represents the HTTP response for a bowel movement
type BowelMovementResponse struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	RecordedAt time.Time `json:"recorded_at"`

	BristolType  int                 `json:"bristol_type"`
	Volume       *shared.Volume      `json:"volume,omitempty"`
	Color        *shared.Color       `json:"color,omitempty"`
	Consistency  *shared.Consistency `json:"consistency,omitempty"`
	Floaters     bool                `json:"floaters"`
	Pain         int                 `json:"pain"`
	Strain       int                 `json:"strain"`
	Satisfaction int                 `json:"satisfaction"`
	PhotoURL     string              `json:"photo_url,omitempty"`
	SmellLevel   *shared.SmellLevel  `json:"smell_level,omitempty"`
	HasDetails   bool                `json:"has_details"`
}

// BowelMovementDetailsResponse represents the HTTP response for bowel movement details
type BowelMovementDetailsResponse struct {
	ID              string    `json:"id"`
	BowelMovementID string    `json:"bowel_movement_id"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`

	Notes             string      `json:"notes,omitempty"`
	DetailedNotes     string      `json:"detailed_notes,omitempty"`
	Environment       string      `json:"environment,omitempty"`
	PreConditions     string      `json:"pre_conditions,omitempty"`
	PostConditions    string      `json:"post_conditions,omitempty"`
	AIAnalysis        interface{} `json:"ai_analysis,omitempty"`
	AIConfidence      *float64    `json:"ai_confidence,omitempty"`
	AIRecommendations string      `json:"ai_recommendations,omitempty"`
	Tags              []string    `json:"tags,omitempty"`
	WeatherCondition  string      `json:"weather_condition,omitempty"`
	StressLevel       *int        `json:"stress_level,omitempty"`
	SleepQuality      *int        `json:"sleep_quality,omitempty"`
	ExerciseIntensity *int        `json:"exercise_intensity,omitempty"`
}

// BowelMovementStatsResponse represents the HTTP response for bowel movement statistics
type BowelMovementStatsResponse struct {
	TotalCount          int64       `json:"total_count"`
	AveragePain         float64     `json:"average_pain"`
	AverageStrain       float64     `json:"average_strain"`
	AverageSatisfaction float64     `json:"average_satisfaction"`
	MostCommonBristol   int         `json:"most_common_bristol"`
	FrequencyPerDay     float64     `json:"frequency_per_day"`
	BristolDistribution map[int]int `json:"bristol_distribution"`
}

// ListBowelMovementsResponse represents the HTTP response for listing bowel movements
type ListBowelMovementsResponse struct {
	BowelMovements []BowelMovementResponse `json:"bowel_movements"`
	Total          int64                   `json:"total"`
	Limit          int                     `json:"limit"`
	Offset         int                     `json:"offset"`
}

// ToCreateInput converts HTTP request to domain input
func (r *CreateBowelMovementRequest) ToCreateInput() *bowelmovement.CreateBowelMovementInput {
	input := &bowelmovement.CreateBowelMovementInput{
		BristolType:  r.BristolType,
		RecordedAt:   r.RecordedAt,
		Floaters:     r.Floaters,
		Pain:         r.Pain,
		Strain:       r.Strain,
		Satisfaction: r.Satisfaction,
		PhotoURL:     r.PhotoURL,
	}

	if r.Volume != nil {
		input.Volume = r.Volume
	}
	if r.Color != nil {
		input.Color = r.Color
	}
	if r.Consistency != nil {
		input.Consistency = r.Consistency
	}
	if r.SmellLevel != nil {
		input.SmellLevel = r.SmellLevel
	}

	return input
}

// ToUpdateInput converts HTTP request to domain input
func (r *UpdateBowelMovementRequest) ToUpdateInput() *bowelmovement.UpdateBowelMovementInput {
	return &bowelmovement.UpdateBowelMovementInput{
		BristolType:  r.BristolType,
		RecordedAt:   r.RecordedAt,
		Volume:       r.Volume,
		Color:        r.Color,
		Consistency:  r.Consistency,
		Floaters:     r.Floaters,
		Pain:         r.Pain,
		Strain:       r.Strain,
		Satisfaction: r.Satisfaction,
		PhotoURL:     r.PhotoURL,
		SmellLevel:   r.SmellLevel,
	}
}

// ToCreateDetailsInput converts HTTP request to domain input
func (r *CreateBowelMovementDetailsRequest) ToCreateDetailsInput() *bowelmovement.CreateBowelMovementDetailsInput {
	return &bowelmovement.CreateBowelMovementDetailsInput{
		Notes:             r.Notes,
		DetailedNotes:     r.DetailedNotes,
		Environment:       r.Environment,
		PreConditions:     r.PreConditions,
		PostConditions:    r.PostConditions,
		Tags:              r.Tags,
		WeatherCondition:  r.WeatherCondition,
		StressLevel:       r.StressLevel,
		SleepQuality:      r.SleepQuality,
		ExerciseIntensity: r.ExerciseIntensity,
	}
}

// ToUpdateDetailsInput converts HTTP request to domain input
func (r *UpdateBowelMovementDetailsRequest) ToUpdateDetailsInput() *bowelmovement.UpdateBowelMovementDetailsInput {
	return &bowelmovement.UpdateBowelMovementDetailsInput{
		Notes:             r.Notes,
		DetailedNotes:     r.DetailedNotes,
		Environment:       r.Environment,
		PreConditions:     r.PreConditions,
		PostConditions:    r.PostConditions,
		Tags:              r.Tags,
		WeatherCondition:  r.WeatherCondition,
		StressLevel:       r.StressLevel,
		SleepQuality:      r.SleepQuality,
		ExerciseIntensity: r.ExerciseIntensity,
	}
}

// FromDomain converts domain model to HTTP response
func FromDomain(bm *bowelmovement.BowelMovement) BowelMovementResponse {
	return BowelMovementResponse{
		ID:           bm.ID,
		UserID:       bm.UserID,
		CreatedAt:    bm.CreatedAt,
		UpdatedAt:    bm.UpdatedAt,
		RecordedAt:   bm.RecordedAt,
		BristolType:  bm.BristolType,
		Volume:       bm.Volume,
		Color:        bm.Color,
		Consistency:  bm.Consistency,
		Floaters:     bm.Floaters,
		Pain:         bm.Pain,
		Strain:       bm.Strain,
		Satisfaction: bm.Satisfaction,
		PhotoURL:     bm.PhotoURL,
		SmellLevel:   bm.SmellLevel,
		HasDetails:   bm.HasDetails,
	}
}

// FromDomainDetails converts domain details to HTTP response
func FromDomainDetails(details *bowelmovement.BowelMovementDetails) BowelMovementDetailsResponse {
	return BowelMovementDetailsResponse{
		ID:                details.ID,
		BowelMovementID:   details.BowelMovementID,
		CreatedAt:         details.CreatedAt,
		UpdatedAt:         details.UpdatedAt,
		Notes:             details.Notes,
		DetailedNotes:     details.DetailedNotes,
		Environment:       details.Environment,
		PreConditions:     details.PreConditions,
		PostConditions:    details.PostConditions,
		AIAnalysis:        details.AIAnalysis,
		AIConfidence:      details.AIConfidence,
		AIRecommendations: details.AIRecommendations,
		Tags:              details.Tags,
		WeatherCondition:  details.WeatherCondition,
		StressLevel:       details.StressLevel,
		SleepQuality:      details.SleepQuality,
		ExerciseIntensity: details.ExerciseIntensity,
	}
}

// FromDomainStats converts domain stats to HTTP response
func FromDomainStats(stats *bowelmovement.UserBowelMovementStats) BowelMovementStatsResponse {
	return BowelMovementStatsResponse{
		TotalCount:          stats.TotalCount,
		AveragePain:         stats.AveragePain,
		AverageStrain:       stats.AverageStrain,
		AverageSatisfaction: stats.AverageSatisfaction,
		MostCommonBristol:   stats.MostCommonBristol,
		FrequencyPerDay:     stats.FrequencyPerDay,
		BristolDistribution: stats.BristolDistribution,
	}
}
