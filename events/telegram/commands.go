package telegram

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"shulamah_bot_golang/downloader"
	"shulamah_bot_golang/lib/e"
	"shulamah_bot_golang/storage"
	"strings"
)

const (
	RndCmd      = "/rnd"
	HelpCmd     = "/help"
	StartCmd    = "/start"
	DownloadCmd = "/download"
)

func (p *Processor) doCmd(text string, chatID int, username string) error {
	text = strings.TrimSpace(text)
	log.Printf("got new command '%s' from '%s'", text, username)

	fields := strings.Fields(text)
	if len(fields) == 0 {
		return p.tg.SendMessage(chatID, msgUnkownCommand)
	}

	cmd := fields[0]

	switch cmd {
	case RndCmd:
		return p.SendRandom(chatID, username)
	case HelpCmd:
		return p.SendHelp(chatID)
	case StartCmd:
		return p.SendHello(chatID)
	case DownloadCmd:
		return p.HandleDownloadCommand(chatID, text)
	default:
		// Если это не команда, но валидный URL — сохранить
		if isAddCmd(text) {
			return p.savePage(chatID, text, username)
		}
		return p.tg.SendMessage(chatID, msgUnkownCommand)
	}
}





func (p *Processor) savePage(chatID int, pageURL string, username string) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: save page", err) }()

	page := &storage.Page{
		URL:      pageURL,
		UserName: username,
	}

	exists, err := p.storage.IsExists(page)
	if err != nil {
		return err
	}

	if exists {
		return p.tg.SendMessage(chatID, msgALreadyExists)
	}

	if err := p.storage.Save(page); err != nil {
		return err
	}

	return p.tg.SendMessage(chatID, msgSaved)
}

func (p *Processor) SendRandom(chatID int, username string) (err error) {
	defer func() { err = e.Wrap("can't do command: send random page", err) }()

	page, err := p.storage.PickRandom(username)
	if errors.Is(err, storage.ErrNoSavedPages) {
		return p.tg.SendMessage(chatID, msgNoSavedPages)
	}
	if err != nil {
		return err
	}

	if err := p.tg.SendMessage(chatID, page.URL); err != nil {
		return err
	}

	return p.storage.Remove(page)
}

func (p *Processor) SendHelp(chatID int) error {
	return p.tg.SendMessage(chatID, msgHelp)
}

func (p *Processor) SendHello(chatID int) error {
	return p.tg.SendMessage(chatID, msgHello)
}

func isAddCmd(text string) bool {
	return isURL(text)
}

func isURL(text string) bool {
	u, err := url.Parse(text)
	return err == nil && u.Host != ""
}

func (p *Processor) HandleDownloadCommand(chatID int, text string) error {
    args := strings.Fields(text)
    if len(args) < 2 {
        return p.tg.SendMessage(chatID, "Usage: /download <video_url>")
    }

    url := args[1]

    // Ensure downloads folder exists
    const downloadDir = "downloads"
    if err := os.MkdirAll(downloadDir, os.ModePerm); err != nil {
        return p.tg.SendMessage(chatID, "Failed to create download directory")
    }

    filePath, err := downloader.DownloadVideo(url, downloadDir)
    if err != nil {
        return p.tg.SendMessage(chatID, fmt.Sprintf("Error downloading video: %v", err))
    }
    return p.tg.SendMessage(chatID, fmt.Sprintf("Video downloaded to:\n%s", filePath))
}
