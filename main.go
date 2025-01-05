package main

import (
	"flag"
	"log"
	"saves-given-link-bot/clients/telegram"
)

const (
	tgBotGost = "api.telegram.org"
)

func main() {
	client := telegram.New(tgBotGost, mustToken())

}

func mustToken() string {
	token := flag.String("bot-token", "", "token for access for telegram bot")

	flag.Parse()

	if *token == "" {
		log.Fatal("token is not specified ")
	}

	return *token
}
