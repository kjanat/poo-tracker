# Analytics Domain

The analytics domain provides comprehensive health analytics functionality for the Poo Tracker application. This package processes data from multiple domains (bowel movements, meals, symptoms, medications) to generate insights, correlations, trends, and personalized recommendations.

## Architecture

```tree
analytics/
├── models.go          # Core data structures for analytics
├── service.go         # Service interface and complex domain models
├── errors.go          # Domain-specific error definitions
└── README.md          # This documentation
```

## Core Functionality

### 1. Health Overview (`GetUserHealthOverview`)

Generates comprehensive health summaries by aggregating data across all health domains:

**Input**:

- User ID
- Date range (start, end)

**Output**:

- `HealthOverview` containing:
  - Bowel movement statistics (frequency, Bristol scale trends, regularity)
  - Meal statistics (calories, fiber intake, dietary patterns)
  - Symptom statistics (frequency, severity trends, common types)
  - Medication statistics (adherence, effectiveness)
  - Overall health score and trend direction

### 2. Correlation Analysis (`GetCorrelationAnalysis`)

Identifies statistical relationships between different health factors:

**Correlations Analyzed**:

- Meal ingredients → Bowel movement quality
- Meal timing → Symptom occurrence
- Medication effectiveness → Symptom reduction
- Dietary triggers → Digestive issues

**Output**:

- `CorrelationAnalysis` with strength, confidence, and sample size for each correlation

### 3. Trend Analysis (`GetTrendAnalysis`)

Performs time-series analysis to identify health trends over time:

**Trends Tracked**:

- Bowel movement regularity and quality over time
- Symptom severity and frequency patterns
- Dietary adherence and nutrition trends
- Overall health trajectory

**Output**:

- `TrendAnalysis` with direction (improving/stable/declining), slope, and confidence

### 4. Behavior Pattern Recognition (`GetBehaviorPatterns`)

Identifies patterns in user behavior and their health impacts:

**Pattern Types**:

- **Eating Patterns**: Meal timing, portion sizes, cuisine preferences, consistency
- **Bowel Patterns**: Preferred timing, regularity, response to meals
- **Symptom Patterns**: Timing, triggers, seasonal variations, cyclical patterns
- **Lifestyle Patterns**: Stress levels, sleep quality, exercise correlation

### 5. Health Insights (`GetHealthInsights`)

Generates actionable insights based on comprehensive data analysis:

**Insight Categories**:

- Key findings from recent data
- Identified risk factors
- Positive health factors to maintain
- Personalized recommendations with evidence
- Alert level assessment (LOW/MEDIUM/HIGH)

### 6. Health Scoring (`GetHealthScore`)

Calculates overall health scores with component breakdown:

**Score Components**:

- Bowel Health (40% weight): Regularity, Bristol scale consistency
- Symptom Control (30% weight): Frequency and severity management
- Nutrition (20% weight): Dietary quality and consistency
- Medication Adherence (10% weight): Compliance and effectiveness

**Output**:

- Overall score (0-100)
- Component scores with explanations
- Trend direction and factors affecting score
- Benchmarks against population averages

### 7. Personalized Recommendations (`GetRecommendations`)

Generates evidence-based health recommendations:

**Recommendation Types**:

- **DIETARY**: Nutrition and meal timing suggestions
- **LIFESTYLE**: Exercise, stress management, routine optimization
- **MEDICAL**: When to consult healthcare providers
- **TRACKING**: Improved data collection strategies

**Recommendation Structure**:

- Priority level (LOW/MEDIUM/HIGH)
- Confidence score (0-1)
- Supporting evidence
- Specific action steps
- Expected impact and timeline

## Data Models

### Core Analytics Models

#### `BowelMovementSummary`

```go
type BowelMovementSummary struct {
    TotalCount          int64   // Total movements in period
    AveragePerDay       float64 // Daily frequency
    MostCommonBristol   int     // Most common Bristol type (1-7)
    AveragePain         float64 // Average pain level (1-10)
    AverageStrain       float64 // Average strain level (1-10)
    AverageSatisfaction float64 // Average satisfaction (1-10)
    RegularityScore     float64 // Regularity metric (0-1)
}
```

#### `MealSummary`

```go
type MealSummary struct {
    TotalMeals       int64   // Number of meals tracked
    AveragePerDay    float64 // Meals per day
    TotalCalories    int     // Total calories consumed
    AverageCalories  float64 // Average calories per meal
    FiberRichPercent float64 // Percentage of fiber-rich meals
    HealthScore      float64 // Overall dietary health score
}
```

#### `SymptomSummary`

```go
type SymptomSummary struct {
    TotalSymptoms      int64   // Total symptoms tracked
    AveragePerDay      float64 // Symptoms per day
    AverageSeverity    float64 // Average severity (1-10)
    MostCommonCategory string  // Most frequent symptom category
    MostCommonType     string  // Most frequent symptom type
    TrendDirection     string  // IMPROVING/STABLE/DECLINING
}
```

### Pattern Recognition Models

#### `EatingPattern`

```go
type EatingPattern struct {
    MealTiming           map[string]float64 // Hour -> frequency
    MealSizeDistribution map[string]float64 // Size -> frequency
    PreferredCuisines    []string           // Preferred food types
    DietaryConsistency   float64            // Consistency score (0-1)
}
```

#### `BowelPattern`

```go
type BowelPattern struct {
    PreferredTiming     map[string]float64 // Hour -> frequency
    RegularityScore     float64            // Regularity metric (0-1)
    ConsistencyPatterns map[string]float64 // Bristol type -> frequency
    ResponseToMeals     float64            // Hours after meals
}
```

### Analysis Results

#### `Correlation`

```go
type Correlation struct {
    Factor      string  // Input factor (e.g., "dairy_intake")
    Outcome     string  // Outcome measure (e.g., "digestive_symptoms")
    Strength    float64 // Correlation strength (-1 to 1)
    Confidence  float64 // Statistical confidence (0 to 1)
    Description string  // Human-readable explanation
    SampleSize  int     // Number of data points
}
```

#### `DataTrend`

```go
type DataTrend struct {
    Direction   string             // IMPROVING/STABLE/DECLINING
    Slope       float64            // Rate of change
    Confidence  float64            // Trend confidence (0-1)
    TimePoints  []time.Time        // Time series data points
    Values      []float64          // Corresponding values
    Seasonality map[string]float64 // Seasonal patterns
}
```

## Usage Examples

### Getting Health Overview

```go
// Get comprehensive health overview for last 30 days
end := time.Now()
start := end.AddDate(0, 0, -30)
overview, err := analyticsService.GetUserHealthOverview(ctx, userID, start, end)
if err != nil {
    return err
}

// Access aggregated statistics
fmt.Printf("Bowel movements per day: %.1f\n", overview.BowelMovementStats.AveragePerDay)
fmt.Printf("Overall health score: %.1f\n", overview.OverallHealthScore)
fmt.Printf("Health trend: %s\n", overview.TrendDirection)
```

### Analyzing Correlations

```go
// Identify meal-symptom correlations
correlations, err := analyticsService.GetCorrelationAnalysis(ctx, userID, start, end)
if err != nil {
    return err
}

// Review significant correlations
for _, correlation := range correlations.MealSymptomCorrelations {
    if correlation.Confidence > 0.7 {
        fmt.Printf("Found correlation: %s → %s (strength: %.2f)\n",
            correlation.Factor, correlation.Outcome, correlation.Strength)
    }
}
```

### Getting Recommendations

```go
// Get personalized recommendations
recommendations, err := analyticsService.GetRecommendations(ctx, userID)
if err != nil {
    return err
}

// Process high-priority recommendations
for _, rec := range recommendations {
    if rec.Priority == "HIGH" {
        fmt.Printf("High priority: %s\n", rec.Title)
        fmt.Printf("Actions: %v\n", rec.ActionSteps)
    }
}
```

## Error Handling

The analytics domain defines specific errors for different failure scenarios:

- `ErrInsufficientData`: Not enough data for meaningful analysis
- `ErrInvalidDateRange`: Invalid date range parameters
- `ErrCorrelationFailed`: Correlation analysis failed
- `ErrTrendAnalysisFailed`: Trend analysis failed
- `ErrUserNotAuthorized`: User lacks access to analytics
- `ErrInvalidAnalysisType`: Unsupported analysis type

## Integration Points

### Data Sources

- **Bowel Movement Domain**: Frequency, Bristol scale, satisfaction metrics
- **Meal Domain**: Nutritional data, timing, ingredients, calories
- **Symptom Domain**: Types, severity, timing, triggers
- **Medication Domain**: Adherence, effectiveness, side effects

### Consumers

- **REST API**: Exposes analytics endpoints for frontend
- **Background Jobs**: Scheduled analysis and insight generation
- **Notification Service**: Alert generation based on health insights
- **Reporting Service**: Health summary reports and trends

## Performance Considerations

### Data Volume

- Designed to handle years of user health data
- Efficient aggregation algorithms for large datasets
- Configurable date ranges to limit analysis scope

### Caching Strategy

- Health scores cached for 24 hours
- Trend analysis cached for 1 hour
- Real-time insights for immediate feedback

### Scalability

- Stateless service design for horizontal scaling
- Database queries optimized with proper indexing
- Background processing for computationally intensive analysis

## Future Enhancements

### Planned Features

1. **Machine Learning Integration**: Advanced pattern recognition and prediction
2. **Comparative Analytics**: Anonymous population-level comparisons
3. **Export Capabilities**: PDF health reports and data export
4. **Real-time Monitoring**: Live health metric tracking
5. **Integration APIs**: Third-party health app connectivity

### Analytics Improvements

1. **Seasonal Analysis**: Long-term seasonal pattern recognition
2. **Cohort Analysis**: Group-based health trend analysis
3. **Predictive Modeling**: Health outcome prediction
4. **Intervention Tracking**: Measure impact of lifestyle changes
