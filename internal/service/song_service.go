package service

import (
	"fmt"
	"log"
	"music-library/internal/models"
	"music-library/internal/repository"
	"strings"
	"time"
)

// SongService представляет сервисный слой для работы с песнями
type SongService struct {
	repo        *repository.SongRepository
	externalAPI *ExternalAPIService
}

// NewSongService создает новый экземпляр сервиса песен
func NewSongService(repo *repository.SongRepository, externalAPI *ExternalAPIService) *SongService {
	return &SongService{
		repo:        repo,
		externalAPI: externalAPI,
	}
}

// GetSongs возвращает список песен с применением фильтрации и пагинации
func (s *SongService) GetSongs(group string, page, limit int) ([]models.Song, error) {
	// Валидация входных параметров
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10 // значение по умолчанию
	}

	log.Printf("Fetching songs for group: %s, page: %d, limit: %d", group, page, limit)

	songs, err := s.repo.GetSongs(group, page, limit)
	if err != nil {
		log.Printf("Error in GetSongs: %v", err)
		return nil, fmt.Errorf("failed to retrieve songs: %w", err)
	}

	return songs, nil
}

// CreateSong создает новую песню с обогащением данными из внешнего API
func (s *SongService) CreateSong(req models.CreateSongRequest) (*models.Song, error) {
	// Нормализация входных данных
	group := strings.TrimSpace(req.Group)
	songName := strings.TrimSpace(req.Song)

	if group == "" || songName == "" {
		return nil, fmt.Errorf("group and song name cannot be empty")
	}

	// Получение дополнительной информации из внешнего API
	songDetails, err := s.externalAPI.GetSongDetails(group, songName)
	if err != nil {
		log.Printf("Error fetching song details from external API: %v", err)
		return nil, fmt.Errorf("failed to fetch song details: %w", err)
	}

	// Обработка даты релиза
	if songDetails.ReleaseDate == "" {
		songDetails.ReleaseDate = time.Now().Format("02.01.2006")
	}

	// Создание объекта песни для сохранения в базе данных
	song := &models.Song{
		Group:       group,
		SongName:    songName,
		ReleaseDate: songDetails.ReleaseDate,
		Text:        songDetails.Text,
		Link:        songDetails.Link,
	}

	// Сохранение песни в репозитории
	createdSong, err := s.repo.CreateSong(song)
	if err != nil {
		log.Printf("Error creating song in repository: %v", err)
		return nil, fmt.Errorf("failed to create song: %w", err)
	}

	log.Printf("Created song: %s by %s", createdSong.SongName, createdSong.Group)
	return createdSong, nil
}

// UpdateSong обновляет информацию о песне
func (s *SongService) UpdateSong(songID int, updateData models.Song) error {
	// Валидация входных данных
	if songID <= 0 {
		return fmt.Errorf("invalid song ID")
	}

	// Нормализация данных
	updateData.ID = songID
	updateData.Group = strings.TrimSpace(updateData.Group)
	updateData.SongName = strings.TrimSpace(updateData.SongName)

	// Проверка обязательных полей
	if updateData.Group == "" || updateData.SongName == "" {
		return fmt.Errorf("group and song name cannot be empty")
	}

	// Обновление даты релиза, если не указана
	if updateData.ReleaseDate == "" {
		updateData.ReleaseDate = time.Now().Format("02.01.2006")
	}

	// Вызов репозитория для обновления
	err := s.repo.UpdateSong(&updateData)
	if err != nil {
		log.Printf("Error updating song: %v", err)
		return fmt.Errorf("failed to update song: %w", err)
	}

	log.Printf("Updated song ID: %d", songID)
	return nil
}

// DeleteSong удаляет песню по идентификатору
func (s *SongService) DeleteSong(songID int) error {
	// Валидация входных данных
	if songID <= 0 {
		return fmt.Errorf("invalid song ID")
	}

	// Вызов репозитория для удаления
	err := s.repo.DeleteSong(songID)
	if err != nil {
		log.Printf("Error deleting song: %v", err)
		return fmt.Errorf("failed to delete song: %w", err)
	}

	log.Printf("Deleted song ID: %d", songID)
	return nil
}

// GetSongText возвращает текст песни постранично
func (s *SongService) GetSongText(songID, page, limit int) (string, error) {
	// Валидация входных параметров
	if songID <= 0 {
		return "", fmt.Errorf("invalid song ID")
	}

	if page < 1 {
		page = 1
	}

	if limit < 1 || limit > 10 {
		limit = 3 // значение по умолчанию
	}

	// Получение текста песни постранично
	songText, err := s.repo.GetSongText(songID, page, limit)
	if err != nil {
		log.Printf("Error getting song text: %v", err)
		return "", fmt.Errorf("failed to retrieve song text: %w", err)
	}

	return songText, nil
}
