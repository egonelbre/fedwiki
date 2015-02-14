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
		"reference": Item{
			"id":    Prop{reflect.String, false},
			"site":  Prop{reflect.String, false},
			"title": Prop{reflect.String, true},
			"text":  Prop{reflect.String, true},
		},
	}
)
