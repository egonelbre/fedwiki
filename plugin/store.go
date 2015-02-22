package plugin

import (
	"path/filepath"

	"github.com/egonelbre/fedwiki/page/globstore"
)

type Store struct {
	Dir string
	globstore.Store
}

func New(folder string) *Store {
	return &Store{
		Dir:   folder,
		Store: globstore.New(filepath.Join(folder, "*", "pages", "*")),
	}
}
