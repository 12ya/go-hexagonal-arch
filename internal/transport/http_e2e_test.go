package transport_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/12ya/go-hexagonal-arch/internal/repository/inmemory"
	"github.com/12ya/go-hexagonal-arch/internal/services"
	"github.com/12ya/go-hexagonal-arch/internal/transport"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type HttpTestSuite struct {
	suite.Suite
	portService transport.PortService
	httpServer  transport.HttpServer
}

func NewHttpTestSuite() *HttpTestSuite {
	suite := &HttpTestSuite{}

	portStoreRepo := inmemory.NewStore()

	suite.portService = services.NewPortService(portStoreRepo)

	// http server with application injected
	suite.httpServer = transport.NewHttpServer(suite.portService)

	return suite
}

func TestHttpTestSuite(t *testing.T) {
	suite.Run(t, NewHttpTestSuite())
}

func (suite *HttpTestSuite) TestUploadPorts() {
	portsRequest, err := os.ReadFile("./ports_request.json")
	require.NoError(suite.T(), err)

	requestPortsTotal := countJSONPorts(suite.T(), portsRequest)

	portsResponse, err := os.ReadFile("./ports_response.json")
	require.NoError(suite.T(), err)

	req := httptest.NewRequest(http.MethodPost, "/ports", bytes.NewBuffer(portsRequest))

	w := httptest.NewRecorder()

	suite.httpServer.UploadPorts(w, req)

	res := w.Result()
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	require.NoError(suite.T(), err)

	require.Equal(suite.T(), http.StatusOK, res.StatusCode)
	require.Equal(suite.T(), portsResponse, data)

	storedPortsTotal, err := suite.portService.CountPorts(context.Background())
	require.NoError(suite.T(), err)

	require.Equal(suite.T(), requestPortsTotal, storedPortsTotal)
}

func countJSONPorts(t *testing.T, portsJSON []byte) int {
	t.Helper()
	var ports map[string]struct{}
	err := json.Unmarshal(portsJSON, &ports)
	require.NoError(t, err)
	return len(ports)
}
