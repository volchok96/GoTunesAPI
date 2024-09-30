package models

import (
    "time"
)

// Song представляет песню в базе данных
type Song struct {
    ID          uint      `gorm:"primaryKey" json:"id"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
    DeletedAt   *time.Time `gorm:"index" json:"deleted_at,omitempty"` 
    Group   string    `json:"group"`      
    Song    string    `json:"song"`       
    ReleaseDate string    `json:"release_date"`
    Text        string    `json:"text"`
    Link        string    `json:"link"`
}

// SongDetail представляет детальную информацию о песне
type SongDetail struct {
    ReleaseDate string `json:"release_date"`
    Text        string `json:"text"`
    Link        string `json:"link"`
}

// NewSongRequest используется при добавлении новой песни
type NewSongRequest struct {
    Group string `json:"group" binding:"required"` 
    Song  string `json:"song" binding:"required"`  
}
