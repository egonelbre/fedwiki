package server

import "net/http"

func (s *Server) handleSystem(rw http.ResponseWriter, r *http.Request) (response interface{}, code int) {
	return Error(http.StatusNotImplemented, "not implemented yet")

	//TODO:
	//   sitemap.json    json array of page.Header
	//   slugs.json      json array of page.Slug
	//   plugins.json    json array of plugin name
	//   factories.json  json array of factory
	//   export.json     ???
}
