package controllers

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "go-tunes/models"
    "go-tunes/database"
    "go-tunes/services"
    "log"
)

// GetSongs retrieves all songs with filtering and pagination
// @Summary Get all songs
// @Description Get all songs with filtering and pagination
// @Produce json
// @Success 200 {array} models.Song
// @Router /songs [get]
func GetSongs(c *gin.Context) {
    var songs []models.Song
    db := database.Connect()
    db.Find(&songs)
    log.Println("INFO: Retrieved songs")
    c.JSON(http.StatusOK, songs)
}

// GetSongText retrieves the text of a song with pagination by verses
// @Summary Get a song by ID
// @Description Get a song by ID
// @Produce json
// @Param id path int true "Song ID"
// @Success 200 {object} models.Song
// @Router /songs/{id} [get]
func GetSongText(c *gin.Context) {
    id := c.Param("id")
    var song models.Song
    db := database.Connect()
    if err := db.First(&song, id).Error; err != nil {
        log.Printf("ERROR: Song with ID %s not found", id)
        c.AbortWithStatus(http.StatusNotFound)
        return
    }
    log.Printf("INFO: Retrieved text for song ID %s", id)
    c.JSON(http.StatusOK, song.Text)
}

// AddSong adds a new song
// @Summary Add a new song
// @Description Add a new song
// @Accept json
// @Produce json
// @Param song body services.NewSongRequest true "Song data"
// @Success 200 {object} models.Song
// @Router /songs [post]
func AddSong(c *gin.Context) {
    var song models.Song
    if err := c.ShouldBindJSON(&song); err != nil {
        log.Printf("ERROR: Invalid song data: %v", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Make a request to the external API
    songDetail, err := services.GetSongDetail(song.Group, song.Name)
    if err != nil {
        log.Printf("ERROR: Failed to get song detail: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // Enrich the song data
    song.ReleaseDate = songDetail.ReleaseDate
    song.Text = songDetail.Text
    song.Link = songDetail.Link

    db := database.Connect()
    db.Create(&song)
    log.Printf("INFO: Added new song: %v", song)
    c.JSON(http.StatusOK, song)
}

// UpdateSong updates an existing song
// @Summary Update a song
// @Description Update an existing song
// @Accept json
// @Produce json
// @Param id path int true "Song ID"
// @Param song body models.Song true "Song data"
// @Success 200 {object} models.Song
// @Router /songs/{id} [put]
func UpdateSong(c *gin.Context) {
    id := c.Param("id")
    var song models.Song
    db := database.Connect()
    if err := db.First(&song, id).Error; err != nil {
        log.Printf("ERROR: Song with ID %s not found", id)
        c.AbortWithStatus(http.StatusNotFound)
        return
    }
    if err := c.ShouldBindJSON(&song); err != nil {
        log.Printf("ERROR: Invalid song data: %v", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    db.Save(&song)
    log.Printf("INFO: Updated song with ID %s", id)
    c.JSON(http.StatusOK, song)
}

// DeleteSong deletes a song by ID
// @Summary Delete a song
// @Description Delete a song by ID
// @Produce json
// @Param id path int true "Song ID"
// @Success 200 {object} map[string]interface{}
// @Router /songs/{id} [delete]
func DeleteSong(c *gin.Context) {
    id := c.Param("id")
    db := database.Connect()
    if err := db.Delete(&models.Song{}, id).Error; err != nil {
        log.Printf("ERROR: Failed to delete song with ID %s", id)
        c.AbortWithStatus(http.StatusNotFound)
        return
    }
    log.Printf("INFO: Deleted song with ID %s", id)
    c.JSON(http.StatusOK, map[string]interface{}{"id #" + id: "deleted"})
}
