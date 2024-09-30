package main

import (
    "log"
    "github.com/gin-gonic/gin"
    "go-tunes/config"
    "go-tunes/controllers"
    "go-tunes/database"
    _ "go-tunes/docs"
    "github.com/swaggo/gin-swagger"
    "github.com/gin-contrib/cors"
    "github.com/swaggo/files"
    "net/http"
)

// @title Music Library API
// @version 1.0
// @description API для управления библиотекой песен.
// @host localhost:8080
// @BasePath /

func main() {
    config.LoadEnv()
    db := database.Connect()
    database.Migrate(db)

    // Основной сервер на порту 8080
    router := gin.Default()
    router.Use(cors.New(cors.Config{
        AllowAllOrigins: true, // Разрешает все источники (для тестирования, не рекомендуется в продакшене)
        AllowMethods:    []string{"GET", "POST", "PUT", "DELETE"},
        AllowHeaders:    []string{"Origin", "Content-Type", "Authorization"},
    }))

    router.GET("/info", controllers.GetSongInfo)
    router.GET("/songs", controllers.GetSongs)
    router.GET("/songs/:id", controllers.GetSongText)
    router.POST("/songs", controllers.AddSong)
    router.PUT("/songs/:id", controllers.UpdateSong)
    router.DELETE("/songs/:id", controllers.DeleteSong)

    // Swagger
    router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

    // Второй сервер на порту 8081
    go func() {
        testRouter := gin.Default()
        testRouter.GET("/info", func(c *gin.Context) {
            group := c.Query("group")
            song := c.Query("song")

            // Возвращаем фиктивные данные для тестирования
            if group != "" && song != "" {
                c.JSON(http.StatusOK, gin.H{
                    "group":        group,
                    "song":         song,
                    "release_date": "2024-01-01",
                    "text":         "This is a sample song text.",
                    "link":         "http://example.com/song",
                })
            } else {
                c.JSON(http.StatusBadRequest, gin.H{
                    "error": "missing parameters",
                })
            }
        })

        // Запускаем сервер на порту 8081
        if err := testRouter.Run(":8081"); err != nil {
            log.Fatalf("Ошибка при запуске тестового сервера: %v", err)
        }
    }()

    // Запускаем основной сервер на порту 8080
    log.Fatal(router.Run(":8080"))
}
