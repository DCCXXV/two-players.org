package logger

import (
	"log/slog"
	"os"
	"time"

	"github.com/lmittmann/tint"
)

var appLogger *slog.Logger

func Init() {
	level := slog.LevelDebug

	handlerOptions := &tint.Options{
		Level:      level,
		TimeFormat: time.Kitchen,
		AddSource:  false,
		NoColor:    false,
	}

	handler := tint.NewHandler(os.Stdout, handlerOptions)
	appLogger = slog.New(handler)

	slog.SetDefault(appLogger)
	appLogger.Info("Logger initialized", "level", level.String(), "colors", true)
}

func Get() *slog.Logger {
	if appLogger == nil {
		slog.Error("Application logger not initialized, using basic slog. Call logger.Init() first.")
		Init()
	}
	return appLogger
}
