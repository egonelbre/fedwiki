package globstore

import (
	"os"
	"path/filepath"

	"github.com/egonelbre/fedwiki"
)

type Store struct {
	Glob string
}

// serves/saves to files matching the glob
func New(glob string) *Store {
	return &Store{glob}
}

func (store *Store) path(slug fedwiki.Slug) (string, error) {
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

	return "", fedwiki.ErrNotExist
}

func (store *Store) Exists(slug fedwiki.Slug) bool {
	path, err := store.path(slug)
	if err != nil {
		return false
	}

	stat, err := os.Stat(path)
	return os.IsExist(err) && !stat.IsDir()
}

func (store *Store) Load(slug fedwiki.Slug) (*fedwiki.Page, error) {
	path, err := store.path(slug)
	if err != nil {
		return nil, err
	}

	return pageutil.Load(path)
}

func (store *Store) Save(slug fedwiki.Slug, page *fedwiki.Page) error {
	path, err := store.path(slug)
	if err != nil {
		return nil, err
	}

	return pageutil.Save(page, path)
}

// Discards any errors that happen in sub-stores
func (store *Store) List() ([]*fedwiki.PageHeader, error) {
	matches, err := filepath.Glob(store.Glob)
	err = pageutil.ConvertOSError(err)
	if err != nil {
		return "", err
	}

	headers := []*fedwiki.PageHeader{}
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
