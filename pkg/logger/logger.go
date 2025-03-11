package logger

import (
	"io"
	"os"

	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
	appctx "infinitoon.dev/infinitoon/pkg/context"
)

type Logger struct {
	cfg    LoggerConfig
	logger zerolog.Logger
}

func NewLogger(ctx *appctx.AppContext, cfg LoggerConfig) *Logger {
	switch cfg.Level {
	case LOG_TRACE:
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	case LOG_DEBUG:
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case LOG_INFO:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case LOG_WARNING:
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case LOG_ERROR:
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case LOG_FATAL:
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	}
	var writers []io.Writer
	if ctx.GetStr(appctx.EnvironmentKey) == "development" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		writers = append(writers, zerolog.ConsoleWriter{Out: os.Stderr})
	}

	if cfg.Output == LOG_OUTPUT_FILE {
		writers = append(writers, initLogRotator(&cfg))
	}

	if len(writers) == 0 {
		writers = append(writers, zerolog.ConsoleWriter{Out: os.Stderr})
	}

	logger := zerolog.New(io.MultiWriter(writers...)).With().
		Str("service", ctx.GetStr(appctx.AppNameKey)).
		Str("ENV", ctx.GetStr(appctx.EnvironmentKey)).
		Timestamp().
		Caller().
		Logger()

	log := &Logger{
		cfg:    cfg,
		logger: logger,
	}
	ctx.Set(appctx.LoggerKey, log)
	return log
}

func initLogRotator(cfg *LoggerConfig) *lumberjack.Logger {
	return &lumberjack.Logger{
		Filename:   cfg.FilePath,
		MaxSize:    cfg.MaxSize,    // Max size in megabytes before log is rotated
		MaxBackups: cfg.MaxBackups, // Max number of old log files to retain
		MaxAge:     cfg.MaxAge,     // Max number of days to retain old log files
		Compress:   cfg.Compress,   // Compress old log files
	}

}

func (l *Logger) Debug() *zerolog.Event {
	return l.logger.Debug()
}

func (l *Logger) Info() *zerolog.Event {
	return l.logger.Info()
}

func (l *Logger) Warn() *zerolog.Event {
	return l.logger.Warn()
}

func (l *Logger) Error() *zerolog.Event {
	return l.logger.Error()
}

func (l *Logger) Fatal() *zerolog.Event {
	return l.logger.Fatal()
}

func (l *Logger) Panic() *zerolog.Event {
	return l.logger.Panic()
}
