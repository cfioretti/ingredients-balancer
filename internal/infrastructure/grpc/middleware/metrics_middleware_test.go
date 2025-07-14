package middleware

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	infraMetrics "github.com/cfioretti/ingredients-balancer/internal/infrastructure/metrics"
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

func TestNewMetricsMiddleware(t *testing.T) {
	mockDomainMetrics := new(MockBalancerMetrics)
	mockPrometheusMetrics := &infraMetrics.PrometheusMetrics{}

	middleware := NewMetricsMiddleware(mockDomainMetrics, mockPrometheusMetrics)

	assert.NotNil(t, middleware)
	assert.Equal(t, mockDomainMetrics, middleware.domainMetrics)
	assert.Equal(t, mockPrometheusMetrics, middleware.prometheusMetrics)
}

func TestUnaryServerInterceptor_Success(t *testing.T) {
	mockMetrics := new(MockBalancerMetrics)
	middleware := NewMetricsMiddleware(mockMetrics, &infraMetrics.PrometheusMetrics{})

	interceptor := middleware.UnaryServerInterceptor()

	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return "response", nil
	}

	info := &grpc.UnaryServerInfo{
		FullMethod: "/ingredients_balancer.IngredientsBalancer/Balance",
	}

	mockMetrics.On("IncrementGRPCRequests", "Balance", "OK").Return()
	mockMetrics.On("RecordGRPCRequestDuration", "Balance", mock.AnythingOfType("time.Duration")).Return()
	mockMetrics.On("IncrementBalanceOperations", "unknown").Return()
	mockMetrics.On("RecordBalanceOperationDuration", "unknown", mock.AnythingOfType("time.Duration")).Return()

	response, err := interceptor(context.Background(), "request", info, handler)

	assert.NoError(t, err)
	assert.Equal(t, "response", response)
	mockMetrics.AssertExpectations(t)
}

func TestUnaryServerInterceptor_Error(t *testing.T) {
	mockMetrics := new(MockBalancerMetrics)
	middleware := NewMetricsMiddleware(mockMetrics, &infraMetrics.PrometheusMetrics{})

	interceptor := middleware.UnaryServerInterceptor()

	expectedError := status.Error(codes.InvalidArgument, "invalid request")
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, expectedError
	}

	info := &grpc.UnaryServerInfo{
		FullMethod: "/ingredients_balancer.IngredientsBalancer/Balance",
	}

	mockMetrics.On("IncrementGRPCRequests", "Balance", "InvalidArgument").Return()
	mockMetrics.On("RecordGRPCRequestDuration", "Balance", mock.AnythingOfType("time.Duration")).Return()
	mockMetrics.On("IncrementBalanceOperationErrors", "unknown", "validation_error").Return()

	response, err := interceptor(context.Background(), "request", info, handler)

	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Equal(t, expectedError, err)
	mockMetrics.AssertExpectations(t)
}

func TestExtractMethodName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"/ingredients_balancer.IngredientsBalancer/Balance", "Balance"},
		{"/service/Method", "Method"},
		{"", "unknown"},
		{"NoSlash", "unknown"},
	}

	for _, test := range tests {
		result := extractMethodName(test.input)
		assert.Equal(t, test.expected, result)
	}
}

func TestGetStatusCode(t *testing.T) {
	tests := []struct {
		input    error
		expected string
	}{
		{nil, "OK"},
		{status.Error(codes.InvalidArgument, "invalid"), "InvalidArgument"},
		{status.Error(codes.NotFound, "not found"), "NotFound"},
		{errors.New("generic error"), "UNKNOWN"},
	}

	for _, test := range tests {
		result := getStatusCode(test.input)
		assert.Equal(t, test.expected, result)
	}
}

func TestIsBalanceOperation(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"Balance", true},
		{"OptimizeIngredients", true},
		{"ValidateRecipe", true},
		{"AnalyzePans", true},
		{"UnknownMethod", false},
		{"", false},
	}

	for _, test := range tests {
		result := isBalanceOperation(test.input)
		assert.Equal(t, test.expected, result)
	}
}

func TestMapGRPCErrorToBusinessError(t *testing.T) {
	tests := []struct {
		input    error
		expected string
	}{
		{nil, "none"},
		{status.Error(codes.InvalidArgument, "invalid"), "validation_error"},
		{status.Error(codes.NotFound, "not found"), "recipe_not_found"},
		{status.Error(codes.Internal, "internal"), "internal_error"},
		{status.Error(codes.Unavailable, "unavailable"), "service_unavailable"},
		{status.Error(codes.DeadlineExceeded, "timeout"), "timeout"},
		{status.Error(codes.Unknown, "unknown"), "unknown_error"},
		{errors.New("generic error"), "unknown_error"},
	}

	for _, test := range tests {
		result := mapGRPCErrorToBusinessError(test.input)
		assert.Equal(t, test.expected, result)
	}
}

func TestSplitString(t *testing.T) {
	tests := []struct {
		input     string
		separator string
		expected  []string
	}{
		{"a/b/c", "/", []string{"a", "b", "c"}},
		{"", "/", []string{}},
		{"no-separator", "/", []string{"no-separator"}},
		{"a//b", "/", []string{"a", "b"}},
	}

	for _, test := range tests {
		result := splitString(test.input, test.separator)
		assert.Equal(t, test.expected, result)
	}
}
