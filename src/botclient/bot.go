package botclient

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	HOST = "https://sandbox.api.sgroup.qq.com"
	GET  = "GET"
	POST = "POST"
)

type BotClient struct {
	botToken   string
	httpClient *http.Client
}

func CreateBotClient(botToken string) *BotClient {
	instance := &BotClient{
		botToken:   botToken,
		httpClient: &http.Client{},
	}

	return instance
}

func (b *BotClient) getRequest(method string, path string, body *string) *http.Request {
	vUrl, _ := url.JoinPath(HOST, path)
	var reader io.Reader
	if body != nil {
		reader = strings.NewReader(*body)
	} else {
		reader = nil
	}
	request, err := http.NewRequest(method, vUrl, reader)
	if err != nil {
		log.Panicln(err.Error())
	}

	request.Header.Add("Authorization", b.botToken)

	return request
}

func (b *BotClient) GetWebSocketUrl() string {
	request := b.getRequest(GET, "/gateway", nil)
	response, err := b.httpClient.Do(request)
	if err != nil {
		log.Panicln(err.Error())
	}
	urlResponse := UrlResponse{}
	getResult(response, &urlResponse)

	return urlResponse.Url
}

func (b *BotClient) EstablishWSConnection() {
	wsUrl := b.GetWebSocketUrl()
	if len(wsUrl) == 0 {
		log.Panicln("WS Url is blank!")
	}

	dialer := websocket.Dialer{}
	conn, _, err := dialer.Dial(wsUrl, nil)
	if err != nil {
		log.Panicln(err.Error())
	}

	// Init
	initMsg := receive(conn, &InitMessage{})
	heartbeatInterval := initMsg.D.HeartbeatInterval

	// Authorization
	send(conn, IdentifyMessage{
		Op: 2,
		D: IdentifyD{
			Token:      b.botToken,
			Intents:    1 << 9,
			Shard:      []int{0, 1},
			Properties: map[string]string{},
		},
	})
	receive(conn, &ReadyEventMessage{})

	var latestMessage *int
	latestMessage = nil
	readCh := make(chan []byte, 10)
	go b.readMessageHandler(conn, readCh)

	interval := time.Millisecond * time.Duration(heartbeatInterval)
	log.Printf("Heartbeat interval: %v\n", interval)
	timer := time.NewTimer(interval)
	for {
		select {
		case <-timer.C:
			send(conn, HeartbeatMessage{Op: 1, D: latestMessage})
			timer.Reset(interval)
		case msg := <-readCh:
			logReceived(msg)
			opMsg := getOp(msg)
			switch opMsg.Op {
			case 11: // Heartbeat
			case 0: // Message
				msgObj := MessageVO{}
				json.Unmarshal(opMsg.D, &msgObj)
				fmt.Printf("%v 发送了一条消息 “%v”\n", msgObj.Author.Username, msgObj.Content)
			}
		}
	}

	conn.Close()
}

func (b *BotClient) readMessageHandler(conn *websocket.Conn, ch chan []byte) {
	for {
		_, bytes, err := conn.ReadMessage()
		if err != nil {
			log.Panicln(err)
		}
		ch <- bytes
	}
}

func getOp(raw []byte) OpMessage {
	msg := OpMessage{}
	json.Unmarshal(raw, &msg)
	return msg
}

func send(conn *websocket.Conn, msg interface{}) {
	sendOpt(conn, msg, true)
}

func sendOpt(conn *websocket.Conn, msg interface{}, printLog bool) {
	bytes, _ := json.Marshal(msg)
	if printLog {
		logSent(bytes)
	}
	_ = conn.WriteMessage(websocket.TextMessage, bytes)
}

func receive[T interface{}](conn *websocket.Conn, result *T) *T {
	_, bytes, _ := conn.ReadMessage()
	logReceived(bytes)
	json.Unmarshal(bytes, result)
	return result
}

func logSent(s []byte) {
	log.Printf("Sent >> %v\n", string(s))
}
func logReceived(s []byte) {
	log.Printf("Received >> %v\n", string(s))
}
