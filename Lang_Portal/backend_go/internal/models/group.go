package models

import "time"

type Group struct {
	ID        int       `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type GroupWithWordCount struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	WordsCount int    `json:"words_count"`
}

type GroupDetail struct {
	ID            int        `json:"id"`
	Name          string     `json:"name"`
	Words         []WordStat `json:"words"`
	Stats         GroupStats `json:"stats"`
}

type GroupStats struct {
	TotalWordCount int `json:"total_word_count"`
}

type WordStat struct {
	ID           int    `json:"id"`
	Bengali      string `json:"bengali"`
	English      string `json:"english"`
	CorrectCount int    `json:"correct_count"`
	WrongCount   int    `json:"wrong_count"`
}

type WordsGroup struct {
	ID        int       `json:"id" db:"id"`
	WordID    int       `json:"word_id" db:"word_id"`
	GroupID   int       `json:"group_id" db:"group_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
