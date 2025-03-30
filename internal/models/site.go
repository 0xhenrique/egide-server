package models

import "time"

type ProtectionMode string

const (
	// SimpleProtection represents basic security measures
	SimpleProtection ProtectionMode = "simple"
	
	// HardenedProtection represents advanced security measures (premium)
	HardenedProtection ProtectionMode = "hardened"
)

type Site struct {
	ID             int64          `json:"id"`
	UserID         int64          `json:"user_id"`
	Domain         string         `json:"domain"`
	ProtectionMode ProtectionMode `json:"protection_mode"`
	Active         bool           `json:"active"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
}

// Data required to create or update a site
type SiteInput struct {
	Domain         string         `json:"domain" validate:"required,fqdn"`
	ProtectionMode ProtectionMode `json:"protection_mode" validate:"required,oneof=simple hardened"`
	Active         *bool          `json:"active,omitempty"`
}
