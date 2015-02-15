package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"path/filepath"
	"strings"
)

const (
	StaticRoute = "/static/"
	ClientRoute = "/client/"
)

func (s *Server) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		s.serveIndex(rw, r)
		return
	} else if strings.HasPrefix(r.URL.Path, StaticRoute) {
		s.serveStatic(rw, r)
		return
	} else if strings.HasPrefix(r.URL.Path, ClientRoute) {
		s.serveClient(rw, r)
		return
	}

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
	if strings.HasPrefix(r.URL.Path, "/system/") {
		response, code = s.handleSystem(rw, r)
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
		s.RenderHTML(rw, response)
	}
}

func (s *Server) serveStatic(rw http.ResponseWriter, r *http.Request) {
	upath := filepath.Join(s.StaticDir, r.URL.Path[len(StaticRoute):])
	http.ServeFile(rw, r, path.Clean(upath))
}

func (s *Server) serveClient(rw http.ResponseWriter, r *http.Request) {
	upath := filepath.Join(s.ClientDir, r.URL.Path[len(ClientRoute):])
	http.ServeFile(rw, r, path.Clean(upath))
}

func (s *Server) serveIndex(rw http.ResponseWriter, r *http.Request) {
	s.RenderTemplate(rw, "client.html", s.Title)
}
