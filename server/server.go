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
