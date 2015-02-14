package server

import (
	"encoding/json"
	"fmt"
	"mime"
	"net/http"
	"path"
)

type ContentFunc func(rw http.ResponseWriter, r *http.Request) (response interface{}, code int)

func HandleAcceptHeader(fn ContentFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		responseType := ""

		if r.Header.Get("Accept") != "" {
			switch {
			case Accepts(r, "application/json"):
				responseType = "application/json"
			case Accepts(r, "text/html"):
				responseType = "text/html"
			default:
				http.Error(rw, fmt.Sprintf(`Unknown Accept header "%s".`, r.Header.Get("Accept")), http.StatusNotAcceptable)
				return
			}
		}

		// back-comp with older clients
		if responseType == "" {
			ext := path.Ext(r.URL.Path)
			switch ext {
			case ".json":
				r.URL.Path = r.URL.Path[:len(r.URL.Path)-len(ext)]
				responseType = "application/json"
			case ".html":
				r.URL.Path = r.URL.Path[:len(r.URL.Path)-len(ext)]
				responseType = "text/html"
			}
		}

		response, code := fn(rw, r)

		switch responseType {
		case "application/json":
			enc := json.NewEncoder(rw)
			rw.Header().Set("Content-Type", "application/json")
			rw.WriteHeader(code)
			enc.Encode(response)
		case "text/html":
			fallthrough
		default:
			rw.Header().Set("Content-Type", "text/html")
			rw.WriteHeader(code)
			fmt.Fprintf(rw, "%#v\n", response)
		}
	}
}

func Accepts(r *http.Request, mimetype string) bool {
	for _, accept := range r.Header[http.CanonicalHeaderKey("Accept")] {
		m, _, err := mime.ParseMediaType(accept)
		if err != nil {
			continue
		}
		if m == mimetype {
			return true
		}
	}
	return false
}
