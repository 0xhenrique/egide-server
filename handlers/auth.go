package handlers

import (
	"database/sql"
	"net/http"
	"strings"
	"time"
	"os"

    "github.com/gin-contrib/sessions"
    "github.com/gin-contrib/sessions/cookie"
    "github.com/markbates/goth"
    "github.com/markbates/goth/gothic"
    "github.com/markbates/goth/providers/github"

	"github.com/0xhenrique/egide-server/database"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	//"golang.org/x/crypto/bcrypt"
)

func InitAuth() {
    store := cookie.NewStore([]byte("secret"))
    gothic.Store = store
	goth.UseProviders(
		github.New(
			os.Getenv("GITHUB_KEY"), 
			os.Getenv("GITHUB_SECRET"), 
			"http://127.0.0.1:8080/api/auth/github/callback",
		),
	)
}

func OAuthLogin(c *gin.Context) {
	c.Request.URL.RawQuery = "provider=github"
	gothic.BeginAuthHandler(c.Writer, c.Request)
}

func OAuthCallback(c *gin.Context) {
	c.Request.URL.RawQuery = "provider=github"
	user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to authenticate user: " + err.Error()})
		return
	}

	db := c.MustGet("db").(*sql.DB)

	// Check if user already exists
	var existingUserID int64
	err = db.QueryRow(`SELECT id FROM users WHERE email = ?`, user.Email).Scan(&existingUserID)

	if err != nil && err != sql.ErrNoRows {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	var userID int64
	if err == sql.ErrNoRows {
		// User doesn't exist, create a new one
		err = db.QueryRow(`
			INSERT INTO users (username, email, github_id) 
			VALUES (?, ?, ?) 
			RETURNING id
		`, user.Name, user.Email, user.UserID).Scan(&userID)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user: " + err.Error()})
			return
		}
	} else {
		userID = existingUserID
		// Update GitHub ID if needed
		_, err = db.Exec(`UPDATE users SET github_id = ? WHERE id = ?`, user.UserID, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
			return
		}
	}

	// Create a session
	token := uuid.New().String()
	expiresAt := time.Now().Add(24 * time.Hour)

	_, err = db.Exec(`
		INSERT INTO sessions (user_id, token, expires_at) 
		VALUES (?, ?, ?)
	`, userID, token, expiresAt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create session"})
		return
	}

	// Return the token in the response
	c.JSON(http.StatusOK, gin.H{
		"token":      token,
		"expires_at": expiresAt,
		"user": gin.H{
			"id":       userID,
			"username": user.Name,
			"email":    user.Email,
			"avatar":   user.AvatarURL,
		},
	})
}

func Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) == 2 && tokenParts[0] == "Bearer" {
			token := tokenParts[1]
			db := c.MustGet("db").(*sql.DB)
			db.Exec("DELETE FROM sessions WHERE token = ?", token)
		}
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

func GetUser(c *gin.Context) {
	userID := c.MustGet("user_id").(int64)
	db := c.MustGet("db").(*sql.DB)

	var user database.User
	err := db.QueryRow(`
		SELECT id, username, email, created_at, updated_at 
		FROM users 
		WHERE id = ?
	`, userID).Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user information"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

func UpdateUser(c *gin.Context) {
	userID := c.MustGet("user_id").(int64)
	
	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := c.MustGet("db").(*sql.DB)
	
	// Start transaction
	tx, err := db.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	
	if req.Username != "" {
		_, err = tx.Exec("UPDATE users SET username = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?", 
			req.Username, userID)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update username"})
			return
		}
	}
	
	if req.Email != "" {
		_, err = tx.Exec("UPDATE users SET email = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?", 
			req.Email, userID)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update email"})
			return
		}
	}
	
	if err = tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user information"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "User information updated successfully"})
}

func DeleteUser(c *gin.Context) {
	userID := c.MustGet("user_id").(int64)
	db := c.MustGet("db").(*sql.DB)

	// cascade delete sessions and websites
	_, err := db.Exec("DELETE FROM users WHERE id = ?", userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user account"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User account deleted successfully"})
}
