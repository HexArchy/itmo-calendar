package main

import (
	"go.uber.org/fx"

	"github.com/hexarchy/itmo-calendar/internal/app"
)

func main() {
	fx.New(app.Module).Run()
}
