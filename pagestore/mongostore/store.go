package mongostore

import (
	"github.com/egonelbre/fedwiki"

	"gopkg.in/mgo.v2"
)

var _ fedwiki.PageStore = (*Store)(nil)

type Store struct {
	session    *mgo.Session
	collection string
}

func translate(err error) error {
	if err == nil {
		return nil
	}
	if err == mgo.ErrNotFound {
		return fedwiki.ErrNotExist
	}
	return err
}

func New(url, collection string) (*Store, error) {
	main, err := mgo.Dial(url)
	if err != nil {
		return nil, err
	}
	return &Store{main, collection}, nil
}

func (store *Store) c() (*mgo.Session, *mgo.Collection) {
	s := store.session.Copy()
	return s, s.DB("").C(store.collection)
}

func (store *Store) Exists(slug fedwiki.Slug) bool {
	session, c := store.c()
	defer session.Close()

	n, err := c.FindId(slug).Count()
	return (n > 0) && (err != nil)
}

func (store *Store) Create(slug fedwiki.Slug, page *fedwiki.Page) error {
	session, c := store.c()
	defer session.Close()

	return translate(c.Insert(page))
}

func (store *Store) Load(slug fedwiki.Slug) (*fedwiki.Page, error) {
	session, c := store.c()
	defer session.Close()

	page := &fedwiki.Page{}
	if err := c.FindId(slug).One(page); err != nil {
		return nil, translate(err)
	}

	return page, nil
}

func (store *Store) Save(slug fedwiki.Slug, page *fedwiki.Page) error {
	session, c := store.c()
	defer session.Close()

	return translate(c.UpdateId(slug, page))
}

func (store *Store) Delete(slug fedwiki.Slug) error {
	session, c := store.c()
	defer session.Close()

	return translate(c.RemoveId(slug))
}

func (store *Store) List() ([]*fedwiki.PageHeader, error) {
	session, c := store.c()
	defer session.Close()

	var headers []*fedwiki.PageHeader
	if err := c.Find(nil).All(&headers); err != nil {
		return nil, translate(err)
	}
	return headers, nil
}
