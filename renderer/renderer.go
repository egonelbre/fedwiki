package renderer

import (
	"fmt"
	"html/template"
	"io"
	"regexp"
	"strings"

	"github.com/egonelbre/fedwiki/page"
)

type Renderer struct {
	Glob string
}

func New(glob string) *Renderer {
	return &Renderer{glob}
}

func (r *Renderer) Render(responseType string, w io.Writer, tname string, data interface{}) error {
	t, err := template.New("").Funcs(helpers).ParseGlob(r.Glob)
	if err != nil {
		return err
	}

	if tname == "" {
		tname = "static"
	}

	return t.ExecuteTemplate(w, tname+".html", data)
}

var (
	rxInternal = regexp.MustCompile(`\[\[([^\]]+)\]\]`)
	rxExternal = regexp.MustCompile(`\[((?:http|https|ftp):.*?) (.*?)\]`)
)

func replaceLinks(s string) template.HTML {
	s = rxInternal.ReplaceAllStringFunc(s, func(s string) string {
		s = strings.Trim(s, "[]")
		return fmt.Sprintf(`<a href="%s">%s</a>`, page.Slugify(s), s)
	})
	s = rxExternal.ReplaceAllString(s, `<a href="$1">$2</a>`)
	return template.HTML(s)
}

var (
	helpers = template.FuncMap{
		"resolve": func(s string) template.HTML {
			s = template.HTMLEscapeString(s)
			return replaceLinks(s)
		},
		"html": func(s string) template.HTML {
			return replaceLinks(s)
		},
	}
)
