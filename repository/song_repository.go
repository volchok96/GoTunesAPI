package repository

import (
    "go-tunes/models"
    "gorm.io/gorm"
)

type SongRepository struct {
    DB *gorm.DB
}

func NewSongRepository(db *gorm.DB) *SongRepository {
    return &SongRepository{DB: db}
}

// Сохранение песни в базу данных
func (repo *SongRepository) SaveSong(song *models.Song) (*models.Song, error) {
    if err := repo.DB.Create(song).Error; err != nil {
        return nil, err
    }
    return song, nil
}

// Получение всех песен с пагинацией
func (repo *SongRepository) GetAllSongs(page int, limit int) ([]models.Song, error) {
    var songs []models.Song
    offset := (page - 1) * limit

    if err := repo.DB.Limit(limit).Offset(offset).Find(&songs).Error; err != nil {
        return nil, err
    }
    return songs, nil
}

// Получение песни по ID
func (repo *SongRepository) GetSongByID(id uint) (*models.Song, error) {
    var song models.Song
    if err := repo.DB.First(&song, id).Error; err != nil {
        return nil, err
    }
    return &song, nil
}

// Обновление песни
func (repo *SongRepository) UpdateSong(song *models.Song) (*models.Song, error) {
    if err := repo.DB.Save(song).Error; err != nil {
        return nil, err
    }
    return song, nil
}

// Удаление песни
func (repo *SongRepository) DeleteSong(id uint) error {
    if err := repo.DB.Delete(&models.Song{}, id).Error; err != nil {
        return err
    }
    return nil
}
