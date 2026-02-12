package config

import (
	"os"
	"strings"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/rawbytes"
	"github.com/knadh/koanf/v2"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
)

// Loader wraps a koanf instance loaded from YAML + env overrides.
type Loader struct {
	k *koanf.Koanf
}

// NewLoader reads configuration from a YAML file specified by --config flag,
// expands ${VAR} environment variables in the file, and applies APP_ prefixed
// env overrides.
func NewLoader() (*Loader, error) {
	k := koanf.New(".")

	fs := pflag.NewFlagSet("config", pflag.ContinueOnError)
	configPath := fs.String("config", "", "Path to configuration file")
	if err := fs.Parse(os.Args[1:]); err != nil {
		return nil, errors.Wrap(err, "parse flags")
	}

	if *configPath == "" {
		return nil, errors.New("--config flag is required")
	}

	data, err := os.ReadFile(*configPath)
	if err != nil {
		return nil, errors.Wrap(err, "read config file")
	}

	if errLoad := k.Load(rawbytes.Provider([]byte(os.ExpandEnv(string(data)))), yaml.Parser()); errLoad != nil {
		return nil, errors.Wrap(errLoad, "load yaml")
	}

	// APP_POSTGRES__POOL__MAX_CONNECTIONS -> postgres.pool.max_connections
	if errEnv := k.Load(env.Provider("APP_", ".", func(s string) string {
		s = strings.TrimPrefix(s, "APP_")
		s = strings.ToLower(s)
		s = strings.ReplaceAll(s, "__", ".")
		return s
	}), nil); errEnv != nil {
		return nil, errors.Wrap(errEnv, "load env overrides")
	}

	return &Loader{k: k}, nil
}

// Unmarshal decodes a config section at the given path into dst.
// Uses yaml struct tags for field mapping.
func (l *Loader) Unmarshal(path string, dst any) error {
	return l.k.UnmarshalWithConf(path, dst, koanf.UnmarshalConf{Tag: "yaml"})
}
