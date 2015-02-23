package fedwiki

import (
	"fmt"
	"reflect"
)

// ActionSpec can validate a Action
type ActionSpec map[string]PropSpec

// ItemSpec can validate a Item
type ItemSpec map[string]PropSpec

// PropSpec is a value definition
type PropSpec struct {
	Kind     reflect.Kind
	Optional bool
}

func (actspec ActionSpec) Validate(action Action) error {
	spec := (map[string]PropSpec)(actspec)
	vals := (map[string]interface{})(action)
	return validate(spec, vals)
}

func (itemspec ItemSpec) Validate(item Item) error {
	spec := (map[string]PropSpec)(itemspec)
	vals := (map[string]interface{})(item)
	return validate(spec, vals)
}

func validate(spec map[string]PropSpec, values map[string]interface{}) error {
	var errs errlist
	for key, PropSpec := range spec {
		v, ok := values[key]
		if PropSpec.Optional && !ok {
			continue
		}
		if !ok {
			errs = append(errs, fmt.Errorf(`did not find property "%v"`, key))
			continue
		}

		actual := reflect.TypeOf(v).Kind()
		if PropSpec.Kind != actual {
			errs = append(errs, fmt.Errorf(`expected "%v" to be "%v", but was "%v"`, key, PropSpec.Kind, actual))
		}
	}

	return errs
}

var (
	ActionSpecs = map[string]ActionSpec{
		"add": ActionSpec{
			"id":    PropSpec{reflect.String, false},
			"item":  PropSpec{reflect.Interface, false},
			"after": PropSpec{reflect.String, true},
		},
		"edit": ActionSpec{
			"id":   PropSpec{reflect.String, false},
			"item": PropSpec{reflect.Interface, false},
		},
		"move": ActionSpec{
			"id":    PropSpec{reflect.String, false},
			"after": PropSpec{reflect.String, true},
			"order": PropSpec{reflect.Slice, false},
		},
		"remove": ActionSpec{
			"id": PropSpec{reflect.String, false},
		},
	}
)

var (
	ItemSpecs = map[string]ItemSpec{
		"paragraph": ItemSpec{
			"id":   PropSpec{reflect.String, false},
			"text": PropSpec{reflect.String, false},
		},
		"html": ItemSpec{
			"id":   PropSpec{reflect.String, false},
			"text": PropSpec{reflect.String, false},
		},
		"reference": ItemSpec{
			"id":    PropSpec{reflect.String, false},
			"site":  PropSpec{reflect.String, false},
			"title": PropSpec{reflect.String, true},
			"text":  PropSpec{reflect.String, true},
		},
	}
)

type errlist []error

func (errs errlist) Error() string {
	s := ""
	for i, err := range errs {
		if i == 0 {
			s = err.Error()
			continue
		}
		s += "; " + err.Error()
	}
	return s
}
