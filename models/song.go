package models

import (
    "time"
    // "gorm.io/gorm"
)

// Song represents a song in the database
type Song struct {
    ID          uint           `gorm:"primaryKey" json:"id"`
    CreatedAt   time.Time      `json:"created_at"`
    UpdatedAt   time.Time      `json:"updated_at"`
    DeletedAt   time.Time `gorm:"index" json:"deleted_at,omitempty"`
    Group       string         `json:"group"`
    Name        string         `json:"name"`
    ReleaseDate string         `json:"release_date"`
    Text        string         `json:"text"`
    Link        string         `json:"link"`
}

// SongDetail represents detailed information about a song
type SongDetail struct {
    ReleaseDate string `json:"releaseDate"`
    Text        string `json:"text"`
    Link        string `json:"link"`
}
