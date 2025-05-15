package app

import (
	"context"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/hexarchy/itmo-calendar/internal/app/container"
	"github.com/hexarchy/itmo-calendar/internal/config"
	"github.com/hexarchy/itmo-calendar/internal/handlers/http"
	api "github.com/hexarchy/itmo-calendar/internal/handlers/http/v1"
	"github.com/hexarchy/itmo-calendar/pkg/shutdown"
)

type App struct {
	Cfg    *config.Config
	Logger *zap.Logger

	Container  *container.Container
	HTTPServer *http.Server
}

func New(ctx context.Context, cfg *config.Config) (*App, error) {
	app := App{
		Cfg: cfg,
	}
	var err error

	app.Logger, err = initLogger(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "init logger")
	}

	app.Container, err = container.New(ctx, cfg, app.Logger)
	if err != nil {
		return nil, errors.Wrap(err, "new container")
	}

	apiHandler, err := api.NewHandler(&app.Container.UseCases, app.Logger)
	if err != nil {
		return nil, errors.Wrap(err, "new api handler")
	}

	app.HTTPServer, err = http.New(
		app.Container,
		cfg.HTTPServer,
		http.WithAPIHandler(apiHandler),
		http.WithLogger(app.Logger),
	)
	if err != nil {
		return nil, errors.Wrap(err, "new http server")
	}

	// Register shutdown callbacks.
	shutdownCallbacks := make([]*shutdown.Callback, 0)
	shutdownCallbacks = append(shutdownCallbacks, gracefulShutdownCallbackZapLogger(app.Logger))

	for _, cb := range shutdownCallbacks {
		shutdown.AddCallback(cb)
	}

	return &app, nil
}
