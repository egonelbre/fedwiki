package page

import "time"

type Header struct {
	Slug     Slug      `json:"slug"`
	Title    string    `json:"title"`
	Date     time.Time `json:"date"`
	Synopsis string    `json:"synopsis"`
}

type Page struct {
	Header  Header
	Story   Story   `json:"story"`
	Journal Journal `json:"journal"`
}

type Story []Item
type Journal []Action

func (page *Page) Apply(action Action) error {
	fn, ok := actionfns[action.Type()]
	if !ok {
		return ErrUnknownAction
	}

	err := fn(page, action)
	if err != nil {
		return err
	}

	if t, err := action.Date(); err == nil {
		page.Header.Date = t
	}
	return nil
}
