package inmemory

import (
	"errors"

	"github.com/12ya/go-hexagonal-arch/internal/domain"
)

type Port struct {
	ID      string
	Name    string
	Code    string
	City    string
	Country string

	Alias       []string
	Regions     []string
	Coordinates []float64
	Province    string
	Timezone    string
	Unlocs      []string
}

func (p *Port) storePortToDomain() (*domain.Port, error) {
	if p == nil {
		return nil, errors.New("store port is nil")
	}

	return domain.NewPort(
		p.ID,
		p.Name,
		p.Code,
		p.City,
		p.Country,
		append([]string(nil), p.Alias...),
		append([]string(nil), p.Regions...),
		append([]float64(nil), p.Coordinates...),
		p.Province,
		p.Timezone,
		append([]string(nil), p.Unlocs...),
	)
}

func (p *Port) toUpdate() *domain.PortUpdate {
	return &domain.PortUpdate{
		Name:        &p.Name,
		Code:        &p.Code,
		City:        &p.City,
		Country:     &p.Country,
		Alias:       &p.Alias,
		Regions:     &p.Regions,
		Coordinates: &p.Coordinates,
		Province:    &p.Province,
		Timezone:    &p.Timezone,
		Unlocs:      &p.Unlocs,
	}
}

func domainPortToStore(dp *domain.Port) (*Port, error) {
	if dp == nil {
		return nil, domain.ErrNil
	}
	return &Port{
		ID:          dp.ID(),
		Name:        dp.Name(),
		Code:        dp.Code(),
		City:        dp.City(),
		Country:     dp.Country(),
		Alias:       dp.Alias(),
		Regions:     dp.Regions(),
		Coordinates: dp.Coordinates(),
		Province:    dp.Province(),
		Timezone:    dp.Timezone(),
		Unlocs:      dp.Unlocs(),
	}, nil
}
