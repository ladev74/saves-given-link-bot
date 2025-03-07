package telegram

import (
	"errors"
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

var (
	ErrUnknowEventType = errors.New("unknow event type")
	ErrUnknowMetaType  = errors.New("unknow meta type")
)

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

	p.offset = updates[len(updates)-1].ID + 1

	return res, nil
}

func (p *Processor) Process(event events.Event) error {
	switch event.Type {
	case events.Message:
		return p.processMessage(event)
	default:
		return e.Wrap("can't process message", ErrUnknowEventType)
	}
}

func (p *Processor) processMessage(event events.Event) error {
	meta, err := meta(event)
	if err != nil {
		return e.Wrap("can't process message", err)
	}

	if err := p.doCmd(event.Text, meta.ChatID, meta.Usermame); err != nil {
		return e.Wrap("can't process message", err)
	}

	return nil
}

func meta(event events.Event) (Meta, error) {
	res, ok := event.Meta.(Meta)
	if !ok {
		return Meta{}, e.Wrap("can't get meta", ErrUnknowMetaType)
	}

	return res, nil
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
