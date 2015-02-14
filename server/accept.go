package server

import (
	"mime"
	"net/http"
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
