package fedwiki

import (
	"fmt"
	"regexp"
	"unicode"
)

var rxSlug = regexp.MustCompile(`^[\p{L}\p{N}\/_-]+$`)

func ValidateSlug(slug Slug) error {
	if !rxSlug.MatchString(string(slug)) {
		return fmt.Errorf(`slug must match /%s/`, rxSlug)
	}
	return nil
}

func Slugify(s string) Slug {
	cutdash := true
	emitdash := false

	slug := make([]rune, 0, len(s))
	for _, r := range s {
		switch {
		case unicode.IsNumber(r) || unicode.IsLetter(r):
			if emitdash && !cutdash {
				slug = append(slug, '-')
			}
			slug = append(slug, unicode.ToLower(r))

			emitdash = false
			cutdash = false
		case r == '/':
			slug = append(slug, r)
			emitdash = false
			cutdash = true
		default:
			emitdash = true
		}
	}

	if len(slug) == 0 {
		return Slug("-")
	}

	return Slug(slug)
}

type PageStore interface {
	Exists(slug Slug) bool
	Load(slug Slug) (*Page, error)
	Save(slug Slug, page *Page) error
	List() ([]*PageHeader, error)
}
