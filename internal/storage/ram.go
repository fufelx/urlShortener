package storage

import "sync"

var (
	Mu                 = &sync.Mutex{}
	InMemory           = false
	OriginalToShortmap = make(map[string]string) // мапа для хранения сокращения и оригинальной ссылки
	ShortToOriginalmap = make(map[string]string) // мапа для хранения оригинальной ссылки и ее сокращения
)
