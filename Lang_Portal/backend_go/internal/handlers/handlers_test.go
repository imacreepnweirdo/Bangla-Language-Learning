package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io"
	"lang-portal/internal/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}

	// Create tables for testing
	migrations := []string{
		`CREATE TABLE words (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			bengali TEXT NOT NULL,
			parts TEXT NOT NULL,
			english TEXT NOT NULL
		)`,
		`CREATE TABLE groups (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL
		)`,
		`CREATE TABLE words_groups (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			word_id INTEGER NOT NULL,
			group_id INTEGER NOT NULL,
			FOREIGN KEY (word_id) REFERENCES words(id),
			FOREIGN KEY (group_id) REFERENCES groups(id)
		)`,
		`CREATE TABLE study_sessions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			group_id INTEGER NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			study_activities_id INTEGER
		)`,
		`CREATE TABLE study_activities (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			study_session_id INTEGER NOT NULL,
			group_id INTEGER NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (study_session_id) REFERENCES study_sessions(id)
		)`,
		`CREATE TABLE word_review_items (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			word_id INTEGER NOT NULL,
			study_session_id INTEGER NOT NULL,
			correct BOOLEAN NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (word_id) REFERENCES words(id),
			FOREIGN KEY (study_session_id) REFERENCES study_sessions(id)
		)`,
	}

	for _, migration := range migrations {
		_, err := db.Exec(migration)
		if err != nil {
			t.Fatal(err)
		}
	}

	return db
}

func seedTestData(t *testing.T, db *sql.DB) {
	// Insert test group
	_, err := db.Exec("INSERT INTO groups (name) VALUES (?)", "Basic Greetings")
	if err != nil {
		t.Fatal(err)
	}

	// Insert test word
	_, err = db.Exec("INSERT INTO words (bengali, parts, english) VALUES (?, ?, ?)",
		"স্বাগতম", "অব্যয়", "welcome")
	if err != nil {
		t.Fatal(err)
	}

	// Link word to group
	_, err = db.Exec("INSERT INTO words_groups (word_id, group_id) VALUES (1, 1)")
	if err != nil {
		t.Fatal(err)
	}

	// Insert study session
	_, err = db.Exec("INSERT INTO study_sessions (group_id, study_activities_id) VALUES (1, 1)")
	if err != nil {
		t.Fatal(err)
	}
}

func TestDashboardHandlers(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	seedTestData(t, db)

	gin.SetMode(gin.TestMode)
	handler := NewDashboardHandler(db)

	t.Run("GetLastStudySession", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/api/dashboard/last_study_session", nil)

		handler.GetLastStudySession(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response models.LastStudySession
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, 1, response.ID)
		assert.Equal(t, "Basic Greetings", response.GroupName)
	})

	t.Run("GetStudyProgress", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/api/dashboard/study_progress", nil)

		handler.GetStudyProgress(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response models.StudyProgress
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, 0, response.TotalWordsStudied) // No reviews yet
		assert.Equal(t, 1, response.TotalAvailableWords)
	})

	t.Run("GetQuickStats", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/api/dashboard/quick-stats", nil)

		handler.GetQuickStats(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response models.QuickStats
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, 1, response.TotalWords)
		assert.Equal(t, 1, response.TotalGroups)
		assert.Equal(t, 1, response.SessionsCompleted)
	})
}

func TestWordHandlers(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	seedTestData(t, db)

	gin.SetMode(gin.TestMode)
	handler := NewWordHandler(db)

	t.Run("GetWords", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/api/words?page=1", nil)

		handler.GetWords(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response models.PaginatedWords
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response.Items, 1)
		assert.Equal(t, 1, response.Pagination.Page)
		assert.Equal(t, "স্বাগতম", response.Items[0].Bengali)
	})

	t.Run("GetWord", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "1"}}
		c.Request, _ = http.NewRequest("GET", "/api/words/1", nil)

		handler.GetWord(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response models.WordDetail
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, 1, response.ID)
		assert.Equal(t, "স্বাগতম", response.Bengali)
		assert.Equal(t, "অব্যয়", response.Parts)
		assert.Equal(t, "welcome", response.English)
	})

	t.Run("GetWordNotFound", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "999"}}
		c.Request, _ = http.NewRequest("GET", "/api/words/999", nil)

		handler.GetWord(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestGroupHandlers(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	seedTestData(t, db)

	gin.SetMode(gin.TestMode)
	handler := NewGroupHandler(db)

	t.Run("GetGroups", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/api/groups?page=1", nil)

		handler.GetGroups(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response models.PaginatedGroups
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response.Items, 1)
		assert.Equal(t, "Basic Greetings", response.Items[0].Name)
		assert.Equal(t, 1, response.Items[0].WordsCount)
	})

	t.Run("GetGroup", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "1"}}
		c.Request, _ = http.NewRequest("GET", "/api/groups/1", nil)

		handler.GetGroup(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response models.GroupDetail
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, 1, response.ID)
		assert.Equal(t, "Basic Greetings", response.Name)
		assert.Equal(t, 1, response.Stats.TotalWordCount)
	})
}

func TestStudyActivityHandlers(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	seedTestData(t, db)

	gin.SetMode(gin.TestMode)
	handler := NewStudyActivityHandler(db)

	t.Run("GetStudyActivity", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{{Key: "id", Value: "1"}}
		c.Request, _ = http.NewRequest("GET", "/api/study_activities/1", nil)

		handler.GetStudyActivity(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response models.StudyActivityDetail
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, 1, response.ID)
		assert.Equal(t, "Vocabulary Practice", response.Name)
	})

	t.Run("CreateStudyActivity", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/api/study_activities", nil)
		c.Request.Header.Set("Content-Type", "application/json")

		requestBody := models.CreateStudySessionRequest{
			GroupID:         1,
			StudyActivityID: 1,
		}
		jsonBody, _ := json.Marshal(requestBody)
		c.Request.Body = io.NopCloser(bytes.NewBuffer(jsonBody))

		handler.CreateStudyActivity(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response models.CreateStudySessionResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, 1, response.GroupID)
	})
}

func TestStudySessionHandlers(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	seedTestData(t, db)

	gin.SetMode(gin.TestMode)
	handler := NewStudySessionHandler(db)

	t.Run("ReviewWord", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{
			{Key: "id", Value: "1"},
			{Key: "word_id", Value: "1"},
		}
		c.Request, _ = http.NewRequest("POST", "/api/study_sessions/1/words/1/review", nil)
		c.Request.Header.Set("Content-Type", "application/json")

		requestBody := map[string]bool{"correct": true}
		jsonBody, _ := json.Marshal(requestBody)
		c.Request.Body = io.NopCloser(bytes.NewBuffer(jsonBody))

		handler.ReviewWord(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response models.WordReviewResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, 1, response.WordID)
		assert.Equal(t, 1, response.StudySessionID)
		assert.True(t, response.Correct)
	})

	t.Run("ReviewWordInvalidSession", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = gin.Params{
			{Key: "id", Value: "999"},
			{Key: "word_id", Value: "1"},
		}
		c.Request, _ = http.NewRequest("POST", "/api/study_sessions/999/words/1/review", nil)
		c.Request.Header.Set("Content-Type", "application/json")

		requestBody := map[string]bool{"correct": true}
		jsonBody, _ := json.Marshal(requestBody)
		c.Request.Body = io.NopCloser(bytes.NewBuffer(jsonBody))

		handler.ReviewWord(c)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestResetHandlers(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	seedTestData(t, db)

	gin.SetMode(gin.TestMode)
	handler := NewResetHandler(db)

	t.Run("ResetHistory", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/api/reset_history", nil)

		handler.ResetHistory(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response models.ResetResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response.Success)
		assert.Equal(t, "Study history reset successfully", response.Message)
	})

	t.Run("FullReset", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/api/full_reset", nil)

		handler.FullReset(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response models.ResetResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response.Success)
		assert.Equal(t, "System has been fully reset", response.Message)
	})
}
