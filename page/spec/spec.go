package spec

import (
	"fmt"
	"reflect"

	"github.com/egonelbre/wiki-go-server/page"
)

// Action can validate a page.Action
type Action map[string]Prop

// Item can validate a page.Item
type Item map[string]Prop

// Prop is a value definition
type Prop struct {
	Kind     reflect.Kind
	Optional bool
}

func (actspec Action) Validate(action page.Action) error {
	spec := (map[string]Prop)(actspec)
	vals := (map[string]interface{})(action)
	return validate(spec, vals)
}

func (itemspec Item) Validate(item page.Item) error {
	spec := (map[string]Prop)(itemspec)
	vals := (map[string]interface{})(item)
	return validate(spec, vals)
}

func validate(spec map[string]Prop, values map[string]interface{}) error {
	var errs errlist
	for key, prop := range spec {
		v, ok := values[key]
		if prop.Optional && !ok {
			continue
		}
		if !ok {
			errs = append(errs, fmt.Errorf(`did not find property "%v"`, key))
			continue
		}

		actual := reflect.TypeOf(v).Kind()
		if prop.Kind != actual {
			errs = append(errs, fmt.Errorf(`expected "%v" to be "%v", but was "%v"`, key, expected, actual))
		}
	}
}
