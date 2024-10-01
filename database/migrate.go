package database

import (
    "log"
    "gorm.io/gorm"
    "go-tunes/models"
)

func Migrate(db *gorm.DB) {
    err := db.AutoMigrate(&models.Song{})
    if err != nil {
        log.Fatal("Migration failed: ", err)
    }
}
