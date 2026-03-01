package handlers

import (
	"database/sql"
	"net/http"
	"lang-portal/internal/models"

	"github.com/gin-gonic/gin"
)

type DashboardHandler struct {
	db *sql.DB
}

func NewDashboardHandler(db *sql.DB) *DashboardHandler {
	return &DashboardHandler{db: db}
}

func (h *DashboardHandler) GetLastStudySession(c *gin.Context) {
	query := `
		SELECT ss.id, ss.group_id, ss.created_at, ss.study_activities_id, g.name as group_name
		FROM study_sessions ss
		JOIN groups g ON ss.group_id = g.id
		ORDER BY ss.created_at DESC
		LIMIT 1
	`
	
	var session models.LastStudySession
	var studyActivitiesID sql.NullInt64
	
	err := h.db.QueryRow(query).Scan(
		&session.ID, 
		&session.GroupID, 
		&session.CreatedAt, 
		&studyActivitiesID,
		&session.GroupName,
	)
	
	if err == sql.ErrNoRows {
		c.JSON(http.StatusOK, gin.H{"message": "No study sessions found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	if studyActivitiesID.Valid {
		session.StudyActivitiesID = int(studyActivitiesID.Int64)
	}
	
	c.JSON(http.StatusOK, session)
}

func (h *DashboardHandler) GetStudyProgress(c *gin.Context) {
	var progress models.StudyProgress
	
	// Get total words studied (distinct words that have been reviewed)
	err := h.db.QueryRow("SELECT COUNT(DISTINCT word_id) FROM word_review_items").Scan(&progress.TotalWordsStudied)
	if err != nil {
		progress.TotalWordsStudied = 0
	}
	
	// Get total available words
	err = h.db.QueryRow("SELECT COUNT(*) FROM words").Scan(&progress.TotalAvailableWords)
	if err != nil {
		progress.TotalAvailableWords = 0
	}
	
	c.JSON(http.StatusOK, progress)
}

func (h *DashboardHandler) GetQuickStats(c *gin.Context) {
	var stats models.QuickStats
	
	// Get total words
	h.db.QueryRow("SELECT COUNT(*) FROM words").Scan(&stats.TotalWords)
	
	// Get total groups
	h.db.QueryRow("SELECT COUNT(*) FROM groups").Scan(&stats.TotalGroups)
	
	// Get words learned (distinct words with at least one correct review)
	h.db.QueryRow(`
		SELECT COUNT(DISTINCT word_id) 
		FROM word_review_items 
		WHERE correct = 1
	`).Scan(&stats.WordsLearned)
	
	// Get sessions completed
	h.db.QueryRow("SELECT COUNT(*) FROM study_sessions").Scan(&stats.SessionsCompleted)
	
	// Calculate accuracy rate
	var correctCount, totalCount int
	h.db.QueryRow("SELECT COUNT(*) FROM word_review_items WHERE correct = 1").Scan(&correctCount)
	h.db.QueryRow("SELECT COUNT(*) FROM word_review_items").Scan(&totalCount)
	
	if totalCount > 0 {
		stats.AccuracyRate = float64(correctCount) / float64(totalCount) * 100
	}
	
	// Current streak (simplified - consecutive days with activity)
	stats.CurrentStreak = 1 // Placeholder
	
	c.JSON(http.StatusOK, stats)
}
