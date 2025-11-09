package main

import (
	"context"
	"log/slog"
	"tip-aggregator/internal/config"
	"tip-aggregator/internal/database"
	"tip-aggregator/internal/events"
	"tip-aggregator/internal/logger"
	"tip-aggregator/internal/providers"
	"tip-aggregator/internal/socket"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx          context.Context
	db           *database.DB
	socket       *socket.Socket
	config       *config.Config
	providers    map[string]providers.Provider
	eventHandler events.EventHandler
}

func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	a.db = database.NewDatabase()
	a.config = config.NewConfig(a.db, a.onConfigUpdate)
	a.providers = providers.InitializeProviders(a.config)

	a.socket = socket.NewSocket(":" + a.config.Settings["socketPort"])

	a.eventHandler = events.NewHandler(a.ctx, a.socket)

	a.startHandling()
}

// shudown to clean up, eg. database close
func (a *App) shutdown(ctx context.Context) {
	a.db.Close()
	for _, provider := range a.providers {
		provider.Stop()
	}
	logger.Info(context.Background(), "APP", "Cleanup finished.")
}

// this restarts a specific provider & also sends an event to frontend with most recent config
// for example, this resets chaturbate's nextUrl.
// it also handles enable/disable, since only .Start() checks for enabled bool
func (a *App) onConfigUpdate(provider string) {
	logger.Info(context.Background(), "APP", "Config updated")
	// send new cfg to frontend app
	runtime.EventsEmit(a.ctx, "config_update", a.config)

	// restart provider if required
	if provider != "" {
		go func() {
			a.providers[provider].Stop()
			a.providers[provider].Start(a.eventHandler)
		}()
	}
}

func (a *App) startHandling() {
	for _, provider := range a.providers {
		go provider.Start(a.eventHandler)
	}
}

// Public Methods (frontend)
func (a *App) GetConfig() config.Config {
	return *a.config
}

func (a *App) SetSettings(settings map[string]any) {
	for key, value := range settings {
		logger.Debug(context.Background(), "APP", "Setting a setting", slog.String("Key", key), slog.Any("Value", value))
		a.config.SetSetting(key, value)
	}
}
func (a *App) SetProviderSettings(provider string, config config.Provider) {
	a.config.SetProviderSettings(provider, &config)
}
