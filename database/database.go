package database

import (
    "log"
    "os"
    "sync"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

var (
    db  *gorm.DB
    once sync.Once
)

// Connect устанавливает соединение с базой данных и возвращает *gorm.DB.
func Connect() *gorm.DB {
    once.Do(func() {
        dsn := os.Getenv("DATABASE_URL")
        var err error
        db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
        if err != nil {
            log.Fatal("Failed to connect to the database:", err)
        }

        sqlDB, err := db.DB()
        if err != nil {
            log.Fatal("Failed to get sql.DB from gorm.DB:", err)
        }

        // Настройка пула подключений
        sqlDB.SetMaxOpenConns(10)  // Максимальное количество открытых подключений
        sqlDB.SetMaxIdleConns(5)   // Максимальное количество "спящих" подключений
        sqlDB.SetConnMaxLifetime(0) // Время жизни соединения (0 - не ограничивать)
    })

    return db
}
