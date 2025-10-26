package providers

import (
	"encoding/json"
	"math/rand"
	"regexp"

	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
	"tip-aggregator/internal/config"
	"tip-aggregator/internal/events"
)

type Chaturbate struct {
	config   *config.Provider
	nextUrl  string
	stopChan chan struct{}
}

func NewChaturbate(cfg *config.Provider) *Chaturbate {
	return &Chaturbate{
		config:  cfg,
		nextUrl: "",
	}
}

func (c Chaturbate) GetName() string {
	return "chaturbate"
}

func (c Chaturbate) Start(handler events.EventHandler) error {
	if !c.config.Enabled || c.config.ApiToken == "" {
		fmt.Println("Provider Chaturbate not enabled or API Token missing")
		return errors.New("Provider Chaturbate not enabled or API Token missing")
	}
	pattern := `^https:\/\/eventsapi\.chaturbate\.com\/events\/[a-zA-Z0-9]+\/[a-zA-Z0-9]+\/$`
	match, err := regexp.MatchString(pattern, c.config.ApiToken)
	if !match || err != nil {
		fmt.Println("Chaturbate API Token Format wrong.")
		return errors.New("Chaturbate API Token Format wrong.")
	}

	c.stopChan = make(chan struct{})
	for {
		select {
		case <-c.stopChan:
			fmt.Println("CB Stopped")
			return nil
		default:
			handler(c.GetName(), c.fetch())
		}
	}
}

func (c Chaturbate) Stop() {
	c.nextUrl = ""
	if c.stopChan != nil {
		close(c.stopChan)
	}
}

func (c Chaturbate) fetch() events.Event {
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
		default:
			return nil
		}
	}

	return nil
}

func (c Chaturbate) _fakeFetch() cbResponse {
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
