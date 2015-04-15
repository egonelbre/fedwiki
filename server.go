package fedwiki

import (
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"
	"path"
	"strings"
)

// errorInfo contains relevant information about an error
type errorInfo struct {
	Status string
	Code   int
	Detail string
}

// ErrorResponse creates a response based on http Error code
func ErrorResponse(ecode int, format string, args ...interface{}) (code int, template string, data interface{}) {
	return ecode, "error", errorInfo{
		Status: http.StatusText(ecode),
		Code:   ecode,
		Detail: fmt.Sprintf(format, args...),
	}
}

// Template is the interface that is used to render pages as HTML
type Template interface {
	// RenderHTML renders data as HTML with the appropriate template
	// template = "" means that it should be rendered as a regular page
	RenderHTML(w io.Writer, template string, data interface{}) error
}

// Handler is the interface for handling federated wiki requests
type Handler interface {
	Handle(r *http.Request) (code int, template string, data interface{})
}

// HandlerFunc type adapts a function to be used as a regular handler
type HandlerFunc func(r *http.Request) (code int, template string, data interface{})

func (fn HandlerFunc) Handle(r *http.Request) (code int, template string, data interface{}) {
	return fn(r)
}

// This implements basic management of request of headers and canonicalizes the requests
type Server struct {
	Handler  Handler
	Template Template
}

func (server *Server) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	responseType := ""

	if r.Header.Get("Accept") != "" {
		h := parseAccept(r)
		switch {
		case h.Accepts("*/*"):
			fallthrough
		case h.Accepts("application/json"):
			responseType = "application/json"
		case h.Accepts("text/html"):
			responseType = "text/html"
		case h.Accepts("text/plain"):
			responseType = "text/plain"
		default:
			http.Error(rw, fmt.Sprintf(`Unknown Accept header "%s".`, r.Header.Get("Accept")), http.StatusNotAcceptable)
			return
		}
	}

	// handle explicit extensions
	ext := path.Ext(r.URL.Path)
	switch ext {
	case ".json":
		r.URL.Path = r.URL.Path[:len(r.URL.Path)-len(ext)]
		responseType = "application/json"
	case ".html":
		r.URL.Path = r.URL.Path[:len(r.URL.Path)-len(ext)]
		responseType = "text/html"
	}

	switch {
	case responseType == "" && server.Template == nil:
		responseType = "application/json"
	case responseType == "":
		responseType = "text/html"
	case responseType == "text/html" && server.Template == nil:
		responseType = "application/json"
	}

	code, template, data := server.Handler.Handle(r)

	rw.Header().Set("Content-Type", responseType)
	rw.WriteHeader(code)

	switch responseType {
	case "application/json":
		json.NewEncoder(rw).Encode(data)
	case "text/plain":
		fmt.Fprintf(rw, "%#v\n", data)
	case "text/html":
		err := server.Template.RenderHTML(rw, template, data)
		if err != nil {
			fmt.Fprintf(rw, err.Error())
		}
	default:
		fmt.Fprintf(rw, fmt.Sprintf("Unknown Content-Type \"%v\"", responseType))
	}
}

type acceptHeaders []string

// Accepts checks whether mimetype is allowed
func (spec acceptHeaders) Accepts(mimetype string) bool {
	for _, mtype := range spec {
		if mtype == mimetype {
			return true
		}
	}
	return false
}

// parseAccept parses request Accept header
func parseAccept(r *http.Request) acceptHeaders {
	var spec acceptHeaders

	accepts := r.Header.Get("Accept")
	params := strings.Split(accepts, ";")
	for _, accept := range strings.Split(params[0], ",") {
		m, _, err := mime.ParseMediaType(accept)
		if err != nil {
			continue
		}
		spec = append(spec, strings.TrimSpace(m))
	}
	return spec
}
