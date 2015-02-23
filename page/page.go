package page

import (
	"fmt"
	"time"
)

type Header struct {
	Slug     Slug   `json:"slug"`
	Title    string `json:"title"`
	Date     Date   `json:"date"`
	Synopsis string `json:"synopsis"`
}

type Page struct {
	Header
	Story   Story   `json:"story"`
	Journal Journal `json:"journal"`
}

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

// if no date is found then it will use the current time!
func (page *Page) LastModified() time.Time {
	if !page.Date.IsZero() {
		return page.Date.Time
	}

	for _, action := range page.Journal {
		if t, err := action.Date(); err == nil && !t.IsZero() {
			return t.Time
		}
	}

	return time.Now()
}

type Story []Item

func (s Story) IndexOf(id string) (index int, ok bool) {
	for i, item := range s {
		if item.Id() == id {
			return i, true
		}
	}
	return -1, false
}

func (s *Story) insertAt(i int, item Item) {
	t := *s
	t = append(t, Item{})
	copy(t[i+1:], t[i:])
	t[i] = item
	*s = t
}

func (s *Story) Prepend(item Item) {
	s.insertAt(0, item)
}

func (s *Story) Append(item Item) {
	*s = append(*s, item)
}

func (s *Story) InsertAfter(after string, item Item) error {
	if i, ok := s.IndexOf(after); ok {
		s.insertAt(i+1, item)
		return nil
	}
	return fmt.Errorf("missing item id \"%v\"", after)
}

func (s Story) SetById(id string, item Item) error {
	if i, ok := s.IndexOf(id); ok {
		s[i] = item
		return nil
	}
	return fmt.Errorf("missing item id \"%v\"", id)
}

func (ps *Story) Move(id string, after string) error {
	item, err := ps.RemoveById(id)
	if err != nil {
		return err
	}
	if after != "" {
		return ps.InsertAfter(after, item)
	}
	ps.Prepend(item)
	return nil
}

func (s *Story) RemoveById(id string) (item Item, err error) {
	if i, ok := s.IndexOf(id); ok {
		t := *s
		item = t[i]
		copy(t[i:], t[i+1:])
		t = t[:len(t)-1]
		*s = t
		return item, nil
	}
	return item, fmt.Errorf("missing item id \"%v\"", id)
}

type Journal []Action
