package metrics

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrometheusMetrics(t *testing.T) {
	metrics := NewPrometheusMetrics()

	assert.NotNil(t, metrics)
	assert.NotNil(t, metrics.balanceOperationsTotal)
	assert.NotNil(t, metrics.balanceOperationDuration)
	assert.NotNil(t, metrics.ingredientProcessingTotal)
	assert.NotNil(t, metrics.grpcRequestsTotal)
	assert.NotNil(t, metrics.balancingAccuracy)

	metrics.IncrementBalanceOperations("napoletana")
	metrics.IncrementIngredientProcessing("flour", true)
	metrics.IncrementGRPCRequests("Balance", "OK")
	metrics.SetActiveBalanceOperations(3)
	metrics.RecordBalancingAccuracy(95.5)
	metrics.RecordIngredientWastage(2.5)
	metrics.RecordPanUtilization(97.8)
	metrics.IncrementOptimizationStrategies("minimize_waste")
	metrics.IncrementQualityChecks("ingredient_ratio", true)
	metrics.SetActiveGRPCConnections(5)

	assert.Equal(t, "true", boolToString(true))
	assert.Equal(t, "false", boolToString(false))
}
