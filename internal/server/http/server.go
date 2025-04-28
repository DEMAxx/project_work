package internalhttp

import (
	"context"
	"errors"
	"fmt"
	"github.com/DEMAxx/project_work/internal/file_search"
	lrucache "github.com/DEMAxx/project_work/internal/lru_cache"
	"github.com/DEMAxx/project_work/pkg/config"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"net/http"
	"strings"
	"time"
)

type Server struct {
	httpServer      *http.Server
	grpcHostAndPort string
	logger          *zerolog.Logger
	cache           lrucache.Cache
}

func NewServer(logger *zerolog.Logger, hostAndPort string, cache lrucache.Cache, cnf *config.Config) *Server {
	mux := http.NewServeMux()

	mux.Handle("/hello", LoggingMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientIP := r.RemoteAddr
		dateTime := time.Now().Format(time.RFC3339)
		method := r.Method
		path := r.URL.Path
		httpVersion := r.Proto
		userAgent := r.Header.Get("User-Agent")

		logger.Info().Msg(
			fmt.Sprintf(
				"Client IP: %s, DateTime: %s, Method: %s, Path: %s, HTTP Version: %s, User Agent: %s",
				clientIP, dateTime, method, path, httpVersion, userAgent,
			),
		)

		write, err := w.Write([]byte("Hello, World!"))
		if err != nil {
			return
		}
		logger.Info().Msg(fmt.Sprintf("response: %d", write))
	}), logger))

	mux.Handle("/fill/", LoggingMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path[len("/fill/"):]
		parts := strings.Split(path, "/")

		if len(parts) < 3 {
			http.Error(w, "Invalid URL format", 400)
			return
		}

		height, width, imageUrl := parts[0], parts[1], strings.Join(parts[2:], "/")

		logger.Info().Msg(
			fmt.Sprintf(
				"Extracted vars - height: %s, width: %s, image url: %s", height, width, imageUrl,
			),
		)

		if !strings.HasSuffix(imageUrl, ".jpg") {
			http.Error(w, "Invalid image URL format. Only .jpg files are supported.", http.StatusBadRequest)
			return
		}

		cachedImageUrl, found := cache.Get(
			lrucache.Key(
				fmt.Sprintf(
					"%s_%s_%s", width, height, imageUrl,
				),
			),
		)

		if found {
			//todo return cached image

			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte(fmt.Sprintf("Image already cached: %s", cachedImageUrl)))
			if err != nil {
				return
			}
		}

		uid := uuid.New()

		if err := file_search.FetchFileFromURL(
			imageUrl,
			fmt.Sprintf("%s/%s_%s_%s.jpg", cnf.UploadPath, width, height, uid),
			logger,
		); err != nil {
			http.Error(w, "Failed to fetch image.", http.StatusInternalServerError)
		}

		if err := cache.Set(
			lrucache.Key(
				fmt.Sprintf(
					"%s_%s_%s", width, height, imageUrl,
				),
			),
			fmt.Sprintf("%s_%s_%s", width, height, uid),
		); err {
			http.Error(w, "Failed to store image.", http.StatusInternalServerError)
		}

		w.WriteHeader(http.StatusOK)
	}), logger))

	return &Server{
		httpServer: &http.Server{
			Addr:              hostAndPort,
			Handler:           mux,
			ReadHeaderTimeout: 5 * time.Second,
		},
		logger: logger,
		cache:  cache,
	}
}

func (s *Server) Start(ctx context.Context) error {
	s.logger.Info().Msg("Starting HTTP server...")

	// Start HTTP server
	go func() {
		s.logger.Info().Msg("HTTP server start...")

		if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.logger.Error().Msg(fmt.Sprintf("HTTP server ListenAndServe: %s", err.Error()))
		}

		s.logger.Info().Msg("HTTP server started")
	}()

	<-ctx.Done()
	return s.Stop(ctx)
}

func (s *Server) Stop(ctx context.Context) error {
	s.logger.Info().Msg("Stopping HTTP server...")

	// Stop HTTP server
	shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(shutdownCtx); err != nil {
		s.logger.Error().Msg(fmt.Sprintf("HTTP server Shutdown: %s", err.Error()))
		return err
	}

	s.logger.Info().Msg("HTTP server stopped")
	return nil
}
