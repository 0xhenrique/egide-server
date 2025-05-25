package service

import (
	"fmt"
	"math/rand"
	"time"

	"egide-server/internal/repository"
	"egide-server/internal/models"
)

// KpiMetric represents a single KPI metric with value and change
type KpiMetric struct {
	Value    interface{} `json:"value"`
	Change   *float64    `json:"change,omitempty"`
	Subvalue *string     `json:"subvalue,omitempty"`
}

// KpiData represents all KPI metrics for the dashboard
type KpiData struct {
	TotalRequests  KpiMetric `json:"totalRequests"`
	BlockedThreats KpiMetric `json:"blockedThreats"`
	ResponseTime   KpiMetric `json:"responseTime"`
	Uptime         KpiMetric `json:"uptime"`
}

// MetricsService handles metrics data operations
type MetricsService struct {
	rand               *rand.Rand
	healthCheckRepo    *repository.HealthCheckRepository
}

// NewMetricsService creates a new metrics service
func NewMetricsService(healthCheckRepo *repository.HealthCheckRepository) *MetricsService {
	source := rand.NewSource(time.Now().UnixNano())
	return &MetricsService{
		rand:            rand.New(source),
		healthCheckRepo: healthCheckRepo,
	}
}

// GetKpiData returns KPI data for the given user's sites
func (s *MetricsService) GetKpiData(userID int64) (*KpiData, error) {
	now := time.Now()
	
	// Get current month data (last 30 days)
	currentStart := now.AddDate(0, 0, -30)
	currentChecks, err := s.healthCheckRepo.GetChecksInRange(currentStart, now)
	if err != nil {
		return nil, fmt.Errorf("failed to get current health checks: %v", err)
	}
	
	// Get previous month data (30-60 days ago)
	previousStart := now.AddDate(0, 0, -60)
	previousEnd := now.AddDate(0, 0, -30)
	previousChecks, err := s.healthCheckRepo.GetChecksInRange(previousStart, previousEnd)
	if err != nil {
		return nil, fmt.Errorf("failed to get previous health checks: %v", err)
	}
	
	// Calculate uptime and response time metrics
	currentUptime, currentAvgResponseTime := s.calculateMetrics(currentChecks)
	previousUptime, previousAvgResponseTime := s.calculateMetrics(previousChecks)
	
	// Calculate change percentages
	var uptimeChange *float64
	var responseTimeChange *float64
	
	if len(previousChecks) > 0 {
		uptimeChangeVal := currentUptime - previousUptime
		uptimeChange = &uptimeChangeVal
		
		if previousAvgResponseTime > 0 {
			responseTimeChangeVal := ((currentAvgResponseTime - previousAvgResponseTime) / previousAvgResponseTime) * 100
			responseTimeChange = &responseTimeChangeVal
		}
	}
	
	// Generate mock data for total requests and blocked threats (as before)
	totalRequests := 500000 + s.rand.Intn(1000000)
	totalRequestsChange := (s.rand.Float64() * 30) - 15
	
	blockedThreats := 2000 + s.rand.Intn(5000)
	blockedThreatsChange := (s.rand.Float64() * 20) - 10
	
	// Calculate downtime information
	downtimeInfo := s.calculateDowntimeInfo(currentChecks)
	
	// Format response time
	responseTimeStr := fmt.Sprintf("%.0fms", currentAvgResponseTime)
	
	// Round values
	totalRequestsChangeRounded := round(totalRequestsChange, 1)
	blockedThreatsChangeRounded := round(blockedThreatsChange, 1)
	
	return &KpiData{
		TotalRequests: KpiMetric{
			Value:  totalRequests,
			Change: &totalRequestsChangeRounded,
		},
		BlockedThreats: KpiMetric{
			Value:  blockedThreats,
			Change: &blockedThreatsChangeRounded,
		},
		ResponseTime: KpiMetric{
			Value:  responseTimeStr,
			Change: responseTimeChange,
		},
		Uptime: KpiMetric{
			Value:    round(currentUptime, 2),
			Change:   uptimeChange,
			Subvalue: &downtimeInfo,
		},
	}, nil
}

// calculateMetrics computes uptime percentage and average response time from health checks
func (s *MetricsService) calculateMetrics(checks []*models.HealthCheck) (uptime float64, avgResponseTime float64) {
	if len(checks) == 0 {
		return 100.0, 0.0 // Default to 100% uptime if no data
	}
	
	successCount := 0
	totalResponseTime := 0
	
	for _, check := range checks {
		if check.Success {
			successCount++
		}
		totalResponseTime += check.ResponseTimeMs
	}
	
	uptime = (float64(successCount) / float64(len(checks))) * 100
	avgResponseTime = float64(totalResponseTime) / float64(len(checks))
	
	return uptime, avgResponseTime
}

// calculateDowntimeInfo creates a human-readable downtime string
func (s *MetricsService) calculateDowntimeInfo(checks []*models.HealthCheck) string {
	if len(checks) == 0 {
		return "No monitoring data available"
	}
	
	// Calculate total downtime in the period
	downtime := 0
	for _, check := range checks {
		if !check.Success {
			// Each failed check represents ~1 minute of downtime (our check interval)
			downtime += 60 // seconds
		}
	}
	
	if downtime == 0 {
		return "No downtime recorded"
	}
	
	// Convert to human readable format
	minutes := downtime / 60
	remainingSeconds := downtime % 60
	
	if minutes > 0 {
		return fmt.Sprintf("Total downtime: %dm %ds", minutes, remainingSeconds)
	}
	return fmt.Sprintf("Total downtime: %ds", remainingSeconds)
}

// round rounds a float64 to the specified number of decimal places
func round(value float64, decimals int) float64 {
	precision := 1.0
	for i := 0; i < decimals; i++ {
		precision *= 10
	}
	return float64(int(value*precision+0.5)) / precision
}
