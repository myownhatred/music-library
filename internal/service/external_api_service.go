package service

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"music-library/internal/config"
	"music-library/internal/models"
	"net/http"
	"time"
)

type ExternalAPIService struct {
	baseURL string
	timeout time.Duration
}

func NewExternalAPIService(cfg *config.Config) *ExternalAPIService {
	return &ExternalAPIService{
		baseURL: cfg.APIBaseURL,
		timeout: time.Duration(cfg.APITimeout) * time.Second,
	}
}

func (s *ExternalAPIService) GetSongDetails(group, song string) (*models.Song, error) {
	client := &http.Client{Timeout: s.timeout}

	url := fmt.Sprintf("%s/info?group=%s&song=%s", s.baseURL, group, song)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error making request: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("external API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return nil, err
	}

	var songDetails models.Song
	if err := json.Unmarshal(body, &songDetails); err != nil {
		log.Printf("Error parsing response: %v", err)
		return nil, err
	}

	return &songDetails, nil
}
