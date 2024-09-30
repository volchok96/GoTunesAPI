package services

import (
    "fmt"
    "net/http"
    "encoding/json"
)

type SongDetail struct {
    ReleaseDate string `json:"releaseDate"`
    Text        string `json:"text"`
    Link        string `json:"link"`
}

// GetSongDetail fetches song details from the external API
func GetSongDetail(group string, song string) (SongDetail, error) {
    // Формирование URL запроса
    url := fmt.Sprintf("http://localhost:8080/info?group=%s&song=%s", group, song)

    // Выполнение GET-запроса
    resp, err := http.Get(url)
    if err != nil {
        return SongDetail{}, fmt.Errorf("failed to fetch song details: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return SongDetail{}, fmt.Errorf("failed to fetch song details: %s", resp.Status)
    }

    // Декодирование JSON ответа в структуру SongDetail
    var songDetail SongDetail
    if err := json.NewDecoder(resp.Body).Decode(&songDetail); err != nil {
        return SongDetail{}, fmt.Errorf("failed to decode response: %w", err)
    }

    return songDetail, nil
}
