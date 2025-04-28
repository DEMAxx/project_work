package file_search

import (
	"fmt"
	"github.com/rs/zerolog"
	"io"
	"net/http"
	"os"
)

func FetchFileFromURL(imageUrl, outputPath string, logger *zerolog.Logger) error {
	resp, err := http.Get(imageUrl)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.Error().Msg("failed to close response body")
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to fetch file: %s", resp.Status)
	}

	// Create the output file
	outFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}

	defer func(outFile *os.File) {
		err := outFile.Close()
		if err != nil {
			logger.Error().Msg("failed to close response body")
		}
	}(outFile)

	_, err = io.Copy(outFile, resp.Body)

	if err != nil {
		return err
	}

	return nil
}
