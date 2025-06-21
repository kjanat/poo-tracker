package insights

import (
	"fmt"
	"time"

	"github.com/kjanat/poo-tracker/backend/internal/domain/analytics"
	"github.com/kjanat/poo-tracker/backend/internal/domain/bowelmovement"
	"github.com/kjanat/poo-tracker/backend/internal/domain/meal"
	"github.com/kjanat/poo-tracker/backend/internal/domain/medication"
	"github.com/kjanat/poo-tracker/backend/internal/domain/symptom"
	"github.com/kjanat/poo-tracker/backend/internal/infrastructure/service/analytics/shared"
)

// InsightEngine generates health insights and recommendations
type InsightEngine struct{}

// NewInsightEngine creates a new insight engine
func NewInsightEngine() *InsightEngine {
	return &InsightEngine{}
}

// GenerateHealthInsights generates comprehensive health insights
func (ie *InsightEngine) GenerateHealthInsights(
	bowelMovements []bowelmovement.BowelMovement,
	meals []meal.Meal,
	symptoms []symptom.Symptom,
	medications []medication.Medication,
	correlations []*analytics.Correlation,
	trends []*shared.TrendLine,
) []*shared.InsightRecommendation {
	insights := make([]*shared.InsightRecommendation, 0)

	// Generate pattern-based insights
	patternInsights := ie.generatePatternInsights(bowelMovements, meals, symptoms)
	insights = append(insights, patternInsights...)

	// Generate correlation-based insights
	correlationInsights := ie.generateCorrelationInsights(correlations)
	insights = append(insights, correlationInsights...)

	// Generate trend-based insights
	trendInsights := ie.generateTrendInsights(trends)
	insights = append(insights, trendInsights...)

	// Generate medication insights
	medicationInsights := ie.generateMedicationInsights(medications, symptoms)
	insights = append(insights, medicationInsights...)

	// Sort by priority (High, Medium, Low)
	ie.sortInsightsByPriority(insights)

	return insights
}

// GetRecommendationStrings converts recommendations to string descriptions
func (ie *InsightEngine) GetRecommendationStrings(
	bowelValues []bowelmovement.BowelMovement,
	mealValues []meal.Meal,
	symptomValues []symptom.Symptom,
) []string {
	recs := ie.GenerateRecommendations(bowelValues, mealValues, symptomValues)
	var result []string
	for _, rec := range recs {
		result = append(result, rec.Message)
	}
	return result
}

// Pattern-based insights
func (ie *InsightEngine) generatePatternInsights(
	bowelMovements []bowelmovement.BowelMovement,
	meals []meal.Meal,
	symptoms []symptom.Symptom,
) []*shared.InsightRecommendation {
	insights := make([]*shared.InsightRecommendation, 0)

	// Bristol stool type consistency
	bristolInsight := ie.analyzeBristolConsistency(bowelMovements)
	if bristolInsight != nil {
		insights = append(insights, bristolInsight)
	}

	// Meal timing patterns
	timingInsight := ie.analyzeMealTimingPatterns(meals)
	if timingInsight != nil {
		insights = append(insights, timingInsight)
	}

	// Symptom clustering
	symptomInsight := ie.analyzeSymptomClustering(symptoms)
	if symptomInsight != nil {
		insights = append(insights, symptomInsight)
	}

	return insights
}

func (ie *InsightEngine) analyzeBristolConsistency(movements []bowelmovement.BowelMovement) *shared.InsightRecommendation {
	if len(movements) < 7 {
		return nil
	}

	// Calculate Bristol type distribution
	bristolCounts := make(map[int]int)
	for _, bm := range movements {
		bristolCounts[bm.BristolType]++
	}

	// Check for concerning patterns
	type1and2Count := bristolCounts[1] + bristolCounts[2] // Constipation
	type6and7Count := bristolCounts[6] + bristolCounts[7] // Diarrhea
	totalCount := len(movements)

	var insight *shared.InsightRecommendation

	if float64(type1and2Count)/float64(totalCount) > 0.5 {
		insight = &shared.InsightRecommendation{
			Type:       "LIFESTYLE",
			Category:   "Bowel Health",
			Message:    "Constipation Pattern Detected - Over 50% of your bowel movements indicate hard stools",
			Evidence:   fmt.Sprintf("%d out of %d movements were Bristol types 1-2, indicating constipation", type1and2Count, totalCount),
			Priority:   "HIGH",
			Confidence: 0.8,
			ActionItems: []string{
				"Increase fiber intake with fruits, vegetables, and whole grains",
				"Drink more water throughout the day",
				"Consider adding physical activity to your routine",
				"Consult with a healthcare provider if pattern persists",
			},
		}
	} else if float64(type6and7Count)/float64(totalCount) > 0.4 {
		insight = &shared.InsightRecommendation{
			Type:       "LIFESTYLE",
			Category:   "Bowel Health",
			Message:    "Loose Stool Pattern Detected - Over 40% of your bowel movements indicate loose stools",
			Evidence:   fmt.Sprintf("%d out of %d movements were Bristol types 6-7, indicating loose stools", type6and7Count, totalCount),
			Priority:   "MEDIUM",
			Confidence: 0.7,
			ActionItems: []string{
				"Keep a food diary to identify potential triggers",
				"Consider reducing dairy, gluten, or spicy foods temporarily",
				"Stay hydrated to replace lost fluids",
				"Consult with a healthcare provider if pattern continues",
			},
		}
	}

	return insight
}

func (ie *InsightEngine) analyzeMealTimingPatterns(meals []meal.Meal) *shared.InsightRecommendation {
	if len(meals) < 7 {
		return nil
	}

	// Analyze meal timing regularity
	hourCounts := make(map[int]int)
	for _, meal := range meals {
		hour := meal.MealTime.Hour()
		hourCounts[hour]++
	}

	// Check for irregular eating patterns
	var irregularHours []int
	var regularHours []int

	for hour, count := range hourCounts {
		if count >= 3 { // Regular if appears 3+ times in sample
			regularHours = append(regularHours, hour)
		} else if count == 1 {
			irregularHours = append(irregularHours, hour)
		}
	}

	if len(irregularHours) > len(regularHours) {
		return &shared.InsightRecommendation{
			ID:       fmt.Sprintf("meal-timing-%d", time.Now().Unix()),
			Type:     "pattern",
			Priority: "MEDIUM",
			Title:    "Irregular Meal Timing Pattern",
			Description: fmt.Sprintf("Your meal times appear inconsistent across %d meals. Regular meal timing can improve digestive health.",
				len(meals)),
			Evidence: []string{
				fmt.Sprintf("Meals recorded at %d different irregular times", len(irregularHours)),
				fmt.Sprintf("Only %d regular meal times identified", len(regularHours)),
				"Irregular eating can affect digestion and gut health",
			},
			ActionSteps: []string{
				"Try to eat meals at consistent times each day",
				"Aim for 3 main meals with 4-6 hour intervals",
				"Set meal reminders if needed",
				"Plan meals in advance to maintain routine",
			},
			CreatedAt: time.Now(),
		}
	}

	return nil
}

func (ie *InsightEngine) analyzeSymptomClustering(symptoms []symptom.Symptom) *shared.InsightRecommendation {
	if len(symptoms) < 5 {
		return nil
	}

	// Group symptoms by day
	dailySymptoms := make(map[string][]symptom.Symptom)
	for _, symptom := range symptoms {
		dayKey := symptom.RecordedAt.Format("2006-01-02")
		dailySymptoms[dayKey] = append(dailySymptoms[dayKey], symptom)
	}

	// Find days with multiple symptoms
	clusterDays := 0
	maxSymptomsInDay := 0
	for _, daySymptoms := range dailySymptoms {
		if len(daySymptoms) >= 3 {
			clusterDays++
		}
		if len(daySymptoms) > maxSymptomsInDay {
			maxSymptomsInDay = len(daySymptoms)
		}
	}

	if clusterDays >= 2 || maxSymptomsInDay >= 4 {
		return &shared.InsightRecommendation{
			ID:       fmt.Sprintf("symptom-clustering-%d", time.Now().Unix()),
			Type:     "pattern",
			Priority: "MEDIUM",
			Title:    "Symptom Clustering Detected",
			Description: fmt.Sprintf("Multiple symptoms occurring together on %d days, with up to %d symptoms in a single day.",
				clusterDays, maxSymptomsInDay),
			Evidence: []string{
				fmt.Sprintf("%d days with 3+ symptoms", clusterDays),
				fmt.Sprintf("Maximum %d symptoms in one day", maxSymptomsInDay),
				"Symptom clustering may indicate trigger events",
			},
			ActionSteps: []string{
				"Look for common triggers on high-symptom days",
				"Track what you ate 6-24 hours before symptom clusters",
				"Note stress levels, sleep quality, and activity on these days",
				"Consider discussing pattern with healthcare provider",
			},
			CreatedAt: time.Now(),
		}
	}

	return nil
}

// Correlation-based insights

func (ie *InsightEngine) generateCorrelationInsights(correlations []*analytics.Correlation) []*shared.InsightRecommendation {
	insights := make([]*shared.InsightRecommendation, 0)

	for _, corr := range correlations {
		if corr.Confidence > 0.7 && (corr.Strength > 0.5 || corr.Strength < -0.5) {
			insight := ie.createCorrelationInsight(corr)
			if insight != nil {
				insights = append(insights, insight)
			}
		}
	}

	return insights
}

func (ie *InsightEngine) createCorrelationInsight(corr *analytics.Correlation) *shared.InsightRecommendation {
	priority := "MEDIUM"
	if corr.Confidence > 0.8 {
		priority = "HIGH"
	}

	return &shared.InsightRecommendation{
		ID:          fmt.Sprintf("correlation-%s-%d", corr.Factor, time.Now().Unix()),
		Type:        "correlation",
		Priority:    priority,
		Title:       fmt.Sprintf("Strong Correlation: %s and %s", corr.Factor, corr.Outcome),
		Description: corr.Description,
		Evidence: []string{
			fmt.Sprintf("Correlation strength: %.3f", corr.Strength),
			fmt.Sprintf("Confidence level: %.1f%%", corr.Confidence*100),
			fmt.Sprintf("Based on %d data points", corr.SampleSize),
		},
		ActionSteps: ie.generateCorrelationActions(corr),
		CreatedAt:   time.Now(),
	}
}

func (ie *InsightEngine) generateCorrelationActions(corr *analytics.Correlation) []string {
	actions := make([]string, 0)

	if corr.Factor == "Spicy Food Level" && corr.Outcome == "Bowel Movement Pain" && corr.Strength > 0.3 {
		actions = append(actions, []string{
			"Consider reducing spicy food intake",
			"Try milder seasonings like herbs instead of hot spices",
			"Monitor pain levels when avoiding spicy foods",
			"Gradually reintroduce spicy foods to test tolerance",
		}...)
	} else if corr.Factor == "Fiber-Rich Meals" && corr.Outcome == "Bristol Stool Type" && corr.Strength > 0.3 {
		actions = append(actions, []string{
			"Continue including fiber-rich foods in your diet",
			"Aim for 25-35g of fiber daily from various sources",
			"Increase water intake with higher fiber consumption",
			"Track which fiber sources work best for you",
		}...)
	} else if corr.Factor == "Dairy Consumption" && corr.Outcome == "Symptom Severity" && corr.Strength > 0.3 {
		actions = append(actions, []string{
			"Consider reducing dairy intake temporarily",
			"Try lactose-free alternatives",
			"Monitor symptoms during dairy elimination",
			"Discuss potential lactose intolerance with healthcare provider",
		}...)
	} else {
		// Generic actions for other correlations
		actions = append(actions, []string{
			fmt.Sprintf("Monitor %s intake and its effects", corr.Factor),
			fmt.Sprintf("Track changes in %s when modifying %s", corr.Outcome, corr.Factor),
			"Keep detailed records to confirm this pattern",
			"Discuss findings with healthcare provider if concerning",
		}...)
	}

	return actions
}

// Trend-based insights

func (ie *InsightEngine) generateTrendInsights(trends []*shared.TrendLine) []*shared.InsightRecommendation {
	insights := make([]*shared.InsightRecommendation, 0)

	for _, trend := range trends {
		if trend.Confidence > 0.6 && trend.Direction != "stable" {
			insight := ie.createTrendInsight(trend)
			if insight != nil {
				insights = append(insights, insight)
			}
		}
	}

	return insights
}

func (ie *InsightEngine) createTrendInsight(trend *shared.TrendLine) *shared.InsightRecommendation {
	priority := "MEDIUM"
	if trend.Confidence > 0.8 {
		priority = "HIGH"
	}

	return &shared.InsightRecommendation{
		ID:       fmt.Sprintf("trend-%s-%d", trend.Name, time.Now().Unix()),
		Type:     "trend",
		Priority: priority,
		Title:    fmt.Sprintf("%s Trend: %s", trend.Name, trend.Direction),
		Description: fmt.Sprintf("Your %s shows a %s trend over time with %s confidence.",
			trend.Name, trend.Direction, trend.Significance),
		Evidence: []string{
			fmt.Sprintf("Trend direction: %s", trend.Direction),
			fmt.Sprintf("Confidence: %.1f%%", trend.Confidence*100),
			fmt.Sprintf("Based on %d data points", len(trend.Points)),
		},
		ActionSteps: ie.generateTrendActions(trend),
		CreatedAt:   time.Now(),
	}
}

func (ie *InsightEngine) generateTrendActions(trend *shared.TrendLine) []string {
	actions := make([]string, 0)

	switch trend.Name {
	case "Pain":
		if trend.Direction == "improving" {
			actions = append(actions, "Continue current pain management strategies")
		} else {
			actions = append(actions, "Consider reviewing pain triggers and management")
		}
	case "Satisfaction":
		if trend.Direction == "declining" {
			actions = append(actions, "Evaluate factors affecting bowel movement satisfaction")
		}
	case "Symptom Severity":
		if trend.Direction == "declining" {
			actions = append(actions, "Continue current symptom management approach")
		} else {
			actions = append(actions, "Review and adjust symptom management strategies")
		}
	default:
		actions = append(actions, fmt.Sprintf("Monitor %s trend and related factors", trend.Name))
	}

	actions = append(actions, "Track progress over the next few weeks")
	return actions
}

// Medication insights

func (ie *InsightEngine) generateMedicationInsights(
	medications []medication.Medication,
	symptoms []symptom.Symptom,
) []*shared.InsightRecommendation {
	insights := make([]*shared.InsightRecommendation, 0)

	// Check for inactive medications
	inactiveCount := 0
	for _, med := range medications {
		if !med.IsActive {
			inactiveCount++
		}
	}

	if inactiveCount > 0 && len(medications) > 0 {
		adherencePercent := float64(len(medications)-inactiveCount) / float64(len(medications)) * 100

		if adherencePercent < 80 {
			insight := &shared.InsightRecommendation{
				ID:       fmt.Sprintf("medication-adherence-%d", time.Now().Unix()),
				Type:     "medication",
				Priority: "HIGH",
				Title:    "Low Medication Adherence Detected",
				Description: fmt.Sprintf("Only %.1f%% of your medications are currently active. Medication adherence is important for treatment effectiveness.",
					adherencePercent),
				Evidence: []string{
					fmt.Sprintf("%d of %d medications are inactive", inactiveCount, len(medications)),
					"Low adherence may affect treatment outcomes",
					"Regular medication use is often necessary for optimal results",
				},
				ActionSteps: []string{
					"Review inactive medications with healthcare provider",
					"Set medication reminders if forgetfulness is an issue",
					"Discuss any side effects or concerns about medications",
					"Consider pill organizers or medication tracking apps",
				},
				CreatedAt: time.Now(),
			}
			insights = append(insights, insight)
		}
	}

	return insights
}

// Utility functions

func (ie *InsightEngine) convertInsightToRecommendation(insight *shared.InsightRecommendation) *analytics.Recommendation {
	return &analytics.Recommendation{
		Type:        insight.Type,
		Priority:    insight.Priority,
		Title:       insight.Title,
		Description: insight.Description,
		ActionSteps: insight.ActionSteps,
		Confidence:  0.8, // Default confidence for high-priority insights
		Evidence:    insight.Evidence,
	}
}

func (ie *InsightEngine) sortInsightsByPriority(insights []*shared.InsightRecommendation) {
	// Simple sorting: HIGH -> MEDIUM -> LOW
	highPriority := make([]*shared.InsightRecommendation, 0)
	mediumPriority := make([]*shared.InsightRecommendation, 0)
	lowPriority := make([]*shared.InsightRecommendation, 0)

	for _, insight := range insights {
		switch insight.Priority {
		case "HIGH":
			highPriority = append(highPriority, insight)
		case "MEDIUM":
			mediumPriority = append(mediumPriority, insight)
		default:
			lowPriority = append(lowPriority, insight)
		}
	}

	// Clear and refill the original slice in place
	n := 0
	for _, insight := range highPriority {
		insights[n] = insight
		n++
	}
	for _, insight := range mediumPriority {
		insights[n] = insight
		n++
	}
	for _, insight := range lowPriority {
		insights[n] = insight
		n++
	}
}
