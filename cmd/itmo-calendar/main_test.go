package main

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/fx"

	"github.com/hexarchy/itmo-calendar/internal/app"
)

func TestAppAssembly(t *testing.T) {
	err := fx.ValidateApp(app.Module)
	require.NoError(t, err)
}
