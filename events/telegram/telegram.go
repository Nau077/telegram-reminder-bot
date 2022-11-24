package telegram

import (
	"github.com/Nau077/golang-tg-bot/clients/telegram"
	"github.com/Nau077/golang-tg-bot/events"
	"github.com/Nau077/golang-tg-bot/lib/e"
	"github.com/Nau077/golang-tg-bot/lib/storage"
)

type Processor struct {
	tg     *telegram.Client
	offset int
	storage storage.Storage
}

func New(client *telegram.Client, storage storage.Storage) *Processor{
	return &Processor{
		tg: client,
		storage: storage,
	}
}

func (p *Processor) Fetch(limit int) ([]events.Event, error) {
	update, err := p.tg.Updates(p.offset, limit)
	if err != nil {
		return nil, e.Wrap("не могу получить события", err)
	}
	// релоцируем сразу память под результат
	res:=make([]events.Event, 0, len(update))
	// нужно обойти все апдейты и преобразовать их в евенты
	for _, u := make(update) {
		res = append(res, event(u))
	}

}

func event(u telegram.Update) events.Event {
	
}