package folderstore

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/egonelbre/wiki-go-server/page"
	"github.com/egonelbre/wiki-go-server/page/pageutil"
)

type Store struct {
	Dir string
}

func New(dir string) *Store {
	return &Store{dir}
}

func (store *Store) path(slug page.Slug) string {
	filename := filepath.FromSlash(string(slug)) + ".json"
	return filepath.Join(store.Dir, filename)
}

func (store *Store) Exists(slug page.Slug) bool {
	stat, err := os.Stat(store.path(slug))
	return os.IsExist(err) && !stat.IsDir()
}

func (store *Store) Load(slug page.Slug) (*page.Page, error) {
	return pageutil.Load(store.path(slug))
}

func (store *Store) Save(slug page.Slug, page *page.Page) error {
	return pageutil.Save(page, store.path(slug))
}

// Discards any errors that happen in sub-stores
func (store *Store) List() ([]*page.Header, error) {
	list, err := ioutil.ReadDir(store.Dir)
	err = pageutil.ConvertOSError(err)
	if err != nil {
		return nil, err
	}

	headers := []*page.Header{}
	for _, info := range list {
		if filepath.Ext(info.Name()) == ".json" {
			filename := filepath.Join(store.Dir, info.Name())

			header, err := pageutil.LoadHeader(filename)
			err = pageutil.ConvertOSError(err)
			if err != nil {
				return nil, err
			}

			headers = append(headers, header)
		}
	}
	return headers, nil
}
