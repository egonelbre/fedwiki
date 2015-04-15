// This package implements PageStore for a folder
package folderstore

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/egonelbre/fedwiki"
	"github.com/egonelbre/fedwiki/pagestore"
)

// Store load/saves pages from a directory
type Store struct {
	Dir string
}

func New(dir string) *Store {
	return &Store{dir}
}

func (store *Store) path(slug fedwiki.Slug) string {
	filename := filepath.FromSlash(string(slug))
	return filepath.Join(store.Dir, filename)
}

func (store *Store) Exists(slug fedwiki.Slug) bool {
	stat, err := os.Stat(store.path(slug))
	return os.IsExist(err) && !stat.IsDir()
}

func (store *Store) Load(slug fedwiki.Slug) (*fedwiki.Page, error) {
	page, err := pagestore.Load(store.path(slug), slug)
	return page, err
}

func (store *Store) Create(slug fedwiki.Slug, page *fedwiki.Page) error {
	return pagestore.Create(page, store.path(slug))
}

func (store *Store) Save(slug fedwiki.Slug, page *fedwiki.Page) error {
	return pagestore.Save(page, store.path(slug))
}

// Discards any errors that happen in sub-stores
func (store *Store) List() ([]*fedwiki.PageHeader, error) {
	list, err := ioutil.ReadDir(store.Dir)
	err = pagestore.ConvertOSError(err)
	if err != nil {
		return nil, err
	}

	headers := []*fedwiki.PageHeader{}
	for _, info := range list {
		filename := filepath.Join(store.Dir, info.Name())

		slug := fedwiki.Slugify(filepath.Base(filename))

		header, err := pagestore.LoadHeader(filename, slug)
		err = pagestore.ConvertOSError(err)
		//TODO: maybe ignore this error?
		if err != nil {
			return nil, err
		}

		headers = append(headers, header)
	}
	return headers, nil
}
