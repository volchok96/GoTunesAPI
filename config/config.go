package config

import (
    "github.com/joho/godotenv"
    "log"
)

func LoadEnv() {
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Ошибка при загрузке .env файла")
    }
}
