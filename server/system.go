package server

import "net/http"

func (s *Server) handleSystem(rw http.ResponseWriter, r *http.Request) (response interface{}, code int) {
	return Error(http.StatusNotImplemented, "not implemented yet")
}
