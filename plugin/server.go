package plugin

import (
	"log"
	"net/http"
	"path/filepath"
	"sync"

	"github.com/egonelbre/fedwiki"
)

type Factory struct {
	Name     string
	Title    string
	Category string
}

type Plugin struct {
	Name      string
	Folder    string
	Factories []*Factory
}

// plugin Server implements interfaces
//   fedwiki.Handler
// 		/system/factories
// 		/system/plugins
//   http.Handler
// 		/plugin/*
type Server struct {
	Glob string

	mu      sync.RWMutex
	plugins []*Plugin
}

func (server *Server) Update() {
	server.mu.Lock()
	defer server.mu.Unlock()

	matches, err := filepath.Glob(server.Glob)
	if err != nil {
		log.Printf("Failed to update plugin list: %v", err)
	}
}

func (server *Server) Handle(r *http.Request) (code int, template string, data interface{}) {
	switch r.URL.Path {
	case "/system/plugins":
		server.mu.RLock()
		defer server.mu.RUnlock()

		plugins := make([]string, 0, 10)
		for _, plugin := range server.plugins {
			plugins = append(plugins, plugin.Name)
		}

		return http.StatusOK, "plugins", plugins
	case "/system/factories":
		server.mu.RLock()
		defer server.mu.RUnlock()

		factories := make([]*Factory, 0, 10)
		for _, plugin := range server.plugins {
			factories = append(factories, plugin.Factories...)
		}

		return http.StatusOK, "factories", factories
	}

	return fedwiki.ErrorResponse(http.StatusNotFound, "Page not found.")
}

func (server *Server) ServeHTTP(rw http.ResponseWriter, r *http.Request) {

}
