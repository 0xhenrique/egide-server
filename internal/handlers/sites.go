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

	active := true
	if input.Active != nil {
		active = *input.Active
	}

	site := &models.Site{
		UserID:         userID,
		Domain:         input.Domain,
		ProtectionMode: input.ProtectionMode,
		Active:         active,
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
	if input.Active != nil {
		existingSite.Active = *input.Active
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
