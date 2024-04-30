package e2e

import (
	"context"
	"encoding/json"
	"github.com/evacchi/envoy-ext-server/extproc"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestE2E_Example(t *testing.T) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	cmd := setupEnvoy(t, "testdata/envoy.yaml")
	defer cmd.Process.Kill()

	go setupEcho(ctx, 8000)
	go extproc.ServeFromConfig(ctx, "testdata/ext-server.yaml")

	time.Sleep(1 * time.Second)

	req, err := http.NewRequest("POST", "http://127.0.0.1:8080/", strings.NewReader(`{"hello":"world"}`))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	// The trivial plugin will append a message to the JSON payload.
	req.Header.Set("X-Trivial-Process", "yes")
	// The wasm plugin will append a message to the JSON payload.
	req.Header.Set("X-Wasm", "append")

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	decoded := &Request{}
	reader, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	err = json.Unmarshal(reader, decoded)
	require.NoError(t, err)

	require.Equal(t, `hello from the trivial plugin`, decoded.Headers["X-Custom-Header"][0])
	require.Equal(t, `{"hello":"world","message":"hello from trivial plugin","trailer":"...and that's all folks"}`, decoded.Body)

}
