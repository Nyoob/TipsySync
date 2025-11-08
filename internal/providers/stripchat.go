package providers

import (
	"encoding/json"

	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"
	"tip-aggregator/internal/config"
	"tip-aggregator/internal/events"

	"github.com/gorilla/websocket"
)

type Stripchat struct {
	config *config.Provider
	socket *websocket.Conn
	stopCh chan struct{}
}

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
	fmt.Println("Starting Stripchat")
	f.stopCh = make(chan struct{})
	// if !f.config.Enabled || f.config.ApiToken == "" { // TODO uncomment this once frontend done
	// 	fmt.Println("Provider Stripchat not enabled or ChatroomID missing")
	// 	return errors.New("Provider Stripchat not enabled or ChatroomID missing")
	// }

	// f._fake() // INFO: uncomment to fake WS data TODO implement this shit

	u := url.URL{Scheme: "wss", Host: "websocket-sp-v6.stripchat.com", Path: "/connection/websocket"}

	headers := http.Header{}
	headers.Add("Origin", "https://www.stripchat.com")
	headers.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:144.0) Gecko/20100101 Firefox/144.0")

	// connect to socket
	var err error
	f.socket, _, err = websocket.DefaultDialer.Dial(u.String(), headers)
	if err != nil {
		fmt.Println("Error connecting to Stripchat Websocket: ", err)
		return errors.New("Error connecting to Stripchat Websocket")
	}

	// chatRoomId := f.getChatRoomId(f.config.ApiToken) // TODO care about this later
	// chatRoomId := f.config.ApiToken;
	chatRoomId := "222404103"
	if chatRoomId == "" {
		return errors.New("SC ChatroomID is undefined. Username is likely wrong.")
	}

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

	// setup 25sec interval ping
	ticker := time.NewTicker(25 * time.Second)
	defer ticker.Stop()

	go func() {
		select {
		case <-f.stopCh:
			fmt.Println("Stripchat stopped successfully")
			return
		default:
			for range ticker.C {
				err := f.socket.WriteMessage(websocket.TextMessage, []byte("p"))
				if err != nil {
					fmt.Println("Error sending stripchat keepalive")
				}
			}
		}
	}()

	// start handling
	for {
		select {
		case <-f.stopCh:
			fmt.Println("Stopped stripchat")
			return nil
		default:
			_, message, err := f.socket.ReadMessage()
			fmt.Println("SC MSG RECIEVED", string(message))
			fmt.Println("SC ERR RECIEVED", err)
			if err != nil {
				fmt.Println("Stripchat Socket ReadMessage Error, likely called Stop()", err)
				return nil
			}
			// handler(f.GetName(), f.onMessage(message))
		}
	}
}

func (f *Stripchat) Stop() { // TODO: fix goofy ahh nil pointer deref
	select {
	case <-f.stopCh:
		fmt.Println("Stopping Stripchat -> Already stopped")
		return
	default:
		fmt.Println("Stopping Stripchat")
		if f.socket != nil {
			f.socket.Close()
		}
		close(f.stopCh)
	}
}

func (f *Stripchat) onMessage(rawMessage []byte) events.Event {
	var msg stripchatResponseMessage
	err := json.Unmarshal(rawMessage, &msg)
	if err != nil {
		fmt.Println("Error unmarshalling stripchat Message", err)
		return nil
	}

	fmt.Println("STRIPCHAT WS RECIEVED:", msg)

	return nil
}

// func (f *Stripchat) getChatRoomId(username string) string {
// 	resp, err := http.Get("https://apiv3.fansly.com/api/v1/account?usernames=" + username + "&ngsw-bypass=true")
// 	if err != nil {
// 		fmt.Println("Error getting Fansly ChatroomID (fetch) ", err)
// 		return ""
// 	}
// 	defer resp.Body.Close()

// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		return ""
// 	}

// 	var data fanslyAccountResponse
// 	err = json.Unmarshal(body, &data)
// 	if err != nil {
// 		fmt.Println("Error parsing Fansly Account Response", err)
// 		return ""
// 	}

// 	if len(data.Response) == 0 {
// 		fmt.Println("Response Data is empty", data, " (For username: ", username, ")")
// 		return ""
// 	}

// 	return data.Response[0].Streaming.Channel.ChatRoomId
// }

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
						Amount                          int         `json:"amount"`
						Body                            string      `json:"body"`
						FanClubNumberMonthsOfSubscribed int         `json:"fanClubNumberMonthsOfSubscribed"`
						FanClubTier                     interface{} `json:"fanClubTier"`
						IsAnonymous                     bool        `json:"isAnonymous"`
						Source                          string      `json:"source"`
						TipData                         struct {
							TipperKey string `json:"tipperKey"`
						} `json:"tipData"`
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
