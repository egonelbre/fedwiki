package server

import (
	"net/http"

	"github.com/egonelbre/fedwiki/page"
)

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
