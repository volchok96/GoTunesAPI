package models

import (
    "gorm.io/gorm"
)

type Song struct {
    gorm.Model
    Group       string `json:"group"`
    Name        string `json:"name"`
    ReleaseDate string `json:"release_date"`
    Text        string `json:"text"`
    Link        string `json:"link"`
}

type SongDetail struct {
    ReleaseDate string `json:"releaseDate"`
    Text        string `json:"text"`
    Link        string `json:"link"`
}
