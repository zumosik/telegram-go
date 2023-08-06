package files

import (
	"encoding/gob"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/zumosik/telegram-go/lib/e"
	"github.com/zumosik/telegram-go/storage"
)

type Storage struct {
	basePath string
}

const defaultPerm = 0774

func New(basePath string) Storage {
	return Storage{
		basePath: basePath,
	}
}

func (s Storage) Save(page *storage.Page) error {
	const errorMsg = "can't save page"

	fPath := filepath.Join(s.basePath, page.UserName)

	if err := os.MkdirAll(fPath, defaultPerm); err != nil {
		return e.Wrap(errorMsg, err)
	}

	fName, err := fileName(page)
	if err != nil {
		return e.Wrap(errorMsg, err)
	}

	fPath = filepath.Join(fPath, fName)

	file, err := os.Create(fPath)
	if err != nil {
		return e.Wrap(errorMsg, err)
	}

	defer file.Close()

	if err := gob.NewEncoder(file).Encode(page); err != nil {
		return e.Wrap(errorMsg, err)

	}

	return nil
}

func (s Storage) PickRandom(username string) (page *storage.Page, err error) {
	const errorMsg = "can't pick random page"

	path := filepath.Join(s.basePath, username)

	files, err := os.ReadDir(path) // ERROR
	// if we havent save pages we will get an error here
	if err != nil {
		return nil, e.Wrap(errorMsg, err)
	}

	if len(files) <= 0 {
		return nil, storage.ErrNoSavedPages
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	n := r.Intn(len(files))

	file := files[n]

	return s.decodePath(filepath.Join(path, file.Name()))

}

func (s Storage) List(username string) ([]*storage.Page, error) {
	const errorMsg = "can't get list of pages"

	path := filepath.Join(s.basePath, username)

	files, err := os.ReadDir(path) // ERROR
	// if we havent save pages we will get an error here
	if err != nil {
		return nil, e.Wrap(errorMsg, err)
	}

	if len(files) <= 0 {
		return nil, storage.ErrNoSavedPages
	}

	var result []*storage.Page

	for _, f := range files {
		p, _ := s.decodePath(filepath.Join(path, f.Name())) // not checking error
		result = append(result, p)
	}

	if len(result) <= 0 {
		return nil, storage.ErrNoSavedPages
	}

	return result, nil
}

func (s Storage) Remove(p *storage.Page) error {
	fName, err := fileName(p)
	if err != nil {
		return e.Wrap("can't remove file", err)
	}
	path := filepath.Join(s.basePath, p.UserName, fName)

	if err := os.Remove(path); err != nil {
		return e.Wrap(fmt.Sprintf("can't remove file %s", path), err)
	}

	return nil
}

func (s Storage) IsExists(p *storage.Page) (bool, error) {
	fName, err := fileName(p)
	if err != nil {
		return false, e.Wrap("can't check if file exists", err)
	}

	path := filepath.Join(s.basePath, p.UserName, fName)

	switch _, err := os.Stat(path); {
	case errors.Is(err, os.ErrNotExist):
		return false, nil
	case err != nil:
		return false, e.Wrap(fmt.Sprintf("can't check if file exists %s", path), err)
	}

	return true, nil
}

func (s Storage) decodePath(filePath string) (*storage.Page, error) {
	const errorMsg = "can't decode path"

	f, err := os.Open(filePath)
	if err != nil {
		return nil, e.Wrap(errorMsg, err)
	}
	defer f.Close()

	var p storage.Page

	if err := gob.NewDecoder(f).Decode(&p); err != nil {
		return nil, e.Wrap(errorMsg, err)
	}

	return &p, nil
}

func fileName(p *storage.Page) (string, error) {
	return p.Hash()
}
