package storage

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"shulamah_bot_golang/lib/e"
)

type Storage interface {
	Save(p *Page) error
	PickRandom(UserName string) (*Page, error)
	Remove(p *Page) error
	IsExists(p *Page) (bool, error)
}

var ErrNoSavedPages = errors.New("no saved page")

type Page struct {
	URL      string
	UserName string
}

func (p Page) Hash() (string, error) {
	h := sha1.New()

	if _, err := io.WriteString(h, p.URL);err != nil{
		return "", e.Wrap("cant calculate", err)
	}

	if _, err := io.WriteString(h, p.UserName);err != nil{
		return "", e.Wrap("cant calculate", err)
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}