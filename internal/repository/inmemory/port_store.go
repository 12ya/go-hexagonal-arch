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

func (s *Store) CountPorts(context.Context) (int, error) {
	s.mu.Lock()
	defer s.mu.RUnlock()

	return len(s.data), nil
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
		s.updatePort(storePort.ID, storePort.toUpdate())
	} else {
		s.data[storePort.ID] = storePort
	}
	return nil
}

func (s *Store) updatePort(id string, update *domain.PortUpdate) {
	storePort := s.data[id]
	// updated := storePort.Copy()

	if update.Name != nil {
		storePort.Name = *update.Name
	}
	if update.Code != nil {
		storePort.Code = *update.Code
	}
	if update.City != nil {
		storePort.City = *update.City
	}
	if update.Country != nil {
		storePort.Country = *update.Country
	}
	if update.Alias != nil {
		storePort.Alias = append([]string(nil), *update.Alias...)
	}
	if update.Regions != nil {
		storePort.Regions = append([]string(nil), *update.Regions...)
	}
	if update.Coordinates != nil {
		storePort.Coordinates = append([]float64(nil), *update.Coordinates...)
	}
	if update.Province != nil {
		storePort.Province = *update.Province
	}
	if update.Timezone != nil {
		storePort.Timezone = *update.Timezone
	}
	if update.Unlocs != nil {
		storePort.Unlocs = append([]string(nil), *update.Unlocs...)
	}
}
