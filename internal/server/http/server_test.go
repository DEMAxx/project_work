package internalhttp

import (
	"fmt"
	"github.com/DEMAxx/project_work/internal/lrucache"
	"github.com/DEMAxx/project_work/pkg/config"
	"github.com/DEMAxx/project_work/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

const testImagesDir = "testdata"

func TestServer(t *testing.T) {
	logs := logger.MustSetupLogger(config.AppName, "test", true, "INFO")
	cnf := config.Config{
		Capability: 10,
		UploadPath: testImagesDir,
	}

	cache := lrucache.NewCache(cnf.Capability, cnf.UploadPath, logs)
	server := NewServer(&logs, "localhost:8080", cache, &cnf)

	t.Run("hello", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/hello", nil)
		require.NoError(t, err)

		rec := httptest.NewRecorder()

		server.httpServer.Handler.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "Hello, World!", rec.Body.String())
	})

	t.Run("fill success", func(t *testing.T) {
		err := os.MkdirAll(testImagesDir, 0755)

		require.NoError(t, err)

		fileUrl := "raw.githubusercontent.com/OtusGolang/final_project/master/examples/image-previewer/_gopher_original_1024x504.jpg"
		path := fmt.Sprintf(
			"/fill/%d/%d/%s",
			100,
			100,
			fileUrl,
		)

		req, err := http.NewRequest(http.MethodGet, path, nil)
		require.NoError(t, err)

		rec := httptest.NewRecorder()

		server.httpServer.Handler.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "image/jpeg", rec.Header().Get("Content-Type"))

		err = os.RemoveAll(testImagesDir)
		require.NoError(t, err)
	})
}
