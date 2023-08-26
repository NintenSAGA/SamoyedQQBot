package botclient

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func getResult(response *http.Response, result any) {
	all, err := io.ReadAll(response.Body)
	if err != nil {
		log.Panicln(err.Error())
	}
	err = json.Unmarshal(all, result)
	if err != nil {
		log.Panicln(err)
	}
}

type UrlResponse struct {
	Url string
}

// ==================== Messages ==================== //
type InitMessage struct {
	Op int
	D  struct {
		HeartbeatInterval int `json:"heartbeat_interval"`
	}
}

type IdentifyMessage struct {
	Op int       `json:"op"`
	D  IdentifyD `json:"d"`
}

type IdentifyD struct {
	Token      string            `json:"token"`
	Intents    int               `json:"intents"`
	Shard      []int             `json:"shard"`
	Properties map[string]string `json:"properties"`
}

type ReadyEventMessage struct {
	Op int
	S  int
	T  string
	D  struct {
		Version   int
		SessionId string `json:"session_id"`
		User      struct {
			Id       string
			Username string
			Bot      bool
		}
		Shard []int
	}
}

type HeartbeatMessage struct {
	Op int  `json:"op"`
	D  *int `json:"d"`
}

type OpMessage struct {
	Op int
	S  int
	T  string
	Id string
	D  json.RawMessage
}

type MessageVO struct {
	Id        string   `json:"id"`
	ChannelId string   `json:"channel_id"`
	GuildId   string   `json:"guild_id"`
	Content   string   `json:"content"`
	Author    UserVO   `json:"author"`
	Mentions  []UserVO `json:"mentions"`
}

type UserVO struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
	Bot      bool   `json:"bot"`
}

type SendMessageRequest struct {
	channelId        string
	Content          string              `json:"content"`
	MessageReference *MessageReferenceVo `json:"message_reference"`
	Image            string              `json:"image"`
	MsgId            string              `json:"msg_id"`
	EventId          string              `json:"event_id"`
}

type MessageReferenceVo struct {
	MessageId string `json:"message_id"`
}
