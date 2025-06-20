package bowelmovement

import (
	"time"

	"github.com/kjanat/poo-tracker/backend/internal/domain/bowelmovement"
	"github.com/kjanat/poo-tracker/backend/internal/domain/shared"
)

// CreateBowelMovementRequest represents the HTTP request for creating a bowel movement
type CreateBowelMovementRequest struct {
	BristolType  int       `json:"bristolType" binding:"required,min=1,max=7"`
	RecordedAt   time.Time `json:"recordedAt" binding:"required"`
	Volume       *string   `json:"volume,omitempty"`
	Color        *string   `json:"color,omitempty"`
	Consistency  *string   `json:"consistency,omitempty"`
	Floaters     bool      `json:"floaters"`
	Pain         int       `json:"pain" binding:"min=1,max=10"`
	Strain       int       `json:"strain" binding:"min=1,max=10"`
	Satisfaction int       `json:"satisfaction" binding:"min=1,max=10"`
	PhotoURL     string    `json:"photoUrl,omitempty"`
	SmellLevel   *string   `json:"smellLevel,omitempty"`
}

// UpdateBowelMovementRequest represents the HTTP request for updating a bowel movement
type UpdateBowelMovementRequest struct {
	BristolType  *int       `json:"bristolType,omitempty" binding:"omitempty,min=1,max=7"`
	RecordedAt   *time.Time `json:"recordedAt,omitempty"`
	Volume       *string    `json:"volume,omitempty"`
	Color        *string    `json:"color,omitempty"`
	Consistency  *string    `json:"consistency,omitempty"`
	Floaters     *bool      `json:"floaters,omitempty"`
	Pain         *int       `json:"pain,omitempty" binding:"omitempty,min=1,max=10"`
	Strain       *int       `json:"strain,omitempty" binding:"omitempty,min=1,max=10"`
	Satisfaction *int       `json:"satisfaction,omitempty" binding:"omitempty,min=1,max=10"`
	PhotoURL     *string    `json:"photoUrl,omitempty"`
	SmellLevel   *string    `json:"smellLevel,omitempty"`
}

// CreateBowelMovementDetailsRequest represents the HTTP request for creating bowel movement details
type CreateBowelMovementDetailsRequest struct {
	Notes             string   `json:"notes,omitempty"`
	DetailedNotes     string   `json:"detailedNotes,omitempty"`
	Environment       string   `json:"environment,omitempty"`
	PreConditions     string   `json:"preConditions,omitempty"`
	PostConditions    string   `json:"postConditions,omitempty"`
	Tags              []string `json:"tags,omitempty"`
	WeatherCondition  string   `json:"weatherCondition,omitempty"`
	StressLevel       *int     `json:"stressLevel,omitempty" binding:"omitempty,min=1,max=10"`
	SleepQuality      *int     `json:"sleepQuality,omitempty" binding:"omitempty,min=1,max=10"`
	ExerciseIntensity *int     `json:"exerciseIntensity,omitempty" binding:"omitempty,min=1,max=10"`
}

// UpdateBowelMovementDetailsRequest represents the HTTP request for updating bowel movement details
type UpdateBowelMovementDetailsRequest struct {
	Notes             *string  `json:"notes,omitempty"`
	DetailedNotes     *string  `json:"detailedNotes,omitempty"`
	Environment       *string  `json:"environment,omitempty"`
	PreConditions     *string  `json:"preConditions,omitempty"`
	PostConditions    *string  `json:"postConditions,omitempty"`
	Tags              []string `json:"tags,omitempty"`
	WeatherCondition  *string  `json:"weatherCondition,omitempty"`
	StressLevel       *int     `json:"stressLevel,omitempty" binding:"omitempty,min=1,max=10"`
	SleepQuality      *int     `json:"sleepQuality,omitempty" binding:"omitempty,min=1,max=10"`
	ExerciseIntensity *int     `json:"exerciseIntensity,omitempty" binding:"omitempty,min=1,max=10"`
}

// BowelMovementResponse represents the HTTP response for a bowel movement
type BowelMovementResponse struct {
	ID         string    `json:"id"`
	UserID     string    `json:"userId"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
	RecordedAt time.Time `json:"recordedAt"`

	BristolType  int                 `json:"bristolType"`
	Volume       *shared.Volume      `json:"volume,omitempty"`
	Color        *shared.Color       `json:"color,omitempty"`
	Consistency  *shared.Consistency `json:"consistency,omitempty"`
	Floaters     bool                `json:"floaters"`
	Pain         int                 `json:"pain"`
	Strain       int                 `json:"strain"`
	Satisfaction int                 `json:"satisfaction"`
	PhotoURL     string              `json:"photoUrl,omitempty"`
	SmellLevel   *shared.SmellLevel  `json:"smellLevel,omitempty"`
	HasDetails   bool                `json:"hasDetails"`
}

// BowelMovementDetailsResponse represents the HTTP response for bowel movement details
type BowelMovementDetailsResponse struct {
	ID              string    `json:"id"`
	BowelMovementID string    `json:"bowelMovementId"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`

	Notes             string      `json:"notes,omitempty"`
	DetailedNotes     string      `json:"detailedNotes,omitempty"`
	Environment       string      `json:"environment,omitempty"`
	PreConditions     string      `json:"preConditions,omitempty"`
	PostConditions    string      `json:"postConditions,omitempty"`
	AIAnalysis        interface{} `json:"aiAnalysis,omitempty"`
	AIConfidence      *float64    `json:"aiConfidence,omitempty"`
	AIRecommendations string      `json:"aiRecommendations,omitempty"`
	Tags              []string    `json:"tags,omitempty"`
	WeatherCondition  string      `json:"weatherCondition,omitempty"`
	StressLevel       *int        `json:"stressLevel,omitempty"`
	SleepQuality      *int        `json:"sleepQuality,omitempty"`
	ExerciseIntensity *int        `json:"exerciseIntensity,omitempty"`
}

// BowelMovementStatsResponse represents the HTTP response for bowel movement statistics
type BowelMovementStatsResponse struct {
	TotalCount          int64       `json:"totalCount"`
	AveragePain         float64     `json:"averagePain"`
	AverageStrain       float64     `json:"averageStrain"`
	AverageSatisfaction float64     `json:"averageSatisfaction"`
	MostCommonBristol   int         `json:"mostCommonBristol"`
	FrequencyPerDay     float64     `json:"frequencyPerDay"`
	BristolDistribution map[int]int `json:"bristolDistribution"`
}

// ListBowelMovementsResponse represents the HTTP response for listing bowel movements
type ListBowelMovementsResponse struct {
	BowelMovements []BowelMovementResponse `json:"bowelMovements"`
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
