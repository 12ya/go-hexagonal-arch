package inmemory

import (
	"context"
	"fmt"
	"sync"

	"github.com/12ya/go-hexagonal-arch/internal/domain"
)

type Store struct {
	data map[string]*Port
	mu   sync.RWMutex
}

func NewStore() *Store {
	return &Store{
		data: make(map[string]*Port),
	}
}

func (s *Store) GetPort(ctx context.Context, id string) (*domain.Port, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	storePort, exists := s.data[id]
	if !exists {
		return nil, domain.ErrNotFound
	}

	port, err := storePort.storePortToDomain()
	if err != nil {
		return nil, fmt.Errorf("portStoreToDomain failed: %w", err)
	}

	return port, nil
}
