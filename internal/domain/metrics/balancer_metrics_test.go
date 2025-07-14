package metrics

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockBalancerMetrics struct {
	mock.Mock
}

func (m *MockBalancerMetrics) IncrementBalanceOperations(recipeType string) {
	m.Called(recipeType)
}

func (m *MockBalancerMetrics) RecordBalanceOperationDuration(recipeType string, duration time.Duration) {
	m.Called(recipeType, duration)
}

func (m *MockBalancerMetrics) IncrementBalanceOperationErrors(recipeType string, errorType string) {
	m.Called(recipeType, errorType)
}

func (m *MockBalancerMetrics) SetActiveBalanceOperations(count int) {
	m.Called(count)
}

func (m *MockBalancerMetrics) IncrementIngredientProcessing(ingredientType string, success bool) {
	m.Called(ingredientType, success)
}

func (m *MockBalancerMetrics) RecordIngredientProcessingDuration(ingredientType string, duration time.Duration) {
	m.Called(ingredientType, duration)
}

func (m *MockBalancerMetrics) IncrementIngredientOptimizations(optimizationType string) {
	m.Called(optimizationType)
}

func (m *MockBalancerMetrics) RecordIngredientWastage(wastePercentage float64) {
	m.Called(wastePercentage)
}

func (m *MockBalancerMetrics) IncrementRecipeAnalysis(recipeComplexity string) {
	m.Called(recipeComplexity)
}

func (m *MockBalancerMetrics) RecordRecipeAnalysisDuration(duration time.Duration) {
	m.Called(duration)
}

func (m *MockBalancerMetrics) IncrementRecipeValidations(validationType string, valid bool) {
	m.Called(validationType, valid)
}

func (m *MockBalancerMetrics) RecordRecipePortions(portionCount int) {
	m.Called(portionCount)
}

func (m *MockBalancerMetrics) IncrementPanDistributions(panSize string) {
	m.Called(panSize)
}

func (m *MockBalancerMetrics) RecordPanDistributionAccuracy(accuracy float64) {
	m.Called(accuracy)
}

func (m *MockBalancerMetrics) IncrementPanOptimizations(optimizationType string) {
	m.Called(optimizationType)
}

func (m *MockBalancerMetrics) RecordPanUtilization(utilization float64) {
	m.Called(utilization)
}

func (m *MockBalancerMetrics) IncrementGRPCRequests(method string, statusCode string) {
	m.Called(method, statusCode)
}

func (m *MockBalancerMetrics) RecordGRPCRequestDuration(method string, duration time.Duration) {
	m.Called(method, duration)
}

func (m *MockBalancerMetrics) SetActiveGRPCConnections(count int) {
	m.Called(count)
}

func (m *MockBalancerMetrics) RecordBalancingAccuracy(accuracy float64) {
	m.Called(accuracy)
}

func (m *MockBalancerMetrics) IncrementOptimizationStrategies(strategy string) {
	m.Called(strategy)
}

func (m *MockBalancerMetrics) RecordIngredientDistribution(distribution float64) {
	m.Called(distribution)
}

func (m *MockBalancerMetrics) IncrementQualityChecks(checkType string, passed bool) {
	m.Called(checkType, passed)
}

func TestNewMetricsRecorder(t *testing.T) {
	mockMetrics := new(MockBalancerMetrics)
	recorder := NewMetricsRecorder(mockMetrics)

	assert.NotNil(t, recorder)
	assert.Equal(t, mockMetrics, recorder.metrics)
}

func TestRecordBalanceOperation_Success(t *testing.T) {
	mockMetrics := new(MockBalancerMetrics)
	recorder := NewMetricsRecorder(mockMetrics)

	result := BalanceOperationResult{
		RecipeType:        "napoletana",
		Duration:          100 * time.Millisecond,
		Success:           true,
		IngredientsCount:  5,
		PansCount:         3,
		PortionCount:      4,
		OptimizationType:  "minimize_waste",
		AccuracyScore:     95.5,
		WastagePercentage: 2.1,
		UtilizationScore:  97.9,
	}

	mockMetrics.On("RecordBalanceOperationDuration", result.RecipeType, result.Duration).Return()
	mockMetrics.On("IncrementBalanceOperations", result.RecipeType).Return()
	mockMetrics.On("RecordBalancingAccuracy", result.AccuracyScore).Return()
	mockMetrics.On("IncrementOptimizationStrategies", result.OptimizationType).Return()
	mockMetrics.On("RecordIngredientWastage", result.WastagePercentage).Return()
	mockMetrics.On("RecordPanUtilization", result.UtilizationScore).Return()
	mockMetrics.On("RecordRecipePortions", result.PortionCount).Return()
	mockMetrics.On("RecordPanDistributionAccuracy", result.AccuracyScore).Return()
	mockMetrics.On("IncrementIngredientProcessing", "balanced", true).Return()

	recorder.RecordBalanceOperation(context.Background(), result)

	mockMetrics.AssertExpectations(t)
}

func TestRecordBalanceOperation_Error(t *testing.T) {
	mockMetrics := new(MockBalancerMetrics)
	recorder := NewMetricsRecorder(mockMetrics)

	result := BalanceOperationResult{
		RecipeType: "napoletana",
		Duration:   50 * time.Millisecond,
		Success:    false,
		ErrorType:  "validation_error",
	}

	mockMetrics.On("RecordBalanceOperationDuration", result.RecipeType, result.Duration).Return()
	mockMetrics.On("IncrementBalanceOperations", result.RecipeType).Return()
	mockMetrics.On("IncrementBalanceOperationErrors", result.RecipeType, result.ErrorType).Return()
	mockMetrics.On("IncrementIngredientProcessing", "balanced", false).Return()

	recorder.RecordBalanceOperation(context.Background(), result)

	mockMetrics.AssertExpectations(t)
}

func TestRecordIngredientProcessing(t *testing.T) {
	mockMetrics := new(MockBalancerMetrics)
	recorder := NewMetricsRecorder(mockMetrics)

	duration := 25 * time.Millisecond

	mockMetrics.On("IncrementIngredientProcessing", "flour", true).Return()
	mockMetrics.On("RecordIngredientProcessingDuration", "flour", duration).Return()

	recorder.RecordIngredientProcessing(context.Background(), "flour", duration, true)

	mockMetrics.AssertExpectations(t)
}

func TestRecordRecipeAnalysis(t *testing.T) {
	mockMetrics := new(MockBalancerMetrics)
	recorder := NewMetricsRecorder(mockMetrics)

	duration := 15 * time.Millisecond

	mockMetrics.On("IncrementRecipeAnalysis", "medium").Return()
	mockMetrics.On("RecordRecipeAnalysisDuration", duration).Return()

	recorder.RecordRecipeAnalysis(context.Background(), "medium", duration)

	mockMetrics.AssertExpectations(t)
}

func TestRecordPanDistribution(t *testing.T) {
	mockMetrics := new(MockBalancerMetrics)
	recorder := NewMetricsRecorder(mockMetrics)

	mockMetrics.On("IncrementPanDistributions", "30cm").Return()
	mockMetrics.On("RecordPanDistributionAccuracy", 96.0).Return()
	mockMetrics.On("IncrementPanOptimizations", "even_distribution").Return()

	recorder.RecordPanDistribution(context.Background(), "30cm", 96.0, "even_distribution")

	mockMetrics.AssertExpectations(t)
}

func TestRecordQualityCheck(t *testing.T) {
	mockMetrics := new(MockBalancerMetrics)
	recorder := NewMetricsRecorder(mockMetrics)

	mockMetrics.On("IncrementQualityChecks", "ingredient_ratio", true).Return()

	recorder.RecordQualityCheck(context.Background(), "ingredient_ratio", true)

	mockMetrics.AssertExpectations(t)
}
