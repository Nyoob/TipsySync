package main

import (
	"context"
	"tip-aggregator/internal/config"
	"tip-aggregator/internal/database"
	"tip-aggregator/internal/events"
	"tip-aggregator/internal/providers"
)

// App struct
type App struct {
	ctx context.Context
	db *database.DB
	config *config.Config
	providers []providers.Provider
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	a.db = database.NewDatabase(); 
	a.config = config.NewConfig(a.db, a.onConfigUpdate);
	a.providers = providers.InitializeProviders(a.config)

	a.startHandling()
}

func (a *App) onConfigUpdate(cfg config.Config) {
	for _, provider := range a.providers {
		go provider.Stop()
	}
	a.startHandling()
}

func (a *App) startHandling() {
	eventHandler := events.NewHandler()
	for _, provider := range a.providers {
		go provider.Start(eventHandler)
	}
}

// Public Methods (frontend)
func (a *App) GetConfig() config.Config {
	return *a.config;
}
