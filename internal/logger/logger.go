package logger

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/phsym/console-slog"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"log/slog"
)

var (
	globalLogger *slog.Logger
	once         sync.Once
)

func initLogger() {
	globalLogger = slog.New(
		console.NewHandler(os.Stderr, &console.HandlerOptions{
			Level: slog.LevelDebug,
      AddSource: true,
		}),
	)
}

func ensureLoggerInitialized() {
	once.Do(initLogger)
}

func Info(ctx context.Context, space string, msg string, attrs ...any) {
	ensureLoggerInitialized()
	globalLogger.InfoContext(getContext(ctx), buildMsg(space, msg), attrs...)
}

func Error(ctx context.Context, space string, msg string, attrs ...any) error {
	ensureLoggerInitialized()
	globalLogger.ErrorContext(getContext(ctx), buildMsg(space, msg), attrs...)
	return errors.New(space + msg)
}

func Debug(ctx context.Context, space string, msg string, attrs ...any) {
	ensureLoggerInitialized()
	globalLogger.DebugContext(getContext(ctx), buildMsg(space, msg), attrs...)
}

func Toast(ctx context.Context, space string, msg string, level slog.Level) {
	runtime.EventsEmit(ctx, "log_toast", LogToast{
		timestamp: time.Now(),
		space:     space,
		msg:       msg,
		level:     level,
	})
}

func buildMsg(space string, msg string) string {
	return fmt.Sprintf("[%s] %s", space, msg)
}

func getContext(ctx context.Context) context.Context {
	if ctx == nil {
		return context.Background()
	}

	return ctx
}

type LogToast struct {
	timestamp time.Time
	space     string
	msg       string
	level     slog.Level
}
