package internalhttp

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/DEMAxx/project_work/internal/lrucache"
	"github.com/DEMAxx/project_work/pkg/config"
	"github.com/DEMAxx/project_work/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testImagesDir = "testdata"

func TestServer(t *testing.T) {
	ctx := context.Background()
	logs := logger.MustSetupLogger(config.AppName, "test", true, "INFO")
	cnf := config.Config{
		Capability: 10,
		UploadPath: testImagesDir,
	}

	cache := lrucache.NewCache(cnf.Capability, cnf.UploadPath, logs)
	server := NewServer(ctx, &logs, "localhost:8080", cache, &cnf)

	t.Run("hello", func(t *testing.T) {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/hello", nil)
		require.NoError(t, err)

		rec := httptest.NewRecorder()

		server.httpServer.Handler.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "Hello, World!", rec.Body.String())
	})

	t.Run("fill success", func(t *testing.T) {
		err := os.MkdirAll(testImagesDir, 0o755)

		require.NoError(t, err)

		fileURL := "raw.githubusercontent.com/OtusGolang/final_project/master/examples/image-previewer/_gopher_original_1024x504.jpg" //nolint
		path := fmt.Sprintf(
			"/fill/%d/%d/%s",
			100,
			100,
			fileURL,
		)

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, path, nil)
		require.NoError(t, err)

		rec := httptest.NewRecorder()

		server.httpServer.Handler.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "image/jpeg", rec.Header().Get("Content-Type"))

		err = os.RemoveAll(testImagesDir)
		require.NoError(t, err)
	})
}
