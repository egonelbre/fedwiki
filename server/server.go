package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"strings"

	"github.com/egonelbre/wiki-go-server/page"
)

type Server struct {
	Pages   page.Store
	Systems []System
}

func New(pages page.Store, systems ...System) *Server {
	return &Server{
		Pages:   pages,
		Systems: systems,
	}
}

func (s *Server) IsAuthorized(r *http.Request) bool {
	return true
}

func (s *Server) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	responseType := ""
	if r.Header.Get("Accept") != "" {
		spec := ParseAccept(r)
		switch {
		case spec.Accepts("application/json"):
			responseType = "application/json"
		case spec.Accepts("text/html"):
			responseType = "text/html"
		default:
			http.Error(rw, fmt.Sprintf(`Unknown Accept header "%s".`, r.Header.Get("Accept")), http.StatusNotAcceptable)
			return
		}
	}

	// back-comp with older clients
	ext := path.Ext(r.URL.Path)
	switch ext {
	case ".json":
		r.URL.Path = r.URL.Path[:len(r.URL.Path)-len(ext)]
		responseType = "application/json"
	case ".html":
		r.URL.Path = r.URL.Path[:len(r.URL.Path)-len(ext)]
		responseType = "text/html"
	}

	var response interface{}
	var code int
	var template string

	if strings.HasPrefix(r.URL.Path, "/system/") {
		response, code, template = s.handleSystem(rw, r)
	} else {
		response, code = s.handlePage(rw, r)
	}

	switch responseType {
	case "application/json":
		enc := json.NewEncoder(rw)
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(code)
		enc.Encode(response)
	case "text/html":
		fallthrough
	default:
		rw.Header().Set("Content-Type", "text/html")
		rw.WriteHeader(code)
		if template == "" {
			s.RenderHTML(rw, response)
		} else {
			s.RenderTemplate(rw, template+".html", response)
		}

	}
}
