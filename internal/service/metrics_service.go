package service

import (
	"math/rand"
	"time"
	"fmt"
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
	rand *rand.Rand
}

// NewMetricsService creates a new metrics service
func NewMetricsService() *MetricsService {
	source := rand.NewSource(time.Now().UnixNano())
	return &MetricsService{
		rand: rand.New(source),
	}
}

// GetKpiData returns KPI data for the given user's sites
// Currently mocks the data, but in the future will fetch from the actual data source
func (s *MetricsService) GetKpiData(userID int64) (*KpiData, error) {
	// In a real implementation, this would fetch actual metrics for the user's sites
	// For now, we'll generate realistic mock data
	
	// Total Requests (typically a large number)
	totalRequests := 500000 + s.rand.Intn(1000000)
	totalRequestsChange := (s.rand.Float64() * 30) - 15 // -15% to +15%
	
	// Blocked Threats (smaller than total requests)
	blockedThreats := 2000 + s.rand.Intn(5000)
	blockedThreatsChange := (s.rand.Float64() * 20) - 10 // -10% to +10%
	
	// Response Time (typically between 30ms and 200ms)
	responseTimeValue := 30 + s.rand.Intn(170)
	responseTimeChange := (s.rand.Float64() * 10) - 5 // -5% to +5%
	
	// Uptime (typically high, between 99.5% and 100%)
	uptimeValue := 99.5 + (s.rand.Float64() * 0.5)
	
	// Calculate a realistic downtime string based on uptime percentage
	// 100% uptime = 0 downtime, 99.5% = ~7.2 hours/month
	downtimeMinutes := int((100 - uptimeValue) * 0.01 * 24 * 60) // downtime in minutes for last day
	downtimeString := fmt.Sprintf("Total downtime: %dm %ds", 
		downtimeMinutes, 
		s.rand.Intn(60)) // random seconds
	
	// Format the response time as a string with "ms" suffix
	responseTimeStr := fmt.Sprintf("%dms", responseTimeValue)
	
	// Round the changes to 1 decimal place
	totalRequestsChangeRounded := round(totalRequestsChange, 1)
	blockedThreatsChangeRounded := round(blockedThreatsChange, 1)
	responseTimeChangeRounded := round(responseTimeChange, 1)
	
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
			Change: &responseTimeChangeRounded,
		},
		Uptime: KpiMetric{
			Value:    round(uptimeValue, 2),
			Subvalue: &downtimeString,
		},
	}, nil
}

// round rounds a float64 to the specified number of decimal places
func round(value float64, decimals int) float64 {
	precision := 1.0
	for i := 0; i < decimals; i++ {
		precision *= 10
	}
	return float64(int(value*precision+0.5)) / precision
}
