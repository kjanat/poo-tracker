package analyzer

import (
	"fmt"
	"math"
	"time"

	"github.com/kjanat/poo-tracker/backend/internal/domain/analytics"
	"github.com/kjanat/poo-tracker/backend/internal/domain/bowelmovement"
	"github.com/kjanat/poo-tracker/backend/internal/domain/meal"
	"github.com/kjanat/poo-tracker/backend/internal/domain/medication"
	"github.com/kjanat/poo-tracker/backend/internal/domain/symptom"
	"github.com/kjanat/poo-tracker/backend/internal/infrastructure/service/analytics/shared"
)

// CorrelationAnalyzer analyzes correlations between different data types
type CorrelationAnalyzer struct{}

// NewCorrelationAnalyzer creates a new correlation analyzer
func NewCorrelationAnalyzer() *CorrelationAnalyzer {
	return &CorrelationAnalyzer{}
}

// CalculateMealBowelCorrelations analyzes correlations between meals and bowel movements
func (ca *CorrelationAnalyzer) CalculateMealBowelCorrelations(
	meals []meal.Meal,
	bowelMovements []bowelmovement.BowelMovement,
) []*analytics.Correlation {
	correlations := make([]*analytics.Correlation, 0)

	// Group data by day for correlation analysis
	dailyData := ca.groupMealBowelDataByDay(meals, bowelMovements)

	// Analyze fiber-Bristol correlation
	fiberBristolCorr := ca.analyzeFiberBristolCorrelation(dailyData)
	if fiberBristolCorr != nil {
		correlations = append(correlations, fiberBristolCorr)
	}

	// Analyze calorie-satisfaction correlation
	calorieSatisfactionCorr := ca.analyzeCalorieSatisfactionCorrelation(dailyData)
	if calorieSatisfactionCorr != nil {
		correlations = append(correlations, calorieSatisfactionCorr)
	}

	// Analyze spicy food-pain correlation
	spicyPainCorr := ca.analyzeSpicyPainCorrelation(dailyData)
	if spicyPainCorr != nil {
		correlations = append(correlations, spicyPainCorr)
	}

	return correlations
}

// CalculateMealSymptomCorrelations analyzes correlations between meals and symptoms
func (ca *CorrelationAnalyzer) CalculateMealSymptomCorrelations(
	meals []meal.Meal,
	symptoms []symptom.Symptom,
) []*analytics.Correlation {
	correlations := make([]*analytics.Correlation, 0)

	// Group data by day for correlation analysis
	dailyData := ca.groupMealSymptomDataByDay(meals, symptoms)

	// Analyze spicy food-symptom correlation
	spicySymptomCorr := ca.analyzeSpicyFoodSymptomCorrelation(dailyData)
	if spicySymptomCorr != nil {
		correlations = append(correlations, spicySymptomCorr)
	}

	// Analyze dairy-symptom correlation
	dairySymptomCorr := ca.analyzeDairySymptomCorrelation(dailyData)
	if dairySymptomCorr != nil {
		correlations = append(correlations, dairySymptomCorr)
	}

	// Analyze gluten-symptom correlation
	glutenSymptomCorr := ca.analyzeGlutenSymptomCorrelation(dailyData)
	if glutenSymptomCorr != nil {
		correlations = append(correlations, glutenSymptomCorr)
	}

	return correlations
}

// CalculateMedicationEffectiveness analyzes medication effectiveness
func (ca *CorrelationAnalyzer) CalculateMedicationEffectiveness(
	medications []medication.Medication,
	symptoms []symptom.Symptom,
	bowelMovements []bowelmovement.BowelMovement,
) []*analytics.MedicationEffect {
	effectiveness := make([]*analytics.MedicationEffect, 0)

	// Group medications by name
	medicationGroups := ca.groupMedicationsByName(medications)

	for medName, meds := range medicationGroups {
		effect := ca.analyzeIndividualMedicationEffectiveness(medName, meds, symptoms, bowelMovements)
		if effect != nil {
			effectiveness = append(effectiveness, effect)
		}
	}

	return effectiveness
}

// Helper functions for meal-bowel correlations

type DailyMealBowelData struct {
	Date            time.Time
	FiberRichMeals  int
	TotalCalories   int
	SpicyLevel      float64
	BristolAvg      float64
	PainAvg         float64
	SatisfactionAvg float64
}

func (ca *CorrelationAnalyzer) groupMealBowelDataByDay(
	meals []meal.Meal,
	bowelMovements []bowelmovement.BowelMovement,
) []DailyMealBowelData {
	dataMap := make(map[string]*DailyMealBowelData)

	// Process meals
	for _, meal := range meals {
		dayKey := meal.MealTime.Format("2006-01-02")
		if data, exists := dataMap[dayKey]; exists {
			data.TotalCalories += meal.Calories
			if meal.FiberRich {
				data.FiberRichMeals++
			}
			if meal.SpicyLevel != nil {
				data.SpicyLevel = math.Max(data.SpicyLevel, float64(*meal.SpicyLevel))
			}
		} else {
			spicyLevel := 0.0
			if meal.SpicyLevel != nil {
				spicyLevel = float64(*meal.SpicyLevel)
			}
			fiberCount := 0
			if meal.FiberRich {
				fiberCount = 1
			}
			dataMap[dayKey] = &DailyMealBowelData{
				Date:           meal.MealTime,
				FiberRichMeals: fiberCount,
				TotalCalories:  meal.Calories,
				SpicyLevel:     spicyLevel,
			}
		}
	}

	// Process bowel movements
	bowelDataByDay := make(map[string][]bowelmovement.BowelMovement)
	for _, bm := range bowelMovements {
		dayKey := bm.RecordedAt.Format("2006-01-02")
		bowelDataByDay[dayKey] = append(bowelDataByDay[dayKey], bm)
	}

	// Calculate averages for each day
	for dayKey, dayBMs := range bowelDataByDay {
		if data, exists := dataMap[dayKey]; exists {
			bristolSum, painSum, satisfactionSum := 0.0, 0.0, 0.0
			for _, bm := range dayBMs {
				bristolSum += float64(bm.BristolType)
				painSum += float64(bm.Pain)
				satisfactionSum += float64(bm.Satisfaction)
			}
			count := float64(len(dayBMs))
			data.BristolAvg = bristolSum / count
			data.PainAvg = painSum / count
			data.SatisfactionAvg = satisfactionSum / count
		}
	}

	// Convert map to slice
	result := make([]DailyMealBowelData, 0, len(dataMap))
	for _, data := range dataMap {
		result = append(result, *data)
	}

	return result
}

func (ca *CorrelationAnalyzer) analyzeFiberBristolCorrelation(data []DailyMealBowelData) *analytics.Correlation {
	if len(data) < 3 {
		return nil
	}

	fiberValues := make([]float64, len(data))
	bristolValues := make([]float64, len(data))

	for i, d := range data {
		fiberValues[i] = float64(d.FiberRichMeals)
		bristolValues[i] = d.BristolAvg
	}

	coefficient := shared.CalculateCorrelation(fiberValues, bristolValues)
	strength := shared.InterpretCorrelationStrength(coefficient)
	confidence := shared.CalculateConfidenceScore(len(data))

	return &analytics.Correlation{
		Factor:      "Fiber-Rich Meals",
		Outcome:     "Bristol Stool Type",
		Strength:    shared.RoundToDecimalPlaces(coefficient, 3),
		Confidence:  confidence,
		Description: ca.generateCorrelationDescription("fiber intake", "stool consistency", coefficient, strength),
		SampleSize:  len(data),
	}
}

func (ca *CorrelationAnalyzer) analyzeCalorieSatisfactionCorrelation(data []DailyMealBowelData) *analytics.Correlation {
	if len(data) < 3 {
		return nil
	}

	calorieValues := make([]float64, len(data))
	satisfactionValues := make([]float64, len(data))

	for i, d := range data {
		calorieValues[i] = float64(d.TotalCalories)
		satisfactionValues[i] = d.SatisfactionAvg
	}

	coefficient := shared.CalculateCorrelation(calorieValues, satisfactionValues)
	strength := shared.InterpretCorrelationStrength(coefficient)
	confidence := shared.CalculateConfidenceScore(len(data))

	return &analytics.Correlation{
		Factor:      "Daily Calorie Intake",
		Outcome:     "Bowel Movement Satisfaction",
		Strength:    shared.RoundToDecimalPlaces(coefficient, 3),
		Confidence:  confidence,
		Description: ca.generateCorrelationDescription("calorie intake", "satisfaction", coefficient, strength),
		SampleSize:  len(data),
	}
}

func (ca *CorrelationAnalyzer) analyzeSpicyPainCorrelation(data []DailyMealBowelData) *analytics.Correlation {
	if len(data) < 3 {
		return nil
	}

	spicyValues := make([]float64, len(data))
	painValues := make([]float64, len(data))

	for i, d := range data {
		spicyValues[i] = d.SpicyLevel
		painValues[i] = d.PainAvg
	}

	coefficient := shared.CalculateCorrelation(spicyValues, painValues)
	strength := shared.InterpretCorrelationStrength(coefficient)
	confidence := shared.CalculateConfidenceScore(len(data))

	return &analytics.Correlation{
		Factor:      "Spicy Food Level",
		Outcome:     "Bowel Movement Pain",
		Strength:    shared.RoundToDecimalPlaces(coefficient, 3),
		Confidence:  confidence,
		Description: ca.generateCorrelationDescription("spicy food intake", "pain levels", coefficient, strength),
		SampleSize:  len(data),
	}
}

// Helper functions for meal-symptom correlations

type DailyMealSymptomData struct {
	Date            time.Time
	SpicyMeals      int
	DairyMeals      int
	GlutenMeals     int
	SymptomSeverity float64
	SymptomCount    int
}

func (ca *CorrelationAnalyzer) groupMealSymptomDataByDay(
	meals []meal.Meal,
	symptoms []symptom.Symptom,
) []DailyMealSymptomData {
	dataMap := make(map[string]*DailyMealSymptomData)

	// Process meals
	for _, meal := range meals {
		dayKey := meal.MealTime.Format("2006-01-02")
		if data, exists := dataMap[dayKey]; exists {
			if meal.SpicyLevel != nil && *meal.SpicyLevel > 2 {
				data.SpicyMeals++
			}
			if meal.Dairy {
				data.DairyMeals++
			}
			if meal.Gluten {
				data.GlutenMeals++
			}
		} else {
			spicyCount := 0
			if meal.SpicyLevel != nil && *meal.SpicyLevel > 2 {
				spicyCount = 1
			}
			dairyCount := 0
			if meal.Dairy {
				dairyCount = 1
			}
			glutenCount := 0
			if meal.Gluten {
				glutenCount = 1
			}

			dataMap[dayKey] = &DailyMealSymptomData{
				Date:        meal.MealTime,
				SpicyMeals:  spicyCount,
				DairyMeals:  dairyCount,
				GlutenMeals: glutenCount,
			}
		}
	}

	// Process symptoms
	for _, symptom := range symptoms {
		dayKey := symptom.RecordedAt.Format("2006-01-02")
		if data, exists := dataMap[dayKey]; exists {
			data.SymptomCount++
			data.SymptomSeverity += float64(symptom.Severity)
		}
	}

	// Calculate averages
	for _, data := range dataMap {
		if data.SymptomCount > 0 {
			data.SymptomSeverity /= float64(data.SymptomCount)
		}
	}

	// Convert map to slice
	result := make([]DailyMealSymptomData, 0, len(dataMap))
	for _, data := range dataMap {
		result = append(result, *data)
	}

	return result
}

func (ca *CorrelationAnalyzer) analyzeSpicyFoodSymptomCorrelation(data []DailyMealSymptomData) *analytics.Correlation {
	if len(data) < 3 {
		return nil
	}

	spicyValues := make([]float64, len(data))
	symptomValues := make([]float64, len(data))

	for i, d := range data {
		spicyValues[i] = float64(d.SpicyMeals)
		symptomValues[i] = d.SymptomSeverity
	}

	coefficient := shared.CalculateCorrelation(spicyValues, symptomValues)
	strength := shared.InterpretCorrelationStrength(coefficient)
	confidence := shared.CalculateConfidenceScore(len(data))

	return &analytics.Correlation{
		Factor:      "Spicy Food Consumption",
		Outcome:     "Symptom Severity",
		Strength:    shared.RoundToDecimalPlaces(coefficient, 3),
		Confidence:  confidence,
		Description: ca.generateCorrelationDescription("spicy food", "symptom severity", coefficient, strength),
		SampleSize:  len(data),
	}
}

func (ca *CorrelationAnalyzer) analyzeDairySymptomCorrelation(data []DailyMealSymptomData) *analytics.Correlation {
	if len(data) < 3 {
		return nil
	}

	dairyValues := make([]float64, len(data))
	symptomValues := make([]float64, len(data))

	for i, d := range data {
		dairyValues[i] = float64(d.DairyMeals)
		symptomValues[i] = d.SymptomSeverity
	}

	coefficient := shared.CalculateCorrelation(dairyValues, symptomValues)
	strength := shared.InterpretCorrelationStrength(coefficient)
	confidence := shared.CalculateConfidenceScore(len(data))

	return &analytics.Correlation{
		Factor:      "Dairy Consumption",
		Outcome:     "Symptom Severity",
		Strength:    shared.RoundToDecimalPlaces(coefficient, 3),
		Confidence:  confidence,
		Description: ca.generateCorrelationDescription("dairy consumption", "symptom severity", coefficient, strength),
		SampleSize:  len(data),
	}
}

func (ca *CorrelationAnalyzer) analyzeGlutenSymptomCorrelation(data []DailyMealSymptomData) *analytics.Correlation {
	if len(data) < 3 {
		return nil
	}

	glutenValues := make([]float64, len(data))
	symptomValues := make([]float64, len(data))

	for i, d := range data {
		glutenValues[i] = float64(d.GlutenMeals)
		symptomValues[i] = d.SymptomSeverity
	}

	coefficient := shared.CalculateCorrelation(glutenValues, symptomValues)
	strength := shared.InterpretCorrelationStrength(coefficient)
	confidence := shared.CalculateConfidenceScore(len(data))

	return &analytics.Correlation{
		Factor:      "Gluten Consumption",
		Outcome:     "Symptom Severity",
		Strength:    shared.RoundToDecimalPlaces(coefficient, 3),
		Confidence:  confidence,
		Description: ca.generateCorrelationDescription("gluten consumption", "symptom severity", coefficient, strength),
		SampleSize:  len(data),
	}
}

// Medication effectiveness analysis

func (ca *CorrelationAnalyzer) groupMedicationsByName(medications []medication.Medication) map[string][]medication.Medication {
	groups := make(map[string][]medication.Medication)

	for _, med := range medications {
		groups[med.Name] = append(groups[med.Name], med)
	}

	return groups
}

func (ca *CorrelationAnalyzer) analyzeIndividualMedicationEffectiveness(
	medName string,
	medications []medication.Medication,
	symptoms []symptom.Symptom,
	bowelMovements []bowelmovement.BowelMovement,
) *analytics.MedicationEffect {
	if len(medications) == 0 {
		return nil
	}

	// Find medication start dates
	var startDate *time.Time
	for _, med := range medications {
		if med.StartDate != nil {
			if startDate == nil || med.StartDate.Before(*startDate) {
				startDate = med.StartDate
			}
		}
	}

	if startDate == nil {
		return nil
	}

	// Analyze symptom improvement before/after medication
	symptomImprovement := ca.analyzeSymptomImprovement(symptoms, *startDate)
	bowelImprovement := ca.analyzeBowelImprovement(bowelMovements, *startDate)

	// Calculate overall effectiveness score
	effectivenessScore := ca.calculateMedicationEffectivenessScore(symptomImprovement, bowelImprovement)

	return &analytics.MedicationEffect{
		MedicationName:     medName,
		SymptomImprovement: symptomImprovement,
		BowelImprovement:   bowelImprovement,
		EffectivenessScore: shared.RoundToDecimalPlaces(effectivenessScore, 1),
	}
}

func (ca *CorrelationAnalyzer) analyzeSymptomImprovement(symptoms []symptom.Symptom, startDate time.Time) map[string]float64 {
	improvement := make(map[string]float64)

	beforeSymptoms := make([]symptom.Symptom, 0)
	afterSymptoms := make([]symptom.Symptom, 0)

	for _, symptom := range symptoms {
		if symptom.RecordedAt.Before(startDate) {
			beforeSymptoms = append(beforeSymptoms, symptom)
		} else {
			afterSymptoms = append(afterSymptoms, symptom)
		}
	}

	// Calculate improvement (simplified)
	beforeSeverity := ca.calculateAverageSymptomSeverity(beforeSymptoms)
	afterSeverity := ca.calculateAverageSymptomSeverity(afterSymptoms)

	improvementPercent := 0.0
	if beforeSeverity > 0 {
		improvementPercent = ((beforeSeverity - afterSeverity) / beforeSeverity) * 100
	}

	improvement["overall"] = shared.RoundToDecimalPlaces(improvementPercent, 1)

	return improvement
}

func (ca *CorrelationAnalyzer) analyzeBowelImprovement(bowelMovements []bowelmovement.BowelMovement, startDate time.Time) float64 {
	beforeMovements := make([]bowelmovement.BowelMovement, 0)
	afterMovements := make([]bowelmovement.BowelMovement, 0)

	for _, bm := range bowelMovements {
		if bm.RecordedAt.Before(startDate) {
			beforeMovements = append(beforeMovements, bm)
		} else {
			afterMovements = append(afterMovements, bm)
		}
	}

	beforeScore := ca.calculateBowelScore(beforeMovements)
	afterScore := ca.calculateBowelScore(afterMovements)

	improvement := afterScore - beforeScore
	return shared.RoundToDecimalPlaces(improvement, 1)
}

func (ca *CorrelationAnalyzer) calculateAverageSymptomSeverity(symptoms []symptom.Symptom) float64 {
	if len(symptoms) == 0 {
		return 0
	}

	total := 0.0
	for _, symptom := range symptoms {
		total += float64(symptom.Severity)
	}

	return total / float64(len(symptoms))
}

func (ca *CorrelationAnalyzer) calculateBowelScore(movements []bowelmovement.BowelMovement) float64 {
	if len(movements) == 0 {
		return 0
	}

	score := 0.0
	for _, bm := range movements {
		// Simple scoring: lower pain and strain, higher satisfaction
		bmScore := float64(bm.Satisfaction) - float64(bm.Pain) - float64(bm.Strain)
		score += bmScore
	}

	return score / float64(len(movements))
}

func (ca *CorrelationAnalyzer) calculateMedicationEffectivenessScore(
	symptomImprovement map[string]float64,
	bowelImprovement float64,
) float64 {
	score := 50.0 // Base score

	// Add symptom improvement
	if overall, exists := symptomImprovement["overall"]; exists {
		score += overall * 0.3 // 30% weight for symptom improvement
	}

	// Add bowel improvement
	score += bowelImprovement * 0.2 // 20% weight for bowel improvement

	return math.Max(0, math.Min(100, score))
}

// Utility functions

func (ca *CorrelationAnalyzer) generateCorrelationDescription(factor, outcome string, coefficient float64, strength string) string {
	direction := "positive"
	if coefficient < 0 {
		direction = "negative"
	}

	return fmt.Sprintf("There is a %s %s correlation between %s and %s (r=%.3f).",
		strength, direction, factor, outcome, coefficient)
}
