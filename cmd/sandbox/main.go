package main

import (
	"context"
	"log"

	"github.com/pkg/errors"

	"github.com/hexarchy/itmo-calendar/internal/config"
	configcore "github.com/hexarchy/itmo-calendar/pkg/config"
)

func main() {
	err := run()
	if err != nil {
		log.Fatalf("error running app: %v\n", err)
	}
}

func run() error {
	cfg := &config.Config{}
	err := configcore.Init(cfg)
	if err != nil {
		log.Fatal("Fail to load config: ", err)
	}

	ctx := context.Background()
	sandbox, err := NewSandbox(ctx, cfg)
	if err != nil {
		return errors.Wrap(err, "new sandbox")
	}

	err = exec(ctx, sandbox, sandbox.Logger)
	if err != nil {
		return errors.Wrap(err, "exec")
	}

	return nil
}
