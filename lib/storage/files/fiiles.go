package files

import (
	"encoding/gob"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/Nau077/golang-tg-bot/lib/e"
	"github.com/Nau077/golang-tg-bot/lib/storage"
)

// реализация с файловой системой

type Storage struct {
	basePath string
}

const defaultPerm = 0774

var ErrNoSavedPages = errors.New("нет сохранённых страниц")

func New(basePath string) Storage {
	return Storage{basePath: basePath}
}

func (s Storage) Save(page *storage.Page) (err error) {
	defer func() {
		err = e.WrapIfErr("не могу сохранить страницу", err)
	}()
	// получаем путь до директории с файлами
	fPath := filepath.Join(s.basePath, page.UserName)

	if err := os.MkdirAll(fPath, defaultPerm); err != nil {
		return err
	}

	fName, err := fileName(page)
	if err != nil {
		return err
	}
	// формируем путь до директории, куда будет сохраняться файл
	fPath = filepath.Join(fPath, fName)
	// создаём файл
	file, err := os.Create(fPath)
	if err != nil {
		return err
	}

	defer func() {
		_ = file.Close()
	}()
	// записываем в файл страницу в нужном формате
	if err := gob.NewEncoder(file).Encode(page); err != nil {
		return err
	}

	return nil

}

func (s Storage) PickRandom(userName string) (page *storage.Page, err error) {
	defer func() {
		err = e.WrapIfErr("не могу сохранить", err)
	}()

	path := filepath.Join(s.basePath, userName)
	// получаем список файлов
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	// нужно, чтобы бот мог сообщить пользователю, что ничего не сохранил
	if len(files) == 0 {
		return nil, ErrNoSavedPages
	}
	// нужно получить число от 0 до номера последнего файла
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(len(files))

	file := files[n]
	// декодируем файл и возвращаем содержимое
	return s.decodePage(filepath.Join(path, file.Name()))
}

func (s Storage) IsExists(p *storage.Page) (bool, error) {
	fileName, err := fileName(p)
	if err != nil {
		return false, e.Wrap("не могу удалить файл", err)
	}

	path := filepath.Join(s.basePath, p.UserName, fileName)
	// проверяем на существование файл

	switch _, err = os.Stat(path); {
	case errors.Is(err, os.ErrNotExist):
		return false, nil
	case err != nil:
		msg := "не могу удалить файл %s"
		return false, e.Wrap(fmt.Sprintf(msg, path), err)
	}

	return true, nil
}

func (s Storage) Remove(p *storage.Page) error {
	fileName, err := fileName(p)
	if err != nil {
		return e.Wrap("не могу удалить файл", err)
	}

	path := filepath.Join(s.basePath, p.UserName, fileName)
	if err := os.Remove(path); err != nil {
		msg := "не могу удалить файл %s"
		return e.Wrap(fmt.Sprintf(msg, path), err)
	}

	return nil
}

func (s Storage) decodePage(filePath string) (*storage.Page, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, e.Wrap("не могу декодировать файл", err)
	}

	defer func() {
		_ = file.Close()
	}()

	var p storage.Page
	if err := gob.NewDecoder(file).Decode(&p); err != nil {
		return nil, e.Wrap("не могу декодировать файл", err)
	}

	return &p, nil
}

func fileName(p *storage.Page) (string, error) {
	return p.Hash()
}
