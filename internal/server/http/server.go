package internalhttp

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/DEMAxx/project_work/internal/filemodifier"
	"github.com/DEMAxx/project_work/internal/lrucache"
	"github.com/DEMAxx/project_work/pkg/config"
	"github.com/rs/zerolog"
)

const TIMEOUT = 5 * time.Second

type Server struct {
	httpServer *http.Server
	logger     *zerolog.Logger
	cache      lrucache.Cache
}

func NewServer(
	logger *zerolog.Logger,
	hostAndPort string,
	cache lrucache.Cache,
	cnf *config.Config,
) *Server {
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

		if path == "" {
			http.Error(w, "Missing URL parameter", http.StatusBadRequest)
			return
		}

		modifier, err := filemodifier.New(
			strings.Split(path, "/"),
			logger,
			cnf,
			cache,
		)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		cachedImage, found := modifier.GetFromCache()

		if found {
			logger.Info().Msg(fmt.Sprintf("Image retrieved from cache"))

			w.Header().Set("Content-Type", "image/jpeg")
			w.WriteHeader(http.StatusOK)
			_, err := w.Write(cachedImage.([]byte))
			if err != nil {
				logger.Error().Msg("Failed to write cached image to response")
			}
			return
		}

		resizedImage, err := modifier.ResizeImage()

		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to modify image: %s", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "image/jpeg")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(resizedImage)
		if err != nil {
			logger.Error().Msg("Failed to write response body")
			return
		}
	}), logger))

	return &Server{
		httpServer: &http.Server{
			Addr:              hostAndPort,
			Handler:           mux,
			ReadHeaderTimeout: TIMEOUT,
		},
		logger: logger,
		cache:  cache,
	}
}

func (s *Server) Start(ctx context.Context) error {
	s.logger.Info().Msg(fmt.Sprintf("Starting HTTP server on %s...", s.httpServer.Addr))

	// Start HTTP server
	go func() {
		s.logger.Info().Msg("HTTP server start...")

		if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.logger.Error().Msg(fmt.Sprintf("HTTP server ListenAndServe: %s", err.Error()))
		}
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
