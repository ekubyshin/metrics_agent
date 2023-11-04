package logger

import (
	"github.com/ekubyshin/metrics_agent/internal/config"
	"go.uber.org/zap"
)

type Logger interface {
	Sync() error
	Info(args ...any)
}

type DefaultLogger struct {
	Writer *zap.SugaredLogger
	logger *zap.Logger
}

func NewLoggerFromEnv(cfg *config.Config) (*DefaultLogger, error) {
	if cfg.Env == "production" {
		return NewProductionLogger()
	}
	return NewDevelopmentLogger()
}

func NewProductionLogger() (*DefaultLogger, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}
	sugar := logger.Sugar()
	return &DefaultLogger{
		Writer: sugar,
		logger: logger,
	}, nil
}

func NewDevelopmentLogger() (*DefaultLogger, error) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		return nil, err
	}
	sugar := logger.Sugar()
	return &DefaultLogger{
		Writer: sugar,
		logger: logger,
	}, nil
}

func (l *DefaultLogger) Sync() error {
	return l.logger.Sync()
}

func (l *DefaultLogger) Info(args ...any) {
	l.Writer.Infoln(args...)
}
