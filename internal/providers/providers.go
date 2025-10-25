package providers

import (
	"tip-aggregator/internal/config"
	"tip-aggregator/internal/events"
)

type Provider interface {
  Start(events.EventHandler) error
  Stop()
  GetName() string
}

func InitializeProviders(cfg *config.Config) []Provider {
  return []Provider {
	  NewChaturbate(cfg.Providers["chaturbate"]),
  }
}

