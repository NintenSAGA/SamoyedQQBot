package main

import (
	"log"
	"os"
	"strings"
)

func main() {
	botToken := os.Getenv("BOT_TOKEN")
	if len(botToken) == 0 {
		log.Panicln("Missing env BOT_TOKEN")
	}

	log.Printf("BOT_TOKEN = %v\n", botToken)

	isOnline := false
	onlineEnv := os.Getenv("ONLINE")
	if len(onlineEnv) != 0 && strings.EqualFold("true", onlineEnv) {
		log.Printf("Using online environment.")
		isOnline = true
	}

	client := CreateBotClient(botToken, isOnline)
	url := client.GetWebSocketUrl()
	println(url)

	client.EstablishWSConnection()
}
