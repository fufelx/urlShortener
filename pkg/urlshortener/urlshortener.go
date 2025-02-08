package urlshortener

import (
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
	res := "http://localhost:3030/" + string(short)
	return res
}
