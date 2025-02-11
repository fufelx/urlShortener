package storage

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Store struct {
	ctx context.Context
	db  pgxpool.Pool
}

type UrlInfo struct {
	Url      string `json:"url"`
	ShortUrl string `json:"shorturl"`
}

type Storage interface {
	AddUrl(info UrlInfo) (string, error)
	GetUrlByShotrurl(shorturl string) (UrlInfo, error)
}
