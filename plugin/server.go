package plugin

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/egonelbre/fedwiki"
)

type Factory struct {
	Name     string
	Title    string
	Category string
}

type Plugin struct {
	Name    string
	Folder  string
	Factory *Factory
}

// plugin Server implements interfaces
//   fedwiki.Handler
// 		/system/factories
// 		/system/plugins
//   http.Handler
// 		/plugin/*
type Server struct {
	Dir string

	mu      sync.RWMutex
	plugins map[string]*Plugin
}

func NewServer(dir string) *Server {
	server := &Server{}
	server.Dir = dir
	server.plugins = make(map[string]*Plugin)
	return server
}

func readPlugin(dirname string) (*Plugin, error) {
	plugin := &Plugin{}

	name := filepath.Base(dirname)
	if !strings.HasPrefix(name, "wiki-plugin-") {
		return nil, errors.New("bad folder name")
	}
	name = name[len("wiki-plugin-"):]

	plugin.Name = name
	plugin.Folder = filepath.Join(dirname, "client")

	factory := &Factory{}
	data, err := ioutil.ReadFile(filepath.Join(dirname, "factory.json"))
	if err != nil {
		return plugin, nil
	}

	if err := json.Unmarshal(data, factory); err != nil {
		return plugin, nil
	}

	plugin.Factory = factory
	return plugin, nil
}

func (server *Server) Update() {
	list, err := ioutil.ReadDir(server.Dir)
	if err != nil {
		log.Printf("Failed to update plugin list: %v", err)
		return
	}

	plugins := make(map[string]*Plugin)
	for _, info := range list {
		if !info.IsDir() {
			continue
		}

		plugin, err := readPlugin(filepath.Join(server.Dir, info.Name()))
		if err != nil {
			continue
		}
		plugins[plugin.Name] = plugin
	}

	server.mu.Lock()
	defer server.mu.Unlock()
	server.plugins = plugins
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

		sort.Strings(plugins)

		return http.StatusOK, "plugins", plugins
	case "/system/factories":
		server.mu.RLock()
		defer server.mu.RUnlock()

		factories := make([]*Factory, 0, 10)
		for _, plugin := range server.plugins {
			if plugin.Factory != nil {
				factories = append(factories, plugin.Factory)
			}
		}

		return http.StatusOK, "factories", factories
	}

	return fedwiki.ErrorResponse(http.StatusNotFound, "Page not found.")
}

func (server *Server) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	tokens := strings.SplitN(r.URL.Path, "/", 4)
	if len(tokens) < 4 {
		http.Error(rw, "Plugin not found.", http.StatusNotFound)
		return
	}
	// tokens[0] == ""
	// tokens[1] == "plugins"
	name := tokens[2]
	file := tokens[3]

	server.mu.RLock()
	defer server.mu.RUnlock()

	plugin, ok := server.plugins[name]
	if !ok {
		http.Error(rw, "Plugin not found.", http.StatusNotFound)
		return
	}

	file = filepath.Join(plugin.Folder, filepath.FromSlash(path.Clean("/"+file)))
	http.ServeFile(rw, r, file)
}
