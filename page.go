package fedwiki

import (
	"fmt"
	"strconv"
	"time"
)

type PageHeader struct {
	Slug     Slug   `json:"slug"`
	Title    string `json:"title"`
	Date     Date   `json:"date"`
	Synopsis string `json:"synopsis,omitempty"`

	// may contain extra information specific to client/server
	Meta Meta `json:"meta,omitempty"`
}

type Meta map[string]interface{}

type Page struct {
	PageHeader
	Story   Story   `json:"story,omitempty"`
	Journal Journal `json:"journal,omitempty"`
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
		page.PageHeader.Date = t
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

func (s Story) IndexOf(id string) (index int, ok bool) {
	for i, item := range s {
		if item.ID() == id {
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

func (s Story) SetByID(id string, item Item) error {
	if i, ok := s.IndexOf(id); ok {
		s[i] = item
		return nil
	}
	return fmt.Errorf("missing item id \"%v\"", id)
}

func (ps *Story) Move(id string, after string) error {
	item, err := ps.RemoveByID(id)
	if err != nil {
		return err
	}
	if after != "" {
		return ps.InsertAfter(after, item)
	}
	ps.Prepend(item)
	return nil
}

func (s *Story) RemoveByID(id string) (item Item, err error) {
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

type Item map[string]interface{}

func (item Item) Val(key string) string {
	if s, ok := item[key].(string); ok {
		return s
	}
	return ""
}

func (item Item) ID() string { return item.Val("id") }

type Date struct{ time.Time }

func NewDate(t time.Time) Date { return Date{t} }

func ParseDate(data string) (Date, error) {
	v, err := strconv.Atoi(data)
	if err != nil {
		return Date{}, err
	}
	return Date{time.Unix(int64(v), 0)}, nil
}

func (d Date) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Itoa(int(d.Unix()))), nil
}

func (d *Date) UnmarshalJSON(data []byte) (err error) {
	*d, err = ParseDate(string(data))
	return
}
