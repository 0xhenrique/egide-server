package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"egide-server/internal/auth"
	"egide-server/internal/models"
	"egide-server/internal/repository"
)

type SiteHandler struct {
	siteRepo  *repository.SiteRepository
	validator *validator.Validate
}

func NewSiteHandler(siteRepo *repository.SiteRepository) *SiteHandler {
	return &SiteHandler{
		siteRepo:  siteRepo,
		validator: validator.New(),
	}
}

func (h *SiteHandler) ListSites(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.UserIDFromContext(r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	sites, err := h.siteRepo.FindByUserID(userID)
	if err != nil {
		http.Error(w, "Failed to fetch sites", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sites)
}

func (h *SiteHandler) GetSite(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.UserIDFromContext(r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	siteID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid site ID", http.StatusBadRequest)
		return
	}

	site, err := h.siteRepo.FindByID(siteID)
	if err != nil {
		http.Error(w, "Site not found", http.StatusNotFound)
		return
	}

	// Ensure the site belongs to the authenticated user
	if site.UserID != userID {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(site)
}

func (h *SiteHandler) CreateSite(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.UserIDFromContext(r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var input models.SiteInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.validator.Struct(input); err != nil {
		http.Error(w, "Validation error: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Always set active=false and verified=false for new sites
	// Ignoring any provided active value as it doesn't make sense for unverified sites
	active := false
	verified := false

	site := &models.Site{
		UserID:         userID,
		Domain:         input.Domain,
		ProtectionMode: input.ProtectionMode,
		Active:         active,
		Verified:       verified,
	}

	siteID, err := h.siteRepo.Create(site)
	if err != nil {
		http.Error(w, "Failed to create site: "+err.Error(), http.StatusInternalServerError)
		return
	}

	site, err = h.siteRepo.FindByID(siteID)
	if err != nil {
		http.Error(w, "Site created but failed to fetch", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(site)
}

func (h *SiteHandler) UpdateSite(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.UserIDFromContext(r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	siteID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid site ID", http.StatusBadRequest)
		return
	}

	existingSite, err := h.siteRepo.FindByID(siteID)
	if err != nil {
		http.Error(w, "Site not found", http.StatusNotFound)
		return
	}

	if existingSite.UserID != userID {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	var input models.SiteInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.validator.Struct(input); err != nil {
		http.Error(w, "Validation error: "+err.Error(), http.StatusBadRequest)
		return
	}

	existingSite.Domain = input.Domain
	existingSite.ProtectionMode = input.ProtectionMode
	
	// If user is trying to set active=true but site is not verified, reject
	if input.Active != nil {
		if *input.Active && !existingSite.Verified {
			http.Error(w, "Cannot activate protection for unverified site. Please verify the site first.", http.StatusBadRequest)
			return
		}
		existingSite.Active = *input.Active
	}
	
	// Handle verified status - only in Update endpoint for admin purposes,
	// normal users should use the dedicated VerifySite endpoint
	if input.Verified != nil {
		existingSite.Verified = *input.Verified
		
		// If setting verified to false, DO NOT forget to set 'active' to false
		if !*input.Verified {
			existingSite.Active = false
		}
	}

	if err := h.siteRepo.Update(existingSite); err != nil {
		http.Error(w, "Failed to update site: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(existingSite)
}

func (h *SiteHandler) DeleteSite(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.UserIDFromContext(r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	siteID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid site ID", http.StatusBadRequest)
		return
	}

	existingSite, err := h.siteRepo.FindByID(siteID)
	if err != nil {
		http.Error(w, "Site not found", http.StatusNotFound)
		return
	}

	if existingSite.UserID != userID {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	if err := h.siteRepo.Delete(siteID); err != nil {
		http.Error(w, "Failed to delete site: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *SiteHandler) VerifySite(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.UserIDFromContext(r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	siteID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid site ID", http.StatusBadRequest)
		return
	}

	existingSite, err := h.siteRepo.FindByID(siteID)
	if err != nil {
		http.Error(w, "Site not found", http.StatusNotFound)
		return
	}

	if existingSite.UserID != userID {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	var input struct {
		Verified bool `json:"verified"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// If unverifying a site, also ensure it's deactivated
	if !input.Verified && existingSite.Active {
		// First update active status to false
		existingSite.Active = false
		if err := h.siteRepo.Update(existingSite); err != nil {
			http.Error(w, "Failed to deactivate site during unverification: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if err := h.siteRepo.UpdateVerificationStatus(siteID, input.Verified); err != nil {
		http.Error(w, "Failed to update verification status: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Get the updated site
	updatedSite, err := h.siteRepo.FindByID(siteID)
	if err != nil {
		http.Error(w, "Site updated but failed to fetch", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedSite)
}

func (h *SiteHandler) ToggleSiteActivation(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.UserIDFromContext(r.Context())
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	siteID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid site ID", http.StatusBadRequest)
		return
	}

	existingSite, err := h.siteRepo.FindByID(siteID)
	if err != nil {
		http.Error(w, "Site not found", http.StatusNotFound)
		return
	}

	if existingSite.UserID != userID {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	var input struct {
		Active bool `json:"active"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if input.Active && !existingSite.Verified {
		http.Error(w, "Cannot activate protection for unverified site. Please verify the site first.", http.StatusBadRequest)
		return
	}

	// Update active status
	existingSite.Active = input.Active
	if err := h.siteRepo.Update(existingSite); err != nil {
		http.Error(w, "Failed to update activation status: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(existingSite)
}
