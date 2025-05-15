package logger

import (
	"strings"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	_defaultSamplingInitial    = 100
	_defaultSamplingThereafter = 100
	_syncTimeout               = 5 * time.Second
)

// New creates a configured zap.Logger.
func New(cfg *Config, app *AppInfo) (*zap.Logger, error) {
	level := zap.InfoLevel
	if err := level.UnmarshalText([]byte(cfg.Level)); err != nil {
		return nil, errors.Wrap(err, "parse log level")
	}

	stacktraceLevel := zap.ErrorLevel
	if err := stacktraceLevel.UnmarshalText([]byte(cfg.Stacktrace)); err != nil {
		return nil, errors.Wrap(err, "parse stacktrace level")
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	if cfg.Development {
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoderConfig.EncodeCaller = zapcore.FullCallerEncoder
	}

	zapCfg := zap.Config{
		Level:             zap.NewAtomicLevelAt(level),
		Development:       cfg.Development,
		DisableCaller:     false,
		DisableStacktrace: false,
		Sampling: &zap.SamplingConfig{
			Initial:    _defaultSamplingInitial,
			Thereafter: _defaultSamplingThereafter,
		},
		Encoding:         cfg.Encoding,
		EncoderConfig:    encoderConfig,
		OutputPaths:      cfg.OutputPaths,
		ErrorOutputPaths: cfg.ErrorOutputPaths,
	}

	if !cfg.Sampling {
		zapCfg.Sampling = nil
	}

	l, err := zapCfg.Build()
	if err != nil {
		return nil, errors.Wrap(err, "build logger")
	}

	l = l.With(
		zap.String("app", app.Name),
		zap.String("environment", app.Environment),
		zap.String("version", app.Version),
		zap.String("instance", app.Instance),
	)

	return l, nil
}

// Sync flushes logger buffers with a timeout.
func Sync(l *zap.Logger) {
	done := make(chan struct{})
	go func() {
		err := l.Sync()
		if err != nil && !strings.Contains(err.Error(), "invalid argument") {
			l.Error("Failed to flush logger buffers", zap.Error(err))
		}
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(_syncTimeout):
		l.Warn("Logger flush timed out after 5 seconds")
	}
}
