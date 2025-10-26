package main

import (
	"context"
	"fmt"
	"tip-aggregator/internal/config"
	"tip-aggregator/internal/database"
	"tip-aggregator/internal/events"
	"tip-aggregator/internal/providers"
)

type App struct {
	ctx context.Context
	db *database.DB
	config *config.Config
	providers []providers.Provider
}

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

// shudown to clean up, eg. database close
func (a *App) shutdown(ctx context.Context) {
	a.db.Close()
	for _, provider := range a.providers {
		provider.Stop()
	}
	fmt.Println("TA Cleanup finished!")
}

// this restarts providers when the config is updated, just to make sure everything works for sure
// for example, this resets chaturbate's nextUrl.
func (a *App) onConfigUpdate() {
	for _, provider := range a.providers {
		provider.Stop()
	}
	a.startHandling()
}

func (a *App) startHandling() {
	eventHandler := events.NewHandler(a.ctx)
	for _, provider := range a.providers {
		go provider.Start(eventHandler)
	}
}

// Public Methods (frontend)
func (a *App) GetConfig() config.Config {
	return *a.config;
}

func (a *App) SetProviderSettings(provider string, config config.Provider) {
	a.config.SetProviderSettings(provider, &config)
}
