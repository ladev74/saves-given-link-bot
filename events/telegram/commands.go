package telegram

import (
	"context"
	"errors"
	"log"
	"net/url"
	"saves-given-link-bot/lib/e"
	"saves-given-link-bot/storage"
	"strings"
)

const (
	RndCmd    = "/rnd"
	HelpCmd   = "/help"
	StartCmd  = "/start"
	PrintList = "/print"
)

func (p *Processor) doCmd(text string, chatID int, userName string) error {
	text = strings.TrimSpace(text)

	log.Printf("got new command '%s' from '%s'", text, userName)

	if isAddCmd(text) {
		return p.savePage(chatID, text, userName)
	}

	switch text {
	case RndCmd:
		return p.sendRandom(chatID, userName)
	case HelpCmd:
		return p.sendHelp(chatID)
	case StartCmd:
		return p.sendHello(chatID)
	case PrintList:
		return p.printList(chatID, userName)
	default:
		return p.tg.SendMessage(chatID, msgUnknownCommand)
	}
}

func (p *Processor) savePage(chatID int, pageURL string, userName string) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: can't save page", err) }()

	page := &storage.Page{
		URL:      pageURL,
		UserName: userName,
	}

	isExists, err := p.storage.IsExists(context.Background(), page)
	if err != nil {
		return err
	}

	if isExists {
		return p.tg.SendMessage(chatID, msgAlreadyExists)
	}

	if err := p.storage.Save(context.Background(), page); err != nil {
		return err
	}

	if err := p.tg.SendMessage(chatID, msgSaved); err != nil {
		return err
	}

	return nil
}

func (p *Processor) sendRandom(chatID int, userName string) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: can't send random", err) }()

	page, err := p.storage.PickRandom(context.Background(), userName)
	if err != nil && !errors.Is(err, storage.ErrNoSavedPages) {
		return err
	}
	if errors.Is(err, storage.ErrNoSavedPages) {
		return p.tg.SendMessage(chatID, msgNoSavedPages)
	}

	if err := p.tg.SendMessage(chatID, page.URL); err != nil {
		return err
	}

	return p.storage.Remove(context.Background(), page)
}

func (p *Processor) printList(chatID int, userName string) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: can't print the list", err) }()

	page, err := p.storage.PrintList(context.Background(), userName)
	if err != nil {
		return err
	}

	if errors.Is(err, storage.ErrNoSavedPages) {
		return p.tg.SendMessage(chatID, msgNoSavedPages)
	}

	if err := p.tg.SendMessage(chatID, page.URL); err != nil {
		return err
	}

	return nil
}

func (p *Processor) sendHelp(chatID int) error {
	return p.tg.SendMessage(chatID, msgHelp)
}

func (p *Processor) sendHello(chatID int) error {
	return p.tg.SendMessage(chatID, msgHello)
}

func isAddCmd(text string) bool {
	return isURL(text)
}

func isURL(text string) bool {
	u, err := url.Parse(text)

	return err == nil && u.Host != ""
}
