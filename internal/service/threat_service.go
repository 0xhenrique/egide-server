package service

import (
	"math/rand"
	"strconv"
	"time"

	"egide-server/internal/models"
)

// ThreatService handles threat data operations
type ThreatService struct {
	// In the future, this could contain a client to communicate with egide-data
	rand *rand.Rand
}

// NewThreatService creates a new threat service
func NewThreatService() *ThreatService {
	// Create a new random source with fixed seed for reproducible mock data
	source := rand.NewSource(time.Now().UnixNano())
	return &ThreatService{
		rand: rand.New(source),
	}
}

// GetRecentThreats returns recent threats for the given sites
// Currently mocks the data, but in the future will fetch from egide-data
func (s *ThreatService) GetRecentThreats(sites []*models.Site) ([]*models.Threat, error) {
	var threats []*models.Threat

	for _, site := range sites {
		// Generate between 2 and 5 threats per site
		numThreats := s.rand.Intn(4) + 2
		for i := 0; i < numThreats; i++ {
			threat := s.generateMockThreat(site.Domain)
			threats = append(threats, threat)
		}
	}

	return threats, nil
}

// GetThreatDistribution returns the distribution of threats by nature
// Currently mocks the data, but in the future will fetch from egide-data
func (s *ThreatService) GetThreatDistribution(sites []*models.Site) ([]*models.ThreatDistribution, error) {
	// In a real implementation, this would query the egide-data service
	// For now, we'll generate mock data with realistic looking numbers
	
	// Create a map to store the count for each threat nature
	distributionMap := make(map[models.ThreatNature]int)
	
	// Initialize with all threat types (1-6)
	for i := 1; i <= 6; i++ {
		distributionMap[models.ThreatNature(i)] = 0
	}
	
	// Generate realistic looking counts
	// AI Crawlers (1) and SQL Injections (5) are usually more common
	distributionMap[1] = 20000 + s.rand.Intn(20000)    // AI Crawler: 20,000 - 40,000
	distributionMap[2] = 500 + s.rand.Intn(1000)       // DDoS: 500 - 1,500
	distributionMap[3] = 20 + s.rand.Intn(100)         // Brute Force: 20 - 120
	distributionMap[4] = 10 + s.rand.Intn(50)          // XSS: 10 - 60
	distributionMap[5] = 2000 + s.rand.Intn(3000)      // SQL Injection: 2,000 - 5,000
	distributionMap[6] = 5 + s.rand.Intn(20)           // Other: 5 - 25
	
	// Convert map to slice of ThreatDistribution
	var distribution []*models.ThreatDistribution
	for nature, count := range distributionMap {
		distribution = append(distribution, &models.ThreatDistribution{
			Nature: nature,
			Count:  count,
		})
	}
	
	return distribution, nil
}

// generateMockThreat creates a mock threat for a given domain
func (s *ThreatService) generateMockThreat(domain string) *models.Threat {
	// Random IP addresses for sources
	ipPrefixes := []string{"192.168", "10.0", "172.16", "8.8"}
	var sources []string
	numSources := s.rand.Intn(3) + 1 // 1 to 3 sources
	for i := 0; i < numSources; i++ {
		prefix := ipPrefixes[s.rand.Intn(len(ipPrefixes))]
		ip := prefix + "." + strconv.Itoa(s.rand.Intn(255)) + "." + strconv.Itoa(s.rand.Intn(255))
		sources = append(sources, ip)
	}

	// Random time within the last 7 days
	maxAge := 7 * 24 * time.Hour
	randomAge := time.Duration(s.rand.Int63n(int64(maxAge)))
	timestamp := time.Now().Add(-randomAge)

	return &models.Threat{
		Nature: models.ThreatNature(s.rand.Intn(5) + 1), // 1 to 5
		Source: sources,
		Time:   timestamp,
		Site:   domain,
		Status: models.ThreatStatus(s.rand.Intn(3) + 1), // 1 to 3
	}
}
