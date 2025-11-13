package logger

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
	"tip-aggregator/internal/helpers"

	"log/slog"

	"github.com/phsym/console-slog"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

var (
	globalLogger *slog.Logger
	once         sync.Once
	logDir       = helpers.GetConfigDir() + "/logs"
)

func initLogger() {
	globalLogger = slog.New(
		console.NewHandler(os.Stderr, &console.HandlerOptions{
			Level:     slog.LevelDebug,
			AddSource: true,
		}),
	)
	go cleanupOldLogs();
}

func ensureLoggerInitialized() {
	once.Do(initLogger)
}

func logWithLevel(ctx context.Context, level slog.Level, space, msg string, attrs ...any) {
	ensureLoggerInitialized()

	builtMsg := fmt.Sprintf("[%s] %s", space, msg)
	actualCtx := getContext(ctx)

	switch level {
	case slog.LevelInfo:
		globalLogger.InfoContext(actualCtx, builtMsg, attrs...)
	case slog.LevelError:
		globalLogger.ErrorContext(actualCtx, builtMsg, attrs...)
	case slog.LevelDebug:
		globalLogger.DebugContext(actualCtx, builtMsg, attrs...)
	}

	go saveLogToFile(level, space, msg, attrs...)
}

func Info(ctx context.Context, space string, msg string, attrs ...any) {
	ensureLoggerInitialized()
	logWithLevel(ctx, slog.LevelInfo, space, msg, attrs...)
}

func Error(ctx context.Context, space string, msg string, attrs ...any) error {
	ensureLoggerInitialized()
	logWithLevel(ctx, slog.LevelError, space, msg, attrs...)
	return errors.New(space + msg)
}

func Debug(ctx context.Context, space string, msg string, attrs ...any) {
	ensureLoggerInitialized()
	logWithLevel(ctx, slog.LevelDebug, space, msg, attrs...)
}

func Toast(ctx context.Context, space string, msg string, level slog.Level) {
	logWithLevel(ctx, level, space, msg)
	runtime.EventsEmit(ctx, "log_toast", LogToast{
		timestamp: time.Now(),
		space:     space,
		msg:       msg,
		level:     level, // LevelDebug should only print to js console, rest should toast.
	})
}

func getContext(ctx context.Context) context.Context {
	if ctx == nil {
		return context.Background()
	}

	return ctx
}

func saveLogToFile(level slog.Level, space, msg string, attrs ...any) {
	os.MkdirAll(logDir, 0755)

	logLine := fmt.Sprintf("%s [%s] [%s] %s\n",
		time.Now().Format(time.RFC3339),
		level.String(),
		space,
		fmt.Sprint(msg),
	)

	// Append the log line to today's log file
	logFileName := filepath.Join(logDir, time.Now().Format("2006-01-02")+".log")
	f, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open log file: %v\n", err)
		return
	}
	defer f.Close()

	_, err = f.WriteString(logLine)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write log to file: %v\n", err)
	}
}

func cleanupOldLogs() {
	files, err := os.ReadDir(logDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read log directory: %v\nLikely first time startup.", err)
		return
	}

	cutoff := time.Now().AddDate(0, 0, -3) // 3 days ago

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		fi, err := file.Info()
		if err != nil {
			continue
		}
		if fi.ModTime().Before(cutoff) {
			os.Remove(filepath.Join(logDir, file.Name()))
		}
	}
}

type LogToast struct {
	timestamp time.Time
	space     string
	msg       string
	level     slog.Level
}
