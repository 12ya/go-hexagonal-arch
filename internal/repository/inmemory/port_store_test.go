package inmemory

import (
	"context"
	"testing"

	"github.com/12ya/go-hexagonal-arch/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestPortStore_Upsert(t *testing.T) {
	store := NewStore()

	t.Run("create prot", func(t *testing.T) {
		randomPort := newRandomDomainPort(t)
		err := store.Upsert(context.Background(), randomPort)
		require.NoError(t, err)

		port, err := store.GetPort(context.Background(), randomPort.ID())
		require.NoError(t, err)

		require.Equal(t, port, randomPort)
	})
}

func newRandomDomainPort(t *testing.T) *domain.Port {
	t.Helper()

	randomID := uuid.New().String()
	port, err := domain.NewPort(
		randomID,
		randomID,
		randomID,
		randomID,
		randomID,
		[]string{randomID},
		[]string{randomID},
		[]float64{1.0, 2.0},
		randomID,
		randomID,
		nil,
	)
	require.NoError(t, err)

	return port
}
