package main

import (
    "log"
    "github.com/gin-gonic/gin"
    "go-tunes/config"
    "go-tunes/controllers"
    "go-tunes/database"
    _ "go-tunes/docs"
    "github.com/swaggo/gin-swagger"
    "github.com/swaggo/files"
    "net/http"
)

// @title Music Library API
// @version 1.0
// @description API для управления библиотекой песен.
// @host localhost:8080
// @BasePath /

func main() {
    // Загрузка переменных окружения
    config.LoadEnv()

    // Подключение к базе данных и выполнение миграций
    db := database.Connect()
    database.Migrate(db)

    // Основной сервер на порту 8080
    router := gin.Default()

    // Определение маршрутов для основного API
    router.GET("/info", controllers.GetSongInfo)       // Информация о песне
    router.GET("/songs", controllers.GetSongs)         // Список песен
    router.GET("/songs/:id/verses", controllers.GetSongTextWithPagination)  // Текст песни по ID
    router.PUT("/songs/:id", controllers.UpdateSong)   // Обновление песни по ID
    router.DELETE("/songs/:id", controllers.DeleteSong) // Удаление песни по ID

    // Swagger для документации
    router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

    // Второй сервер на порту 8081 - Эмуляция внешнего API
    go startMockServer()

    // Запускаем основной сервер на порту 8080
    log.Fatal(router.Run(":8080"))
}

// startMockServer запускает тестовый сервер на порту 8081 для эмуляции внешнего API
func startMockServer() {
    testRouter := gin.Default()

    // Определяем маршрут для эмуляции внешнего API
    testRouter.GET("/info", func(c *gin.Context) {
        group := c.Query("group")
        song := c.Query("song")

        // Проверка параметров запроса
        if group == "" || song == "" {
            c.JSON(http.StatusBadRequest, gin.H{"error": "missing parameters"})
            return
        }

        // Используем отдельную функцию для получения данных из JSON без необходимости использования *gin.Context
        songDetail, err := controllers.GetSongDetailFromJSON(group, song)
        if err != nil {
            c.JSON(http.StatusNotFound, gin.H{"error": "song not found"})
            return
        }

        c.JSON(http.StatusOK, songDetail)
    })

    // Запускаем сервер на порту 8081
    if err := testRouter.Run(":8081"); err != nil {
        log.Fatalf("Ошибка при запуске тестового сервера: %v", err)
    }
}
