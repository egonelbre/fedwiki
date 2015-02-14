package spec

import "reflect"

var (
	Items = map[string]Action{
		"paragraph": spec{
			"id":   Prop{reflect.String, false},
			"text": Prop{reflect.String, false},
		},
		"html": spec{
			"id":   Prop{reflect.String, false},
			"text": Prop{reflect.String, false},
		},
	}
)
