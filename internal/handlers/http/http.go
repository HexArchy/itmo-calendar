package http

import (
	"context"
	"net"
	"net/http"
	"net/http/pprof"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/hexarchy/itmo-calendar/internal/handlers/http/v1/gen"

	api "github.com/hexarchy/itmo-calendar/internal/handlers/http/v1"
)

const _defaultServerShutdownTimeout = 10 * time.Second

// Server is the HTTP server for ITMO Calendar API.
type Server struct {
	server *http.Server
	logger *zap.Logger
	config *Config
}

// NewServer creates a new HTTP server with ogen handler, Scalar docs and debug routes.
func NewServer(ogenSrv *gen.Server, cfg *Config, logger *zap.Logger) *Server {
	mux := http.NewServeMux()

	// ogen API routes (already prefixed with /api/v1)
	mux.Handle("/api/v1/", ogenSrv)

	// OpenAPI spec and documentation
	mux.Handle("/openapi.yaml", api.SpecHandler())
	mux.Handle("/docs", api.ScalarDocsHandler())

	// Debug/pprof routes
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	mux.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
	mux.Handle("/debug/pprof/heap", pprof.Handler("heap"))
	mux.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
	mux.Handle("/debug/pprof/block", pprof.Handler("block"))
	mux.Handle("/debug/pprof/allocs", pprof.Handler("allocs"))
	mux.Handle("/debug/pprof/mutex", pprof.Handler("mutex"))

	srvLogger := logger.With(zap.String("component", "http_server"))

	addr := net.JoinHostPort(cfg.Host, strconv.Itoa(cfg.Port))

	httpSrv := &http.Server{
		Addr:         addr,
		Handler:      NewLoggingMiddleware(srvLogger)(mux),
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	return &Server{
		server: httpSrv,
		logger: srvLogger,
		config: cfg,
	}
}

// Start starts the HTTP server (blocking).
func (s *Server) Start() error {
	tlsEnabled := s.config.TLS != nil && s.config.TLS.Enabled
	s.logger.Info("Starting HTTP server", zap.String("address", s.server.Addr), zap.Bool("tls", tlsEnabled))

	if tlsEnabled {
		tlsConfig, err := s.config.TLS.BuildTLSConfig(s.config.Host)
		if err != nil {
			return errors.Wrap(err, "build TLS config")
		}
		s.server.TLSConfig = tlsConfig

		err = s.server.ListenAndServeTLS(s.config.TLS.CertFile, s.config.TLS.KeyFile)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return errors.Wrap(err, "server failed to start (TLS)")
		}
	} else {
		err := s.server.ListenAndServe()
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
