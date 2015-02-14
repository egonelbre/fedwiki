package server

import (
	"encoding/json"
	"fmt"
	"mime"
	"net/http"
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
			fmt.Fprintf(response)
		}
	}
}

type ErrorResponse struct {
	Status string `json:"status"`
	Code   int    `json:"code"`
	Detail string `json:"detail"`
}

func Error(code int, detail string) (r interface{}, rcode int) {
	return ErrorResponse{
		Status: http.StatusText(code),
		Code:   code,
		Detail: detail,
	}, code
}

func Errorf(code int, format string, args ...interface{}) (r interface{}, rcode int) {
	return ErrorResponse{
		Status: http.StatusText(code),
		Code:   code,
		Detail: fmt.Sprintf(format, args...),
	}, code
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
