package main

import (
	"flag"
	"fmt"
	"log"
)

func main() {
	t := mustToken()
	fmt.Println(t)
}

func mustToken() string {
	token := flag.String("tg-api-token", "", "telegram bot api token")

	flag.Parse()

	if *token == "" {
		log.Fatal("bot token is not specified")
	}

	return *token
}
