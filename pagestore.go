package fedwiki

import (
	"fmt"
	"regexp"
	"strings"
)

var rxSlug = regexp.MustCompile(`^[a-zA-Z0-9\.\/_-]+$`)

func ValidateSlug(slug Slug) error {
	if !rxSlug.MatchString(string(slug)) {
		return fmt.Errorf(`slug must match /%s/`, rxSlug)
	}
	return nil
}

var (
	rxName     = regexp.MustCompile(`([A-Z][a-z]+)`)
	rxNumber   = regexp.MustCompile(`([0-9]+)`)
	rxToDash   = regexp.MustCompile(`[\-\s]+`)
	rxSlashes  = regexp.MustCompile(`\s*(\/\s*)+`)
	rxTrailing = regexp.MustCompile(`[\/ ]+$`)
	rxRemove   = regexp.MustCompile(`[^a-zA-Z0-9\-\/_ ]`)
)

func Slugify(s string) Slug {
	s = rxName.ReplaceAllString(s, " $1 ")
	s = rxNumber.ReplaceAllString(s, " $1 ")
	s = rxRemove.ReplaceAllString(s, "")
	s = strings.TrimSpace(s)
	s = rxSlashes.ReplaceAllString(s, "/")
	s = rxTrailing.ReplaceAllString(s, "")
	s = rxToDash.ReplaceAllString(s, "-")
	s = strings.ToLower(s)
	return Slug(s)
}

type PageStore interface {
	Exists(slug Slug) bool
	Load(slug Slug) (*Page, error)
	Save(slug Slug, page *Page) error
	List() ([]*PageHeader, error)
}
