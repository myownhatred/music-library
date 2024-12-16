package main

import (
	"log"
	"music-library/internal/config"
	"music-library/internal/handlers"
	"music-library/internal/repository"
	"music-library/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Загрузка переменных окружения
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Инициализация конфигурации
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Cannot load config: %v", err)
	}

	// Подключение к базе данных
	db, err := repository.InitPostgresDB(cfg)
	if err != nil {
		log.Fatalf("Cannot connect to database: %v", err)
	}
	defer db.Close()

	// Создание репозитория, сервисов и обработчиков
	songRepo := repository.NewSongRepository(db)
	externalAPI := service.NewExternalAPIService(cfg)
	songService := service.NewSongService(songRepo, externalAPI)
	songHandler := handlers.NewSongHandler(songService)

	// Настройка роутера Gin
	router := gin.Default()

	// Группировка маршрутов
	v1 := router.Group("/api")
	{
		v1.GET("/songs", songHandler.GetSongs)
		v1.GET("/song/text", songHandler.GetSongText)
		v1.POST("/song", songHandler.CreateSong)
		v1.PUT("/song/:id", songHandler.UpdateSong)
		v1.DELETE("/song/:id", songHandler.DeleteSong)
	}

	// Запуск сервера
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
