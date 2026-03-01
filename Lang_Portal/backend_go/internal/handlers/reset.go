package handlers

import (
	"database/sql"
	"net/http"
	"lang-portal/internal/models"

	"github.com/gin-gonic/gin"
)

type ResetHandler struct {
	db *sql.DB
}

func NewResetHandler(db *sql.DB) *ResetHandler {
	return &ResetHandler{db: db}
}

func (h *ResetHandler) ResetHistory(c *gin.Context) {
	// Delete all study-related data but keep words and groups
	tx, err := h.db.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer tx.Rollback()

	// Delete in correct order due to foreign key constraints
	_, err = tx.Exec("DELETE FROM word_review_items")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	_, err = tx.Exec("DELETE FROM study_activities")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	_, err = tx.Exec("DELETE FROM study_sessions")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := models.ResetResponse{
		Success: true,
		Message: "Study history reset successfully",
	}

	c.JSON(http.StatusOK, response)
}

func (h *ResetHandler) FullReset(c *gin.Context) {
	// Delete all data including words and groups
	tx, err := h.db.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer tx.Rollback()

	// Delete in correct order due to foreign key constraints
	tables := []string{
		"word_review_items",
		"study_activities", 
		"study_sessions",
		"words_groups",
		"words",
		"groups",
	}

	for _, table := range tables {
		_, err = tx.Exec("DELETE FROM " + table)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := models.ResetResponse{
		Success: true,
		Message: "System has been fully reset",
	}

	c.JSON(http.StatusOK, response)
}
