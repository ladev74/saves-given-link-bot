package main

import (
	"flag"
	"log"
	tgClient "saves-given-link-bot/clients/telegram"
	eventcomsumer "saves-given-link-bot/consumer/event-comsumer"
	"saves-given-link-bot/events/telegram"
	"saves-given-link-bot/storage/files"
)

const (
	tgBotGost   = "api.telegram.org"
	storagePath = "file_storage"
	bathSize    = 100
)

func main() {
	eventProcessor := telegram.New(tgClient.New(tgBotGost, mustToken()), files.New(storagePath))

	log.Print("service started")

	consumer := eventcomsumer.New(eventProcessor, eventProcessor, bathSize)

	if err := consumer.Start(); err != nil {
		log.Fatal("service is stopped", err)
	}
}

func mustToken() string {
	token := flag.String("tg-bot-token", "", "token for access for telegram bot")

	flag.Parse()

	if *token == "" {
		log.Fatal("token is not specified ")
	}

	return *token
}
