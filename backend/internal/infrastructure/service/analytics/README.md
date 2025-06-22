# Infrastructure Analytics Services

The infrastructure analytics services provide the implementation layer for analytics functionality, containing specialized components for data aggregation, analysis, and insight generation.

## Architecture

```graphql
infrastructure/service/analytics/
├── analytics_service.go    # Main analytics service implementation
├── analytics.go           # Type aliases and shared definitions
├── aggregator/            # Data aggregation components
│   ├── data_aggregator.go # Core data aggregation logic
│   └── statistics.go      # Statistical calculation utilities
├── analyzer/              # Analysis and correlation engines
│   ├── correlation.go     # Correlation analysis implementation
│   ├── patterns.go        # Pattern recognition algorithms
│   ├── trend_analysis.go  # Trend analysis implementation
│   └── trends.go          # Trend calculation utilities
├── calculator/            # Health scoring and metrics
│   ├── calculators.go     # Legacy calculation utilities
│   └── health_scores.go   # Health score calculation engine
├── insights/              # Insight generation and recommendations
│   ├── insight_engine.go  # Main insight generation engine
│   └── recommendation_generator.go # Personalized recommendations
└── shared/                # Shared utilities and models
    ├── models.go          # Common data structures
    └── patterns.go        # Pattern analysis utilities
```

## Component Overview

### AnalyticsService (`analytics_service.go`)

**Purpose**: Main orchestrator that implements the domain Service interface
**Responsibilities**:

- Coordinates data retrieval from multiple repositories
- Delegates specialized analysis to subcomponents
- Converts between infrastructure and domain models
- Manages configuration and error handling

**Key Methods**:

- `GetUserHealthOverview()`: Aggregates data and generates comprehensive health summary
- `GetCorrelationAnalysis()`: Orchestrates correlation analysis across data types
- `GetTrendAnalysis()`: Performs time-series trend analysis
- `GetBehaviorPatterns()`: Identifies behavioral patterns in health data
- `GetHealthInsights()`: Generates actionable health insights
- `GetHealthScore()`: Calculates overall health scores
- `GetRecommendations()`: Creates personalized recommendations

### Aggregator Package

**Purpose**: Handles data aggregation and statistical calculations

#### `DataAggregator`

- Aggregates raw health data into summary statistics
- Handles large datasets efficiently
- Provides caching for expensive aggregations

#### `Statistics`

- Statistical utility functions
- Data normalization and standardization
- Basic statistical measures (mean, median, variance)

### Analyzer Package

**Purpose**: Performs complex analysis operations

#### `CorrelationAnalyzer`

- Calculates statistical correlations between health metrics
- Meal-bowel movement correlation analysis
- Meal-symptom correlation analysis
- Medication effectiveness analysis

#### `TrendAnalyzer`

- Time-series trend analysis
- Seasonal pattern detection
- Change point detection
- Confidence interval calculations

#### `PatternAnalyzer`

- Behavioral pattern recognition
- Eating pattern analysis
- Bowel movement pattern analysis
- Symptom pattern analysis
- Lifestyle pattern analysis

### Calculator Package

**Purpose**: Health scoring and metric calculations

#### `HealthScoreCalculator`

- Overall health score calculation (0-100)
- Component score breakdown
- Weighted scoring algorithms
- Benchmark comparisons

#### `LegacyCalculator`

- Backward compatibility utilities
- Deprecated calculation methods
- Migration support functions

### Insights Package

**Purpose**: Generates actionable insights and recommendations

#### `InsightEngine`

- Analyzes health data for actionable insights
- Identifies risk factors and positive factors
- Generates evidence-based recommendations
- Prioritizes insights by importance and confidence

#### `RecommendationGenerator`

- Creates personalized health recommendations
- Evidence-based recommendation logic
- Action step generation
- Timeline and impact estimation

### Shared Package

**Purpose**: Common utilities and data structures

#### `Models`

- Shared data structures between components
- Type conversion utilities
- Common interfaces

#### `Patterns`

- Pattern analysis utilities
- Time-based pattern structures
- Behavioral pattern models

## Data Flow

```
1. API Request → AnalyticsService
2. AnalyticsService → Domain Repositories (fetch raw data)
3. AnalyticsService → DataAggregator (aggregate data)
4. AnalyticsService → Specialized Analyzers (perform analysis)
5. Analyzers → Calculators (compute metrics)
6. AnalyticsService → InsightEngine (generate insights)
7. AnalyticsService → Domain Models (convert results)
8. Domain Models → API Response
```

## Configuration

### AnalyticsServiceConfig

```go
type AnalyticsServiceConfig struct {
    DefaultMedicationLimit int           // Max medications to analyze
    DefaultDataWindow      time.Duration // Default analysis window
    // Additional configuration options...
}
```

## Error Handling

Each component follows the domain error model:

- Uses domain-specific errors from `analytics.errors`
- Wraps infrastructure errors with context
- Provides detailed error messages for debugging
- Maintains error chain for traceability

## Performance Optimizations

### Caching Strategy

- Aggregated data cached for repeated analysis
- Health scores cached with TTL
- Pattern analysis results cached by date range

### Database Optimization

- Efficient date range queries
- Proper indexing on user_id and recorded_at
- Pagination for large datasets
- Connection pooling for concurrent analysis

### Memory Management

- Streaming data processing for large datasets
- Garbage collection friendly data structures
- Memory pooling for frequent allocations

## Testing Strategy

### Unit Tests

- Each component has comprehensive unit tests
- Mock repositories for isolated testing
- Test data generators for consistent test scenarios

### Integration Tests

- End-to-end analytics pipeline testing
- Real database integration tests
- Performance benchmarking tests

### Test Data

- Synthetic health data generators
- Edge case scenario testing
- Load testing with realistic data volumes

## Monitoring and Observability

### Metrics

- Analysis execution time tracking
- Cache hit/miss ratios
- Error rate monitoring
- Data volume processed

### Logging

- Structured logging with context
- Performance metrics logging
- Error logging with stack traces
- Audit logging for sensitive operations

## Future Enhancements

### Performance Improvements

1. **Parallel Processing**: Concurrent analysis of independent data streams
2. **Advanced Caching**: Redis-based distributed caching
3. **Data Streaming**: Real-time analysis capabilities
4. **Query Optimization**: Advanced database query optimization

### Analysis Capabilities

1. **Machine Learning Integration**: TensorFlow/PyTorch model integration
2. **Advanced Statistics**: More sophisticated statistical analysis
3. **Anomaly Detection**: Automated health anomaly detection
4. **Predictive Analytics**: Health outcome prediction models

### Scalability

1. **Microservice Architecture**: Split analytics into specialized services
2. **Event-Driven Processing**: Asynchronous analysis processing
3. **Horizontal Scaling**: Multi-instance analytics processing
4. **Data Partitioning**: Efficient data distribution strategies

## 🧹 Code Quality & Pre-commit Hooks

All linting, formatting, and type-checking is managed via [pre-commit](https://pre-commit.com) and the `.pre-commit-config.yaml` in the project root. Husky and lint-staged are no longer used.

**Setup:**

```bash
uv tool install pre-commit  # or pip install pre-commit
pre-commit install
```

Hooks run on every commit, or manually with:

```bash
pre-commit run --all-files
```
