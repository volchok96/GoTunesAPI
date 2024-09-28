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

    router := gin.Default()

    router.GET("/songs", controllers.GetSongs)
    router.GET("/songs/:id", controllers.GetSongs)
    router.POST("/songs", controllers.AddSong)
    router.PUT("/songs/:id", controllers.UpdateSong)
    router.DELETE("/songs/:id", controllers.DeleteSong)

    // Swagger
    router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

    log.Fatal(router.Run(":8080"))
}
