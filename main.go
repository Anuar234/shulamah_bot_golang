package main

import (
	"flag"
	"log"

	tgClient "shulamah_bot_golang/clients/telegram"
	event_consumer "shulamah_bot_golang/consumer/event-consumer"
	"shulamah_bot_golang/events/telegram"
	"shulamah_bot_golang/storage/files"
)


const (
	tgBotHost = "api.telegram.org"
	storagePath = "storage"
	batchSize = 100
	downloadDir = "downloads"
)

func main() {

	eventsProcessor := telegram.New(
		tgClient.New(tgBotHost, mustToken()),
		files.New(storagePath),
	)
	log.Print("service started")

	consumer := event_consumer.New(eventsProcessor, eventsProcessor, batchSize)

	if err := consumer.Start();err!=nil{
		log.Fatal("service is stopped", err)
	}


}

func mustToken() string {
	token := flag.String(
	"tg-bot-token", 
	"", 
	"token for access to telegram bot",
	)


	flag.Parse()

	if *token == "" {
		log.Fatal("token is not specified")
	}

	return *token
}