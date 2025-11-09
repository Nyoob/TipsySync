package config

import (
	"context"
	"fmt"
	"tip-aggregator/internal/database"
	"tip-aggregator/internal/logger"
)

type Config struct {
  Settings map[string]string
  Providers map[string]*Provider
  db *database.DB
  updateCallback func(string)
}

type Provider struct {
  Enabled bool
  ApiToken string // 
  FetchInterval int // in seconds
}

func NewConfig(db *database.DB, updateCallback func(string)) *Config {
  defaultConfig := getDefaultConfig()

  defaultConfig.db = db
  defaultConfig.updateCallback = updateCallback

  for id, value := range defaultConfig.Settings {
    // setup setting in DB if not exist
    _, err := db.Exec(`INSERT OR IGNORE INTO settings (id, value) VALUES (?, ?)`, 
      id, value);
    if err != nil { fmt.Println("Error setting default settings: ", err)}

    // get from DB and set to local cfg
    newSettings := defaultConfig.GetSetting(id)
    if newSettings != "" {
      defaultConfig.Settings[id] = newSettings
    }
  }

  for providerName, providerSettings := range defaultConfig.Providers {
    // setup providers in DB if not exist
    _, err := db.Exec(`INSERT OR IGNORE INTO settings_providers (id, enabled, apiToken, fetchInterval) VALUES (?, ?, ?, ?)`, 
      providerName, providerSettings.Enabled, providerSettings.ApiToken, providerSettings.FetchInterval);
    if err != nil { fmt.Println("Error setting default provider settings: ", err)}

    // get from DB and set to local cfg
    newPSettings := defaultConfig.GetProviderSettings(providerName)
    if newPSettings != nil {
      defaultConfig.Providers[providerName] = newPSettings
    }
  }

  return defaultConfig
}

func (c *Config) SetSetting(key string, value any) {
  valueString, ok := value.(string)
  if !ok {
    logger.Error(context.Background(), "Config", "Couldn't set Setting, can't parse to string")
  }
  c.Settings[key] = valueString
  c.updateCallback("")

  if c.db != nil {
    _, err :=  c.db.Exec(`UPDATE settings
      SET value = ?
      WHERE id = ?`, valueString, key)
    if err != nil { fmt.Println("Error setting provider settings: ", err)}
  }
}

func (c *Config) GetSetting(setting string) string {
	row := c.db.QueryRow(`SELECT value FROM settings WHERE id = ?`, setting);
	var value string;
	err := row.Scan(&value)
	if err != nil {
    fmt.Println("Error getting provider from DB:: ", err)
		return ""
	}
	return value
}

func (c *Config) SetProviderSettings(provider string, settings *Provider) {
  *c.Providers[provider] = *settings 
  c.updateCallback(provider)

  if c.db != nil {
    _, err := c.db.Exec(`UPDATE settings_providers
      SET enabled = ?, apiToken = ?, fetchInterval = ?
      WHERE id = ?`, settings.Enabled, settings.ApiToken, settings.FetchInterval, provider)
    if err != nil { fmt.Println("Error setting provider settings: ", err)}
  }
}

// gets setting from DB
func (c *Config) GetProviderSettings(provider string) *Provider {
	row := c.db.QueryRow(`SELECT enabled, apiToken, fetchInterval FROM settings_providers WHERE id = ?`, provider);
	var cfg Provider;
	err := row.Scan(&cfg.Enabled, &cfg.ApiToken, &cfg.FetchInterval)
	if err != nil {
    fmt.Println("Error getting provider from DB:: ", err)
		return nil
	}
	return &cfg
}

func getDefaultConfig() *Config {
  return &Config{
    Settings: map[string]string{
      "eventListMaxItems": "100",
      "socketPort": "6769",
    },
    Providers: map[string]*Provider{
      "chaturbate": {
        Enabled: false,
        ApiToken: "",
        FetchInterval: 10,
      },
      "fansly": {
        Enabled: false,
        ApiToken: "",
        FetchInterval: 0,
      },
      "stripchat": {
        Enabled: false,
        ApiToken: "",
        FetchInterval: 0,
      },
    },
  }
}

