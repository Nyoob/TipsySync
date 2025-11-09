package providers

import (
	"context"
	"encoding/json"
	"math/rand"
	"regexp"

	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
	"tip-aggregator/internal/config"
	"tip-aggregator/internal/events"
	"tip-aggregator/internal/logger"
)

type Chaturbate struct {
	config  *config.Provider
	nextUrl string
	ctx context.Context
	cancel  context.CancelFunc
	done chan struct{}
}

func NewChaturbate(cfg *config.Provider) *Chaturbate {
	return &Chaturbate{
		config:  cfg,
		nextUrl: "",
	}
}

func (c *Chaturbate) GetName() string {
	return "chaturbate"
}

func (c *Chaturbate) Start(handler events.EventHandler) error {
	if !c.config.Enabled {
		return nil
	}

	logger.Info(context.Background(), "CHATURBATE", "🆕 Starting Service...")
	c.done = make(chan struct{})
	defer close(c.done)
	c.ctx, c.cancel = context.WithCancel(context.Background())

	if c.config.ApiToken == "" {
		return logger.Error(context.Background(), "CHATURBATE", "❌ API Token Missing")
	}

	pattern := `^https:\/\/eventsapi\.chaturbate\.com\/events\/[a-zA-Z0-9]+\/[a-zA-Z0-9]+\/$`
	match, err := regexp.MatchString(pattern, c.config.ApiToken)
	if !match || err != nil {
		fmt.Println("Chaturbate API Token Format wrong.")
		return logger.Error(context.Background(), "CHATURBATE", "❌ API Token Format wrong!")
	}

	for {
		select {
		case <-c.ctx.Done():
			logger.Info(context.Background(), "CHATURBATE", "💯 Stopped successfully")
			return nil
		default:
			handler(c.GetName(), c.fetch())
		}
	}
}

func (c *Chaturbate) Stop() {
	logger.Info(context.Background(), "CHATURBATE", "🛑 Stopping Service...")
	if c.cancel != nil {
		c.cancel()
		<-c.done
	}
}

func (c *Chaturbate) fetch() events.Event {
	var fetchUrl string
	if c.nextUrl != "" {
		fetchUrl = c.nextUrl
	} else {
		fetchUrl = strings.TrimSpace(c.config.ApiToken)
	}
	resp, err := http.Get(fetchUrl + "?timeout=" + strconv.Itoa(c.config.FetchInterval))
	if err != nil {
		fmt.Println("Error getting chaturbate data", err)
		return nil
	}
	var cbResponse cbResponse
	if err := json.NewDecoder(resp.Body).Decode(&cbResponse); err != nil {
		fmt.Println("Error decoding chaturbate data", err)
		return nil
	}
	// cbResponse := c._fakeFetch() // INFO comment code above, uncomment this, to test with a faked cb response

	c.nextUrl = cbResponse.NextUrl

	for _, event := range cbResponse.Events {
		// tip logic
		if event.Method == "tip" {
			var obj cbEventTip
			if err := json.Unmarshal(event.Object, &obj); err != nil {
				return nil
			}

			return events.TipEvent{
				Id:                "cb_" + event.Id,
				User:              c.buildUser(obj.User),
				TipValue:          float64(obj.Tip.Tokens),
				TipValueInDollars: float64(obj.Tip.Tokens) * 0.05,
				TipMessage:        obj.Tip.Message,
				Timestamp:         time.Now(),
			}
		}
		// chat message logic
		if event.Method == "chatMessage" {
			var obj cbEventChatMessage
			if err := json.Unmarshal(event.Object, &obj); err != nil {
				return nil
			}

			return events.ChatMessageEvent{
				Id:          "cb_" + event.Id,
				User:        c.buildUser(obj.User),
				ChatMessage: obj.Message.Message,
				Timestamp:   time.Now(),
			}
		}

		// generic (no additional data, just follow, unfollow, subscribe)
		var obj cbEventGeneric
		if err := json.Unmarshal(event.Object, &obj); err != nil {
			return nil
		}

		user := c.buildUser(obj.User)

		switch event.Method {
		case "fanclubJoin":
			return events.SubscribeEvent{Id: "cb_" + event.Id, TierId: "fanclub", TierName: "Fanclub", User: user, Timestamp: time.Now()}
		case "follow":
			return events.FollowEvent{Id: "cb_" + event.Id, User: user, Timestamp: time.Now()}
		case "unfollow":
			return events.UnfollowEvent{Id: "cb" + event.Id, User: user, Timestamp: time.Now()}
		default:
			return nil
		}
	}

	return nil
}

func (c *Chaturbate) buildUser(user cbUser) events.User {
	tierName := ""
	tierColor := ""
	if user.InFanclub {
		tierName = "Fanclub"
		tierColor = "#009900"
	}
	return events.User{
		Username:            user.Username,
		Gender:              user.Gender,
		IsMod:               user.IsMod,
		HasTks:              user.HasTokens,
		Subscribed:          user.InFanclub,
		SubscribedTierName:  tierName,
		SubscribedTierColor: tierColor,
	}
}

// used for testing without an api token
func (c *Chaturbate) _fakeFetch() cbResponse {
	time.Sleep(10 * time.Second)
	tip := cbEventTip{
		Broadcaster: "tadakonyanko",
		Tip: struct {
			Tokens  int    `json:"tokens"`
			IsAnon  bool   `json:"isAnon"`
			Message string `json:"message"`
		}{
			Tokens:  500,
			IsAnon:  false,
			Message: "Heya! Lorem ipsum dolor sit amet, qui minim labore adipisicing minim sint cillum sint consectetur cupidatat.",
		},
		User: cbUser{
			Gender:     "m",
			IsMod:      false,
			Username:   "nyoob",
			InFanclub:  false,
			HasTokens:  true,
			RecentTips: "many",
		},
	}

	tipBytes, _ := json.Marshal(tip)

	return cbResponse{
		Events: []cbEvents{
			{Method: "tip", Id: fmt.Sprint(rand.Int63()), Object: tipBytes},
		},
		NextUrl: "",
	}
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
type cbEventChatMessage struct {
	Message     cbMessage `json:"message"`
	Broadcaster string    `json:"broadcaster"`
	User        cbUser    `json:"user"`
}
type cbEventGeneric struct {
	Broadcaster string `json:"broadcaster"`
	User        cbUser `json:"user"`
}

type cbMessage struct {
	Color   string `json:"color"`
	BgColor string `json:"bgColor"`
	Message string `json:"message"`
	Font    string `json:"font"`
}

type cbUser struct {
	Username   string `json:"username"`
	InFanclub  bool   `json:"inFanclub"`
	Gender     string `json:"gender"`
	HasTokens  bool   `json:"hasTokens"`
	RecentTips string `json:"recentTips"`
	IsMod      bool   `json:"isMod"`
}
