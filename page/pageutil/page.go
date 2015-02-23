package pageutil

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/egonelbre/fedwiki/page"
)

func Load(filename string, slug page.Slug) (*page.Page, error) {
	data, err := ioutil.ReadFile(filename)
	err = ConvertOSError(err)
	if err != nil {
		return nil, err
	}

	p := &page.Page{}
	err = json.Unmarshal(data, p)
	if err != nil {
		return nil, err
	}

	if p.Header.Date.IsZero() {
		if info, err := os.Stat(filename); err == nil {
			p.Header.Date = page.Date{info.ModTime()}
		} else {
			p.Header.Date = page.Date{p.LastModified()}
		}
	}

	p.Header.Slug = slug
	return p, nil
}

func LoadHeader(filename string, slug page.Slug) (*page.Header, error) {
	data, err := ioutil.ReadFile(filename)
	err = ConvertOSError(err)
	if err != nil {
		return nil, err
	}

	header := &page.Header{}
	err = json.Unmarshal(data, header)
	if err != nil {
		return nil, err
	}

	if header.Date.IsZero() {
		if info, err := os.Stat(filename); err == nil {
			header.Date = page.Date{info.ModTime()}
		}
	}

	header.Slug = slug
	return header, nil
}

func Save(page *page.Page, filename string) error {
	data, err := json.Marshal(page)
	if err != nil {
		return err
	}

	return ConvertOSError(ioutil.WriteFile(filename, data, 0755))
}
