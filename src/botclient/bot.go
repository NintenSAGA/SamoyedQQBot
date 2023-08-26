package botclient

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
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
	_, initMsgRaw, _ := conn.ReadMessage()
	log.Println(string(initMsgRaw))
	initMsg := InitMessage{}
	json.Unmarshal(initMsgRaw, &initMsg)
	heartbeatInterval := initMsg.D.HeartbeatInterval

	// Authorization

	for {

	}

	conn.Close()
}
