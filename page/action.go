package page

import (
	"fmt"
	"time"
)

type Action map[string]interface{}

func (action Action) Str(key string) string {
	if s, ok := action[key].(string); ok {
		return s
	}
	return ""
}

func (action Action) Type() string { return action.Str("type") }

func (action Action) Date() (t time.Time, err error) {
	val, ok := action["date"]
	if !ok {
		return time.Unix(0, 0), fmt.Errorf("date not found")
	}

	switch val := val.(type) {
	case string:
		return time.Parse(time.RFC3339, val)
	case int: // assume date
		return time.Unix(int64(val), 0), nil
	case int64: // assume date
		return time.Unix(val, 0), nil
	}

	return time.Unix(0, 0), fmt.Errorf("unknown date format")
}

var actionfns = map[string]func(p *Page, a Action) error{
	"add": func(p *Page, action Action) error {
		props := action["item"]
		item, ok := props.(Item)
		if !ok {
			return fmt.Errorf("invalid item")
		}
		return p.Story.InsertAfter(action.Str("after"), item)
	},
	"edit": func(p *Page, action Action) error {
		props := action["item"]
		item, ok := props.(Item)
		if !ok {
			return fmt.Errorf("invalid item")
		}
		return p.Story.SetById(action.Str("id"), item)
	},
	"remove": func(p *Page, action Action) error {
		_, err := p.Story.RemoveById(action.Str("id"))
		return err
	},
	"move": func(p *Page, action Action) error {
		return p.Story.Move(action.Str("id"), action.Str("after"))
	},
}
