package handlers

import (
	"database/sql"
	"lang-portal/internal/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type StudyActivityHandler struct {
	db *sql.DB
}

func NewStudyActivityHandler(db *sql.DB) *StudyActivityHandler {
	return &StudyActivityHandler{db: db}
}

func (h *StudyActivityHandler) GetStudyActivity(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	// Return mock data as per spec since study_activities table doesn't have name/thumbnail/description fields
	activity := models.StudyActivityDetail{
		ID:           id,
		Name:         "Vocabulary Practice",
		ThumbnailURL: "https://example.com/thumbnail.jpg",
		Description:  "Practice your vocabulary with flashcards",
	}

	c.JSON(http.StatusOK, activity)
}

func (h *StudyActivityHandler) GetStudySessions(c *gin.Context) {
	idStr := c.Param("id")
	activityID, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	pageStr := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit := 100
	offset := (page - 1) * limit

	// Count total sessions
	var totalItems int
	countQuery := `
		SELECT COUNT(DISTINCT ss.id)
		FROM study_sessions ss
		JOIN study_activities sa ON ss.id = sa.study_session_id
		WHERE sa.id = ?
	`
	err = h.db.QueryRow(countQuery, activityID).Scan(&totalItems)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get sessions with pagination
	query := `
		SELECT DISTINCT ss.id, g.name as group_name, ss.created_at
		FROM study_sessions ss
		JOIN study_activities sa ON ss.id = sa.study_session_id
		JOIN groups g ON ss.group_id = g.id
		WHERE sa.id = ?
		ORDER BY ss.created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := h.db.Query(query, activityID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var sessions []models.StudySessionWithDetails
	for rows.Next() {
		var session models.StudySessionWithDetails
		err := rows.Scan(&session.ID, &session.GroupName, &session.StartTime)
		if err != nil {
			continue
		}

		// Mock additional fields as per spec
		session.ActivityName = "Vocabulary Quiz"
		session.EndTime = session.StartTime
		session.ReviewItemsCount = 10

		sessions = append(sessions, session)
	}

	totalPages := (totalItems + limit - 1) / limit
	pagination := models.Pagination{
		Page:         page,
		TotalPages:   totalPages,
		TotalItems:   totalItems,
		ItemsPerPage: limit,
	}

	response := models.PaginatedStudySessions{
		Items:      sessions,
		Pagination: pagination,
	}

	c.JSON(http.StatusOK, response)
}

func (h *StudyActivityHandler) CreateStudyActivity(c *gin.Context) {
	var req models.CreateStudySessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create study session
	tx, err := h.db.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer tx.Rollback()

	result, err := tx.Exec("INSERT INTO study_sessions (group_id) VALUES (?)", req.GroupID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	sessionID, err := result.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Create study activity
	_, err = tx.Exec("INSERT INTO study_activities (study_session_id, group_id) VALUES (?, ?)", sessionID, req.GroupID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := models.CreateStudySessionResponse{
		ID:      int(sessionID),
		GroupID: req.GroupID,
	}

	c.JSON(http.StatusOK, response)
}
