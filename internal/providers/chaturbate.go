package providers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"
	"tip-aggregator/internal/config"
	"tip-aggregator/internal/events"
)

type Chaturbate struct {
	config   *config.Provider
	nextUrl  string
	stopChan chan struct{}
}

func NewChaturbate(cfg config.Provider) Chaturbate {
	return Chaturbate{
		config:  &cfg,
		nextUrl: "",
	}
}

func (c Chaturbate) GetName() string {
	return "chaturbate"
}

func (c Chaturbate) Start(handler events.EventHandler) error {
	if !c.config.Enabled || c.config.ApiToken == "" {
		return errors.New("Provider Chaturbate not enabled or API Token missing")
	}
	c.stopChan = make(chan struct{})
	for {
		select {
		case <-c.stopChan:
			return nil
		default:
			handler(c.GetName(), c.fetch())
		}
	}
}

func (c Chaturbate) Stop() {
	if c.stopChan != nil {
		close(c.stopChan)
	}
}

func (c Chaturbate) fetch() events.Event {
	var fetchUrl string
	if c.nextUrl != "" {
		fetchUrl = c.nextUrl
	} else {
		fetchUrl = c.config.ApiToken
	}
	resp, err := http.Get(fetchUrl + "?timeout=" + strconv.Itoa(c.config.FetchInterval))
	if err != nil {
		fmt.Println("Error getting chaturbate data", err)
	}
	var cbResponse cbResponse
	if err := json.NewDecoder(resp.Body).Decode(&cbResponse); err != nil {
		fmt.Println("Error decoding chaturbate data", err)
	}

	c.nextUrl = cbResponse.NextUrl

	for _, event := range cbResponse.Events {
		if event.Method == "tip" {
			var obj cbEventTip
			if err := json.Unmarshal(event.Object, &obj); err != nil {
				return nil
			}
			return events.TipEvent{
				Id: event.Id,
				User: events.User{
					Username:   obj.User.Username,
					Gender:     obj.User.Gender,
					IsMod:      obj.User.IsMod,
					HasTks:     obj.User.HasTokens,
					Subscribed: obj.User.InFanclub,
				},
				TipValue:          obj.Tip.Tokens,
				TipValueInDollars: float64(obj.Tip.Tokens) * 0.05,
				TipMessage:        obj.Tip.Message,
				Timestamp:         time.Now(),
			}
		}

		var obj cbEventGeneric
		if err := json.Unmarshal(event.Object, &obj); err != nil {
			return nil
		}

		user := events.User{
			Username:   obj.User.Username,
			Gender:     obj.User.Gender,
			IsMod:      obj.User.IsMod,
			HasTks:     obj.User.HasTokens,
			Subscribed: obj.User.InFanclub,
		}

		switch event.Method {
		case "fanclubJoin":
			return events.SubscribeEvent{Id: event.Id, User: user, Timestamp: time.Now()}
		case "follow":
			return events.FollowEvent{Id: event.Id, User: user, Timestamp: time.Now()}
		case "unfollow":
			return events.UnfollowEvent{Id: event.Id, User: user, Timestamp: time.Now()}
		}
	}

	return nil
}

type cbResponse struct {
	Events  []cbEvents `json:"events"`
	NextUrl string     `json:"nextUrl"`
}
type cbEvents struct {
	Method string          `json:"method"`
	Id     string          `json:"id"`
	Object json.RawMessage `json:"object"`
}
type cbEvent interface{}

type cbEventTip struct {
	Broadcaster string `json:"broadcaster"`
	Tip         struct {
		Tokens  int    `json:"tokens"`
		IsAnon  bool   `json:"isAnon"`
		Message string `json:"message"`
	} `json:"tip"`
	User cbUser `json:"user"`
}
type cbEventGeneric struct {
	Broadcaster string `json:"broadcaster"`
	User        cbUser `json:"user"`
}

type cbUser struct {
	Username   string `json:"username"`
	InFanclub  bool   `json:"inFanclub"`
	Gender     string `json:"gender"`
	HasTokens  bool   `json:"hasTokens"`
	RecentTips string `json:"recentTips"`
	IsMod      bool   `json:"isMod"`
}
