package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path"
	"strings"

	"github.com/egonelbre/fedwiki/page"
)

type Renderer interface {
	Render(responseType string, w io.Writer, template string, data interface{}) error
}

type System interface {
	// if the system doesn't want to handle the request it should
	//   return nil, http.StatusNotFound, ""
	// template is an optional argument to give to the renderer
	Handle(rw http.ResponseWriter, r *http.Request) (response interface{}, code int, template string)
}

// PageSystem is a System that cares about page changes
//   note that stores can change even if there isn't an explicit notification
type PageSystem interface {
	System
	PageChanged(p *page.Page, err error)
}

type Server struct {
	Pages    page.Store
	Renderer Renderer
	Systems  []System
}

func New(pages page.Store, renderer Renderer, systems ...System) *Server {
	return &Server{
		Pages:    pages,
		Renderer: renderer,
		Systems:  systems,
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
		if responseType == "" {
			responseType = "text/html"
		}

		if s.Renderer == nil {
			rw.Header().Set("Content-Type", "text/plain")
			rw.WriteHeader(code)
			fmt.Fprintf(rw, "%#v\n", response)
			return
		}

		rw.Header().Set("Content-Type", "text/html")
		rw.WriteHeader(code)

		if template == "" {
			if _, ok := response.(ErrorResponse); ok {
				template = "error"
			}
		}

		err := s.Renderer.Render(responseType, rw, template, response)
		if err != nil {
			fmt.Fprintf(rw, err.Error())
		}
	}
}

func (s *Server) handleSystem(rw http.ResponseWriter, r *http.Request) (response interface{}, code int, template string) {
	for _, sys := range s.Systems {
		resp, code, template := sys.Handle(rw, r)
		if code != http.StatusNotFound {
			return resp, code, template
		}
	}
	err, code := Error(http.StatusNotFound, "system was not found")
	return err, code, ""
}

func (s *Server) PageChanged(p *page.Page, err error) {
	for _, sys := range s.Systems {
		if sys, ok := sys.(PageSystem); ok {
			sys.PageChanged(p, err)
		}
	}
}
