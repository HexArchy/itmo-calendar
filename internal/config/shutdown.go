package config

import "time"

type Shutdown struct {
	Delay           time.Duration `path:"delay" default:"5s"`
	Timeout         time.Duration `path:"timeout" default:"30s"`
	CallbackTimeout time.Duration `path:"callback_timeout" default:"10s"`
}
