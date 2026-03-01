package models

import "time"

type LastStudySession struct {
	ID                int       `json:"id"`
	GroupID           int       `json:"group_id"`
	CreatedAt         time.Time `json:"created_at"`
	StudyActivitiesID int       `json:"study_activities_id"`
	GroupName         string    `json:"group_name"`
}

type StudyProgress struct {
	TotalWordsStudied   int `json:"total_words_studied"`
	TotalAvailableWords int `json:"total_available_words"`
}

type QuickStats struct {
	TotalWords        int     `json:"total_words"`
	TotalGroups       int     `json:"total_groups"`
	WordsLearned      int     `json:"words_learned"`
	SessionsCompleted int     `json:"sessions_completed"`
	CurrentStreak     int     `json:"current_streak"`
	AccuracyRate      float64 `json:"accuracy_rate"`
}

type Pagination struct {
	Page         int `json:"page"`
	TotalPages   int `json:"total_pages"`
	TotalItems   int `json:"total_items"`
	ItemsPerPage int `json:"items_per_page"`
}

type PaginatedWords struct {
	Items      []WordWithStats `json:"items"`
	Pagination Pagination      `json:"pagination"`
}

type PaginatedGroups struct {
	Items      []GroupWithWordCount `json:"items"`
	Pagination Pagination           `json:"pagination"`
}

type PaginatedStudySessions struct {
	Items      []StudySessionWithDetails `json:"items"`
	Pagination Pagination                `json:"pagination"`
}

type PaginatedGroupWords struct {
	Words      []WordStat `json:"words"`
	Pagination Pagination `json:"pagination"`
}

type ResetResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
