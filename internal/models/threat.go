package models

import (
	"time"
)

// ThreatNature represents the type of the threat
type ThreatNature int

const (
	AICrawler    ThreatNature = 1
	DDoS         ThreatNature = 2
	BruteForce   ThreatNature = 3
	XSS          ThreatNature = 4
	SQLInjection ThreatNature = 5
)

// ThreatStatus represents the status of the threat
type ThreatStatus int

const (
	Blocked    ThreatStatus = 1
	Detected   ThreatStatus = 2
	InAnalysis ThreatStatus = 3
)

// Threat represents a security threat detected for a site
type Threat struct {
	Nature ThreatNature `json:"nature"`
	Source []string     `json:"source"`
	Time   time.Time    `json:"time"`
	Site   string       `json:"site"`
	Status ThreatStatus `json:"status"`
}

// GetNatureName returns the string representation of the threat nature
func (t *Threat) GetNatureName() string {
	switch t.Nature {
	case AICrawler:
		return "AI Crawler"
	case DDoS:
		return "DDoS"
	case BruteForce:
		return "Brute Force"
	case XSS:
		return "XSS"
	case SQLInjection:
		return "SQL Injection"
	default:
		return "Unknown"
	}
}

// GetStatusName returns the string representation of the threat status
func (t *Threat) GetStatusName() string {
	switch t.Status {
	case Blocked:
		return "Blocked"
	case Detected:
		return "Detected"
	case InAnalysis:
		return "In Analysis"
	default:
		return "Unknown"
	}
}
