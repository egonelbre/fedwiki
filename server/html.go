package server

import (
	"fmt"
	"html/template"
	"io"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/egonelbre/wiki-go-server/page"
)

var (
	rxInternal = regexp.MustCompile(`\[\[([^\]]+)\]\]`)
	rxExternal = regexp.MustCompile(`\[((?:http|https|ftp):.*?) (.*?)\]`)
)

func resolve(s string) template.HTML {
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
			return resolve(s)
		},
		"html": func(s string) template.HTML {
			return resolve(s)
		},
	}
)

func (s *Server) RenderTemplate(w io.Writer, tname string, data interface{}) {
	t, err := template.New("").Funcs(helpers).ParseGlob(filepath.Join("views", "*.html"))
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}

	if err := t.ExecuteTemplate(w, tname, data); err != nil {
		fmt.Fprintf(w, err.Error())
	}
}

func (s *Server) RenderHTML(w io.Writer, data interface{}) {
	if _, ok := data.(ErrorResponse); ok {
		s.RenderTemplate(w, "error.html", data)
		return
	}
	s.RenderTemplate(w, "static.html", data)
}
