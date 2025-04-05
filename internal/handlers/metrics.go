package handlers

import (
	"encoding/json"
	"net/http"

	"egide-server/internal/auth"
	"egide-server/internal/service"
)

// MetricsHandler handles metrics-related requests
type MetricsHandler struct {
	metricsService *service.MetricsService
}

// NewMetricsHandler creates a new metrics handler
func NewMetricsHandler(metricsService *service.MetricsService) *MetricsHandler {
	return &MetricsHandler{
		metricsService: metricsService,
	}
}

// GetKpi handles GET /api/metrics/kpi
func (h *MetricsHandler) GetKpi(w http.ResponseWriter, r *http.Request) {
	// Get user ID from authenticated context
	userID, err := auth.UserIDFromContext(r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get KPI data for the user's sites
	kpiData, err := h.metricsService.GetKpiData(userID)
	if err != nil {
		http.Error(w, "Error fetching KPI data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return KPI data as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(kpiData)
}
