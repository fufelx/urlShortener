package service

import (
	"errors"
	"example.com/m/internal/storage"
	"example.com/m/pkg/urlshortener"
)

func NewLinik(db *storage.Store, link string) (string, error) {
	if storage.InMemory {
		// Если ссылка уже существует, возвращаем её
		if shorturlexist, exist := storage.OriginalToShortmap[link]; exist {
			return shorturlexist, nil
		}

		shortlink := urlshortener.MakeUrlShort()
		storage.Mu.Lock()
		storage.OriginalToShortmap[link] = shortlink
		storage.ShortToOriginalmap[shortlink] = link
		storage.Mu.Unlock()

		return shortlink, nil
	} else {
		shortlink := urlshortener.MakeUrlShort()

		//Из бд может вернуться уже существующая короткая ссылка, если ее нет, то вставится сгенерированная
		shortlinktmp, err := db.AddUrl(link, shortlink)
		if err != nil {
			return "", err
		}

		return shortlinktmp, nil
	}
}

func GetOriginalLink(db *storage.Store, shortlink string) (string, error) {
	if storage.InMemory {
		if originalexist, exist := storage.ShortToOriginalmap[shortlink]; exist {
			return originalexist, nil
		} else {
			return "", errors.New("оригинал ссылки отсутствует")
		}
	} else {
		originallink, err := db.GetUrlByShotrurl(shortlink)
		if err != nil {
			return "", errors.New("оригинал ссылки отсутствует")
		}
		return originallink, nil
	}
}
