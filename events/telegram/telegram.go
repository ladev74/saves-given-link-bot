package telegram

import (
	"saves-given-link-bot/clients/telegram"
	"saves-given-link-bot/events"
	"saves-given-link-bot/lib/e"
	"saves-given-link-bot/storage"
)

type Processor struct {
	tg      *telegram.Client
	offset  int
	storage storage.Storage
}

type Meta struct {
	ChatID   int
	Usermame string
}

func New(client *telegram.Client, storage storage.Storage) *Processor {
	return &Processor{
		tg:      client,
		storage: storage,
	}
}

func (p *Processor) Fetch(limit int) ([]events.Event, error) {
	updates, err := p.tg.Updates(p.offset, limit)
	if err != nil {
		return nil, e.Wrap("can't get events", err)
	}

	if len(updates) == 0 {
		return nil, nil
	}

	res := make([]events.Event, 0, len(updates))

	for _, u := range updates {
		res = append(res, event(u))
	}
}

func event(upd telegram.Update) events.Event {
	udpType := fetchType(upd)

	res := events.Event{
		Type: udpType,
		Text: fetchText(upd),
	}

	if udpType == events.Message {
		res.Meta = Meta{
			ChatID:   upd.Message.Chat.ID,
			Usermame: upd.Message.From.UserName,
		}
	}

	return res
}

func fetchText(upd telegram.Update) string {
	if upd.Message == nil {
		return ""
	}

	return upd.Message.Text
}

func fetchType(upd telegram.Update) events.Type {
	if upd.Message == nil {
		return events.Unknow
	}

	return events.Message
}
