package main

import (
	"flag"
	"log"

	"github.com/cmczk/links-bot/clients/telegram"
)

const (
	tgBotHost = "api.telegram.org"
)

func main() {
	tgClient := telegram.New(tgBotHost, mustToken())
	_ = tgClient
}

func mustToken() string {
	token := flag.String("tg-api-token", "", "telegram bot api token")

	flag.Parse()

	if *token == "" {
		log.Fatal("bot token is not specified")
	}

	return *token
}
