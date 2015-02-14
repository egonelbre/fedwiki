package server

import (
	"net/http"

	"github.com/egonelbre/wiki-go-server/page"
)

type Server struct {
	Pages page.Store
}

func New(pages page.Store) *Server {
	return &Server{pages}
}

func (s *Server) IsAuthorized(r *http.Request) bool {
	return true
}

func (s *Server) ServePage(rw http.ResponseWriter, r *http.Request) (response interface{}, code int) {
	slug := page.Slug(r.URL.Path)
	if err := page.ValidateSlug(slug); err != nil {
		return Error(http.StatusBadRequest, err.Error())
	}

	if r.Method != "GET" && !s.IsAuthorized(r) {
		return Error(http.StatusUnauthorized, "Unauthorized request.")
	}

	switch r.Method {
	case "GET":
		p, err := s.Pages.Get(slug)
		if err != nil {
			if err == page.IsNotExist {
				return Errorf(http.StatusNotFound, `Page "%s" does not exist.`, slug)
			}
			return Error(http.StatusInternalServerError, err.Error())
		}
		return p, http.StatusOK
	case "PUT":
		p, err := s.Pages.Set(slug)
		return p, http.StatusOK
	case "PATCH":
		if err != nil {
			if err == page.IsNotExist {
				return Errorf(http.StatusNotFound, `Page "%s" does not exist.`, slug)
			}
			return Error(http.StatusInternalServerError, err.Error())
		}
		p, err := s.Pages.Get(slug)
		return p, http.StatusOK
	default:
		return Errorf(http.StatusNotAcceptable, `Unknown request Method "%s".`, r.Method)
	}
}
