package main

import (
	"database/sql"
	"log"
	"net/http"

	"lang-portal/internal/handlers"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Initialize database connection
	db, err := sql.Open("sqlite3", "words.db")
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}
	defer db.Close()

	// Test database connection
	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	// Initialize Gin router
	router := gin.Default()

	// Enable CORS
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		
		c.Next()
	})

	// Initialize handlers
	dashboardHandler := handlers.NewDashboardHandler(db)
	studyActivityHandler := handlers.NewStudyActivityHandler(db)
	wordHandler := handlers.NewWordHandler(db)
	groupHandler := handlers.NewGroupHandler(db)
	studySessionHandler := handlers.NewStudySessionHandler(db)
	resetHandler := handlers.NewResetHandler(db)

	// API Routes
	api := router.Group("/api")
	{
		// Dashboard routes
		api.GET("/dashboard/last_study_session", dashboardHandler.GetLastStudySession)
		api.GET("/dashboard/study_progress", dashboardHandler.GetStudyProgress)
		api.GET("/dashboard/quick-stats", dashboardHandler.GetQuickStats)

		// Study activities routes
		api.GET("/study_activities/:id", studyActivityHandler.GetStudyActivity)
		api.GET("/study_activities/:id/study_sessions", studyActivityHandler.GetStudySessions)
		api.POST("/study_activities", studyActivityHandler.CreateStudyActivity)

		// Words routes
		api.GET("/words", wordHandler.GetWords)
		api.GET("/words/:id", wordHandler.GetWord)

		// Groups routes
		api.GET("/groups", groupHandler.GetGroups)
		api.GET("/groups/:id", groupHandler.GetGroup)
		api.GET("/groups/:id/words", groupHandler.GetGroupWords)

		// Study sessions routes
		api.POST("/study_sessions/:id/words/:word_id/review", studySessionHandler.ReviewWord)

		// Reset routes
		api.POST("/reset_history", resetHandler.ResetHistory)
		api.POST("/full_reset", resetHandler.FullReset)
	}

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
			"database": "connected",
		})
	})

	// Start server
	log.Println("Server starting on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
