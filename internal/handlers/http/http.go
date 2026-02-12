package http

import (
	"context"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
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
func NewServer(ogenSrv *gen.Server, calHandler *CalHandler, cfg *Config, logger *zap.Logger) *Server {
	srvLogger := logger.With(zap.String("component", "http_server"))

	r := chi.NewRouter()

	// Global middleware stack
	r.Use(middleware.RealIP)
	r.Use(middleware.RequestID)
	r.Use(NewLoggingMiddleware(srvLogger))
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type", "Authorization", "X-Request-ID"},
	}))

	// Rate-limited routes
	r.Group(func(r chi.Router) {
		r.Use(httprate.LimitByIP(cfg.RateLimitRPM, time.Minute))
		r.Mount("/api/v1", ogenSrv)
		r.Handle("/cal", calHandler)
	})

	// Documentation (no rate limiting)
	r.Handle("/openapi.yaml", api.SpecHandler())
	r.Handle("/docs", api.ScalarDocsHandler())

	// Debug/pprof routes (only if enabled)
	if cfg.PprofEnabled {
		r.Mount("/debug", middleware.Profiler())
	}

	addr := net.JoinHostPort(cfg.Host, strconv.Itoa(cfg.Port))

	httpSrv := &http.Server{
		Addr:         addr,
		Handler:      r,
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
