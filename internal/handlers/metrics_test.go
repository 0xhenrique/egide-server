package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"egide-server/internal/auth"
	"egide-server/internal/service"
)

func TestGetKpi(t *testing.T) {
	// Create metrics service
	metricsService := service.NewMetricsService()

	// Create handler with service
	handler := NewMetricsHandler(metricsService)

	// Create request
	req, err := http.NewRequest("GET", "/api/metrics/kpi", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Add user ID to context
	ctx := auth.WithUserID(context.Background(), 123)
	req = req.WithContext(ctx)

	// Create response recorder
	rr := httptest.NewRecorder()

	// Call handler
	handler.GetKpi(rr, req)

	// Check status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check that we got valid JSON data back
	var kpiData service.KpiData
	err = json.Unmarshal(rr.Body.Bytes(), &kpiData)
	if err != nil {
		t.Errorf("could not parse response as JSON: %v", err)
	}

	// Check that all required fields are present
	if _, ok := kpiData.TotalRequests.Value.(float64); !ok {
		t.Errorf("totalRequests.value should be a number")
	}

	if kpiData.TotalRequests.Change == nil {
		t.Errorf("totalRequests.change should not be nil")
	}

	if _, ok := kpiData.BlockedThreats.Value.(float64); !ok {
		t.Errorf("blockedThreats.value should be a number")
	}

	if kpiData.BlockedThreats.Change == nil {
		t.Errorf("blockedThreats.change should not be nil")
	}

	if _, ok := kpiData.ResponseTime.Value.(string); !ok {
		t.Errorf("responseTime.value should be a string")
	}

	if kpiData.ResponseTime.Change == nil {
		t.Errorf("responseTime.change should not be nil")
	}

	if _, ok := kpiData.Uptime.Value.(float64); !ok {
		t.Errorf("uptime.value should be a number")
	}

	if kpiData.Uptime.Subvalue == nil {
		t.Errorf("uptime.subvalue should not be nil")
	}
}
