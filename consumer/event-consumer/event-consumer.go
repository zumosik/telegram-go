package event_consumer

import (
	"log"
	"time"

	"github.com/zumosik/telegram-go/events"
)

type Consumer struct {
	fetcher   events.Fetcher
	processor events.Processor
	batchSize int
}

func New(fetcher events.Fetcher, processor events.Processor, batchSize int) Consumer {
	return Consumer{
		fetcher:   fetcher,
		processor: processor,
		batchSize: batchSize,
	}
}

func (c Consumer) Start() error {
	for {
		gotEvents, err := c.fetcher.Fetch(c.batchSize)
		if err != nil {
			// this isnt so good
			// fetcher needs retry
			log.Printf("[ERR] consumer: %s", err.Error())
			continue
		}

		if len(gotEvents) == 0 {
			time.Sleep(time.Second)
			continue
		}

		if err := c.handleEvents(gotEvents); err != nil {
			log.Print(err)

			continue
		}
	}
}

func (c *Consumer) handleEvents(events []events.Event) error {
	// TODO: обработать параллельно
	for _, e := range events {
		log.Printf("got new event: %s", e)

		if err := c.processor.Process(e); err != nil {
			log.Printf("error processing msg: %s", err)
			continue
		}
	}

	return nil
}
