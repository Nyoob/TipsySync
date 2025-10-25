package config

import (
	"fmt"
	"tip-aggregator/internal/database"
)

type Config struct {
  Providers map[string]Provider
}

type Provider struct {
  Enabled bool
  ApiToken string // 
  FetchInterval int // in seconds
}

func NewConfig(db *database.DB, updateCb func(Config)) *Config {
  defaultConfig := getDefaultConfig()

  // here do the SQL stuff, idea is:
  // * loop through defaultConfig providers, saving them to db if not exists (for first time load)
  // * have a method "setProviderSettings" to set a provider setting in db

  for providerName, providerSettings := range defaultConfig.Providers {
    // setup providers in DB if not exist
		db.Exec(`INSERT OR IGNORE INTO settings_providers (id, enabled, apiToken, fetchInterval) VALUES (?, ?, ?, ?)`, 
      providerName, providerSettings.Enabled, providerSettings.ApiToken, providerSettings.FetchInterval);

    // get from DB and set to local cfg
    newPSettings := defaultConfig.GetProviderSettings(db, providerName)
    if newPSettings != nil {
      defaultConfig.Providers[providerName] = *newPSettings
    }
  }

  return defaultConfig
}

func (c Config) SetProviderSettings(db *database.DB, provider string, settings Provider) {
  c.Providers[provider] = settings

  if db != nil {
    db.Exec(`UPDATE settings_providers
      SET enabled = ?, apiToken = ?, fetchInterval = ?
      WHERE id = ?`, settings.Enabled, settings.ApiToken, settings.FetchInterval, provider)
  }
}

// gets setting from DB
func (c Config) GetProviderSettings(db *database.DB, provider string) *Provider {
	row := db.QueryRow(`SELECT enabled, apiToken, fetchInterval FROM settings_providers WHERE id = ?`, provider);
	var cfg Provider;
	err := row.Scan(&cfg.ApiToken, &cfg.FetchInterval)
	if err != nil {
		fmt.Println("Error getting provider from DB:", err)
		return nil
	}
	return &cfg
}

func getDefaultConfig() *Config {
  return &Config{
    Providers: map[string]Provider{
      "chaturbate": {
        Enabled: false,
        ApiToken: "", // db field
        FetchInterval: 10, // db field
      },
    },
  }
}

