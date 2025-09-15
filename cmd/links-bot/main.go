package main

import (
	"flag"
	"log"

	tgClient "github.com/cmczk/links-bot/clients/telegram"
	event_consumer "github.com/cmczk/links-bot/consumer/event-consumer"
	"github.com/cmczk/links-bot/events/telegram"
	"github.com/cmczk/links-bot/storage/files"
)

const (
	tgBotHost   = "api.telegram.org"
	storagePath = "storage"
	batchSize   = 100
)

func main() {
	tgClient := tgClient.New(tgBotHost, mustToken())
	eventsProcessor := telegram.New(tgClient, files.New(storagePath))

	log.Println("service started")

	consumer := event_consumer.New(eventsProcessor, eventsProcessor, batchSize)

	if err := consumer.Start(); err != nil {
		log.Fatalf("service stopped: %s", err.Error())
	}
}

func mustToken() string {
	token := flag.String("tg-api-token", "", "telegram bot api token")

	flag.Parse()

	if *token == "" {
		log.Fatal("bot token is not specified")
	}

	return *token
}
