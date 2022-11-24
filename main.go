package main

import (
	"flag"
	"log"

	"github.com/Nau077/golang-tg-bot/clients/telegram"
)

const (
	tgBotHost = "api.telegram.org"
)

func main() {

	tgClient := telegram.New(mustToken())

	// fetcher = fetcher.New()

	// processor = processor.New()

	// consumer.Start(fetcher, processor)
}

func mustToken() string {
	token := flag.String(
		"token-bot-token",
		"",
		"token for access to telegram bot",
	)

	if *token == "" {
		log.Fatal("token is not specified")
	}

	return *token
}
