package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	TEST_HOST   = "https://sandbox.api.sgroup.qq.com"
	ONLINE_HOST = "https://api.sgroup.qq.com"
	GET         = "GET"
	POST        = "POST"
)

type BotClient struct {
	botToken   string
	httpClient *http.Client
	rcvMsgChan chan MessageVO
	sndMsgChan chan SendMessageRequest
	solver     *Solver
	mineId     string
	isOnline   bool
}

func CreateBotClient(botToken string, isOnline bool) *BotClient {
	instance := &BotClient{
		botToken:   botToken,
		httpClient: &http.Client{},
		rcvMsgChan: make(chan MessageVO, 10),
		sndMsgChan: make(chan SendMessageRequest, 10),
		solver:     createSolver(),
		isOnline:   isOnline,
	}

	go instance.sendMessageHandler()
	go instance.receiveMessageHandler()

	return instance
}

func (b *BotClient) getRequest(method string, path string, body *string) *http.Request {
	var host string
	if b.isOnline {
		host = ONLINE_HOST
	} else {
		host = TEST_HOST
	}
	vUrl, _ := url.JoinPath(host, path)
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
	readyEventMessage := receive(conn, &ReadyEventMessage{})
	b.mineId = readyEventMessage.D.User.Id

	var latestMessage *int
	latestMessage = nil
	readCh := make(chan []byte, 10)
	go b.wsReadMessageHandler(conn, readCh)

	interval := time.Millisecond * time.Duration(heartbeatInterval)
	log.Printf("Heartbeat interval: %v\n", interval)
	timer := time.NewTimer(interval)

outer:
	for {
		select {
		case <-timer.C:
			send(conn, HeartbeatMessage{Op: 1, D: latestMessage})
			timer.Reset(interval)
		case msg := <-readCh:
			if msg == nil {
				break outer
			}
			logReceived(msg)
			opMsg := getOp(msg)
			latestMessage = &opMsg.S
			switch opMsg.Op {
			case 11: // Heartbeat
			case 0: // Message
				msgObj := MessageVO{}
				json.Unmarshal(opMsg.D, &msgObj)
				b.rcvMsgChan <- msgObj
			}
		}
	}

	conn.Close()
	b.EstablishWSConnection()
}

func (b *BotClient) wsReadMessageHandler(conn *websocket.Conn, ch chan []byte) {
	for {
		_, bytes, err := conn.ReadMessage()
		if err != nil {
			log.Printf(err.Error())
			ch <- nil
			break
		}
		ch <- bytes
	}
}

func (b *BotClient) receiveMessageHandler() {
	for {
		msg := <-b.rcvMsgChan
		username := msg.Author.Username
		content := msg.Content

		mentionedMe := false
		for _, user := range msg.Mentions {
			if user.Bot {
				mentionedMe = true
				break
			}
		}

		if mentionedMe {
			content = strings.ReplaceAll(content, fmt.Sprintf("\u003c@!%v\u003e", b.mineId), "")
			content = strings.TrimSpace(content)
			sentence := b.solver.generateAnswer(username, content)
			request := SendMessageRequest{
				channelId:        msg.ChannelId,
				Content:          sentence,
				MessageReference: &MessageReferenceVo{MessageId: msg.Id},
			}

			b.sndMsgChan <- request
		} else if rand.Intn(1000) >= 600 {
			request := SendMessageRequest{
				channelId: msg.ChannelId,
				Content:   "(小狗正在偷听你们说话)",
			}
			b.sndMsgChan <- request
		}

	}
}

func (b *BotClient) sendMessageHandler() {
	for {
		request := <-b.sndMsgChan
		b.sendMessage(request)
	}
}

func (b *BotClient) sendMessage(request SendMessageRequest) {
	bytes, _ := json.Marshal(request)
	body := string(bytes)
	req := b.getRequest(POST, fmt.Sprintf("/channels/%v/messages", request.channelId), &body)
	req.Header.Add("Content-Type", "application/json")
	response, err := b.httpClient.Do(req)
	if err != nil {
		log.Printf("sendMessage Error: %v. Request: %v", err.Error(), body)
	}

	raw, _ := io.ReadAll(response.Body)
	log.Printf("sendMessage request: %v response: %v %v\n", request, response, string(raw))
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
