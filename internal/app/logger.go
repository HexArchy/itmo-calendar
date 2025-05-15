package app

import (
	"context"
	"strings"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/hexarchy/itmo-calendar/internal/config"
	"github.com/hexarchy/itmo-calendar/pkg/shutdown"
)

// initLogger configures and initializes the zap logger.
func initLogger(cfg *config.Config) (*zap.Logger, error) {
	level := zap.InfoLevel
	err := level.UnmarshalText([]byte(cfg.Logger.Level))
	if err != nil {
		return nil, errors.Wrap(err, "parse log level")
	}

	stacktraceLevel := zap.ErrorLevel
	err = stacktraceLevel.UnmarshalText([]byte(cfg.Logger.Stacktrace))
	if err != nil {
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

	if !cfg.Logger.Development {
		encoderConfig.EncodeLevel = zapcore.LowercaseLevelEncoder
	} else {
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoderConfig.EncodeCaller = zapcore.FullCallerEncoder
	}

	config := zap.Config{
		Level:             zap.NewAtomicLevelAt(level),
		Development:       cfg.Logger.Development,
		DisableCaller:     false,
		DisableStacktrace: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         cfg.Logger.Encoding,
		EncoderConfig:    encoderConfig,
		OutputPaths:      cfg.Logger.OutputPaths,
		ErrorOutputPaths: cfg.Logger.ErrorOutputPaths,
	}

	if !cfg.Logger.Sampling {
		config.Sampling = nil
	}

	logger, err := config.Build()
	if err != nil {
		return nil, errors.Wrap(err, "build logger")
	}

	logger = logger.With(
		zap.String("app", cfg.App.Name),
		zap.String("environment", cfg.App.Environment),
		zap.String("version", cfg.App.Version),
		zap.String("instance", cfg.App.Instance),
	)

	return logger, nil
}

// gracefulShutdownCallbackZapLogger creates a shutdown callback for zap logger.
func gracefulShutdownCallbackZapLogger(logger *zap.Logger) *shutdown.Callback {
	return &shutdown.Callback{
		Name: "ZapLogger",
		FnCtx: func(ctx context.Context) error {
			logger.Info("Flushing logger buffers")

			done := make(chan struct{})
			go func() {
				err := logger.Sync()
				if err != nil {
					if !strings.Contains(err.Error(), "invalid argument") {
						logger.Error("Failed to flush logger buffers", zap.Error(err))
					}
					logger.Error("Failed to flush logger buffers", zap.Error(err))
				}
				close(done)
			}()

			select {
			case <-done:
				return nil
			case <-ctx.Done():
				logger.Warn("Logger flush timed out")
				return ctx.Err()
			case <-time.After(5 * time.Second):
				logger.Warn("Logger flush timed out after 5 seconds")
				return nil
			}
		},
	}
}
