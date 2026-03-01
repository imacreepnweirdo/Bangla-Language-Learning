package handlers

import (
	"database/sql"
	"lang-portal/internal/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type GroupHandler struct {
	db *sql.DB
}

func NewGroupHandler(db *sql.DB) *GroupHandler {
	return &GroupHandler{db: db}
}

func (h *GroupHandler) GetGroups(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit := 100
	offset := (page - 1) * limit

	// Count total groups
	var totalItems int
	err = h.db.QueryRow("SELECT COUNT(*) FROM groups").Scan(&totalItems)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get groups with word counts
	query := `
		SELECT 
			g.id, g.name,
			COUNT(wg.word_id) as words_count
		FROM groups g
		LEFT JOIN words_groups wg ON g.id = wg.group_id
		GROUP BY g.id, g.name
		ORDER BY g.id
		LIMIT ? OFFSET ?
	`

	rows, err := h.db.Query(query, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var groups []models.GroupWithWordCount
	for rows.Next() {
		var group models.GroupWithWordCount
		err := rows.Scan(&group.ID, &group.Name, &group.WordsCount)
		if err != nil {
			continue
		}
		groups = append(groups, group)
	}

	totalPages := (totalItems + limit - 1) / limit
	pagination := models.Pagination{
		Page:         page,
		TotalPages:   totalPages,
		TotalItems:   totalItems,
		ItemsPerPage: limit,
	}

	response := models.PaginatedGroups{
		Items:      groups,
		Pagination: pagination,
	}

	c.JSON(http.StatusOK, response)
}

func (h *GroupHandler) GetGroup(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	// Get group details
	query := `
		SELECT g.id, g.name
		FROM groups g
		WHERE g.id = ?
	`

	var group models.GroupDetail
	err = h.db.QueryRow(query, id).Scan(&group.ID, &group.Name)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Group not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get total word count
	countQuery := `
		SELECT COUNT(*)
		FROM words_groups
		WHERE group_id = ?
	`

	err = h.db.QueryRow(countQuery, id).Scan(&group.Stats.TotalWordCount)
	if err != nil {
		group.Stats.TotalWordCount = 0
	}

	// Get words (empty array as per spec)
	group.Words = []models.WordStat{}

	c.JSON(http.StatusOK, group)
}

func (h *GroupHandler) GetGroupWords(c *gin.Context) {
	idStr := c.Param("id")
	groupID, err := strconv.Atoi(idStr)
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

	// Count total words in group
	var totalItems int
	countQuery := `
		SELECT COUNT(*)
		FROM words_groups wg
		JOIN words w ON wg.word_id = w.id
		WHERE wg.group_id = ?
	`
	err = h.db.QueryRow(countQuery, groupID).Scan(&totalItems)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get words with stats
	query := `
		SELECT 
			w.id, w.bengali, w.english,
			COALESCE(SUM(CASE WHEN wri.correct = 1 THEN 1 ELSE 0 END), 0) as correct_count,
			COALESCE(SUM(CASE WHEN wri.correct = 0 THEN 1 ELSE 0 END), 0) as wrong_count
		FROM words w
		JOIN words_groups wg ON w.id = wg.word_id
		LEFT JOIN word_review_items wri ON w.id = wri.word_id
		WHERE wg.group_id = ?
		GROUP BY w.id, w.bengali, w.english
		ORDER BY w.id
		LIMIT ? OFFSET ?
	`

	rows, err := h.db.Query(query, groupID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var words []models.WordStat
	for rows.Next() {
		var word models.WordStat
		err := rows.Scan(
			&word.ID,
			&word.Bengali,
			&word.English,
			&word.CorrectCount,
			&word.WrongCount,
		)
		if err != nil {
			continue
		}
		words = append(words, word)
	}

	totalPages := (totalItems + limit - 1) / limit
	pagination := models.Pagination{
		Page:         page,
		TotalPages:   totalPages,
		TotalItems:   totalItems,
		ItemsPerPage: limit,
	}

	response := models.PaginatedGroupWords{
		Words:      words,
		Pagination: pagination,
	}

	c.JSON(http.StatusOK, response)
}
