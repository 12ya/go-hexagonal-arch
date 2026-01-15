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
