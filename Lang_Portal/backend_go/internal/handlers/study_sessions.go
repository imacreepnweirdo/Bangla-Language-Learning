package handlers

import (
	"database/sql"
	"lang-portal/internal/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type StudySessionHandler struct {
	db *sql.DB
}

func NewStudySessionHandler(db *sql.DB) *StudySessionHandler {
	return &StudySessionHandler{db: db}
}

func (h *StudySessionHandler) ReviewWord(c *gin.Context) {
	sessionIDStr := c.Param("id")
	sessionID, err := strconv.Atoi(sessionIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID"})
		return
	}

	wordIDStr := c.Param("word_id")
	wordID, err := strconv.Atoi(wordIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid word ID"})
		return
	}

	var req struct {
		Correct bool `json:"correct" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify study session exists
	var exists bool
	err = h.db.QueryRow("SELECT EXISTS(SELECT 1 FROM study_sessions WHERE id = ?)", sessionID).Scan(&exists)
	if err != nil || !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Study session not found"})
		return
	}

	// Verify word exists
	err = h.db.QueryRow("SELECT EXISTS(SELECT 1 FROM words WHERE id = ?)", wordID).Scan(&exists)
	if err != nil || !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Word not found"})
		return
	}

	// Insert word review item
	query := `
		INSERT INTO word_review_items (word_id, study_session_id, correct, created_at)
		VALUES (?, ?, ?, ?)
	`

	result, err := h.db.Exec(query, wordID, sessionID, req.Correct, time.Now())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := models.WordReviewResponse{
		ID:             int(id),
		Success:        true,
		WordID:         wordID,
		StudySessionID: sessionID,
		Correct:        req.Correct,
		CreatedAt:      time.Now(),
	}

	c.JSON(http.StatusOK, response)
}
