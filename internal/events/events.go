package events

import (
	"context"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type EventHandler func(providerName string, event Event) error

func NewHandler() EventHandler {
	return func(providerName string, event Event) error {
		// todo: send to JS
		runtime.EventsEmit(context.Background(), event.EventType(), event)
		// todo: send to websocket handler
		return nil
	}
}

// event types
type Event interface {
	EventType() string
}

type TipEvent struct {
	Id                string
	User              User
	TipValue          int
	TipValueInDollars float64
	TipMessage        string
	Timestamp         time.Time
}

func (t TipEvent) EventType() string { return "tip" }

type FollowEvent struct {
	Id        string
	User      User
	Timestamp time.Time
}

func (f FollowEvent) EventType() string { return "follow" }

type UnfollowEvent struct {
	Id        string
	User      User
	Timestamp time.Time
}

func (f UnfollowEvent) EventType() string { return "unfollow" }

type SubscribeEvent struct { // eg. cb fanclub
	Id        string
	User      User
	Timestamp time.Time
}

func (s SubscribeEvent) EventType() string { return "subscribe" }

type User struct {
	Username   string
	Subscribed bool
	Gender     string // m(ale), f(emale), t(rans), c(ouple)
	HasTks     bool
	IsMod      bool
}
