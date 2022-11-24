package storage

import (
	"crypto/sha1"
	"fmt"
	"io"
	_ "time"

	"github.com/Nau077/golang-tg-bot/lib/e"
)

type Storage interface {
	Save(p *Page) error
	PickRandom(userName string) (*Page, error)
	Remove(p *Page) error
	IsExists(p *Page) (bool, error)
}

// основной тип данных, с которым будет работать storage

type Page struct {
	URL      string
	UserName string
	// Created time.Time
}

func (p Page) Hash() (string, error) {
	h := sha1.New()

	if _, err := io.WriteString(h, p.URL); err != nil {
		return "", e.Wrap("не могу записать хэш", err)
	}

	if _, err := io.WriteString(h, p.UserName); err != nil {
		return "", e.Wrap("не могу записать хэш", err)
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
