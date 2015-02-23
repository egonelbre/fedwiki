package page

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/egonelbre/fedwiki/server"
)

type Service struct {
	Store
}

func (pages Service) Handle(rw http.ResponseWriter, r *http.Request) *server.Response {
	slug := Slug(r.URL.Path)
	if err := ValidateSlug(slug); err != nil {
		return server.Errorf(http.StatusBadRequest, err.Error())
	}

	switch r.Method {
	case "GET":
		p, err := pages.Load(slug)
		if err != nil {
			if err == ErrNotExist {
				return server.Errorf(http.StatusNotFound, `Page "%s" does not exist.`, slug)
			}
			return server.Errorf(http.StatusInternalServerError, err.Error())
		}
		p.Slug = slug
		return server.StatusOK(p)
	case "PUT":
		var p *Page
		var err error

		switch r.Header.Get("Content-Type") {
		case "":
			fallthrough
		case "application/json":
			p, err = readJSONPage(r.Body)
			r.Body.Close()
			if err != nil {
				return server.Errorf(http.StatusBadRequest, err.Error())
			}
		default:
			return server.Errorf(http.StatusBadRequest, `Invalid request Content-Type "%s".`, r.Header.Get("Content-Type"))
		}

		err = pages.Save(slug, p)
		if err != nil {
			return server.Errorf(http.StatusInternalServerError, err.Error())
		}
		return server.StatusOK(nil)
	case "PATCH":
		var action Action
		var err error

		switch r.Header.Get("Content-Type") {
		case "":
			fallthrough
		case "application/json":
			action, err = readJSONAction(r.Body)
			r.Body.Close()
			if err != nil {
				return server.Errorf(http.StatusBadRequest, err.Error())
			}
		default:
			return server.Errorf(http.StatusBadRequest, `Invalid request Content-Type "%s".`, r.Header.Get("Content-Type"))
		}

		p, err := pages.Load(slug)
		if err != nil {
			if err == ErrNotExist {
				return server.Errorf(http.StatusNotFound, `Page "%s" does not exist.`, slug)
			}
			return server.Errorf(http.StatusInternalServerError, err.Error())
		}

		if err := p.Apply(action); err != nil {
			return server.Errorf(http.StatusInternalServerError, err.Error())
		}
		return server.StatusOK(p)
	default:
		return server.Errorf(http.StatusNotAcceptable, `Unknown request Method "%s".`, r.Method)
	}
}

func readJSONPage(r io.Reader) (*Page, error) {
	dec := json.NewDecoder(r)
	page := &Page{}
	err := dec.Decode(page)
	if err != nil {
		return nil, err
	}
	return page, nil
}

func readJSONAction(r io.Reader) (Action, error) {
	dec := json.NewDecoder(r)
	action := make(Action)
	err := dec.Decode(&action)
	if err != nil {
		return nil, err
	}

	if validator, ok := ActionSpecs[action.Type()]; ok {
		if err := validator.Validate(action); err != nil {
			return nil, err
		}
	}

	return action, nil
}
