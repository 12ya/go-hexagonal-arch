package transport

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/12ya/go-hexagonal-arch/internal/domain"
)

// reads Ports from provided Reader and sends them to portChan
func readPorts(ctx context.Context, r io.Reader, portChan chan Port) error {
	decoder := json.NewDecoder(r)

	// Read the opening delimiter
	token, err := decoder.Token()
	if err != nil {
		return fmt.Errorf("failed to read the opening delimiter: %w", err)
	}

	if token != json.Delim('{') {
		return fmt.Errorf("expected {, got %v", token)
	}

	for decoder.More() {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		token, err := decoder.Token()
		if err != nil {
			return fmt.Errorf("expected string, got %v", token)
		}

		portID, ok := token.(string)
		if !ok {
			return fmt.Errorf("expected string, got %v", token)
		}

		var port Port
		if err := decoder.Decode(&port); err != nil {
			return fmt.Errorf("failed to decode port: %w", err)
		}

		port.ID = portID
		portChan <- port
	}

	return nil
}

func httpPortToDomain(p *Port) (*domain.Port, error) {
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
