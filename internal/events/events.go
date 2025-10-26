package events

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type EventHandler func(providerName string, event Event) error

func NewHandler(ctx context.Context) EventHandler {
	return func(providerName string, event Event) error {
		// TODO: add an error param that emits an error event to display errors to the user
		if event == nil {
			return errors.New("Eventhandler saw empty event")
		}
		wrappedEvent := WrappedEvent{
			Provider: providerName,
			Type:     event.EventType(),
			Event:    event,
		}
		fmt.Println("handling event!!!")
		// send to JS
		runtime.EventsEmit(ctx, "platform_event", wrappedEvent)
		// TODO: send to websocket handler

		return nil
	}
}

type WrappedEvent struct {
	Provider string
	Type     string
	Event    Event
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
