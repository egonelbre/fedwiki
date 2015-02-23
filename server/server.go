package server

import "net/http"

type Handler interface {
	Handle(rw http.ResponseWriter, r *http.Request) *Response
}

type Server struct {
	Handler  Handler
	Renderer Renderer
}

func New(handler Handler, renderer Renderer) *Server {
	return &Server{
		Handler:  handler,
		Renderer: renderer,
	}
}

func (s *Server) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	responseType, code, err := HandleAccept(r)
	if err != nil {
		http.Error(rw, err.Error(), code)
		return
	}

	response := s.Handler.Handle(rw, r)
	RenderCommon(rw, s.Renderer, responseType, response)
}
