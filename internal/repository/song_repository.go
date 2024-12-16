package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"music-library/internal/config"
	"music-library/internal/models"
	"strings"

	_ "github.com/lib/pq"
)

type SongRepository struct {
	db *sql.DB
}

// InitPostgresDB устанавливает соединение с базой данных PostgreSQL
func InitPostgresDB(cfg *config.Config) (*sql.DB, error) {
	// Формирование строки подключения с параметрами из конфигурации
	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
	)

	// Открытие соединения с базой данных
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Printf("Error opening database connection: %v", err)
		return nil, err
	}

	// Проверка соединения с базой данных
	if err = db.Ping(); err != nil {
		log.Printf("Error pinging database: %v", err)
		return nil, err
	}

	log.Println("Successfully connected to the database")
	return db, nil
}

// NewSongRepository создает новый экземпляр репозитория песен
func NewSongRepository(db *sql.DB) *SongRepository {
	return &SongRepository{db: db}
}

// GetSongs возвращает список песен с фильтрацией и пагинацией
func (r *SongRepository) GetSongs(group string, page, limit int) ([]models.Song, error) {
	// Базовый запрос с динамическим условием фильтрации по группе
	query := `SELECT id, "group", song_name, release_date, text, link 
              FROM songs 
              WHERE 1=1`

	// Параметры для запроса
	var args []interface{}
	var conditions []string

	// Добавление фильтра по группе, если указана
	if group != "" {
		conditions = append(conditions, fmt.Sprintf("group ILIKE $%d", len(args)+1))
		args = append(args, "%"+group+"%")
	}

	// Добавление условий к запросу
	if len(conditions) > 0 {
		query += " AND " + strings.Join(conditions, " AND ")
	}

	// Добавление пагинации
	query += " LIMIT $" + fmt.Sprintf("%d", len(args)+1)
	args = append(args, limit)

	query += " OFFSET $" + fmt.Sprintf("%d", len(args)+1)
	args = append(args, (page-1)*limit)

	// Выполнение запроса
	rows, err := r.db.Query(query, args...)
	if err != nil {
		log.Printf("Error querying songs: %v", err)
		return nil, err
	}
	defer rows.Close()

	var songs []models.Song
	for rows.Next() {
		var song models.Song
		err := rows.Scan(
			&song.ID, &song.Group, &song.SongName,
			&song.ReleaseDate, &song.Text, &song.Link,
		)
		if err != nil {
			log.Printf("Error scanning song row: %v", err)
			return nil, err
		}
		songs = append(songs, song)
	}

	return songs, nil
}

// CreateSong добавляет новую песню в базу данных
func (r *SongRepository) CreateSong(song *models.Song) (*models.Song, error) {
	query := `
		INSERT INTO songs (group, song_name, release_date, text, link)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	err := r.db.QueryRow(
		query,
		song.Group,
		song.SongName,
		song.ReleaseDate,
		song.Text,
		song.Link,
	).Scan(&song.ID)

	if err != nil {
		log.Printf("Error creating song: %v", err)
		return nil, err
	}

	return song, nil
}

// UpdateSong обновляет информацию о песне
func (r *SongRepository) UpdateSong(song *models.Song) error {
	query := `
		UPDATE songs 
		SET group = $1, song_name = $2, 
		    release_date = $3, text = $4, link = $5
		WHERE id = $6
	`

	result, err := r.db.Exec(
		query,
		song.Group,
		song.SongName,
		song.ReleaseDate,
		song.Text,
		song.Link,
		song.ID,
	)

	if err != nil {
		log.Printf("Error updating song: %v", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("no song found with the given ID")
	}

	return nil
}

// DeleteSong удаляет песню по идентификатору
func (r *SongRepository) DeleteSong(songID int) error {
	query := `DELETE FROM songs WHERE id = $1`

	result, err := r.db.Exec(query, songID)
	if err != nil {
		log.Printf("Error deleting song: %v", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("no song found with the given ID")
	}

	return nil
}

// GetSongText получает текст песни постранично
func (r *SongRepository) GetSongText(songID, page, limit int) (string, error) {
	query := `
		SELECT text 
		FROM songs 
		WHERE id = $1
	`

	var fullText string
	err := r.db.QueryRow(query, songID).Scan(&fullText)
	if err != nil {
		log.Printf("Error fetching song text: %v", err)
		return "", err
	}

	// Разбиение текста на страницы (куплеты)
	verses := strings.Split(fullText, "\n\n")

	// Вычисление диапазона страниц
	start := (page - 1) * limit
	end := start + limit

	if start >= len(verses) {
		return "", errors.New("page out of range")
	}

	if end > len(verses) {
		end = len(verses)
	}

	return strings.Join(verses[start:end], "\n\n"), nil
}
