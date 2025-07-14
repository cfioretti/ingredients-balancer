package middleware

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	domainMetrics "github.com/cfioretti/ingredients-balancer/internal/domain/metrics"
	infraMetrics "github.com/cfioretti/ingredients-balancer/internal/infrastructure/metrics"
)

type MetricsMiddleware struct {
	domainMetrics     domainMetrics.BalancerMetrics
	prometheusMetrics *infraMetrics.PrometheusMetrics
}

func NewMetricsMiddleware(
	domainMetrics domainMetrics.BalancerMetrics,
	prometheusMetrics *infraMetrics.PrometheusMetrics,
) *MetricsMiddleware {
	return &MetricsMiddleware{
		domainMetrics:     domainMetrics,
		prometheusMetrics: prometheusMetrics,
	}
}

func (m *MetricsMiddleware) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		start := time.Now()

		resp, err := handler(ctx, req)

		duration := time.Since(start)
		method := extractMethodName(info.FullMethod)
		statusCode := getStatusCode(err)

		m.domainMetrics.IncrementGRPCRequests(method, statusCode)
		m.domainMetrics.RecordGRPCRequestDuration(method, duration)

		if isBalanceOperation(method) {
			m.recordBalanceMetrics(ctx, method, duration, err)
		}

		return resp, err
	}
}

func (m *MetricsMiddleware) StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		start := time.Now()

		err := handler(srv, stream)

		duration := time.Since(start)
		method := extractMethodName(info.FullMethod)
		statusCode := getStatusCode(err)

		m.domainMetrics.IncrementGRPCRequests(method, statusCode)
		m.domainMetrics.RecordGRPCRequestDuration(method, duration)

		return err
	}
}

func (m *MetricsMiddleware) recordBalanceMetrics(ctx context.Context, method string, duration time.Duration, err error) {
	switch method {
	case "Balance":
		if err == nil {
			m.domainMetrics.IncrementBalanceOperations("unknown")
			m.domainMetrics.RecordBalanceOperationDuration("unknown", duration)
		} else {
			errorType := mapGRPCErrorToBusinessError(err)
			m.domainMetrics.IncrementBalanceOperationErrors("unknown", errorType)
		}
	}
}

func extractMethodName(fullMethod string) string {
	parts := splitString(fullMethod, "/")
	if len(parts) >= 2 {
		return parts[len(parts)-1]
	}
	return "unknown"
}

func splitString(s, sep string) []string {
	if s == "" {
		return []string{}
	}

	var result []string
	current := ""

	for _, char := range s {
		if string(char) == sep {
			if current != "" {
				result = append(result, current)
				current = ""
			}
		} else {
			current += string(char)
		}
	}

	if current != "" {
		result = append(result, current)
	}

	return result
}

func getStatusCode(err error) string {
	if err == nil {
		return "OK"
	}

	if st, ok := status.FromError(err); ok {
		return st.Code().String()
	}

	return "UNKNOWN"
}

func isBalanceOperation(method string) bool {
	balanceOperations := []string{
		"Balance",
		"OptimizeIngredients",
		"ValidateRecipe",
		"AnalyzePans",
	}

	for _, op := range balanceOperations {
		if method == op {
			return true
		}
	}
	return false
}

func mapGRPCErrorToBusinessError(err error) string {
	if err == nil {
		return "none"
	}

	if st, ok := status.FromError(err); ok {
		switch st.Code() {
		case codes.InvalidArgument:
			return "validation_error"
		case codes.NotFound:
			return "recipe_not_found"
		case codes.Internal:
			return "internal_error"
		case codes.Unavailable:
			return "service_unavailable"
		case codes.DeadlineExceeded:
			return "timeout"
		default:
			return "unknown_error"
		}
	}

	return "unknown_error"
}
