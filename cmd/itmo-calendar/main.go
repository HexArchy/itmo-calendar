package main

import (
	"context"
	"log"

	"github.com/hexarchy/itmo-calendar/internal/app"
	"github.com/hexarchy/itmo-calendar/internal/config"
	configcore "github.com/hexarchy/itmo-calendar/pkg/config"
	"github.com/hexarchy/itmo-calendar/pkg/shutdown"
)

func main() {
	ctx := shutdown.WithContext(context.Background())

	cfg := &config.Config{}
	err := configcore.Init(cfg)
	if err != nil {
		log.Fatal("Fail to load config: ", err)
	}

	application, err := app.New(ctx, cfg)
	if err != nil {
		log.Fatal("Fail to create app: ", err)
	}

	err = application.Start(ctx)
	if err != nil {
		log.Fatal("Fail to start app: ", err)
	}
}
