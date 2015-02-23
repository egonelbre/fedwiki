package server

import "net/http"

type Service interface {
	Handle(rw http.ResponseWriter, r *http.Request) *Response
}

type Server struct {
	Service  Service
	Renderer Renderer
}

func New(service Service, renderer Renderer) *Server {
	return &Server{
		Service:  service,
		Renderer: renderer,
	}
}

func (s *Server) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	responseType, code, err := HandleAccept(r)
	if err != nil {
		http.Error(rw, err.Error(), code)
		return
	}

	response := s.Service.Handle(rw, r)
	RenderCommon(rw, s.Renderer, responseType, response)
}
