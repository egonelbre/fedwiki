package page

import (
	"fmt"
	"time"
)

type Action map[string]interface{}

func (action Action) Val(key string) string {
	if s, ok := action[key].(string); ok {
		return s
	}
	return ""
}

func (action Action) Type() string { return action.Val("type") }

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
	"edit": func(p *Page, a Action) error {
		return nil
	},
	"add": func(p *Page, a Action) error {
		return nil
	},
	"remove": func(p *Page, a Action) error {
		return nil
	},
	"move": func(p *Page, a Action) error {
		return nil
	},
}
