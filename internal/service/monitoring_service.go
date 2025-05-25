package service

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"egide-server/internal/models"
	"egide-server/internal/repository"
)

const (
	// URL to monitor - your personal blog protected by Egide
	MonitorURL = "https://0xhenrique.neocities.org/"
	
	// Monitoring interval
	CheckInterval = 60 * time.Second
	
	// Request timeout - anything above this is considered "down"
	RequestTimeout = 10 * time.Second
	
	// Data retention period
	DataRetention = 60 * 24 * time.Hour // 60 days
)

type MonitoringService struct {
	healthCheckRepo *repository.HealthCheckRepository
	httpClient      *http.Client
	stopChan        chan struct{}
}

func NewMonitoringService(healthCheckRepo *repository.HealthCheckRepository) *MonitoringService {
	return &MonitoringService{
		healthCheckRepo: healthCheckRepo,
		httpClient: &http.Client{
			Timeout: RequestTimeout,
		},
		stopChan: make(chan struct{}),
	}
}

// Start begins the monitoring process
func (s *MonitoringService) Start() {
	log.Println("Starting monitoring service...")
	
	// Run cleanup immediately on start
	s.cleanup()
	
	// Start monitoring loop
	ticker := time.NewTicker(CheckInterval)
	go func() {
		defer ticker.Stop()
		
		// Perform initial check
		s.performHealthCheck()
		
		for {
			select {
			case <-ticker.C:
				s.performHealthCheck()
			case <-s.stopChan:
				log.Println("Monitoring service stopped")
				return
			}
		}
	}()
	
	// Start cleanup routine (runs every 6 hours)
	cleanupTicker := time.NewTicker(6 * time.Hour)
	go func() {
		defer cleanupTicker.Stop()
		for {
			select {
			case <-cleanupTicker.C:
				s.cleanup()
			case <-s.stopChan:
				return
			}
		}
	}()
}

// Stop gracefully stops the monitoring service
func (s *MonitoringService) Stop() {
	close(s.stopChan)
}

// performHealthCheck executes a single health check
func (s *MonitoringService) performHealthCheck() {
	start := time.Now()
	
	req, err := http.NewRequestWithContext(context.Background(), "GET", MonitorURL, nil)
	if err != nil {
		s.recordFailure(start, nil, fmt.Sprintf("Failed to create request: %v", err))
		return
	}
	
	// Add a user agent to identify monitoring requests
	req.Header.Set("User-Agent", "Egide-Monitor/1.0")
	
	resp, err := s.httpClient.Do(req)
	responseTime := time.Since(start)
	
	if err != nil {
		s.recordFailure(start, nil, fmt.Sprintf("Request failed: %v", err))
		return
	}
	defer resp.Body.Close()
	
	// Consider 5xx responses as failures
	success := resp.StatusCode < 500
	
	check := &models.HealthCheck{
		Timestamp:      start,
		ResponseTimeMs: int(responseTime.Milliseconds()),
		StatusCode:     &resp.StatusCode,
		Success:        success,
	}
	
	if !success {
		errorMsg := fmt.Sprintf("Server error: HTTP %d", resp.StatusCode)
		check.Error = &errorMsg
	}
	
	_, err = s.healthCheckRepo.Create(check)
	if err != nil {
		log.Printf("Failed to save health check: %v", err)
		return
	}
	
	if success {
		log.Printf("Health check OK: %dms (HTTP %d)", check.ResponseTimeMs, resp.StatusCode)
	} else {
		log.Printf("Health check FAILED: %dms (HTTP %d)", check.ResponseTimeMs, resp.StatusCode)
	}
}

// recordFailure records a failed health check
func (s *MonitoringService) recordFailure(timestamp time.Time, statusCode *int, errorMsg string) {
	check := &models.HealthCheck{
		Timestamp:      timestamp,
		ResponseTimeMs: int(RequestTimeout.Milliseconds()), // Use timeout as response time for failures
		StatusCode:     statusCode,
		Success:        false,
		Error:          &errorMsg,
	}
	
	_, err := s.healthCheckRepo.Create(check)
	if err != nil {
		log.Printf("Failed to save failed health check: %v", err)
		return
	}
	
	log.Printf("Health check FAILED: %s", errorMsg)
}

// cleanup removes old health check data
func (s *MonitoringService) cleanup() {
	log.Println("Running health check data cleanup...")
	
	err := s.healthCheckRepo.DeleteOldChecks(DataRetention)
	if err != nil {
		log.Printf("Failed to cleanup old health checks: %v", err)
		return
	}
	
	log.Println("Health check cleanup completed")
}
