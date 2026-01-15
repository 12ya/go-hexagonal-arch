package inmemory

import "sync"

type Store struct {
	data map[string]*Port
	mu   sync.RWMutex
}
