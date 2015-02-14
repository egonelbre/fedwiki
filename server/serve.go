package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"

	"github.com/egonelbre/wiki-go-server/page"
	"github.com/egonelbre/wiki-go-server/page/pageutil"
)

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
	if responseType == "" {
		ext := path.Ext(r.URL.Path)
		switch ext {
		case ".json":
			r.URL.Path = r.URL.Path[:len(r.URL.Path)-len(ext)]
			responseType = "application/json"
		case ".html":
			r.URL.Path = r.URL.Path[:len(r.URL.Path)-len(ext)]
			responseType = "text/html"
		}
	}

	response, code := s.handlePage(rw, r)

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

func (s *Server) handlePage(rw http.ResponseWriter, r *http.Request) (response interface{}, code int) {
	slug := page.Slug(r.URL.Path)
	if err := page.ValidateSlug(slug); err != nil {
		return Error(http.StatusBadRequest, err.Error())
	}

	if r.Method != "GET" && !s.IsAuthorized(r) {
		return Error(http.StatusUnauthorized, "Unauthorized request.")
	}

	switch r.Method {
	case "GET":
		p, err := s.Pages.Load(slug)
		if err != nil {
			if err == page.ErrNotExist {
				return Errorf(http.StatusNotFound, `Page "%s" does not exist.`, slug)
			}
			return Error(http.StatusInternalServerError, err.Error())
		}
		p.Slug = slug
		return p, http.StatusOK
	case "PUT":
		var p *page.Page
		var err error

		switch r.Header.Get("Content-Type") {
		case "":
			fallthrough
		case "application/json":
			p, err = pageutil.Read(r.Body)
			r.Body.Close()
			if err != nil {
				return Error(http.StatusBadRequest, err.Error())
			}
		default:
			return Errorf(http.StatusBadRequest, `Invalid request Content-Type "%s".`, r.Header.Get("Content-Type"))
		}

		err = s.Pages.Save(slug, p)
		if err != nil {
			return Error(http.StatusInternalServerError, err.Error())
		}
		return nil, http.StatusOK
	case "PATCH":
		var action page.Action
		var err error

		switch r.Header.Get("Content-Type") {
		case "":
			fallthrough
		case "application/json":
			action, err = pageutil.ReadAction(r.Body)
			r.Body.Close()
			if err != nil {
				return Error(http.StatusBadRequest, err.Error())
			}
		default:
			return Errorf(http.StatusBadRequest, `Invalid request Content-Type "%s".`, r.Header.Get("Content-Type"))
		}

		p, err := s.Pages.Load(slug)
		if err != nil {
			if err == page.ErrNotExist {
				return Errorf(http.StatusNotFound, `Page "%s" does not exist.`, slug)
			}
			return Error(http.StatusInternalServerError, err.Error())
		}

		if err := p.Apply(action); err != nil {
			return Error(http.StatusInternalServerError, err.Error())
		}
		return nil, http.StatusOK
	default:
		return Errorf(http.StatusNotAcceptable, `Unknown request Method "%s".`, r.Method)
	}
}
