/*
 * Copyright (C) 2025 Henrique Marques (0xhenrique)
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program. If not, see <http://www.gnu.org/licenses/>.
 */

package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

    "github.com/gin-contrib/sessions"
    "github.com/gin-contrib/sessions/cookie"

	//"github.com/0xhenrique/egide-server/config"
	"github.com/0xhenrique/egide-server/database"
	"github.com/0xhenrique/egide-server/handlers"
	"github.com/0xhenrique/egide-server/middleware"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

    if os.Getenv("SESSION_SECRET") == "" {
        log.Println("Warning: SESSION_SECRET not set, using a default value")
        os.Setenv("SESSION_SECRET", "egide-secret-key")
    }

	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	router := gin.Default()
	router.Use(middleware.CORS())
	router.Use(middleware.Database(db))
	store := cookie.NewStore([]byte(os.Getenv("SESSION_SECRET")))
	router.Use(sessions.Sessions("github_auth_session", store))
    handlers.InitAuth()

	public := router.Group("/api")
	{
		//public.POST("/register", handlers.Register)
		public.POST("/login", handlers.OAuthLogin)
		public.GET("/auth/github", handlers.OAuthLogin)
		public.GET("/auth/github/callback", handlers.OAuthCallback)
	}

	// Protected routes (needs to be logged in)
	protected := router.Group("/api")
	protected.Use(middleware.Auth())
	{
		protected.GET("/user", handlers.GetUser)
		protected.PUT("/user", handlers.UpdateUser)
		protected.DELETE("/user", handlers.DeleteUser)

		protected.GET("/websites", handlers.GetWebsites)
		protected.POST("/websites", handlers.AddWebsite)
		protected.GET("/websites/:id", handlers.GetWebsite)
		protected.PUT("/websites/:id", handlers.UpdateWebsite)
		protected.DELETE("/websites/:id", handlers.DeleteWebsite)

		protected.PUT("/websites/:id/protection", handlers.UpdateProtectionMode)

		//protected.POST("/logout", handlers.Logout)
		public.GET("/logout", handlers.Logout)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
