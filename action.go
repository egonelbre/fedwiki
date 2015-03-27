package fedwiki

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

func (action Action) Item() (Item, bool) {
	item, ok := action["item"]
	if !ok {
		return nil, false
	}
	m, ismap := (item).(map[string]interface{})
	if !ismap {
		return nil, false
	}
	return (Item)(m), true
}

func (action Action) Date() (t Date, err error) {
	val, ok := action["date"]
	if !ok {
		return Date{time.Unix(0, 0)}, fmt.Errorf("date not found")
	}

	switch val := val.(type) {
	case string:
		t, err := time.Parse(time.RFC3339, val)
		return Date{t}, err
	case int: // assume date
		return Date{time.Unix(int64(val), 0)}, nil
	case int64: // assume date
		return Date{time.Unix(val, 0)}, nil
	}

	return Date{time.Unix(0, 0)}, fmt.Errorf("unknown date format")
}

var actionfns = map[string]func(p *Page, a Action) error{
	"add": func(p *Page, action Action) error {
		item, ok := action.Item()
		if !ok {
			return fmt.Errorf("no item in action")
		}

		after := action.Str("after")
		if after == "" {
			p.Story.Append(item)
			return nil
		}
		return p.Story.InsertAfter(after, item)
	},
	"edit": func(p *Page, action Action) error {
		item, ok := action.Item()
		if !ok {
			return fmt.Errorf("no item in action")
		}
		return p.Story.SetByID(action.Str("id"), item)
	},
	"remove": func(p *Page, action Action) error {
		_, err := p.Story.RemoveByID(action.Str("id"))
		return err
	},
	"move": func(p *Page, action Action) error {
		return p.Story.Move(action.Str("id"), action.Str("after"))
	},
	"create": func(p *Page, action Action) error {
		return nil
	},
	"fork": func(p *Page, action Action) error {
		return nil
	},
}
