package pageutil

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"os"

	"github.com/egonelbre/wiki-go-server/page"
)

func Read(r io.Reader) (*page.Page, error) {
	dec := json.NewDecoder(r)
	page := &page.Page{}
	err := dec.Decode(page)
	if err != nil {
		return nil, err
	}
	return page, nil
}

func Load(filename string) (*page.Page, error) {
	data, err := ioutil.ReadFile(filename)
	err = ConvertOSError(err)
	if err != nil {
		return nil, err
	}

	page := &page.Page{}
	err = json.Unmarshal(data, page)
	if err != nil {
		return nil, err
	}

	if page.Header.Date.IsZero() {
		if info, err := os.Stat(filename); err == nil {
			page.Header.Date = info.ModTime()
		} else {
			page.Header.Date = page.LastModified()
		}
	}

	return page, nil
}

func LoadHeader(filename string) (*page.Header, error) {
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
	return header, nil
}

func Save(page *page.Page, filename string) error {
	data, err := json.Marshal(page)
	if err != nil {
		return err
	}

	return ConvertOSError(ioutil.WriteFile(filename, data, 0755))
}
