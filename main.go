package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	tgClient "github.com/zumosik/telegram-go/clients/telegram"
	event_consumer "github.com/zumosik/telegram-go/consumer/event-consumer"
	"github.com/zumosik/telegram-go/events/telegram"
	"github.com/zumosik/telegram-go/storage/files"
)

const (
	tgBotHost   = "api.telegram.org"
	storagePath = "filestorage"
	batchSize   = 100
)

func main() {

	eventsProcessor := telegram.New(tgClient.New(tgBotHost, mustToken()), files.New(storagePath))

	log.Println("service started")

	consumer := event_consumer.New(eventsProcessor, eventsProcessor, batchSize)

	if err := consumer.Start(); err != nil {
		log.Fatalf("service is stopped: %w", err)
	}

	// fetcher = fetcher.New(tgClient)
	// processor = processor.New(tgClient)

	// consumer.Start(fetcher, processor)
}

func mustToken() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv("TELEGRAM_TOKEN")
}
