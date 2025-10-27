package providers

import (
	"encoding/json"
	"io"

	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"
	"tip-aggregator/internal/config"
	"tip-aggregator/internal/events"

	"github.com/gorilla/websocket"
)

const (
	MESSAGE_TYPE_SINGLE     = 10000
	MESSAGE_TYPE_BULK       = 10001
	EVENT_TYPE_CHAT         = 10
	EVENT_TYPE_SUBSCRIPTION = 53
	ATTACHMENT_TYPE_TIP     = 7
)

type Fansly struct {
	config *config.Provider
	socket *websocket.Conn
	stopCh chan struct{}
}

func NewFansly(cfg *config.Provider) *Fansly {
	return &Fansly{
		config: cfg,
		socket: nil,
	}
}

func (f *Fansly) GetName() string {
	return "fansly"
}

func (f *Fansly) Start(handler events.EventHandler) error {
	fmt.Println("Starting Fansly")
	f.stopCh = make(chan struct{})

	if !f.config.Enabled || f.config.ApiToken == "" {
		fmt.Println("Provider Fansly not enabled or ChatroomID missing")
		return errors.New("Provider Fansly not enabled or ChatroomID missing")
	}

	// f._fake() // INFO: uncomment to fake WS data TODO implement this shit

	u := url.URL{Scheme: "wss", Host: "chatws.fansly.com", Path: "/?v=3"}

	// connect to socket
	var err error
	f.socket, _, err = websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		fmt.Println("Error connecting to Fansly Websocket")
	}

	chatRoomId := f.getChatRoomId(f.config.ApiToken)
	if chatRoomId == "" {
		return errors.New("ChatroomID is undefined. Username is likely wrong.")
	}

	// init chatroom id
	initializer := `{"t":46001,"d":"{\"chatRoomId\":\"` + chatRoomId + `\",\"v\":3}"}`
	f.socket.WriteMessage(websocket.TextMessage, []byte(initializer))

	// setup 22sec interval ping
	ticker := time.NewTicker(22 * time.Second)
	defer ticker.Stop()

	go func() {
		select {
		case <-f.stopCh:
			fmt.Println("Fansly Stopped successfully")
			return
		default:
			for range ticker.C {
				err := f.socket.WriteMessage(websocket.TextMessage, []byte("p"))
				if err != nil {
					fmt.Println("Error sending fansly keepalive")
				}
			}
		}
	}()

	// start handling
	for {
		select {
		case <-f.stopCh:
			fmt.Println("Stopped Fansly")
			return nil
		default:
			_, message, err := f.socket.ReadMessage()
			if err != nil {
				fmt.Println("Fansly Socket ReadMessage Error, likely called Stop()", err)
				return nil
			}
			handler(f.GetName(), f.onMessage(message))
		}
	}
}

func (f *Fansly) Stop() { // TODO: fix goofy ahh nil pointer deref
	select {
	case <-f.stopCh:
		fmt.Println("Stopping Fansly -> Already stopped")
		return
	default:
		fmt.Println("Stopping Fansly")
		f.socket.Close()
		close(f.stopCh)
	}
}

func (f *Fansly) onMessage(rawMessage []byte) events.Event {
	var msg fanslyResponseMessage
	err := json.Unmarshal(rawMessage, &msg)
	if err != nil {
		fmt.Println("Error unmarshalling fansly Message", err)
		return nil
	}

	if msg.Type != MESSAGE_TYPE_BULK && msg.Type != MESSAGE_TYPE_SINGLE {
		fmt.Println("Fansly: Unknown message type")
		return nil
	}

	if msg.Type == MESSAGE_TYPE_BULK { // bulk event
		// TODO handle bulk events
		fmt.Println("BULK EVENT RECEIVED FANSLY")
	}

	f.parseMsg(&msg)

	var event = msg.InnerMessage.Event

	// build chat & tip
	if event.Type == EVENT_TYPE_CHAT {
		// build user
		var subTierName string
		subscriptionInfo, hasSubscriptionInfo := event.ChatRoomMessage.Metadata["senderSubscription"].(map[string]interface{})
		if hasSubscriptionInfo {
			subTierName = subscriptionInfo["tierName"].(string)
		}

		username := event.ChatRoomMessage.Displayname
		if event.ChatRoomMessage.Username != event.ChatRoomMessage.Displayname {
			username += " (" + event.ChatRoomMessage.Username + ")"
		}

		var evtUser = events.User{
			Username:            username,
			Subscribed:          subTierName != "",
			SubscribedTierName:  subTierName,
			SubscribedTierColor: "",
			Gender:              "u",
			IsMod:               event.ChatRoomMessage.Metadata["senderIsStaff"].(bool),
			HasTks:              false,
		}

		jsts := int64(event.ChatRoomMessage.CreatedAt)
		timestamp := time.Unix(jsts/1000, jsts%1000*int64(time.Millisecond))

		// tip
		hasTip, tip := f.attachmentHasTip(event.ChatRoomMessage.Attachments)
		if hasTip {
			return events.TipEvent{
				Id:                event.ChatRoomMessage.ID,
				User:              evtUser,
				TipValue:          tip / 1000,
				TipValueInDollars: tip / 1000,
				TipMessage:        event.ChatRoomMessage.Content,
				Timestamp:         timestamp,
			}
		} else { // regular chat msg, no tip
			return events.ChatMessageEvent{
				Id:          event.ChatRoomMessage.ID,
				ChatMessage: event.ChatRoomMessage.Content,
				User:        evtUser,
				Timestamp:   timestamp,
			}
		}
	}

	if event.Type == EVENT_TYPE_SUBSCRIPTION {
		username := event.SubAlert.Displayname
		if event.SubAlert.Username != event.SubAlert.Displayname {
			username += " (" + event.SubAlert.Username + ")"
		}

		return events.SubscribeEvent{
			Id:       event.SubAlert.ID,
			TierId:   event.SubAlert.SubscriptionTierID,
			TierName: event.SubAlert.SubscriptionTierName,
			User: events.User{
				Username:            username,
				Subscribed:          true,
				SubscribedTierName:  event.SubAlert.SubscriptionTierName,
				SubscribedTierColor: event.SubAlert.SubscriptionTierColor,
				IsMod:               false,
				HasTks:              false,
				Gender:              "u",
			},
			Timestamp: time.Now(),
		}
	}

	return nil
}

func (f *Fansly) parseMsg(msg *fanslyResponseMessage) {
	err := json.Unmarshal([]byte(msg.RawInnerMessage), &msg.InnerMessage)
	if err != nil {
		fmt.Println("Fansly: Error unmarshalling inner message ", err)
	}

	// event
	err = json.Unmarshal([]byte(msg.InnerMessage.RawEvent), &msg.InnerMessage.Event)
	if err != nil {
		fmt.Println("Fansly: Error unmarshalling msg.InnerMessage.Event ", err)
	}

	if msg.InnerMessage.Event.Type == EVENT_TYPE_CHAT {
		// metadata
		err = json.Unmarshal([]byte(msg.InnerMessage.Event.ChatRoomMessage.RawMetadata), &msg.InnerMessage.Event.ChatRoomMessage.Metadata)
		if err != nil {
			fmt.Println("Fansly: Error unmarshalling chatRoomMessage Meta ", err)
		}

		// attachments
		for i := range msg.InnerMessage.Event.ChatRoomMessage.Attachments {
			err = json.Unmarshal([]byte(msg.InnerMessage.Event.ChatRoomMessage.Attachments[i].RawMetadata), &msg.InnerMessage.Event.ChatRoomMessage.Attachments[i].Metadata)
			if err != nil {
				fmt.Println("Fansly: Error unmarshalling attachment Meta ", err)
			}
		}
	}
}

func (f *Fansly) getChatRoomId(username string) string {
	resp, err := http.Get("https://apiv3.fansly.com/api/v1/account?usernames=" + username + "&ngsw-bypass=true")
	if err != nil {
		fmt.Println("Error getting Fansly ChatroomID (fetch) ", err)
		return ""
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ""
	}

	var data fanslyAccountResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println("Error parsing Fansly Account Response", err)
		return ""
	}

	if len(data.Response) == 0 {
		fmt.Println("Response Data is empty", data, " (For username: ", username, ")")
		return ""
	}

	return data.Response[0].Streaming.Channel.ChatRoomId
}

func (f *Fansly) attachmentHasTip(attachments []fanslyResponseAttachment) (bool, float64) {
	for _, attachment := range attachments {
		if attachment.ContentType == 7 {
			return true, attachment.Metadata["amount"].(float64)
		}
	}
	return false, 0
}

type fanslyAccountResponse struct {
	Response []struct {
		Streaming struct {
			Channel struct {
				ChatRoomId string `json:"chatRoomId"`
			} `json:"channel"`
		} `json:"streaming"`
	} `json:"response"`
}

type fanslyResponseMessage struct {
	Type            int    `json:"t"`
	RawInnerMessage string `json:"d"`
	InnerMessage    fanslyResponseInnerMessage
}
type fanslyResponseInnerMessage struct {
	ServiceId int    `json:"serviceId"`
	RawEvent  string `json:"event"`
	Event     fanslyResponseEvent
}

type fanslyResponseEvent struct {
	Type            int                           `json:"type"`
	ChatRoomMessage fanslyResponseChatRoomMessage `json:"chatRoomMessage"`
	SubAlert        fanslyResponseSubAlert        `json:"subAlert"`
}

type fanslyResponseChatRoomMessage struct {
	ChatRoomID        string                     `json:"chatRoomId"`
	SenderID          string                     `json:"senderId"`
	Content           string                     `json:"content"`
	Type              int                        `json:"type"`
	Private           int                        `json:"private"`
	Attachments       []fanslyResponseAttachment `json:"attachments"`
	AccountFlags      int                        `json:"accountFlags"`
	MessageTip        interface{}                `json:"messageTip"`
	RawMetadata       string                     `json:"metadata"`
	Metadata          map[string]interface{}
	ChatRoomAccountID string        `json:"chatRoomAccountId"`
	ID                string        `json:"id"`
	CreatedAt         int64         `json:"createdAt"`
	Embeds            []interface{} `json:"embeds"`
	UsernameColor     string        `json:"usernameColor"`
	Username          string        `json:"username"`
	Displayname       string        `json:"displayname"`
}

type fanslyResponseSubAlert struct {
	ChatRoomID            string `json:"chatRoomId"`
	SenderID              string `json:"senderId"`
	HistoryID             string `json:"historyId"`
	SubscriberID          string `json:"subscriberId"`
	SubscriptionTierID    string `json:"subscriptionTierId"`
	SubscriptionTierName  string `json:"subscriptionTierName"`
	SubscriptionTierColor string `json:"subscriptionTierColor"`
	SubscriptionStreak    int    `json:"subscriptionStreak"`
	SubscriptionTotalDays int    `json:"subscriptionTotalDays"`
	ID                    string `json:"id"`
	UsernameColor         string `json:"usernameColor"`
	Username              string `json:"username"`
	Displayname           string `json:"displayname"`
}

type fanslyResponseAttachment struct {
	ContentType       int    `json:"contentType"`
	ContentID         string `json:"contentId"`
	RawMetadata       string `json:"metadata"`
	Metadata          map[string]interface{}
	ChatRoomMessageID string `json:"chatRoomMessageId"`
}

// INFO lil js snippet do decode the crappy fansly websocket format
// var rawSub = {
//   "t": 10000,
//   "d": "{\"serviceId\":46,\"event\":\"{\\\"type\\\":53,\\\"subAlert\\\":{\\\"chatRoomId\\\":\\\"408830844350771200\\\",\\\"senderId\\\":\\\"281038385793998848\\\",\\\"historyId\\\":\\\"797165341313605638\\\",\\\"subscriberId\\\":\\\"281038385793998848\\\",\\\"subscriptionTierId\\\":\\\"795201999690801152\\\",\\\"subscriptionTierName\\\":\\\"Plus\\\",\\\"subscriptionTierColor\\\":\\\"#F73838\\\",\\\"subscriptionStreak\\\":1,\\\"subscriptionTotalDays\\\":30,\\\"id\\\":\\\"797165385945198592\\\",\\\"usernameColor\\\":\\\"#0066ff\\\",\\\"username\\\":\\\"ZerGo0\\\",\\\"displayname\\\":\\\"ZerGo0\\\"}}\"}"
// }
// var rawTip = {
//   "t": 10000,
//   "d": "{\"serviceId\":46,\"event\":\"{\\\"type\\\":10,\\\"chatRoomMessage\\\":{\\\"chatRoomId\\\":\\\"408830844350771200\\\",\\\"senderId\\\":\\\"281038385793998848\\\",\\\"content\\\":\\\"tip test\\\",\\\"type\\\":0,\\\"private\\\":0,\\\"attachments\\\":[{\\\"contentType\\\":7,\\\"contentId\\\":\\\"797163797902008322\\\",\\\"metadata\\\":\\\"{\\\\\\\"amount\\\\\\\":100}\\\",\\\"chatRoomMessageId\\\":\\\"797163798283694080\\\"}],\\\"accountFlags\\\":6,\\\"messageTip\\\":null,\\\"metadata\\\":\\\"{\\\\\\\"senderIsCreator\\\\\\\":false,\\\\\\\"senderIsStaff\\\\\\\":false,\\\\\\\"senderIsFollowing\\\\\\\":true,\\\\\\\"senderSubscription\\\\\\\":{\\\\\\\"tierId\\\\\\\":\\\\\\\"795201999690801152\\\\\\\",\\\\\\\"tierColor\\\\\\\":\\\\\\\"#F73838\\\\\\\",\\\\\\\"tierName\\\\\\\":\\\\\\\"Plus\\\\\\\"}}\\\",\\\"chatRoomAccountId\\\":\\\"407996034761891840\\\",\\\"id\\\":\\\"797163798283694080\\\",\\\"createdAt\\\":1751553019637,\\\"embeds\\\":[],\\\"usernameColor\\\":\\\"#0066ff\\\",\\\"username\\\":\\\"ZerGo0\\\",\\\"displayname\\\":\\\"ZerGo0\\\"}}\"}"
// }

// var raw = rawSub

// raw.d = JSON.parse(raw.d);
// raw.d.event = JSON.parse(raw.d.event)
// try {
// raw.d.event.chatRoomMessage.metadata = JSON.parse(raw.d.event?.chatRoomMessage?.metadata) ?? {}
// } catch {}

// console.log(raw)
// console.log(JSON.stringify(raw)) // for struct generator
