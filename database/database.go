package database

import (
    "log"
    "os"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

func Connect() *gorm.DB {
    dsn := os.Getenv("DATABASE_URL")
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal("Failed to connect to the database:", err)
    }
    return db
}