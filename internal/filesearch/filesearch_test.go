package filesearch

import (
	"context"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/DEMAxx/project_work/pkg/logger"
	"github.com/stretchr/testify/require"
)

func TestFileSearch(t *testing.T) {
	outputPath := filepath.Join(os.TempDir(), "output")
	logs := logger.MustSetupLogger("previewer", "Test", true, "info")
	ctx := context.Background()

	client := NewClient(
		ctx,
		&http.Request{
			Method: "GET",
			URL: &url.URL{
				Scheme: "http",
				Host:   "localhost",
				Path:   outputPath,
			},
		},
	)

	t.Run("success", func(t *testing.T) {
		r, err := client.FetchFileFromURL(
			"https://raw.githubusercontent.com/OtusGolang/final_project/master/examples/image-previewer/_gopher_original_1024x504.jpg", //nolint
			outputPath,
			&logs,
		)

		require.NoError(t, err)
		require.NotNil(t, r)
		require.True(t, r.StatusCode == http.StatusOK)
		err = r.Body.Close()
		require.NoError(t, err)
	})

	t.Run("wrong address", func(t *testing.T) {
		r, err := client.FetchFileFromURL(
			"https://raw.githubusercontent.com/OtusGolang/final_project/master/examples/image-previewer/not_gopher_original.jpg",
			outputPath,
			&logs,
		)

		require.Error(t, err)

		if err == nil {
			err = r.Body.Close()
			require.NoError(t, err)
		}
	})

	t.Run("not found", func(t *testing.T) {
		r, err := client.FetchFileFromURL("localhost:9999/image.png", outputPath, &logs)
		require.Error(t, err)
		require.ErrorContains(t, err, "connection refused")

		if err == nil {
			err = r.Body.Close()
			require.NoError(t, err)
		}
	})
}
