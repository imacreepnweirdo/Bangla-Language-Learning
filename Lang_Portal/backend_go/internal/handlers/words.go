package handlers

import (
	"database/sql"
	"lang-portal/internal/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type WordHandler struct {
	db *sql.DB
}

func NewWordHandler(db *sql.DB) *WordHandler {
	return &WordHandler{db: db}
}

func (h *WordHandler) GetWords(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit := 100
	offset := (page - 1) * limit

	// Count total words
	var totalItems int
	err = h.db.QueryRow("SELECT COUNT(*) FROM words").Scan(&totalItems)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get words with stats
	query := `
		SELECT 
			w.id, w.bengali, w.parts, w.english,
			COALESCE(SUM(CASE WHEN wri.correct = 1 THEN 1 ELSE 0 END), 0) as correct_count,
			COALESCE(SUM(CASE WHEN wri.correct = 0 THEN 1 ELSE 0 END), 0) as wrong_count
		FROM words w
		LEFT JOIN word_review_items wri ON w.id = wri.word_id
		GROUP BY w.id, w.bengali, w.parts, w.english
		ORDER BY w.id
		LIMIT ? OFFSET ?
	`

	rows, err := h.db.Query(query, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var words []models.WordWithStats
	for rows.Next() {
		var word models.WordWithStats
		err := rows.Scan(
			&word.ID,
			&word.Bengali,
			&word.PartsOfSpeech,
			&word.English,
			&word.CorrectCount,
			&word.WrongCount,
		)
		if err != nil {
			continue
		}
		words = append(words, word)
	}

	// Ensure words is never nil
	if words == nil {
		words = []models.WordWithStats{}
	}

	totalPages := (totalItems + limit - 1) / limit
	pagination := models.Pagination{
		Page:         page,
		TotalPages:   totalPages,
		TotalItems:   totalItems,
		ItemsPerPage: limit,
	}

	response := models.PaginatedWords{
		Items:      words,
		Pagination: pagination,
	}

	c.JSON(http.StatusOK, response)
}

func (h *WordHandler) GetWord(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	// Get word details
	query := `
		SELECT w.id, w.bengali, w.english, w.parts
		FROM words w
		WHERE w.id = ?
	`

	var word models.WordDetail
	err = h.db.QueryRow(query, id).Scan(&word.ID, &word.Bengali, &word.English, &word.Parts)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Word not found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get word stats
	statsQuery := `
		SELECT 
			COALESCE(SUM(CASE WHEN correct = 1 THEN 1 ELSE 0 END), 0) as correct_count,
			COALESCE(SUM(CASE WHEN correct = 0 THEN 1 ELSE 0 END), 0) as wrong_count
		FROM word_review_items
		WHERE word_id = ?
	`

	err = h.db.QueryRow(statsQuery, id).Scan(&word.Stats.CorrectCount, &word.Stats.WrongCount)
	if err != nil {
		word.Stats.CorrectCount = 0
		word.Stats.WrongCount = 0
	}

	// Get word groups
	groupsQuery := `
		SELECT g.id, g.name
		FROM groups g
		JOIN words_groups wg ON g.id = wg.group_id
		WHERE wg.word_id = ?
	`

	rows, err := h.db.Query(groupsQuery, id)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var group models.Group
			err := rows.Scan(&group.ID, &group.Name)
			if err != nil {
				continue
			}
			word.Groups = append(word.Groups, group)
		}
	}

	c.JSON(http.StatusOK, word)
}
