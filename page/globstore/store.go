package globstore

import (
	"os"
	"path/filepath"

	"github.com/egonelbre/wiki-go-server/page"
	"github.com/egonelbre/wiki-go-server/page/pageutil"
)

type Store struct {
	Glob string
}

// serves/saves to files matching the glob
func New(glob string) *Store {
	return &Store{glob}
}

func (store *Store) path(slug page.Slug) (string, error) {
	// todo, cache matches
	matches, err := filepath.Glob(store.Glob)
	err = pageutil.ConvertOSError(err)
	if err != nil {
		return "", err
	}

	for _, filename := range matches {
		if filepath.Base(filename) == string(slug) {
			return filename
		}
	}

	return "", page.ErrNotExist
}

func (store *Store) Exists(slug page.Slug) bool {
	path, err := store.path(slug)
	if err != nil {
		return false
	}

	stat, err := os.Stat(path)
	return os.IsExist(err) && !stat.IsDir()
}

func (store *Store) Load(slug page.Slug) (*page.Page, error) {
	path, err := store.path(slug)
	if err != nil {
		return nil, err
	}

	return pageutil.Load(path)
}

func (store *Store) Save(slug page.Slug, page *page.Page) error {
	path, err := store.path(slug)
	if err != nil {
		return nil, err
	}

	return pageutil.Save(page, path)
}

// Discards any errors that happen in sub-stores
func (store *Store) List() ([]*page.Header, error) {
	matches, err := filepath.Glob(store.Glob)
	err = pageutil.ConvertOSError(err)
	if err != nil {
		return "", err
	}

	headers := []*page.Header{}
	for _, filename := range matches {
		header, err := pageutil.LoadHeader(filename)
		err = pageutil.ConvertOSError(err)
		//TODO: maybe ignore this error?
		if err != nil {
			return nil, err
		}

		headers = append(headers, header)
	}
	return headers, nil
}
