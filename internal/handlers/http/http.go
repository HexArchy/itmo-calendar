package http

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/hexarchy/itmo-calendar/internal/app/container"
	"github.com/hexarchy/itmo-calendar/internal/config"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const (
	_defaultServerShutdownTimeout = 10 * time.Second
)

// Server is HTTP server for ITMO Calendar API.
type Server struct {
	server   *http.Server
	logger   *zap.Logger
	config   *config.HTTPServer
	handlers []APIHandler
}

// APIHandler defines the interface for API handlers.
type APIHandler interface {
	AddRoutes(r *mux.Router)
	GetVersion() string
}

// Option defines a functional option for configuring the server.
type Option func(*Server)

// WithLogger sets a custom logger for the server.
func WithLogger(logger *zap.Logger) Option {
	return func(s *Server) {
		s.logger = logger.With(zap.String("component", "http_server"))
	}
}

// WithAPIHandler adds an API handler to the server.
func WithAPIHandler(handler APIHandler) Option {
	return func(s *Server) {
		s.handlers = append(s.handlers, handler)
	}
}

// New creates a new HTTP server with the provided configuration and options.
func New(c *container.Container, cfg *config.HTTPServer, opts ...Option) (*Server, error) {
	s := &Server{
		logger: c.Logger.With(zap.String("component", "http_server")),
		config: &config.HTTPServer{
			Host:         cfg.Host,
			Port:         cfg.Port,
			TLS:          cfg.TLS,
			ReadTimeout:  cfg.ReadTimeout,
			WriteTimeout: cfg.WriteTimeout,
			IdleTimeout:  cfg.IdleTimeout,
			EnableHTTP2:  cfg.EnableHTTP2,
		},
	}

	for _, opt := range opts {
		opt(s)
	}

	if len(s.handlers) == 0 {
		return nil, errors.New("no API handlers provided")
	}

	return s, nil
}

// Start initializes and starts the HTTP server.
func (s *Server) Start() error {
	router := mux.NewRouter()

	for _, handler := range s.handlers {
		version := handler.GetVersion()
		prefix := fmt.Sprintf("/api/%s", version)
		subrouter := router.PathPrefix(prefix).Subrouter()
		handler.AddRoutes(subrouter)
		s.logger.Info("Registered API handler", zap.String("version", version))
	}

	s.registerDebugRoutes(router)

	addr := net.JoinHostPort(s.config.Host, strconv.Itoa(s.config.Port))

	server := &http.Server{
		Addr:         addr,
		Handler:      NewLoggingMiddleware(s.logger)(router),
		ReadTimeout:  s.config.ReadTimeout,
		WriteTimeout: s.config.WriteTimeout,
		IdleTimeout:  s.config.IdleTimeout,
	}

	s.server = server

	s.logger.Info("Starting HTTP server", zap.String("address", addr), zap.Bool("tls", s.config.TLS.Enabled))

	if s.config.TLS.Enabled {
		tlsConfig, err := s.config.TLS.BuildTLSConfig(s.config.Host)
		if err != nil {
			return errors.Wrap(err, "build TLS config")
		}
		server.TLSConfig = tlsConfig

		err = server.ListenAndServeTLS(s.config.TLS.CertFile, s.config.TLS.KeyFile)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return errors.Wrap(err, "server failed to start (TLS)")
		}
	} else {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			return errors.Wrap(err, "server failed to start")
		}
	}

	return nil
}

// Stop gracefully shuts down the HTTP server.
func (s *Server) Stop(ctx context.Context) error {
	s.logger.Info("Stopping HTTP server")

	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), _defaultServerShutdownTimeout)
		defer cancel()
	}

	err := s.server.Shutdown(ctx)
	if err != nil {
		return errors.Wrap(err, "server shutdown failed")
	}

	s.logger.Info("HTTP server stopped")

	return nil
}
