package files

import (
	"encoding/gob"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"

	"github.com/cmczk/links-bot/lib/e"
	"github.com/cmczk/links-bot/storage"
)

const (
	defaultPerm = 0774
)

type Storage struct {
	basePath string
}

func New(basePath string) Storage {
	return Storage{basePath: basePath}
}

func (s Storage) Save(p *storage.Page) (err error) {
	defer func() { err = e.WrapIfErr("can't save new page", err) }()

	fPath := filepath.Join(s.basePath, p.UserName)

	if err := os.MkdirAll(fPath, defaultPerm); err != nil {
		return err
	}

	fName, err := fileName(p)
	if err != nil {
		return err
	}

	fPath = filepath.Join(fPath, fName)

	f, err := os.Create(fPath)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()

	if err := gob.NewEncoder(f).Encode(p); err != nil {
		return err
	}

	return nil
}

func (s Storage) PickRandom(userName string) (p *storage.Page, err error) {
	defer func() { err = e.WrapIfErr("cannot pick random page", err) }()

	path := filepath.Join(s.basePath, userName)

	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, storage.ErrNoSavedPages
	}

	n := rand.Intn(len(files))

	file := files[n]

	return s.decodePage(filepath.Join(path, file.Name()))
}

func (s Storage) Remove(p *storage.Page) error {
	fName, err := fileName(p)
	if err != nil {
		return e.Wrap("cannot generate file hash", err)
	}

	path := filepath.Join(s.basePath, p.UserName, fName)

	if err := os.Remove(path); err != nil {
		return e.Wrap(fmt.Sprintf("cannot remove file: %s", path), err)
	}

	return nil
}

func (s Storage) Exists(p *storage.Page) (bool, error) {
	fName, err := fileName(p)
	if err != nil {
		return false, e.Wrap("cannot generate file hash", err)
	}

	path := filepath.Join(s.basePath, p.UserName, fName)

	switch _, err = os.Stat(path); {
	case errors.Is(err, os.ErrNotExist):
		return false, nil
	case err != nil:
		return false, e.Wrap(fmt.Sprintf("cannot check if file exists: %s", path), err)
	}

	return true, nil
}

func (s Storage) decodePage(filePath string) (*storage.Page, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, e.Wrap("cannot open file", err)
	}
	defer func() { _ = f.Close() }()

	var p storage.Page

	if err := gob.NewDecoder(f).Decode(&p); err != nil {
		return nil, e.Wrap("cannot decode file", err)
	}

	return &p, nil
}

func fileName(p *storage.Page) (string, error) {
	return p.Hash()
}
