CREATE TABLE IF NOT EXISTS songs (
    id SERIAL PRIMARY KEY,
    group VARCHAR(255) NOT NULL,
    song_name VARCHAR(255) NOT NULL,
    release_date DATE,
    text TEXT,
    link VARCHAR(512)
);