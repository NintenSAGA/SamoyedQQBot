package main

import (
	"log"
	"os"
)

func main() {
	botToken := os.Getenv("BOT_TOKEN")
	if len(botToken) == 0 {
		log.Panicln("Missing env BOT_TOKEN")
	}

	log.Printf("BOT_TOKEN = %v\n", botToken)

	client := CreateBotClient(botToken)
	url := client.GetWebSocketUrl()
	println(url)

	client.EstablishWSConnection()
}
