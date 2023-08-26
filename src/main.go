package main

import (
	"github.com/NintenSAGA/SamoyedQQBot/src/botclient"
	"log"
	"os"
)

func main() {
	botToken := os.Getenv("BOT_TOKEN")
	if len(botToken) == 0 {
		log.Panicln(botToken)
	}

	log.Printf("BOT_TOKEN = %v\n", botToken)

	client := botclient.CreateBotClient(botToken)
	url := client.GetWebSocketUrl()
	println(url)

	client.EstablishWSConnection()
}
