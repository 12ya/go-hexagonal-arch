package services

import (
	"context"

	"github.com/12ya/go-hexagonal-arch/internal/domain"
)

type PortRepository interface {
	Upsert(context.Context, *domain.Port) error
	GetPort(ctx context.Context, id string) (*domain.Port, error)
}

type PortService struct {
	repo PortRepository
}

func NewPortService(repo PortRepository) *PortService {
	return &PortService{repo: repo}
}

func (ps *PortService) GetPort(ctx context.Context, id string) (*domain.Port, error) {
	return ps.repo.GetPort(ctx, id)
}

func (ps *PortService) Upsert(ctx context.Context, port *domain.Port) error {
	return ps.repo.Upsert(ctx, port)
}
