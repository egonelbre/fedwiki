package server

import (
	"fmt"
	"mime"
	"net/http"
	"path"
	"strings"
)

type AcceptSpec []string

func (spec AcceptSpec) Accepts(mimetype string) bool {
	for _, mtype := range spec {
		if mtype == mimetype {
			return true
		}
	}
	return false
}

func ParseAccept(r *http.Request) AcceptSpec {
	var spec AcceptSpec

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

func HandleAccept(r *http.Request) (responseType string, code int, err error) {
	if r.Header.Get("Accept") != "" {
		spec := ParseAccept(r)
		switch {
		case spec.Accepts("application/json"):
			responseType = "application/json"
		case spec.Accepts("text/html"):
			responseType = "text/html"
		default:
			return "", http.StatusNotAcceptable, fmt.Errorf(`Unknown Accept header "%s".`, r.Header.Get("Accept"))
		}
	}

	// back-comp with older clients
	ext := path.Ext(r.URL.Path)
	switch ext {
	case ".json":
		r.URL.Path = r.URL.Path[:len(r.URL.Path)-len(ext)]
		responseType = "application/json"
	case ".html":
		r.URL.Path = r.URL.Path[:len(r.URL.Path)-len(ext)]
		responseType = "text/html"
	}

	return responseType, 0, nil
}
