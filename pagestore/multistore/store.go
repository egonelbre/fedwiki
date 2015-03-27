package multistore

import "github.com/egonelbre/fedwiki"

// Store serves/saves in the stores in specified order.
// A store in the order list will handle the page only if it contains the Page.
// Requests not with a specific will be handled by Fallback.
//
// It is usually beneficial to also have Fallback inside the Order.
type Store struct {
	Fallback fedwiki.PageStore
	Order    []fedwiki.PageStore
}

func New(fallback fedwiki.PageStore, order ...fedwiki.PageStore) *Store {
	return &Store{fallback, order}
}

func (store *Store) Exists(slug fedwiki.Slug) bool {
	for _, store := range store.Order {
		if store.Exists(slug) {
			return true
		}
	}
	return false
}

func (store *Store) Load(slug fedwiki.Slug) (*fedwiki.Page, error) {
	for _, store := range store.Order {
		if store.Exists(slug) {
			return store.Load(slug)
		}
	}

	return store.Fallback.Load(slug)
}

func (store *Store) Create(slug fedwiki.Slug, page *fedwiki.Page) error {
	for _, store := range store.Order {
		if store.Exists(slug) {
			return store.Create(slug, page)
		}
	}

	return store.Fallback.Create(slug)
}

func (store *Store) Save(slug fedwiki.Slug, page *fedwiki.Page) error {
	for _, store := range store.Order {
		if store.Exists(slug) {
			return store.Save(slug, page)
		}
	}

	return store.Fallback.Save(slug)
}

// Discards any errors that happen in sub-stores
func (store *Store) List() ([]*fedwiki.PageHeader, error) {
	headers := []*fedwiki.PageHeader{}
	for _, store := range store.Order {
		sub, _ := store.List()
		headers = append(headers, sub)
	}
	return headers, nil
}
