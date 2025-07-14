package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	domainMetrics "github.com/cfioretti/ingredients-balancer/internal/domain/metrics"
)

type PrometheusMetrics struct {
	// Balance Operations
	balanceOperationsTotal      *prometheus.CounterVec
	balanceOperationDuration    *prometheus.HistogramVec
	balanceOperationErrorsTotal *prometheus.CounterVec
	activeBalanceOperations     prometheus.Gauge

	// Ingredient Processing
	ingredientProcessingTotal    *prometheus.CounterVec
	ingredientProcessingDuration *prometheus.HistogramVec
	ingredientOptimizationsTotal *prometheus.CounterVec
	ingredientWastage            prometheus.Histogram

	// Recipe Analysis
	recipeAnalysisTotal    *prometheus.CounterVec
	recipeAnalysisDuration prometheus.Histogram
	recipeValidationsTotal *prometheus.CounterVec
	recipePortions         prometheus.Histogram

	// Pan Distribution
	panDistributionsTotal   *prometheus.CounterVec
	panDistributionAccuracy prometheus.Histogram
	panOptimizationsTotal   *prometheus.CounterVec
	panUtilization          prometheus.Histogram

	// gRPC Metrics
	grpcRequestsTotal     *prometheus.CounterVec
	grpcRequestDuration   *prometheus.HistogramVec
	activeGRPCConnections prometheus.Gauge

	// Business Quality Metrics
	balancingAccuracy      prometheus.Histogram
	optimizationStrategies *prometheus.CounterVec
	ingredientDistribution prometheus.Histogram
	qualityChecksTotal     *prometheus.CounterVec
}

func NewPrometheusMetrics() *PrometheusMetrics {
	return &PrometheusMetrics{
		// Balance Operations
		balanceOperationsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "ingredients_balancer_balance_operations_total",
				Help: "Total number of balance operations by recipe type",
			},
			[]string{"recipe_type"},
		),
		balanceOperationDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "ingredients_balancer_balance_operation_duration_seconds",
				Help:    "Duration of balance operations by recipe type",
				Buckets: []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1.0, 2.5, 5.0, 10.0},
			},
			[]string{"recipe_type"},
		),
		balanceOperationErrorsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "ingredients_balancer_balance_operation_errors_total",
				Help: "Total number of balance operation errors by recipe type and error type",
			},
			[]string{"recipe_type", "error_type"},
		),
		activeBalanceOperations: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "ingredients_balancer_active_balance_operations",
				Help: "Number of active balance operations",
			},
		),

		// Ingredient Processing
		ingredientProcessingTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "ingredients_balancer_ingredient_processing_total",
				Help: "Total number of ingredient processing operations",
			},
			[]string{"ingredient_type", "success"},
		),
		ingredientProcessingDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "ingredients_balancer_ingredient_processing_duration_seconds",
				Help:    "Duration of ingredient processing operations",
				Buckets: []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1.0, 2.5, 5.0},
			},
			[]string{"ingredient_type"},
		),
		ingredientOptimizationsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "ingredients_balancer_ingredient_optimizations_total",
				Help: "Total number of ingredient optimizations by type",
			},
			[]string{"optimization_type"},
		),
		ingredientWastage: promauto.NewHistogram(
			prometheus.HistogramOpts{
				Name:    "ingredients_balancer_ingredient_wastage_percentage",
				Help:    "Ingredient wastage percentage",
				Buckets: []float64{0, 1, 2, 3, 4, 5, 7, 10, 15, 20, 25, 30},
			},
		),

		// Recipe Analysis
		recipeAnalysisTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "ingredients_balancer_recipe_analysis_total",
				Help: "Total number of recipe analysis operations",
			},
			[]string{"recipe_complexity"},
		),
		recipeAnalysisDuration: promauto.NewHistogram(
			prometheus.HistogramOpts{
				Name:    "ingredients_balancer_recipe_analysis_duration_seconds",
				Help:    "Duration of recipe analysis operations",
				Buckets: []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1.0, 2.5, 5.0},
			},
		),
		recipeValidationsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "ingredients_balancer_recipe_validations_total",
				Help: "Total number of recipe validations",
			},
			[]string{"validation_type", "valid"},
		),
		recipePortions: promauto.NewHistogram(
			prometheus.HistogramOpts{
				Name:    "ingredients_balancer_recipe_portions",
				Help:    "Number of recipe portions",
				Buckets: []float64{1, 2, 3, 4, 5, 6, 8, 10, 12, 15, 20, 25, 30, 40, 50},
			},
		),

		// Pan Distribution
		panDistributionsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "ingredients_balancer_pan_distributions_total",
				Help: "Total number of pan distributions by pan size",
			},
			[]string{"pan_size"},
		),
		panDistributionAccuracy: promauto.NewHistogram(
			prometheus.HistogramOpts{
				Name:    "ingredients_balancer_pan_distribution_accuracy",
				Help:    "Pan distribution accuracy percentage",
				Buckets: []float64{70, 75, 80, 85, 90, 92, 94, 96, 98, 99, 99.5, 100},
			},
		),
		panOptimizationsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "ingredients_balancer_pan_optimizations_total",
				Help: "Total number of pan optimizations by type",
			},
			[]string{"optimization_type"},
		),
		panUtilization: promauto.NewHistogram(
			prometheus.HistogramOpts{
				Name:    "ingredients_balancer_pan_utilization_percentage",
				Help:    "Pan utilization percentage",
				Buckets: []float64{50, 60, 70, 75, 80, 85, 90, 92, 94, 96, 98, 99, 100},
			},
		),

		// gRPC Metrics
		grpcRequestsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "ingredients_balancer_grpc_requests_total",
				Help: "Total number of gRPC requests",
			},
			[]string{"method", "status_code"},
		),
		grpcRequestDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "ingredients_balancer_grpc_request_duration_seconds",
				Help:    "Duration of gRPC requests",
				Buckets: []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1.0, 2.5, 5.0, 10.0},
			},
			[]string{"method"},
		),
		activeGRPCConnections: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "ingredients_balancer_active_grpc_connections",
				Help: "Number of active gRPC connections",
			},
		),

		// Business Quality Metrics
		balancingAccuracy: promauto.NewHistogram(
			prometheus.HistogramOpts{
				Name:    "ingredients_balancer_balancing_accuracy",
				Help:    "Balancing accuracy percentage",
				Buckets: []float64{70, 75, 80, 85, 90, 92, 94, 96, 98, 99, 99.5, 100},
			},
		),
		optimizationStrategies: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "ingredients_balancer_optimization_strategies_total",
				Help: "Total number of optimization strategies used",
			},
			[]string{"strategy"},
		),
		ingredientDistribution: promauto.NewHistogram(
			prometheus.HistogramOpts{
				Name:    "ingredients_balancer_ingredient_distribution",
				Help:    "Ingredient distribution efficiency",
				Buckets: []float64{0.5, 0.6, 0.7, 0.8, 0.85, 0.9, 0.92, 0.94, 0.96, 0.98, 0.99, 1.0},
			},
		),
		qualityChecksTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "ingredients_balancer_quality_checks_total",
				Help: "Total number of quality checks",
			},
			[]string{"check_type", "passed"},
		),
	}
}

func (p *PrometheusMetrics) IncrementBalanceOperations(recipeType string) {
	p.balanceOperationsTotal.WithLabelValues(recipeType).Inc()
}

func (p *PrometheusMetrics) RecordBalanceOperationDuration(recipeType string, duration time.Duration) {
	p.balanceOperationDuration.WithLabelValues(recipeType).Observe(duration.Seconds())
}

func (p *PrometheusMetrics) IncrementBalanceOperationErrors(recipeType string, errorType string) {
	p.balanceOperationErrorsTotal.WithLabelValues(recipeType, errorType).Inc()
}

func (p *PrometheusMetrics) SetActiveBalanceOperations(count int) {
	p.activeBalanceOperations.Set(float64(count))
}

func (p *PrometheusMetrics) IncrementIngredientProcessing(ingredientType string, success bool) {
	p.ingredientProcessingTotal.WithLabelValues(ingredientType, boolToString(success)).Inc()
}

func (p *PrometheusMetrics) RecordIngredientProcessingDuration(ingredientType string, duration time.Duration) {
	p.ingredientProcessingDuration.WithLabelValues(ingredientType).Observe(duration.Seconds())
}

func (p *PrometheusMetrics) IncrementIngredientOptimizations(optimizationType string) {
	p.ingredientOptimizationsTotal.WithLabelValues(optimizationType).Inc()
}

func (p *PrometheusMetrics) RecordIngredientWastage(wastePercentage float64) {
	p.ingredientWastage.Observe(wastePercentage)
}

func (p *PrometheusMetrics) IncrementRecipeAnalysis(recipeComplexity string) {
	p.recipeAnalysisTotal.WithLabelValues(recipeComplexity).Inc()
}

func (p *PrometheusMetrics) RecordRecipeAnalysisDuration(duration time.Duration) {
	p.recipeAnalysisDuration.Observe(duration.Seconds())
}

func (p *PrometheusMetrics) IncrementRecipeValidations(validationType string, valid bool) {
	p.recipeValidationsTotal.WithLabelValues(validationType, boolToString(valid)).Inc()
}

func (p *PrometheusMetrics) RecordRecipePortions(portionCount int) {
	p.recipePortions.Observe(float64(portionCount))
}

func (p *PrometheusMetrics) IncrementPanDistributions(panSize string) {
	p.panDistributionsTotal.WithLabelValues(panSize).Inc()
}

func (p *PrometheusMetrics) RecordPanDistributionAccuracy(accuracy float64) {
	p.panDistributionAccuracy.Observe(accuracy)
}

func (p *PrometheusMetrics) IncrementPanOptimizations(optimizationType string) {
	p.panOptimizationsTotal.WithLabelValues(optimizationType).Inc()
}

func (p *PrometheusMetrics) RecordPanUtilization(utilization float64) {
	p.panUtilization.Observe(utilization)
}

func (p *PrometheusMetrics) IncrementGRPCRequests(method string, statusCode string) {
	p.grpcRequestsTotal.WithLabelValues(method, statusCode).Inc()
}

func (p *PrometheusMetrics) RecordGRPCRequestDuration(method string, duration time.Duration) {
	p.grpcRequestDuration.WithLabelValues(method).Observe(duration.Seconds())
}

func (p *PrometheusMetrics) SetActiveGRPCConnections(count int) {
	p.activeGRPCConnections.Set(float64(count))
}

func (p *PrometheusMetrics) RecordBalancingAccuracy(accuracy float64) {
	p.balancingAccuracy.Observe(accuracy)
}

func (p *PrometheusMetrics) IncrementOptimizationStrategies(strategy string) {
	p.optimizationStrategies.WithLabelValues(strategy).Inc()
}

func (p *PrometheusMetrics) RecordIngredientDistribution(distribution float64) {
	p.ingredientDistribution.Observe(distribution)
}

func (p *PrometheusMetrics) IncrementQualityChecks(checkType string, passed bool) {
	p.qualityChecksTotal.WithLabelValues(checkType, boolToString(passed)).Inc()
}

func boolToString(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

var _ domainMetrics.BalancerMetrics = (*PrometheusMetrics)(nil)
