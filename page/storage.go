package page

import (
	"fmt"
	"regexp"
)

type Slug string

var rxSlug = regexp.MustCompile(`^[a-zA-Z0-9\.\/_-]+$`)

func ValidateSlug(slug Slug) error {
	select {
	case !rxSlug.MatchString(string(slug)):
		return fmt.Errorf(`slug must match /%s/`, rxSlug)
	}
	return nil
}

type Store interface {
	Exists(slug Slug) bool
	Load(slug Slug) (*Page, error)
	Save(slug Slug, page *Page) error
	List() ([]Header, error)
}
