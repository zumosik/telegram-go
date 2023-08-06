package telegram

import (
	"errors"

	"github.com/zumosik/telegram-go/clients/telegram"
	"github.com/zumosik/telegram-go/events"
	"github.com/zumosik/telegram-go/lib/e"
	"github.com/zumosik/telegram-go/storage"
)

type Meta struct {
	ChatID   int
	Username string
}

type Processor struct {
	tgClient *telegram.Client
	offset   int
	storage  storage.Storage
}

var (
	ErrUnknownEventType = errors.New("unknown event type")
	ErrUnknownMetaType  = errors.New("unknown meta type")
)

func New(client *telegram.Client, strorage storage.Storage) *Processor {
	return &Processor{
		tgClient: client,
		offset:   0,
		storage:  strorage,
	}
}

func (p *Processor) Fetch(limit int) ([]events.Event, error) {
	updates, err := p.tgClient.Updates(p.offset, limit)
	if err != nil {
		return nil, e.Wrap("can't get events", err)
	}

	if len(updates) <= 0 {
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
		return p.processMsg(event)
	default:
		return e.Wrap("can't process msg", ErrUnknownEventType)
	}
}

func (p *Processor) processMsg(event events.Event) error {
	meta, err := meta(event)
	if err != nil {
		return e.Wrap("can't process msg", err)
	}

	if err := p.doCmd(event.Text, meta.ChatID, meta.Username); err != nil {
		return p.tgClient.SendMessages(meta.ChatID, msgSomethingWentWrong)
	}

	return nil
}

func meta(event events.Event) (Meta, error) {
	res, ok := event.Meta.(Meta)
	if !ok {
		return Meta{}, e.Wrap("can't get meta", ErrUnknownMetaType)
	}
	return res, nil
}

func event(u telegram.Update) events.Event {
	uType := fetchType(u)
	res := events.Event{
		Type: uType,
		Text: u.Message.Text,
	}

	if uType == events.Message {
		res.Meta = Meta{
			ChatID:   u.Message.Chat.ID,
			Username: u.Message.From.Username,
		}
	}

	return res
}

func fetchType(u telegram.Update) events.Type {
	if u.Message.Text == "" {
		return events.Unknown
	}
	return events.Message
}
