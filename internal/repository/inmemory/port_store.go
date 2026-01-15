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

func (s *Store) Upsert(ctx context.Context, p *domain.Port) error {
	if p == nil {
		return domain.ErrNil
	}

	storePort, err := domainPortToStore(p)
	if err != nil {
		return fmt.Errorf("error converting domain port to store one: %w", err)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.data[storePort.ID]; exists {
		return s.updatePort(ctx, storePort)
	} else {
		return s.createPort(ctx, storePort)
	}
}

func (s *Store) createPort(ctx context.Context, port *Port) error
func (s *Store) updatePort(ctx context.Context, port *Port) error
