package metrics

import (
	"context"
	"time"
)

type BalancerMetrics interface {
	IncrementBalanceOperations(recipeType string)
	RecordBalanceOperationDuration(recipeType string, duration time.Duration)
	IncrementBalanceOperationErrors(recipeType string, errorType string)
	SetActiveBalanceOperations(count int)

	IncrementIngredientProcessing(ingredientType string, success bool)
	RecordIngredientProcessingDuration(ingredientType string, duration time.Duration)
	IncrementIngredientOptimizations(optimizationType string)
	RecordIngredientWastage(wastePercentage float64)

	IncrementRecipeAnalysis(recipeComplexity string)
	RecordRecipeAnalysisDuration(duration time.Duration)
	IncrementRecipeValidations(validationType string, valid bool)
	RecordRecipePortions(portionCount int)

	IncrementPanDistributions(panSize string)
	RecordPanDistributionAccuracy(accuracy float64)
	IncrementPanOptimizations(optimizationType string)
	RecordPanUtilization(utilization float64)

	IncrementGRPCRequests(method string, statusCode string)
	RecordGRPCRequestDuration(method string, duration time.Duration)
	SetActiveGRPCConnections(count int)

	RecordBalancingAccuracy(accuracy float64)
	IncrementOptimizationStrategies(strategy string)
	RecordIngredientDistribution(distribution float64)
	IncrementQualityChecks(checkType string, passed bool)
}

type BalanceOperationResult struct {
	RecipeType        string
	Duration          time.Duration
	Success           bool
	ErrorType         string
	IngredientsCount  int
	PansCount         int
	PortionCount      int
	OptimizationType  string
	AccuracyScore     float64
	WastagePercentage float64
	UtilizationScore  float64
}

type MetricsRecorder struct {
	metrics BalancerMetrics
}

func NewMetricsRecorder(metrics BalancerMetrics) *MetricsRecorder {
	return &MetricsRecorder{
		metrics: metrics,
	}
}

func (m *MetricsRecorder) RecordBalanceOperation(ctx context.Context, result BalanceOperationResult) {
	m.metrics.RecordBalanceOperationDuration(result.RecipeType, result.Duration)
	m.metrics.IncrementBalanceOperations(result.RecipeType)

	if result.Success {
		m.metrics.RecordBalancingAccuracy(result.AccuracyScore)
		m.metrics.IncrementOptimizationStrategies(result.OptimizationType)
		m.metrics.RecordIngredientWastage(result.WastagePercentage)
		m.metrics.RecordPanUtilization(result.UtilizationScore)
		m.metrics.RecordRecipePortions(result.PortionCount)

		if result.PansCount > 0 {
			m.metrics.RecordPanDistributionAccuracy(result.AccuracyScore)
		}

		m.metrics.IncrementIngredientProcessing("balanced", true)
	} else {
		m.metrics.IncrementBalanceOperationErrors(result.RecipeType, result.ErrorType)
		m.metrics.IncrementIngredientProcessing("balanced", false)
	}
}

func (m *MetricsRecorder) RecordIngredientProcessing(ctx context.Context, ingredientType string, duration time.Duration, success bool) {
	m.metrics.IncrementIngredientProcessing(ingredientType, success)
	m.metrics.RecordIngredientProcessingDuration(ingredientType, duration)
}

func (m *MetricsRecorder) RecordRecipeAnalysis(ctx context.Context, complexity string, duration time.Duration) {
	m.metrics.IncrementRecipeAnalysis(complexity)
	m.metrics.RecordRecipeAnalysisDuration(duration)
}

func (m *MetricsRecorder) RecordPanDistribution(ctx context.Context, panSize string, accuracy float64, optimizationType string) {
	m.metrics.IncrementPanDistributions(panSize)
	m.metrics.RecordPanDistributionAccuracy(accuracy)
	m.metrics.IncrementPanOptimizations(optimizationType)
}

func (m *MetricsRecorder) RecordQualityCheck(ctx context.Context, checkType string, passed bool) {
	m.metrics.IncrementQualityChecks(checkType, passed)
}
