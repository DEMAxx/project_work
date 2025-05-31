package filesearch

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/rs/zerolog"
)

func NewClient(ctx context.Context, r *http.Request) *Client {
	return &Client{
		ctx: ctx,
		r:   r,
	}
}

type Client struct {
	ctx context.Context
	r   *http.Request
}

func (p *Client) newHTTPRequest(url string) *http.Request {
	prxReq, _ := http.NewRequestWithContext(p.ctx, p.r.Method, url, p.r.Body)
	prxQuery := prxReq.URL.Query()

	for key, values := range p.r.URL.Query() {
		for _, value := range values {
			prxQuery.Add(key, value)
		}
	}

	for key, values := range p.r.Header {
		for _, value := range values {
			prxReq.Header.Set(key, value)
		}
	}

	prxReq.URL.RawQuery = prxQuery.Encode()

	return prxReq
}

func (p *Client) FetchFileFromURL(imageURL, outputPath string, logger *zerolog.Logger) (*http.Response, error) {
	if strings.HasPrefix(imageURL, "http:/") {
		imageURL = strings.Trim(strings.Replace(imageURL, "http:/", "", 1), "/")
	}

	if strings.HasPrefix(imageURL, "https:/") {
		imageURL = strings.Trim(strings.Replace(imageURL, "https:/", "", 1), "/")
	}

	imageURL = fmt.Sprintf("https://%s", imageURL)

	req := p.newHTTPRequest(imageURL)

	slog.Debug(fmt.Sprintf("Proxy: IN='%s %s' -> OUT='%s %s'", p.r.Method, p.r.URL.String(), req.Method, req.URL.String())) //nolint

	resp, err := http.DefaultClient.Do(req)
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
