package server

import (
	"net/http"

	"github.com/egonelbre/wiki-go-server/page"
)

type Server struct {
	Title string

	Pages   page.Store
	Sitemap *Sitemap

	StaticDir string
}

func New(pages page.Store) *Server {
	return &Server{
		Title:     "",
		Pages:     pages,
		Sitemap:   NewSitemap(pages),
		StaticDir: "static",
	}
}

func (s *Server) IsAuthorized(r *http.Request) bool {
	return true
}
