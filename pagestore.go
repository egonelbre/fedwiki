package fedwiki

type PageStore interface {
	// Exists checks whether a page with `slug` exists
	Exists(slug Slug) bool
	// Create adds a new `page` with `slug`
	Create(slug Slug, page *Page) error
	// Load loads the page with identified by `slug`
	Load(slug Slug) (*Page, error)
	// Save saves the new page to `slug`
	Save(slug Slug, page *Page) error
	// List lists all the page headers
	List() ([]*PageHeader, error)
}
