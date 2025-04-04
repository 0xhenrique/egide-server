package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"egide-server/internal/auth"
	"egide-server/internal/models"
	"egide-server/internal/repository"
	"egide-server/internal/service"
)

// MockSiteRepository is a mock implementation of SiteRepository
type MockSiteRepository struct {
	sites []*models.Site
}

func (m *MockSiteRepository) FindByUserID(userID int64) ([]*models.Site, error) {
	return m.sites, nil
}

// Other methods needed to implement the repository interface
func (m *MockSiteRepository) Create(site *models.Site) (int64, error) { return 0, nil }
func (m *MockSiteRepository) FindByID(id int64) (*models.Site, error) { return nil, nil }
func (m *MockSiteRepository) FindByDomain(userID int64, domain string) (*models.Site, error) { return nil, nil }
func (m *MockSiteRepository) Update(site *models.Site) error { return nil }
func (m *MockSiteRepository) Delete(id int64) error { return nil }

func TestGetRecentThreats(t *testing.T) {
	// Create mock site repository with test data
	mockSiteRepo := &MockSiteRepository{
		sites: []*models.Site{
			{
				ID:             1,
				UserID:         123,
				Domain:         "example.com",
				ProtectionMode: models.SimpleProtection,
				Active:         true,
			},
			{
				ID:             2,
				UserID:         123,
				Domain:         "another-example.com",
				ProtectionMode: models.HardenedProtection,
				Active:         true,
			},
		},
	}

	// Create threat service
	threatService := service.NewThreatService()

	// Create handler with mock repo and service
	handler := NewThreatHandler(mockSiteRepo, threatService)

	// Create request
	req, err := http.NewRequest("GET", "/api/threats", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Add user ID to context
	ctx := auth.WithUserID(context.Background(), 123)
	req = req.WithContext(ctx)

	// Create response recorder
	rr := httptest.NewRecorder()

	// Call handler
	handler.GetRecentThreats(rr, req)

	// Check status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check that we got valid JSON threats back
	var threats []*models.Threat
	err = json.Unmarshal(rr.Body.Bytes(), &threats)
	if err != nil {
		t.Errorf("could not parse response as JSON: %v", err)
	}

	// We should have between 4 and 10 threats (2-5 per site, 2 sites)
	if len(threats) < 4 || len(threats) > 10 {
		t.Errorf("unexpected number of threats: got %d, want between 4 and 10", len(threats))
	}

	// Check that all threats have valid fields
	for _, threat := range threats {
		if threat.Site != "example.com" && threat.Site != "another-example.com" {
			t.Errorf("threat has unexpected site: %s", threat.Site)
		}

		if threat.Nature < 1 || threat.Nature > 6 {
			t.Errorf("threat has invalid nature: %d", threat.Nature)
		}

		if threat.Status < 1 || threat.Status > 3 {
			t.Errorf("threat has invalid status: %d", threat.Status)
		}

		if len(threat.Source) < 1 {
			t.Errorf("threat has no sources")
		}
	}
}

func TestGetThreatDistribution(t *testing.T) {
	// Create mock site repository with test data
	mockSiteRepo := &MockSiteRepository{
		sites: []*models.Site{
			{
				ID:             1,
				UserID:         123,
				Domain:         "example.com",
				ProtectionMode: models.SimpleProtection,
				Active:         true,
			},
			{
				ID:             2,
				UserID:         123,
				Domain:         "another-example.com",
				ProtectionMode: models.HardenedProtection,
				Active:         true,
			},
		},
	}

	// Create threat service
	threatService := service.NewThreatService()

	// Create handler with mock repo and service
	handler := NewThreatHandler(mockSiteRepo, threatService)

	// Create request
	req, err := http.NewRequest("GET", "/api/threats/distribution", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Add user ID to context
	ctx := auth.WithUserID(context.Background(), 123)
	req = req.WithContext(ctx)

	// Create response recorder
	rr := httptest.NewRecorder()

	// Call handler
	handler.GetThreatDistribution(rr, req)

	// Check status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check that we got valid JSON distribution back
	var distribution []*models.ThreatDistribution
	err = json.Unmarshal(rr.Body.Bytes(), &distribution)
	if err != nil {
		t.Errorf("could not parse response as JSON: %v", err)
	}

	// We should have exactly 6 threat types
	if len(distribution) != 6 {
		t.Errorf("unexpected number of threat types: got %d, want 6", len(distribution))
	}

	// Check that all threat types are represented
	natureMap := make(map[models.ThreatNature]bool)
	for _, dist := range distribution {
		natureMap[dist.Nature] = true
		
		// Check that the count is positive
		if dist.Count <= 0 {
			t.Errorf("threat count should be positive, got %d for nature %d", dist.Count, dist.Nature)
		}
		
		// Verify that the nature is valid
		if dist.Nature < 1 || dist.Nature > 6 {
			t.Errorf("invalid threat nature: %d", dist.Nature)
		}
	}
	
	// Verify all 6 threat types are included
	for i := 1; i <= 6; i++ {
		if !natureMap[models.ThreatNature(i)] {
			t.Errorf("missing threat nature: %d", i)
		}
	}
}
