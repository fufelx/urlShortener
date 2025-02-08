package pgsql

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
	"log"
	"os"
)

// UrlInfo - информация о URL.
type UrlInfo struct {
	Url      string `json:"url"`
	ShortUrl string `json:"shorturl"`
}

type Store struct {
	ctx context.Context
	db  pgxpool.Pool
}

// New функция для подключения к БД.
func New() (*Store, error) {
	err := godotenv.Load()
	if err != nil {
		log.Println("файл .env не найден, используются переменные окружения")
	}

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
func (s *Store) AddUrl(info UrlInfo) (string, error) {
	var existingShort string
	err := s.db.QueryRow(s.ctx, "SELECT shorturl FROM urls WHERE url = $1", info.Url).Scan(&existingShort)
	if err == nil {
		return existingShort, nil // Если ссылка уже существует, возвращаем её
	}

	tx, err := s.db.Begin(s.ctx)
	if err != nil {
		return "", err
	}
	defer tx.Rollback(s.ctx)

	_, err = tx.Exec(s.ctx, `INSERT INTO urls(url, shorturl) VALUES ($1,$2)`, info.Url, info.ShortUrl)
	if err != nil {
		return "", err
	} else {
		tx.Commit(s.ctx)
		return info.ShortUrl, nil
	}
}

// GetUrlByShotrurl возвращает URL по ShortURL.
func (s *Store) GetUrlByShotrurl(shorturl string) (UrlInfo, error) {
	rows, err := s.db.Query(s.ctx, `SELECT url, shorturl FROM urls WHERE shorturl = $1`, shorturl)
	if err != nil {
		return UrlInfo{}, err
	} else {
		var info UrlInfo
		for rows.Next() {
			var t UrlInfo
			rows.Scan(&t.Url, &t.ShortUrl)
			info = UrlInfo{t.Url, t.ShortUrl}
			if rows.Err() != nil {
				return UrlInfo{}, err
			}
		}
		return info, nil
	}
}
