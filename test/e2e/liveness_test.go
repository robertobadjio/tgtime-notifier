package e2e

import (
	"io"
	"net/http"
	"os"
	"os/exec"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestE2ELiveness(t *testing.T) {
	const servicePort = "8081"

	cmd := exec.Command("../../tgtime-notifier")
	cmd.Env = append(
		os.Environ(),
		"HTTP_PORT="+servicePort,
		"BOT_TOKEN=8132539422:AAHb_lH0dHEfmxcjS6RqfVOh3egX__t3lU4",
		"WEBHOOK_PATH=telegram",
		"WEBHOOK_LINK=https://tgtime.ru/telegram",
		"TGTIME_AGGREGATOR_HOST=",
		"TGTIME_AGGREGATOR_PORT=1080",
		"TGTIME_API_HOST=",
		"TGTIME_API_PORT=1080",
	)
	require.NoError(t, cmd.Start())

	time.Sleep(time.Second)

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, "http://127.0.0.1:8081/liveness", nil)
	require.NoError(t, err)

	res, err := client.Do(req)
	require.NoError(t, err)
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		require.NoError(t, err)
	}(res.Body)

	body, err := io.ReadAll(res.Body)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "{\"status\":200}\n", string(body))

	require.NoError(t, cmd.Process.Signal(syscall.SIGTERM))

	time.Sleep(time.Second)
}
