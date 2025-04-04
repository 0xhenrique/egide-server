package handlers

import (
	"fmt"
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

    // Get the redirect URL (or use a default)
    frontendURL := h.config.FrontendURL
    if frontendURL == "" {
        frontendURL = "http://localhost:3000" // Default frontend URL
    }

    // Redirect to frontend with the token
    // Using hash fragment is more secure for tokens
    redirectURL := fmt.Sprintf("%s/auth/callback#token=%s&expiry=%d&userId=%d", 
        frontendURL, token, 86400, user.ID)
    
    http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
}
