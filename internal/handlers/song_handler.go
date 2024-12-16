package handlers

import (
	"log"
	"net/http"
	"strconv"

	"music-library/internal/models"
	"music-library/internal/service"

	"github.com/gin-gonic/gin"
)

// SongHandler структура для обработки HTTP-запросов, связанных с песнями
type SongHandler struct {
	songService *service.SongService
}

// NewSongHandler создает новый экземпляр обработчика песен
func NewSongHandler(songService *service.SongService) *SongHandler {
	return &SongHandler{songService: songService}
}

// @title Music Library API
// @version 1.0
// @description API для управления библиотекой песен

// GetSongsHandler godoc
// @Summary Получение списка песен
// @Description Возвращает список песен с фильтрацией и пагинацией
// @Tags songs
// @Accept json
// @Produce json
// @Param group query string false "Название группы"
// @Param song query string false "Название песни"
// @Param page query int false "Номер страницы" default(1)
// @Param limit query int false "Количество записей на странице" default(10)
// @Success 200 {array} models.Song
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /songs [get]
func (h *SongHandler) GetSongs(c *gin.Context) {
	// Извлечение параметров из запроса с значениями по умолчанию
	group := c.Query("group")
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		page = 1
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		limit = 10
	}

	// Вызов сервисного слоя для получения списка песен
	songs, err := h.songService.GetSongs(group, page, limit)
	if err != nil {
		log.Printf("Error fetching songs: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve songs",
			"details": err.Error(),
		})
		return
	}

	// Возвращение успешного ответа
	c.JSON(http.StatusOK, songs)
}

// GetSongTextHandler godoc
// @Summary Получение текста песни
// @Description Возвращает текст песни с пагинацией по куплетам
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "ID песни"
// @Param page query int false "Номер страницы" default(1)
// @Param limit query int false "Количество куплетов на странице" default(10)
// @Success 200 {object} models.SongText
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /songs/{id}/text [get]
func (h *SongHandler) GetSongText(c *gin.Context) {
	// Извлечение параметров из запроса
	songID, err := strconv.Atoi(c.Query("song_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid song ID",
		})
		return
	}

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		page = 1
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "3"))
	if err != nil {
		limit = 3
	}

	// Получение текста песни
	songText, err := h.songService.GetSongText(songID, page, limit)
	if err != nil {
		log.Printf("Error getting song text: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve song text",
			"details": err.Error(),
		})
		return
	}

	// Возвращение текста песни
	c.JSON(http.StatusOK, gin.H{
		"song_id": songID,
		"page":    page,
		"text":    songText,
	})
}

// CreateSongHandler godoc
// @Summary Добавление новой песни
// @Description Создает новую запись о песне
// @Tags songs
// @Accept json
// @Produce json
// @Param song body models.CreateSongRequest true "Информация о песне"
// @Success 201 {object} models.Song
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /songs [post]
func (h *SongHandler) CreateSong(c *gin.Context) {
	// Структура для привязки входящих данных
	var req models.CreateSongRequest

	// Валидация входящего JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Создание песни через сервисный слой
	song, err := h.songService.CreateSong(req)
	if err != nil {
		log.Printf("Error creating song: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create song",
			"details": err.Error(),
		})
		return
	}

	// Возвращение созданной песни
	c.JSON(http.StatusCreated, song)
}

// UpdateSongHandler godoc
// @Summary Обновление информации о песне
// @Description Обновляет существующую запись о песне
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "ID песни"
// @Param song body models.UpdateSongRequest true "Обновленная информация о песне"
// @Success 200 {object} models.Song
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /songs/{id} [put]
func (h *SongHandler) UpdateSong(c *gin.Context) {
	// Извлечение ID песни из параметров пути
	songID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid song ID",
		})
		return
	}

	// Структура для привязки данных обновления
	var updateData models.Song
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Вызов сервисного метода обновления
	if err := h.songService.UpdateSong(songID, updateData); err != nil {
		log.Printf("Error updating song: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update song",
			"details": err.Error(),
		})
		return
	}

	// Возвращение успешного ответа
	c.JSON(http.StatusOK, gin.H{
		"message": "Song updated successfully",
		"song_id": songID,
	})
}

// DeleteSongHandler godoc
// @Summary Удаление песни
// @Description Удаляет запись о песне по идентификатору
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "ID песни"
// @Success 204
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /songs/{id} [delete]
func (h *SongHandler) DeleteSong(c *gin.Context) {
	// Извлечение ID песни из параметров пути
	songID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid song ID",
		})
		return
	}

	// Вызов сервисного метода удаления
	if err := h.songService.DeleteSong(songID); err != nil {
		log.Printf("Error deleting song: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to delete song",
			"details": err.Error(),
		})
		return
	}

	// Возвращение успешного ответа
	c.JSON(http.StatusOK, gin.H{
		"message": "Song deleted successfully",
		"song_id": songID,
	})
}
