package telegram

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/zumosik/telegram-go/lib/e"
	"github.com/zumosik/telegram-go/storage"
)

// add page: https://
// rnd page: /rnd
// help: /help
// start: /start: hello + help

const (
	RandCmd  = "/rnd"
	ListCmd  = "/list"
	HelpCmd  = "/help"
	StartCmd = "/start"
)

func (p *Processor) doCmd(text string, chatID int, username string) error {
	text = strings.TrimSpace(text)

	log.Printf("got new command '%s' from '%s'", text, username)

	if isURL(text) {
		return p.savePage(chatID, text, username)
	}

	switch text {
	case RandCmd:
		return p.sendRandom(chatID, username)
	case ListCmd:
		return p.sendList(chatID, username)
	case HelpCmd:
		return p.sendHelp(chatID)
	case StartCmd:
		return p.sendHello(chatID)
	default:
		return p.tgClient.SendMessages(chatID, msgUnknownCommand)
	}
}

func (p *Processor) savePage(chatID int, pageURL string, username string) error {
	const errorMsg = "can't do command: save page"

	page := &storage.Page{
		URL:      pageURL,
		UserName: username,
	}

	isExists, err := p.storage.IsExists(page)
	if err != nil {
		return e.Wrap(errorMsg, err)
	}

	if isExists {
		return p.tgClient.SendMessages(chatID, msgAlreadyExists)
	}

	if err := p.storage.Save(page); err != nil {
		return e.Wrap(errorMsg, err)
	}

	return p.tgClient.SendMessages(chatID, msgSaved)
}

func (p *Processor) sendRandom(chatID int, username string) error {
	const errorMsg = "can't do command: save page"

	page, err := p.storage.PickRandom(username)
	if errors.Is(err, storage.ErrNoSavedPages) {
		return p.tgClient.SendMessages(chatID, msgNoSavedPages)
	}
	if err != nil {
		return e.Wrap(errorMsg, err)
	}

	return p.tgClient.SendMessages(chatID, page.URL)
}

func (p *Processor) sendList(chatID int, username string) error {
	const errorMsg = "can't do command: get list"
	pages, err := p.storage.List(username)
	if err != nil {
		return e.Wrap(errorMsg, err)
	}
	var text string
	for i, p := range pages {
		text = text + fmt.Sprintf("\n%d. %s", i+1, p.URL)
	}
	return p.tgClient.SendMessages(chatID, text)
}

func (p *Processor) sendHelp(chatID int) error {
	return p.tgClient.SendMessages(chatID, msgHelp)
}

func (p *Processor) sendHello(chatID int) error {
	return p.tgClient.SendMessages(chatID, msgHello)
}

func isURL(text string) bool {
	u, err := url.Parse(text)
	return err == nil && u.Host != ""
}
