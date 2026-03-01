package models

import "time"

type StudySession struct {
	ID                int       `json:"id" db:"id"`
	GroupID           int       `json:"group_id" db:"group_id"`
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
	StudyActivitiesID *int      `json:"study_activities_id" db:"study_activities_id"`
}

type StudySessionWithDetails struct {
	ID               int       `json:"id"`
	ActivityName     string    `json:"activity_name"`
	GroupName        string    `json:"group_name"`
	StartTime        time.Time `json:"start_time"`
	EndTime          time.Time `json:"end_time"`
	ReviewItemsCount int       `json:"review_items_count"`
}

type StudyActivity struct {
	ID             int       `json:"id" db:"id"`
	StudySessionID int       `json:"study_session_id" db:"study_session_id"`
	GroupID        int       `json:"group_id" db:"group_id"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}

type StudyActivityDetail struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	ThumbnailURL string `json:"thumbnail_url"`
	Description  string `json:"description"`
}

type WordReviewItem struct {
	ID             int       `json:"id" db:"id"`
	WordID         int       `json:"word_id" db:"word_id"`
	StudySessionID int       `json:"study_session_id" db:"study_session_id"`
	Correct        bool      `json:"correct" db:"correct"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}

type WordReviewResponse struct {
	ID             int       `json:"id"`
	Success        bool      `json:"success"`
	WordID         int       `json:"word_id"`
	StudySessionID int       `json:"study_session_id"`
	Correct        bool      `json:"correct"`
	CreatedAt      time.Time `json:"created_at"`
}

type CreateStudySessionRequest struct {
	GroupID         int `json:"group_id" binding:"required"`
	StudyActivityID int `json:"study_activity_id" binding:"required"`
}

type CreateStudySessionResponse struct {
	ID      int `json:"id"`
	GroupID int `json:"group_id"`
}
