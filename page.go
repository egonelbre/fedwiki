package fedwiki

import (
	"fmt"
	"strconv"
	"time"
)

// PageHeader represents minimal useful information about the page
type PageHeader struct {
	Slug     Slug   `json:"slug" bson:"_id"`
	Title    string `json:"title"`
	Date     Date   `json:"date"`
	Synopsis string `json:"synopsis,omitempty"`

	// may contain extra information specific to client/server
	Meta Meta `json:"meta,omitempty"`
}

// Meta is used for additional page properties not in fedwiki spec
type Meta map[string]interface{}

// Page represents a federated wiki page
type Page struct {
	PageHeader `bson:",inline"`
	Story      Story   `json:"story,omitempty"`
	Journal    Journal `json:"journal,omitempty"`
}

// Story is the viewable content of the page
type Story []Item

// Journal contains the history of the Page
type Journal []Action

// Apply modifies the page with an action
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

// LastModified returns the date when the page was last modified
// if there is no such date it will return a zero time
func (page *Page) LastModified() time.Time {
	if !page.Date.IsZero() {
		return page.Date.Time
	}

	for i := len(page.Journal) - 1; i >= 0; i-- {
		if t, err := page.Journal[i].Date(); err == nil && !t.IsZero() {
			return t.Time
		}
	}

	return page.Date.Time
}

// IndexOf returns the index of an item with `id`
// ok = false, if that item doesn't exist
func (s Story) IndexOf(id string) (index int, ok bool) {
	for i, item := range s {
		if item.ID() == id {
			return i, true
		}
	}
	return -1, false
}

// insertAt adds an `item` after position `i`
func (s *Story) insertAt(i int, item Item) {
	t := *s
	t = append(t, Item{})
	copy(t[i+1:], t[i:])
	t[i] = item
	*s = t
}

// Prepend adds the `item` as the first item in story
func (s *Story) Prepend(item Item) {
	s.insertAt(0, item)
}

// Appends adds the `item` as the last item in story
func (s *Story) Append(item Item) {
	*s = append(*s, item)
}

// InsertAfter adds the `item` after the item with `id`
func (s *Story) InsertAfter(id string, item Item) error {
	if i, ok := s.IndexOf(id); ok {
		s.insertAt(i+1, item)
		return nil
	}
	return fmt.Errorf("invalid item id '%v'", after)
}

// SetByID replaces item with `id` with `item`
func (s Story) SetByID(id string, item Item) error {
	if i, ok := s.IndexOf(id); ok {
		s[i] = item
		return nil
	}
	return fmt.Errorf("invalid item id '%v'", id)
}

// Move moves the item with `id` after the item with `afterId`
func (ps *Story) Move(id string, afterId string) error {
	item, err := ps.RemoveByID(id)
	if err != nil {
		return err
	}
	if afterId != "" {
		return ps.InsertAfter(afterId, item)
	}
	ps.Prepend(item)
	return nil
}

// Removes item with `id`
func (s *Story) RemoveByID(id string) (item Item, err error) {
	if i, ok := s.IndexOf(id); ok {
		t := *s
		item = t[i]
		copy(t[i:], t[i+1:])
		t = t[:len(t)-1]
		*s = t
		return item, nil
	}
	return item, fmt.Errorf("missing item id '%v'", id)
}

// Item represents a federated wiki Story item
type Item map[string]interface{}

// Val returns a string value from key
func (item Item) Val(key string) string {
	if s, ok := item[key].(string); ok {
		return s
	}
	return ""
}

// ID returns the `item` identificator
func (item Item) ID() string { return item.Val("id") }

// Date represents a federated wiki time
// it's represented by unix time in JSON
type Date struct{ time.Time }

// NewDate returns a federated wiki Date
func NewDate(t time.Time) Date { return Date{t} }

// ParseDate parses federated wiki time format
func ParseDate(data string) (Date, error) {
	v, err := strconv.Atoi(data)
	if err != nil {
		return Date{}, err
	}
	return Date{time.Unix(int64(v), 0)}, nil
}

// MarshalJSON marshals Date as unix timestamp
func (d Date) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Itoa(int(d.Unix()))), nil
}

// UnmarshalJSON unmarshals Date from an unix timestamp
func (d *Date) UnmarshalJSON(data []byte) (err error) {
	*d, err = ParseDate(string(data))
	return
}
