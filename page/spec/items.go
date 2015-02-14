package spec

import "reflect"

var (
	Items = map[string]Item{
		"paragraph": Item{
			"id":   Prop{reflect.String, false},
			"text": Prop{reflect.String, false},
		},
		"html": Item{
			"id":   Prop{reflect.String, false},
			"text": Prop{reflect.String, false},
		},
	}
)
