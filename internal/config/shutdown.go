package config

import "time"

type Shutdown struct {
	Delay           time.Duration `yaml:"delay"`
	Timeout         time.Duration `yaml:"timeout"`
	CallbackTimeout time.Duration `yaml:"callback_timeout"`
}
