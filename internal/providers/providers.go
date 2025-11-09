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

func InitializeProviders(cfg *config.Config) map[string]Provider {
  return map[string]Provider {
    "chaturbate": NewChaturbate(cfg.Providers["chaturbate"]),
    "fansly": NewFansly(cfg.Providers["fansly"]),
    "stripchat": NewStripchat(cfg.Providers["stripchat"]),
  }
}

