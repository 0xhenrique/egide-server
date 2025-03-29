package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/0xhenrique/egide-server/database"
)

type WebsiteRequest struct {
	Domain       string `json:"domain" binding:"required"`
	Description  string `json:"description"`
	ProtectionMode string `json:"protection_mode" binding:"required,oneof=simple hardened"`
}

func GetWebsites(c *gin.Context) {
	userID := c.MustGet("user_id").(int64)
	db := c.MustGet("db").(*sql.DB)

	rows, err := db.Query(`
		SELECT id, user_id, domain, description, protection_mode, created_at, updated_at 
		FROM websites 
		WHERE user_id = ?
	`, userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch websites"})
		return
	}
	defer rows.Close()

	var websites []database.Website
	for rows.Next() {
		var website database.Website
		err := rows.Scan(
			&website.ID,
			&website.UserID,
			&website.Domain,
			&website.Description,
			&website.ProtectionMode,
			&website.CreatedAt,
			&website.UpdatedAt,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse website data"})
			return
		}
		websites = append(websites, website)
	}

	c.JSON(http.StatusOK, gin.H{"websites": websites})
}

func GetWebsite(c *gin.Context) {
	userID := c.MustGet("user_id").(int64)
	websiteID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid website ID"})
		return
	}

	db := c.MustGet("db").(*sql.DB)

	var website database.Website
	err = db.QueryRow(`
		SELECT id, user_id, domain, description, protection_mode, created_at, updated_at 
		FROM websites 
		WHERE id = ? AND user_id = ?
	`, websiteID, userID).Scan(
		&website.ID,
		&website.UserID,
		&website.Domain,
		&website.Description,
		&website.ProtectionMode,
		&website.CreatedAt,
		&website.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Website not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch website"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"website": website})
}

func AddWebsite(c *gin.Context) {
	userID := c.MustGet("user_id").(int64)

	var req WebsiteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := c.MustGet("db").(*sql.DB)

	var websiteID int64
	err := db.QueryRow(`
		INSERT INTO websites (user_id, domain, description, protection_mode) 
		VALUES (?, ?, ?, ?) 
		RETURNING id
	`, userID, req.Domain, req.Description, req.ProtectionMode).Scan(&websiteID)

	if err != nil {
		if err.Error() == "UNIQUE constraint failed: websites.user_id, websites.domain" {
			c.JSON(http.StatusConflict, gin.H{"error": "Website domain already registered"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add website"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Website added successfully",
		"website_id": websiteID,
	})
}

func UpdateWebsite(c *gin.Context) {
	userID := c.MustGet("user_id").(int64)
	websiteID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid website ID"})
		return
	}

	var req WebsiteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := c.MustGet("db").(*sql.DB)

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM websites WHERE id = ? AND user_id = ?", 
		websiteID, userID).Scan(&count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	if count == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Website not found"})
		return
	}

	_, err = db.Exec(`
		UPDATE websites 
		SET domain = ?, description = ?, protection_mode = ?, updated_at = CURRENT_TIMESTAMP 
		WHERE id = ? AND user_id = ?
	`, req.Domain, req.Description, req.ProtectionMode, websiteID, userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update website"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Website updated successfully"})
}

func DeleteWebsite(c *gin.Context) {
	userID := c.MustGet("user_id").(int64)
	websiteID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid website ID"})
		return
	}

	db := c.MustGet("db").(*sql.DB)

	result, err := db.Exec("DELETE FROM websites WHERE id = ? AND user_id = ?", 
		websiteID, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete website"})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Website not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Website deleted successfully"})
}

func UpdateProtectionMode(c *gin.Context) {
	userID := c.MustGet("user_id").(int64)
	websiteID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid website ID"})
		return
	}

	var req struct {
		ProtectionMode string `json:"protection_mode" binding:"required,oneof=simple hardened"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := c.MustGet("db").(*sql.DB)

	result, err := db.Exec(`
		UPDATE websites 
		SET protection_mode = ?, updated_at = CURRENT_TIMESTAMP 
		WHERE id = ? AND user_id = ?
	`, req.ProtectionMode, websiteID, userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update protection mode"})
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Website not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Protection mode updated successfully",
		"protection_mode": req.ProtectionMode,
	})
}
