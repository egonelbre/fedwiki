package fedwiki

import (
	"fmt"
	"unicode"
)

func ValidateSlug(slug Slug) error {
	if len(slug) == 0 {
		return fmt.Errorf("slug cannot be empty")
	}

	for _, r := range slug {
		switch {
		case r == '/':
		case r == '-':
		case unicode.IsNumber(r):
		case unicode.IsLetter(r):
		default:
			return fmt.Errorf(`slug must only contain letters, numbers, - or /`)
		}
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
