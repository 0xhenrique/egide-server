package models

// ThreatDistribution represents the count of threats by nature
type ThreatDistribution struct {
	Nature ThreatNature `json:"nature"`
	Count  int          `json:"count"`
}
