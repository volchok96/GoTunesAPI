package controllers

import (
	"encoding/json"
	"fmt"
	"go-tunes/database"
	"go-tunes/models"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

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

	// Ищем песню в базе данных
	var songRecord models.Song
	if err := db.Where("\"group\" = ? AND song = ?", group, song).First(&songRecord).Error; err != nil {
		log.Printf("INFO: Song with group '%s' and song '%s' not found in database. Attempting to add it.", group, song)

		// Если песни нет в базе данных, попробуем добавить её, используя внешний API
		songDetail, shouldReturn := GetSongDetailFromAPI(group, song, c)
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
func GetSongDetailFromAPI(group, song string, c *gin.Context) (models.SongDetail, bool) {
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

// GetSongDetailFromJSON читает данные о песне из JSON-файла
func GetSongDetailFromJSON(group, song string) (models.SongDetail, error) {
	jsonFile, err := os.Open("song_enrichment.json")
	if err != nil {
		log.Printf("ERROR: Could not open JSON file: %v", err)
		return models.SongDetail{}, fmt.Errorf("could not open JSON file")
	}
	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)

	// Парсинг JSON в структуру
	var enrichmentData SongEnrichment

	if err := json.Unmarshal(byteValue, &enrichmentData); err != nil {
		log.Printf("ERROR: Could not parse JSON file: %v", err)
		return models.SongDetail{}, fmt.Errorf("could not parse JSON file")
	}

	if enrichmentData.Group == group && enrichmentData.Song == song {
		return models.SongDetail{
			ReleaseDate: enrichmentData.ReleaseDate,
			Text:        enrichmentData.Text,
			Link:        enrichmentData.Link,
		}, nil
	}

	return models.SongDetail{}, fmt.Errorf("song not found")
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
// @Param group query string false "Group"
// @Param song query string false "Song"
// @Param release_date query string false "Release Date"
// @Param text query string false "Text"
// @Param link query string false "Link"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Results per page" default(10)
// @Success 200 {array} models.Song
// @Failure 500 {string} string "internal server error"
// @Router /songs [get]
func GetSongs(c *gin.Context) {
	db := database.Connect()
	var songs []models.Song

	// Получение параметров фильтрации
	group := c.Query("group")
	song := c.Query("song")
	releaseDate := c.Query("release_date")
	text := c.Query("text")
	link := c.Query("link")

	// Получение параметров пагинации
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "10")

	// Преобразование значений параметров пагинации в int
	pageNumber, err := strconv.Atoi(page)
	if err != nil || pageNumber < 1 {
		pageNumber = 1
	}

	limitNumber, err := strconv.Atoi(limit)
	if err != nil || limitNumber < 1 {
		limitNumber = 10
	}

	// Построение запроса с учетом фильтров
	query := db.Model(&models.Song{})

	if group != "" {
		query = query.Where("\"group\" ILIKE ?", "%"+group+"%")
	}
	if song != "" {
		query = query.Where("song ILIKE ?", "%"+song+"%")
	}
	if releaseDate != "" {
		query = query.Where("release_date = ?", releaseDate)
	}
	if text != "" {
		query = query.Where("text ILIKE ?", "%"+text+"%")
	}
	if link != "" {
		query = query.Where("link ILIKE ?", "%"+link+"%")
	}

	// Пагинация
	offset := (pageNumber - 1) * limitNumber
	query = query.Offset(offset).Limit(limitNumber)

	// Выполнение запроса
	if err := query.Find(&songs).Error; err != nil {
		log.Printf("ERROR: Failed to retrieve songs: %v", err)
		c.String(http.StatusInternalServerError, "internal server error")
		return
	}

	// Возвращение результатов
	log.Println("INFO: Retrieved songs with filtering and pagination")
	c.JSON(http.StatusOK, songs)
}

// GetSongTextWithPagination retrieves the text of a song with pagination by verses
// @Summary Get a song by ID with pagination
// @Description Retrieve the text of a song by its ID with pagination by verses
// @Produce json
// @Param id path int true "Song ID"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Verses per page" default(1)
// @Success 200 {object} map[string]interface{}
// @Failure 404 {string} string "not found"
// @Failure 500 {string} string "internal server error"
// @Router /songs/{id}/verses [get]
func GetSongTextWithPagination(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Printf("ERROR: Invalid song ID %s", c.Param("id"))
		c.String(http.StatusBadRequest, "invalid song id")
		return
	}
	// Подключение к базе данных
	db := database.Connect()

	// Поиск песни по ID
	var song models.Song
	if err := db.Unscoped().First(&song, id).Error; err != nil {
		log.Printf("ERROR: Song with ID %d not found", id)
		c.String(http.StatusNotFound, "not found")
		return
	}

	// Получение параметров пагинации
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "1"))
	if err != nil || limit < 1 {
		limit = 1
	}

	// Разделение текста песни на куплеты (предполагается, что куплеты разделены "\n\n")
	verses := strings.Split(song.Text, "\n\n")

	// Подсчет общего количества куплетов
	totalVerses := len(verses)
	if totalVerses == 0 {
		log.Printf("ERROR: No verses found for song with ID %d", id)
		c.String(http.StatusNotFound, "not found")
		return
	}

	// Определение начального и конечного индекса для пагинации
	startIndex := (page - 1) * limit
	endIndex := startIndex + limit

	// Проверка, что стартовый индекс находится в пределах доступного диапазона
	if startIndex >= totalVerses {
		log.Printf("ERROR: Page %d out of range for song ID %d", page, id)
		c.String(http.StatusNotFound, "no verses found for the requested page")
		return
	}

	// Ограничение конечного индекса до общего количества куплетов
	if endIndex > totalVerses {
		endIndex = totalVerses
	}

	// Извлечение нужных куплетов
	selectedVerses := verses[startIndex:endIndex]

	// Формирование ответа
	response := map[string]interface{}{
		"song_id":   id,
		"page":      page,
		"limit":     limit,
		"total":     totalVerses,
		"verses":    selectedVerses,
		"total_pages": (totalVerses + limit - 1) / limit, // Подсчет общего количества страниц
	}

	// Логирование и отправка ответа
	log.Printf("INFO: Retrieved verses for song ID %d, page %d", id, page)
	c.JSON(http.StatusOK, response)
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
