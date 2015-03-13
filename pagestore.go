package fedwiki

type PageStore interface {
	Exists(slug Slug) bool
	Load(slug Slug) (*Page, error)
	Save(slug Slug, page *Page) error
	List() ([]*PageHeader, error)
}
