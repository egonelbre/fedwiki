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

	helpers = template.FuncMap{
		"resolve": func(s string) template.HTML {
			s = template.HTMLEscapeString(s)
			s = rxInternal.ReplaceAllStringFunc(s, func(s string) string {
				s = strings.Trim(s, "[]")
				return fmt.Sprintf(`<a href="%s">%s</a>`, page.Slugify(s), s)
			})
			s = rxExternal.ReplaceAllString(s, `<a href="$1">$2</a>`)
			return template.HTML(s)
		},
	}
)

func (s *Server) RenderHTML(w io.Writer, data interface{}) {
	t, err := template.New("").Funcs(helpers).ParseGlob(filepath.Join("views", "*.html"))
	if err != nil {
		fmt.Fprintf(w, "<h1>Error</h1><p>%v</p>", err)
		return
	}

	if err := t.ExecuteTemplate(w, "index.html", data); err != nil {
		fmt.Fprintf(w, "<h1>Error</h1><p>%v</p>", err)
	}
}
