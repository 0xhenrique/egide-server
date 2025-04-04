package handlers

import (
	"encoding/json"
	"net/http"

	"egide-server/internal/auth"
	"egide-server/internal/repository"
	"egide-server/internal/service"
)

type ThreatHandler struct {
	siteRepo      *repository.SiteRepository
	threatService *service.ThreatService
}

func NewThreatHandler(siteRepo *repository.SiteRepository, threatService *service.ThreatService) *ThreatHandler {
	return &ThreatHandler{
		siteRepo:      siteRepo,
		threatService: threatService,
	}
}

// GetRecentThreats handles GET /api/threats
func (h *ThreatHandler) GetRecentThreats(w http.ResponseWriter, r *http.Request) {
	// Get user ID from authenticated context
	userID, err := auth.UserIDFromContext(r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get all sites for the user
	sites, err := h.siteRepo.FindByUserID(userID)
	if err != nil {
		http.Error(w, "Error fetching sites", http.StatusInternalServerError)
		return
	}

	// Get threats for all sites
	threats, err := h.threatService.GetRecentThreats(sites)
	if err != nil {
		http.Error(w, "Error fetching threats: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return threats as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(threats)
}
