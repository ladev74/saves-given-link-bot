package eventcomsumer

import (
	"log"
	"saves-given-link-bot/events"
	"time"
)

type Consumer struct {
	fetcher   events.Fetcher
	processor events.Processor
	bathSize  int
}

func New(fetcher events.Fetcher, processor events.Processor, bathSize int) Consumer {
	return Consumer{
		fetcher:   fetcher,
		processor: processor,
		bathSize:  bathSize,
	}
}

func (c Consumer) Start() error {
	for {
		gotEvents, err := c.fetcher.Fetch(c.bathSize)
		if err != nil {
			log.Printf("[ERR] comsumer: %s", err.Error())

			continue
		}

		if len(gotEvents) == 0 {
			time.Sleep(1 * time.Second)

			continue
		}

		if err := c.handleEvents(gotEvents); err != nil {
			log.Print(err)

			continue
		}
	}
}

// распараллелить
func (c *Consumer) handleEvents(events []events.Event) error {
	for _, event := range events {
		//sync.WaitGroup{}
		log.Printf("got new events: %s", event.Text)

		if err := c.processor.Process(event); err != nil {
			log.Printf("can't handle event: %s", event.Text)

			continue
		}
	}

	return nil
}
