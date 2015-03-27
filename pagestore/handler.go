package pagestore

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/egonelbre/fedwiki"
)

type Handler struct {
	fedwiki.PageStore
}

func (pages Handler) Handle(r *http.Request) (code int, template string, data interface{}) {
	slug := fedwiki.Slug(r.URL.Path)
	if err := fedwiki.ValidateSlug(slug); err != nil {
		return fedwiki.ErrorResponse(http.StatusBadRequest, err.Error())
	}

	switch r.Method {
	case "GET":
		page, err := pages.Load(slug)
		if err != nil {
			if err == fedwiki.ErrNotExist {
				return fedwiki.ErrorResponse(http.StatusNotFound, `Page "%s" does not exist.`, slug)
			}
			return fedwiki.ErrorResponse(http.StatusInternalServerError, err.Error())
		}
		page.Slug = slug
		return http.StatusOK, "", page
	case "PUT":
		var page *fedwiki.Page
		var err error

		switch r.Header.Get("Content-Type") {
		case "":
			fallthrough
		case "application/json":
			page, err = readJSONPage(r.Body)
			r.Body.Close()
			if err != nil {
				return fedwiki.ErrorResponse(http.StatusBadRequest, err.Error())
			}
		default:
			return fedwiki.ErrorResponse(http.StatusBadRequest, `Invalid request Content-Type "%s".`, r.Header.Get("Content-Type"))
		}

		page.Slug = slug
		page.Date = fedwiki.NewDate(time.Now())

		err = pages.Create(slug, page)
		if err != nil {
			return fedwiki.ErrorResponse(http.StatusInternalServerError, err.Error())
		}
		return http.StatusOK, "", page
	case "PATCH":
		var action fedwiki.Action
		var err error

		switch r.Header.Get("Content-Type") {
		case "":
			fallthrough
		case "application/json":
			action, err = readJSONAction(r.Body)
			r.Body.Close()
			if err != nil {
				return fedwiki.ErrorResponse(http.StatusBadRequest, err.Error())
			}
		default:
			return fedwiki.ErrorResponse(http.StatusBadRequest, `Invalid request Content-Type "%s".`, r.Header.Get("Content-Type"))
		}

		page, err := pages.Load(slug)
		if err != nil {
			if err == fedwiki.ErrNotExist {
				return fedwiki.ErrorResponse(http.StatusNotFound, `Page "%s" does not exist.`, slug)
			}
			return fedwiki.ErrorResponse(http.StatusInternalServerError, err.Error())
		}

		if err := page.Apply(action); err != nil {
			return fedwiki.ErrorResponse(http.StatusInternalServerError, err.Error())
		}

		if err := pages.Save(slug, page); err != nil {
			return fedwiki.ErrorResponse(http.StatusInternalServerError, err.Error())
		}

		return http.StatusOK, "", page
	default:
		return fedwiki.ErrorResponse(http.StatusNotAcceptable, `Unknown request Method "%s".`, r.Method)
	}
}

func readJSONPage(r io.Reader) (*fedwiki.Page, error) {
	dec := json.NewDecoder(r)
	page := &fedwiki.Page{}
	err := dec.Decode(page)
	if err != nil {
		return nil, err
	}
	return page, nil
}

func readJSONAction(r io.Reader) (fedwiki.Action, error) {
	dec := json.NewDecoder(r)
	action := make(fedwiki.Action)
	err := dec.Decode(&action)
	if err != nil {
		return nil, err
	}

	return action, nil
}
