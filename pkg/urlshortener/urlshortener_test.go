package urlshortener

import (
	"example.com/m/internal/storage"
	"strconv"
	"testing"
)

// Проверяем, что на 10000(можно изменить) сокращенных ссылок не создается дубликат.
func TestMakeUrlShort(t *testing.T) {
	storage.InMemory = true
	shortmap := make(map[string]int)
	for i := 0; i < 10000; i++ {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			short := MakeUrlShort()
			if _, exist := shortmap[short]; exist {
				t.Errorf("Создан дубликат в %v и %v тестах", i, shortmap[short])
				return
			}
			shortmap[short] = i
		})
	}
}
