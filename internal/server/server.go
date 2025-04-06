package server

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"egide-server/internal/auth"
	"egide-server/internal/config"
	"egide-server/internal/handlers"
	"egide-server/internal/repository"
	"egide-server/internal/service"
)

type Server struct {
	server *http.Server
	config *config.Config
}

func New(cfg *config.Config, db *sql.DB) *Server {
	// Init repos
	userRepo := repository.NewUserRepository(db)
	siteRepo := repository.NewSiteRepository(db)

	// Init services
	authService := auth.NewGitHubService(cfg)
	threatService := service.NewThreatService()
	metricsService := service.NewMetricsService()

	authMiddleware := auth.NewMiddleware(cfg.JWTSecret)
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"}, // This should be restricted in production
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	authHandler := handlers.NewAuthHandler(authService, userRepo, cfg)
	siteHandler := handlers.NewSiteHandler(siteRepo)
	userHandler := handlers.NewUserHandler(userRepo)
	threatHandler := handlers.NewThreatHandler(siteRepo, threatService)
	metricsHandler := handlers.NewMetricsHandler(metricsService)

	// Public routes
	r.Group(func(r chi.Router) {
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		})
		
		// Auth routes
		r.Route("/auth", func(r chi.Router) {
			r.Get("/github", authHandler.GitHubLogin)
			r.Get("/callback", authHandler.GitHubCallback)
		})
	})

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(authMiddleware.Authenticate)
		
		// User routes
		r.Route("/api/users", func(r chi.Router) {
			r.Get("/me", userHandler.GetCurrentUser)
		})
		
		// Site routes
		r.Route("/api/sites", func(r chi.Router) {
			r.Get("/", siteHandler.ListSites)
			r.Post("/", siteHandler.CreateSite)
			r.Get("/{id}", siteHandler.GetSite)
			r.Put("/{id}", siteHandler.UpdateSite)
			r.Delete("/{id}", siteHandler.DeleteSite)
			r.Post("/{id}/verify", siteHandler.VerifySite)
			r.Post("/{id}/activate", siteHandler.ToggleSiteActivation)
		})
		
		// Threat routes
		r.Route("/api/threats", func(r chi.Router) {
			r.Get("/", threatHandler.GetRecentThreats)
			r.Get("/distribution", threatHandler.GetThreatDistribution)
		})
		
		// Metrics routes
		r.Route("/api/metrics", func(r chi.Router) {
			r.Get("/kpi", metricsHandler.GetKpi)
		})
	})

	return &Server{
		server: &http.Server{
			Addr:    fmt.Sprintf(":%d", cfg.ServerPort),
			Handler: r,
		},
		config: cfg,
	}
}

func (s *Server) Start() error {
	log.Printf("Starting server on port %d\n", s.config.ServerPort)
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
