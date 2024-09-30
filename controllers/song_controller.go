package controllers

import (
	"encoding/json"
	"fmt"
	"go-tunes/database"
	"go-tunes/models"
	"go-tunes/services"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
)

type SongEnrichment struct {
	Group       string `json:"group"`
	Song        string `json:"song"`
	ReleaseDate string `json:"release_date"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}

// GetSongInfo обрабатывает запросы для получения информации о песне и добавляет её в базу данных при отсутствии
// @Summary Get song details
// @Description Retrieve detailed information about a song, add to database if not present
// @Produce json
// @Param group query string true "Group"
// @Param song query string true "Song"
// @Success 200 {object} models.SongDetail
// @Failure 400 {string} string "bad request"
// @Failure 500 {string} string "internal server error"
// @Router /info [get]
func GetSongInfo(c *gin.Context) {
	group := c.Query("group")
	song := c.Query("song")

	// Проверка, что параметры не пусты
	if group == "" || song == "" {
		log.Printf("ERROR: Bad request, missing 'group' or 'song' query parameters")
		c.String(http.StatusBadRequest, "bad request: missing required parameters")
		return
	}

	// Подключаемся к базе данных
	db := database.Connect()
	_, err := db.DB() // Получение объекта *sql.DB для дальнейшего закрытия
	if err != nil {
		log.Printf("ERROR: Failed to get sql.DB from gorm.DB: %v", err)
		c.String(http.StatusInternalServerError, "internal server error")
		return
	}
	//defer sqlDB.Close() // Закрытие соединения после выполнения

	// Ищем песню в базе данных
	var songRecord models.Song
	if err := db.Where("\"group\" = ? AND song = ?", group, song).First(&songRecord).Error; err != nil {
		log.Printf("INFO: Song with group '%s' and song '%s' not found in database. Attempting to add it.", group, song)

		// Если песни нет в базе данных, попробуем добавить её, используя внешний API
		songDetail, shouldReturn := getSongDetailFromAPI(group, song, c)
		if shouldReturn {
			return
		}

		// Добавляем песню в базу данных
		newSong := models.Song{
			Group:       group,
			Song:        song,
			ReleaseDate: songDetail.ReleaseDate,
			Text:        songDetail.Text,
			Link:        songDetail.Link,
		}

		if err := db.Create(&newSong).Error; err != nil {
			log.Printf("ERROR: Failed to add new song to the database: %v", err)
			c.String(http.StatusInternalServerError, "internal server error")
			return
		}

		log.Printf("INFO: Added new song to the database: %v", newSong)
		songRecord = newSong
	}

	// Формируем объект ответа
	songDetail := models.SongDetail{
		ReleaseDate: songRecord.ReleaseDate,
		Text:        songRecord.Text,
		Link:        songRecord.Link,
	}

	// Дополнительное обогащение данных с использованием JSON-файла
	enrichSongFromJSON(&songDetail, group, song)

	// Возвращаем результат
	c.JSON(http.StatusOK, songDetail)
}

// getSongDetailFromAPI выполняет запрос к внешнему API для получения данных о песне
func getSongDetailFromAPI(group, song string, c *gin.Context) (models.SongDetail, bool) {
	encodedGroup := url.QueryEscape(group)
	encodedSong := url.QueryEscape(song)
	apiURL := fmt.Sprintf("http://localhost:8081/info?group=%s&song=%s", encodedGroup, encodedSong)
	response, err := http.Get(apiURL)
	if err != nil {
		log.Printf("ERROR: Failed to request external API: %v", err)
		c.String(http.StatusInternalServerError, "internal server error")
		return models.SongDetail{}, true
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		log.Printf("WARNING: External API returned status code %d", response.StatusCode)
		c.String(http.StatusInternalServerError, "failed to retrieve song details from external API")
		return models.SongDetail{}, true
	}

	var apiData models.SongDetail
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Printf("ERROR: Failed to read API response: %v", err)
		c.String(http.StatusInternalServerError, "internal server error")
		return models.SongDetail{}, true
	}

	if err := json.Unmarshal(body, &apiData); err != nil {
		log.Printf("ERROR: Failed to parse API response: %v", err)
		c.String(http.StatusInternalServerError, "internal server error")
		return models.SongDetail{}, true
	}

	return apiData, false
}

// enrichSongFromJSON обогащает данные песни из локального JSON-файла
func enrichSongFromJSON(songDetail *models.SongDetail, group, song string) {
	jsonFile, err := os.Open("song_enrichment.json")
	if err != nil {
		log.Printf("ERROR: Could not open JSON file: %v", err)
		return
	}
	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)

	// Парсинг JSON в структуру
	var enrichmentData SongEnrichment
	if err := json.Unmarshal(byteValue, &enrichmentData); err != nil {
		log.Printf("ERROR: Could not parse JSON file: %v", err)
	} else {
		// Здесь работа с объектом enrichmentData
		if enrichmentData.Group == group && enrichmentData.Song == song {
			songDetail.ReleaseDate = enrichmentData.ReleaseDate
			songDetail.Text = enrichmentData.Text
			songDetail.Link = enrichmentData.Link
		}
	}
}

// GetSongs retrieves all songs with filtering and pagination
// @Summary Get all songs
// @Description Retrieve all songs with optional filtering and pagination
// @Produce json
// @Success 200 {array} models.Song
// @Failure 500 {string} string "internal server error"
// @Router /songs [get]
func GetSongs(c *gin.Context) {
	var songs []models.Song
	db := database.Connect()
	if err := db.Find(&songs).Error; err != nil {
		log.Printf("ERROR: Failed to retrieve songs: %v", err)
		c.String(http.StatusInternalServerError, "internal server error")
		return
	}
	log.Println("INFO: Retrieved songs")
	c.JSON(http.StatusOK, songs)
}

// GetSongText retrieves the text of a song with pagination by verses
// @Summary Get a song by ID
// @Description Retrieve the text of a song by its ID
// @Produce json
// @Param id path int true "Song ID"
// @Success 200 {string} string "Text of the song"
// @Failure 404 {string} string "not found"
// @Failure 500 {string} string "internal server error"
// @Router /songs/{id} [get]
func GetSongText(c *gin.Context) {
	id := c.Param("id")
	var song models.Song
	db := database.Connect()
	if err := db.First(&song, id).Error; err != nil {
		log.Printf("ERROR: Song with ID %s not found", id)
		c.String(http.StatusNotFound, "not found")
		return
	}
	log.Printf("INFO: Retrieved text for song ID %s", id)
	c.JSON(http.StatusOK, song.Text)
}

// AddSong handles adding a new song
// @Summary Add a new song
// @Description Add a new song to the library
// @Accept json
// @Produce json
// @Param song body models.NewSongRequest true "Song data"
// @Success 200 {object} models.Song
// @Failure 400 {object} string "invalid input"
// @Failure 500 {object} string "internal server error"
// @Router /songs [post]
func AddSong(c *gin.Context) {
	// 1. Parse the incoming JSON to get the group and song names
	var newSongRequest models.NewSongRequest
	if err := c.ShouldBindJSON(&newSongRequest); err != nil {
		log.Printf("ERROR: Invalid song data: %v", err)
		c.String(http.StatusBadRequest, "invalid input")
		return
	}

	// 2. Fetch song details from the external API
	songDetail, shouldReturn := newFunction(newSongRequest, c)
	if shouldReturn {
		return
	}

	// 3. Enrich the song data
	song := models.Song{
		Group:       newSongRequest.Group,
		Song:        newSongRequest.Song,
		ReleaseDate: songDetail.ReleaseDate,
		Text:        songDetail.Text,
		Link:        songDetail.Link,
	}

	// 4. Save enriched song data to the database
	db := database.Connect()
	if err := db.Create(&song).Error; err != nil {
		log.Printf("ERROR: Failed to add new song: %v", err)
		c.String(http.StatusInternalServerError, "internal server error")
		return
	}
	log.Printf("INFO: Added new song: %v", song)
	c.JSON(http.StatusOK, song)
}

func newFunction(newSongRequest models.NewSongRequest, c *gin.Context) (services.SongDetail, bool) {
	songDetail, err := services.GetSongDetail(newSongRequest.Group, newSongRequest.Song)
	if err != nil {
		log.Printf("ERROR: Failed to get song detail for group '%s' and song '%s': %v", newSongRequest.Group, newSongRequest.Song, err)
		c.String(http.StatusInternalServerError, "internal server error")
		return services.SongDetail{}, true
	}
	return songDetail, false
}

// UpdateSong updates an existing song
// @Summary Update a song
// @Description Update an existing song by its ID
// @Accept json
// @Produce json
// @Param id path int true "Song ID"
// @Param song body models.Song true "Updated song data"
// @Success 200 {object} models.Song
// @Failure 404 {string} string "not found"
// @Failure 400 {string} string "invalid input"
// @Router /songs/{id} [put]
func UpdateSong(c *gin.Context) {
	id := c.Param("id")
	var song models.Song
	db := database.Connect()
	if err := db.First(&song, id).Error; err != nil {
		log.Printf("ERROR: Song with ID %s not found", id)
		c.String(http.StatusNotFound, "not found")
		return
	}
	if err := c.ShouldBindJSON(&song); err != nil {
		log.Printf("ERROR: Invalid song data: %v", err)
		c.String(http.StatusBadRequest, "invalid input")
		return
	}
	if err := db.Save(&song).Error; err != nil {
		log.Printf("ERROR: Failed to update song with ID %s: %v", id, err)
		c.String(http.StatusInternalServerError, "internal server error")
		return
	}
	log.Printf("INFO: Updated song with ID %s", id)
	c.JSON(http.StatusOK, song)
}

// DeleteSong deletes a song by ID
// @Summary Delete a song
// @Description Delete a song by its ID
// @Produce json
// @Param id path int true "Song ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {string} string "not found"
// @Failure 500 {string} string "internal server error"
// @Router /songs/{id} [delete]
func DeleteSong(c *gin.Context) {
	id := c.Param("id")
	db := database.Connect()
	if err := db.Delete(&models.Song{}, id).Error; err != nil {
		log.Printf("ERROR: Failed to delete song with ID %s: %v", id, err)
		c.String(http.StatusInternalServerError, "internal server error")
		return
	}
	log.Printf("INFO: Deleted song with ID %s", id)
	c.JSON(http.StatusOK, map[string]interface{}{"id #" + id: "deleted"})
}
