package providers

import (
	"context"
	"encoding/json"
	"log/slog"
	"strconv"
	"strings"

	"fmt"
	"net/http"
	"net/url"
	"time"
	"tip-aggregator/internal/config"
	"tip-aggregator/internal/events"
	"tip-aggregator/internal/logger"

	"github.com/PuerkitoBio/goquery"
	"github.com/gorilla/websocket"
)

type Stripchat struct {
	config *config.Provider
	socket *websocket.Conn
	ctx    context.Context
	cancel context.CancelFunc
	done   chan struct{}
}

const (
	userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:144.0) Gecko/20100101 Firefox/144.0"
)

func NewStripchat(cfg *config.Provider) *Stripchat {
	return &Stripchat{
		config: cfg,
		socket: nil,
	}
}

func (f *Stripchat) GetName() string {
	return "stripchat"
}

func (f *Stripchat) Start(handler events.EventHandler) error {
	if !f.config.Enabled {
		return nil
	}

	logger.Info(context.Background(), "STRIPCHAT", "🆕 Starting Service...")
	f.done = make(chan struct{})
	defer close(f.done)
	f.ctx, f.cancel = context.WithCancel(context.Background())

	if f.config.ApiToken == "" {
		return logger.Error(context.Background(), "STRIPCHAT", "❌ ChatroomID not set")
	}

	// f._fake() // INFO: uncomment to fake WS data TODO implement this shit

	u := url.URL{Scheme: "wss", Host: "websocket-sp-v6.stripchat.com", Path: "/connection/websocket"}

	headers := http.Header{}
	headers.Add("Origin", "https://www.stripchat.com")
	headers.Add("User-Agent", userAgent)

	// connect to socket
	var err error
	f.socket, _, err = websocket.DefaultDialer.Dial(u.String(), headers)
	if err != nil {
		return logger.Error(context.Background(), "STRIPCHAT", "Error connecting to Stripchat Websocket", err)
	}

	chatRoomId := f.getChatRoomId(f.config.ApiToken)
	// chatRoomId := f.config.ApiToken;
	// chatRoomId := "222404103"
	if chatRoomId == "" {
		return logger.Error(context.Background(), "STRIPCHAT", "ChatroomID is undefined.", err)
	}

	logger.Debug(context.Background(), "STRIPCHAT", "Found ChatroomID", slog.Any("ChatroomID", chatRoomId))

	// init chatroom id
	initializer := fmt.Sprintf(`{"connect":{"token":"eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJzdWIiOiItNDkxOSIsImluZm8iOnsiaXNHdWVzdCI6dHJ1ZSwidXNlcklkIjotNDkxOX19.EdTdrVX-wIxOt122ih2fMcodHZniOlc0YWSxAtDRIwI","name":"js"},"id":1}
  {"subscribe":{"channel":"changeConfigFeature"},"id":2}
  {"subscribe":{"channel":"userBanned@%[1]v"},"id":6}
  {"subscribe":{"channel":"broadcastChanged@%[1]v"},"id":8}
  {"subscribe":{"channel":"modelStatusChanged@%[1]v"},"id":9}
  {"subscribe":{"channel":"groupShow@%[1]v"},"id":11}
  {"subscribe":{"channel":"broadcastSettingsChanged@%[1]v"},"id":12}
  {"subscribe":{"channel":"topicChanged@%[1]v"},"id":13}
  {"subscribe":{"channel":"tipMenuUpdated@%[1]v"},"id":14}
  {"subscribe":{"channel":"goalChanged@%[1]v"},"id":15}
  {"subscribe":{"channel":"userUpdated@%[1]v"},"id":16}
  {"subscribe":{"channel":"interactiveToyStatusChanged@%[1]v"},"id":17}
  {"subscribe":{"channel":"deleteChatMessages@%[1]v"},"id":18}
  {"subscribe":{"channel":"tipMenuLanguageDetected@%[1]v"},"id":19}
  {"subscribe":{"channel":"userBroadcastServerChanged@%[1]v"},"id":20}
  {"subscribe":{"channel":"fanClubUpdated@%[1]v"},"id":21}
  {"subscribe":{"channel":"modelAppUpdated@%[1]v"},"id":22}
  {"subscribe":{"channel":"newKing@%[1]v"},"id":23}
  {"subscribe":{"channel":"newChatMessage@%[1]v"},"id":24}`, chatRoomId)
	f.socket.WriteMessage(websocket.TextMessage, []byte(initializer))

	// start handling
	for {
		select {
		case <-f.ctx.Done():
			logger.Info(context.Background(), "STRIPCHAT", "💯 Stopped successfully")
			return nil
		default:
			_, message, err := f.socket.ReadMessage()
			if err != nil {
				fmt.Println("Stripchat Socket ReadMessage Error, likely called Stop()", err)
				return nil
			}
			if string(message) == "{}" { // pong
				f.socket.WriteMessage(websocket.TextMessage, []byte("{}"))
			} else {
				handler(f.GetName(), f.onMessage(message))
			}
		}
	}
}

func (f *Stripchat) Stop() {
	logger.Info(context.Background(), "FANSLY", "🛑 Stopping Service...")
	if f.socket != nil {
		f.socket.Close()
	}
	if f.cancel != nil {
		f.cancel()
		<-f.done
	}
}

func (f *Stripchat) onMessage(rawMessage []byte) events.Event {
	var msg stripchatResponseMessage
	err := json.Unmarshal(rawMessage, &msg)
	if err != nil {
		fmt.Println("Error unmarshalling stripchat Message", err)
		return nil
	}

	innerMsg := msg.Push.Pub.Data.Message

	var evtUser = events.User{
		Username:            innerMsg.UserData.Username,
		Subscribed:          innerMsg.Details.FanClubTier != "",
		SubscribedTierName:  innerMsg.Details.FanClubTier,
		SubscribedTierColor: "",
		Gender:              "u",
		IsMod:               innerMsg.UserData.IsAdmin,
		HasTks:              false,
		StripchatIsKing:     innerMsg.AdditionalData.IsKing,
		StripchatIsKnight:   innerMsg.AdditionalData.IsKnight,
		StripchatIsUltimate: innerMsg.UserData.IsUltimate,
	}

	if innerMsg.Type == "tip" {
		return events.TipEvent{
			Id:                "sc_" + strconv.FormatInt(innerMsg.ID, 10),
			User:              evtUser,
			TipValue:          float64(innerMsg.Details.Amount),
			TipValueInDollars: float64(innerMsg.Details.Amount) * 0.05,
			TipMessage:        innerMsg.Details.Body,
			Timestamp:         time.Now(),
		}
	}

	if innerMsg.Type == "text" {
		return events.ChatMessageEvent{
			Id:          "sc_" + strconv.FormatInt(innerMsg.ID, 10),
			User:        evtUser,
			ChatMessage: innerMsg.Details.Body,
			Timestamp:   time.Now(),
		}
	}

	return nil
}

func (f *Stripchat) getChatRoomId(username string) string {
	req, err := http.NewRequest("GET", "https://www.stripchat.com/"+username, nil)
	if err != nil {
		logger.Error(context.Background(), "STRIPCHAT", "Error creating request", slog.Any("Error", err))
		return ""
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Referer", "https://google.com")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Error(context.Background(), "STRIPCHAT", "Error fetching page for chatroomID", slog.Any("Error", err))
		return ""
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		logger.Error(context.Background(), "STRIPCHAT", "Error reading doc for chatroomID", slog.Any("Error", err))
		return ""
	}

	content, exists := doc.Find(`meta[property="og:image"]`).Attr("content")
	if !exists {
		logger.Error(context.Background(), "STRIPCHAT", "Couldn't find og:image in doc")
		return ""
	}

	parts := strings.Split(content, "/")
	uid := parts[len(parts)-1]
	if uid == "" { // handle trailing slash if exist
		uid = parts[len(parts)-2]
	}

	return uid
}

type stripchatResponseMessage struct {
	Push struct {
		Channel string `json:"channel"`
		Pub     struct {
			Data struct {
				AdditionalData struct {
				} `json:"additionalData"`
				Message struct {
					AdditionalData struct {
						IsKing            bool `json:"isKing"`
						IsKnight          bool `json:"isKnight"`
						IsStudioAdmin     bool `json:"isStudioAdmin"`
						IsStudioModerator bool `json:"isStudioModerator"`
					} `json:"additionalData"`
					CacheID   string    `json:"cacheId"`
					CreatedAt time.Time `json:"createdAt"`
					Details   struct {
						Amount                          int    `json:"amount"`
						Body                            string `json:"body"`
						FanClubNumberMonthsOfSubscribed int    `json:"fanClubNumberMonthsOfSubscribed"`
						FanClubTier                     string `json:"fanClubTier"`
						IsAnonymous                     bool   `json:"isAnonymous"`
						Source                          string `json:"source"`
					} `json:"details"`
					ID       int64  `json:"id"`
					ModelID  int    `json:"modelId"`
					Type     string `json:"type"`
					UserData struct {
						HasAdminBadge        bool `json:"hasAdminBadge"`
						HasVrDevice          bool `json:"hasVrDevice"`
						ID                   int  `json:"id"`
						IsAdmin              bool `json:"isAdmin"`
						IsBlocked            bool `json:"isBlocked"`
						IsDeleted            bool `json:"isDeleted"`
						IsExGreen            bool `json:"isExGreen"`
						IsGreen              bool `json:"isGreen"`
						IsModel              bool `json:"isModel"`
						IsOnline             bool `json:"isOnline"`
						IsPermanentlyBlocked bool `json:"isPermanentlyBlocked"`
						IsRegular            bool `json:"isRegular"`
						IsStudio             bool `json:"isStudio"`
						IsSupport            bool `json:"isSupport"`
						IsUltimate           bool `json:"isUltimate"`
						UserRanking          struct {
							IsEx   bool   `json:"isEx"`
							League string `json:"league"`
							Level  int    `json:"level"`
						} `json:"userRanking"`
						Username string `json:"username"`
					} `json:"userData"`
				} `json:"message"`
			} `json:"data"`
			Offset int `json:"offset"`
		} `json:"pub"`
	} `json:"push"`
}
