package filesearch

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/rs/zerolog"
)

func FetchFileFromURL(imageURL, outputPath string, logger *zerolog.Logger) (*http.Response, error) {
	if !strings.HasPrefix(imageURL, "http://") && !strings.HasPrefix(imageURL, "https://") {
		imageURL = fmt.Sprintf("https://%s", imageURL)
	}

	resp, err := http.Get(imageURL) //nolint
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.Error().Msg("failed to close response body")
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch file: %s", resp.Status)
	}

	// Create the output file
	outFile, err := os.Create(outputPath)
	if err != nil {
		return nil, err
	}

	defer func(outFile *os.File) {
		err := outFile.Close()
		if err != nil {
			logger.Error().Msg("failed to close response body")
		}
	}(outFile)

	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
