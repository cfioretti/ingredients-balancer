package http

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type MetricsHandler struct {
	handler http.Handler
}

func NewMetricsHandler() *MetricsHandler {
	return &MetricsHandler{
		handler: promhttp.Handler(),
	}
}

func (h *MetricsHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.Handle("/metrics", h.handler)
}

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/health", h.handleHealth)
}

func (h *HealthHandler) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := `{
		"status": "healthy",
		"service": "ingredients-balancer",
		"version": "1.0.0"
	}`
	w.Write([]byte(response))
}
