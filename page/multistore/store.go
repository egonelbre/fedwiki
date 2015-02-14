package multistore

import "github.com/egonelbre/wiki-go-server/page"

// Store serves/saves in the stores in specified order.
// A store in the order list will handle the page only if it contains the page.
// Requests not with a specific will be handled by Fallback.
//
// It is usually beneficial to also have Fallback inside the Order.
type Store struct {
	Fallback page.Store
	Order    []page.Store
}

func New(fallback page.Store, order ...page.Store) *Store {
	return &Store{fallback, order}
}

func (store *Store) Exists(slug page.Slug) bool {
	for _, store := range store.Order {
		if store.Exists(slug) {
			return true
		}
	}
	return false
}

func (store *Store) Load(slug page.Slug) (*page.Page, error) {
	for _, store := range store.Order {
		if store.Exists(slug) {
			return store.Load(slug)
		}
	}

	return store.Fallback.Load(slug)
}

func (store *Store) Save(slug page.Slug, page *page.Page) error {
	for _, store := range store.Order {
		if store.Exists(slug) {
			return store.Save(slug, page)
		}
	}

	return store.Fallback.Save(slug)
}

// Discards any errors that happen in sub-stores
func (store *Store) List() ([]*page.Header, error) {
	headers := []*page.Header{}
	for _, store := range store.Order {
		sub, _ := store.List()
		headers = append(headers, sub)
	}
	return headers, nil
}
