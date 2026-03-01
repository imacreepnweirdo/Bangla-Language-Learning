package models

import "time"

type Word struct {
	ID        int       `json:"id" db:"id"`
	Bengali   string    `json:"bengali" db:"bengali"`
	Parts     string    `json:"parts" db:"parts"`
	English   string    `json:"english" db:"english"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type WordWithStats struct {
	ID            int    `json:"id"`
	Bengali       string `json:"bengali"`
	PartsOfSpeech string `json:"parts_of_speech"`
	English       string `json:"english"`
	CorrectCount  int    `json:"correct_count"`
	WrongCount    int    `json:"wrong_count"`
}

type WordDetail struct {
	ID      int       `json:"id"`
	Bengali string    `json:"bengali"`
	English string    `json:"english"`
	Parts   string    `json:"parts"`
	Stats   WordStats `json:"stats"`
	Groups  []Group   `json:"items"`
}

type WordStats struct {
	CorrectCount int `json:"correct_count"`
	WrongCount   int `json:"wrong_count"`
}
