package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	tgClient "github.com/zumosik/telegram-go/clients/telegram"
	event_consumer "github.com/zumosik/telegram-go/consumer/event-consumer"
	"github.com/zumosik/telegram-go/events/telegram"
	"github.com/zumosik/telegram-go/storage/postgres"
)

type Config struct {
	Token       string
	DatabaseURL string
}

const (
	tgBotHost   = "api.telegram.org"
	storagePath = "filestorage"
	batchSize   = 100
)

func main() {
	cfg := mustConfig()

	// eventsProcessor := telegram.New(tgClient.New(tgBotHost, cfg.Token), files.New(storagePath))
	storage, err := postgres.New(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Error creating storage: %s", err)
	}
	eventsProcessor := telegram.New(tgClient.New(tgBotHost, cfg.Token), storage)

	log.Println("service started")

	consumer := event_consumer.New(eventsProcessor, eventsProcessor, batchSize)

	if err := consumer.Start(); err != nil {
		log.Fatalf("service is stopped: %s", err)
	}

	// fetcher = fetcher.New(tgClient)
	// processor = processor.New(tgClient)

	// consumer.Start(fetcher, processor)
}

func mustConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return Config{
		Token:       os.Getenv("TELEGRAM_TOKEN"),
		DatabaseURL: os.Getenv("DB_URL"),
	}
}
