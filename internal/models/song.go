package models

type Song struct {
	ID          int    `json:"id" db:"id"`
	Group       string `json:"group" db:"group"`
	SongName    string `json:"song" db:"song_name"`
	ReleaseDate string `json:"release_date" db:"release_date"`
	Text        string `json:"text" db:"text"`
	Link        string `json:"link" db:"link"`
}

type CreateSongRequest struct {
	Group string `json:"group" binding:"required"`
	Song  string `json:"song" binding:"required"`
}
