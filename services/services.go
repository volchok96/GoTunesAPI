package services

import (
    "go-tunes/repository"
    "go-tunes/models"
    "go-tunes/database"
    "net/http"
    "encoding/json"
    "fmt"
)

type NewSongRequest struct {
    Group string `json:"group" binding:"required"`
    Song  string `json:"song" binding:"required"`
}

func CreateSong(newSong NewSongRequest) (*models.Song, error) {
    songDetails, err := GetSongDetail(newSong.Group, newSong.Song)
    if err != nil {
        return nil, err
    }

    song := models.Song{
        Group:       newSong.Group,
        Name:        newSong.Song,
        ReleaseDate: songDetails.ReleaseDate,
        Text:        songDetails.Text,
        Link:        songDetails.Link,
    }

    repo := repository.NewSongRepository(database.Connect())

    return repo.SaveSong(&song)
}

func GetSongDetail(group string, song string) (*models.SongDetail, error) {
    url := fmt.Sprintf("https://external.api/info?group=%s&song=%s", group, song)

    resp, err := http.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("failed to fetch song details: %s", resp.Status)
    }

    var songDetail models.SongDetail
    if err := json.NewDecoder(resp.Body).Decode(&songDetail); err != nil {
        return nil, err
    }

    return &songDetail, nil
}
