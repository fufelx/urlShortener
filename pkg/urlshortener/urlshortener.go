package urlshortener

import (
	"example.com/m/internal/storage"
	"math/rand"
	"time"
)

// Cимволы для коротких ссылок
const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"

func MakeUrlShort() string {
	rand.Seed(time.Now().UnixNano())
	short := make([]byte, 10)
	for i := range short {
		short[i] = chars[rand.Intn(len(chars))]
	}

	if storage.InMemory {
		if _, exist := storage.ShortToOriginalmap[string(short)]; exist {
			// Если сгенерированная короткая ссылка уже существует, генерируем новую
			return MakeUrlShort()
		}
	}

	return string(short)
}
