package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"egide-server/internal/auth"
	"egide-server/internal/config"
	"egide-server/internal/models"
	"egide-server/internal/repository"
)

type AuthHandler struct {
	authService *auth.GitHubService
	userRepo    *repository.UserRepository
	config      *config.Config
}

func NewAuthHandler(authService *auth.GitHubService, userRepo *repository.UserRepository, config *config.Config) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		userRepo:    userRepo,
		config:      config,
	}
}

func (h *AuthHandler) GitHubLogin(w http.ResponseWriter, r *http.Request) {
	url := h.authService.GetAuthURL()
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *AuthHandler) GitHubCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Missing authorization code", http.StatusBadRequest)
		return
	}

	accessToken, err := h.authService.Exchange(code)
	if err != nil {
		http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	githubUser, err := h.authService.GetUser(accessToken)
	if err != nil {
		http.Error(w, "Failed to get user info: "+err.Error(), http.StatusInternalServerError)
		return
	}

	user, err := h.userRepo.FindByGitHubID(strconv.FormatInt(githubUser.ID, 10))
	if err != nil {
		user = &models.User{
			GitHubID: strconv.FormatInt(githubUser.ID, 10),
			Username: githubUser.Login,
			Email:    githubUser.Email,
		}

		userID, err := h.userRepo.Create(user)
		if err != nil {
			http.Error(w, "Failed to create user: "+err.Error(), http.StatusInternalServerError)
			return
		}
		user.ID = userID
	}

	token, err := h.authService.GenerateToken(user.ID)
	if err != nil {
		http.Error(w, "Failed to generate token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := struct {
		Token  string           `json:"token"`
		User   *models.User     `json:"user"`
		Expiry int              `json:"expiry"`
	}{
		Token:  token,
		User:   user,
		Expiry: 86400, // 24 hours in seconds
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
