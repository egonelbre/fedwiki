package page

import (
	"fmt"
	"time"
)

type Journal []Action

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
		return time.Unix(val, 0), nil
	}

	return time.Unix(0, 0), fmt.Errorf("unknown date format")
}
