package repository

import (
    "log"
    "go-tunes/models"
    "gorm.io/gorm"
)

type SongRepository struct {
    DB *gorm.DB
}

func NewSongRepository(db *gorm.DB) *SongRepository {
    log.Println("INFO: Creating new SongRepository.")
    return &SongRepository{DB: db}
}

// SaveSong saves a song to the database
func (repo *SongRepository) SaveSong(song *models.Song) (*models.Song, error) {
    if err := repo.DB.Create(song).Error; err != nil {
        return nil, err
    }
    log.Printf("INFO: Successfully saved song with ID: %d\n", song.ID)
    return song, nil
}

// GetAllSongs retrieves all songs with pagination
func (repo *SongRepository) GetAllSongs(page int, limit int) ([]models.Song, error) {
    log.Printf("INFO: Retrieving all songs. Page: %d, Limit: %d\n", page, limit)
    var songs []models.Song
    offset := (page - 1) * limit

    if err := repo.DB.Limit(limit).Offset(offset).Find(&songs).Error; err != nil {
        log.Printf("ERROR: Failed to retrieve songs. Page: %d, Limit: %d, error: %v\n", page, limit, err)
        return nil, err
    }
    log.Printf("INFO: Successfully retrieved %d songs.\n", len(songs))
    return songs, nil
}

// GetSongByID retrieves a song by its ID
func (repo *SongRepository) GetSongByID(id uint) (*models.Song, error) {
    log.Printf("INFO: Retrieving song with ID: %d\n", id)
    var song models.Song
    if err := repo.DB.First(&song, id).Error; err != nil {
        log.Printf("ERROR: Failed to retrieve song with ID: %d, error: %v\n", id, err)
        return nil, err
    }
    log.Printf("INFO: Successfully retrieved song with ID: %d\n", id)
    return &song, nil
}

// UpdateSong updates an existing song
func (repo *SongRepository) UpdateSong(song *models.Song) (*models.Song, error) {
    log.Printf("INFO: Updating song with ID: %d\n", song.ID)
    if err := repo.DB.Save(song).Error; err != nil {
        log.Printf("ERROR: Failed to update song with ID: %d, error: %v\n", song.ID, err)
        return nil, err
    }
    log.Printf("INFO: Successfully updated song with ID: %d\n", song.ID)
    return song, nil
}

// DeleteSong deletes a song by its ID
func (repo *SongRepository) DeleteSong(id uint) error {
    log.Printf("INFO: Deleting song with ID: %d\n", id)
    if err := repo.DB.Delete(&models.Song{}, id).Error; err != nil {
        log.Printf("ERROR: Failed to delete song with ID: %d, error: %v\n", id, err)
        return err
    }
    log.Printf("INFO: Successfully deleted song with ID: %d\n", id)
    return nil
}
