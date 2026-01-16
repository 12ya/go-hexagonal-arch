package inmemory

import (
	"context"
	"testing"

	"github.com/12ya/go-hexagonal-arch/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestPortStore_Upsert(t *testing.T) {
	t.Parallel()

	store := NewStore()
	t.Run("create port", func(t *testing.T) {
		t.Parallel()
		randomPort := newRandomDomainPort(t)
		err := store.Upsert(context.Background(), randomPort)
		require.NoError(t, err)

		port, err := store.GetPort(context.Background(), randomPort.ID())
		require.NoError(t, err)

		require.Equal(t, port, randomPort)
	})

	t.Run("update port", func(t *testing.T) {
		t.Parallel()
		randomPort := newRandomDomainPort(t)
		err := store.Upsert(context.Background(), randomPort)
		require.NoError(t, err)

		beforeUpdate, err := store.GetPort(context.Background(), randomPort.ID())
		require.NoError(t, err)

		require.Equal(t, beforeUpdate, randomPort)

		err = randomPort.SetName("updated name")
		require.NoError(t, err)

		err = store.Upsert(context.Background(), randomPort)
		require.NoError(t, err)

		updatedPort, err := store.GetPort(context.Background(), randomPort.ID())
		require.NoError(t, err)

		require.NotEqual(t, beforeUpdate.Name(), updatedPort.Name())
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
