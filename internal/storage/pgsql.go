package storage

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"os"
)

// New функция для подключения к БД.
func New() (*Store, error) {
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")
	connStr := fmt.Sprintf("postgres://%s:%s@%s/%s",
		dbUser, dbPassword, dbHost, dbName)
	var ctx context.Context = context.Background()
	db, err := pgxpool.Connect(ctx, connStr)
	if err != nil {
		return nil, err
	}
	result := Store{ctx: ctx, db: *db}
	return &result, nil
}

// AddUrl добавляет URL и ShortURL в БД.
func (s *Store) AddUrl(originallink, shortlink string) (string, error) {
	var existingShort string
	err := s.db.QueryRow(s.ctx, "SELECT shorturl FROM urls WHERE url = $1", originallink).Scan(&existingShort)
	if err == nil {
		return existingShort, nil // Если ссылка уже существует, возвращаем её
	}

	_, err = s.db.Exec(s.ctx, `INSERT INTO urls(url, shorturl) VALUES ($1,$2)`, originallink, shortlink)
	if err != nil {
		return "", err
	}

	return shortlink, nil
}

// GetUrlByShotrurl возвращает URL по ShortURL.
func (s *Store) GetUrlByShotrurl(shorturl string) (string, error) {
	rows, err := s.db.Query(s.ctx, `SELECT url FROM urls WHERE shorturl = $1`, shorturl)
	if err != nil {
		return "", err
	} else {
		var link string
		for rows.Next() {
			rows.Scan(&link)
			if rows.Err() != nil {
				return "", err
			}
		}
		return link, nil
	}
}
