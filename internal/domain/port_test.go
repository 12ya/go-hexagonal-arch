package domain

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewPort(t *testing.T) {
	t.Parallel()

	portID := "port id"
	portName := "port name"
	portCode := "port code"
	portCity := "port city"
	portCountry := "port country"

	t.Run("valid", func(t *testing.T) {
		port, err := NewPort(portID, portName, portCode, portCity, portCountry, nil, nil, nil, "", "", nil)
		require.NoError(t, err)

		require.Equal(t, portID, port.ID())
		require.Equal(t, portCode, port.Code())
		require.Equal(t, portName, port.Name())
		require.Equal(t, portCity, port.City())
		require.Equal(t, portCountry, port.Country())
	})
}
