package file_search

import (
	"github.com/DEMAxx/project_work/pkg/logger"
	"github.com/stretchr/testify/require"
	"net/http"
	"os"
	"path/filepath"
	"testing"
)

func TestFileSearch(t *testing.T) {
	tmpDir := os.TempDir()
	outputPath := filepath.Join(tmpDir, "output")
	logs := logger.MustSetupLogger("previewer", "Test", true, "info")

	t.Run("success", func(t *testing.T) {
		r, err := FetchFileFromURL("https://raw.githubusercontent.com/OtusGolang/final_project/master/examples/image-previewer/_gopher_original_1024x504.jpg", outputPath, &logs)

		require.NoError(t, err)
		require.NotNil(t, r)
		require.True(t, r.StatusCode == http.StatusOK)
	})

	t.Run("wrong address", func(t *testing.T) {
		_, err := FetchFileFromURL("https://raw.githubusercontent.com/OtusGolang/final_project/master/examples/image-previewer/not_gopher_original.jpg", outputPath, &logs)

		require.Error(t, err)
	})

	t.Run("not found", func(t *testing.T) {
		_, err := FetchFileFromURL("localhost:9999/image.png", outputPath, &logs)
		require.Error(t, err)
		require.ErrorContains(t, err, "connection refused")
	})

}
