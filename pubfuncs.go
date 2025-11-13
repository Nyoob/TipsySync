package main

import (
	"context"
	"log/slog"
	"tip-aggregator/internal/config"
	"tip-aggregator/internal/logger"
)

// Public Methods (frontend)
func (a *App) GetInfo() map[string]any {
  name, ok := a.wailsData["name"].(string)
  info, ok := a.wailsData["info"].(map[string]any)
  info["build"] = "development"
  if(BuildNumber != "") {
    info["build"] = BuildNumber
  }
  author, ok := a.wailsData["author"].(map[string]any)
  if !ok {
    return nil
  }
  return map[string]any{
    "name": name,
    "info": info,
    "author": author,
  }
}

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
