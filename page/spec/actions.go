package spec

import "reflect"

var (
	Actions = map[string]Action{
		"add": spec{
			"id":    Prop{reflect.String, false},
			"item":  Prop{reflect.Interface, false},
			"after": Prop{reflect.String, true},
		},
		"edit": spec{
			"id":   Prop{reflect.String, false},
			"item": Prop{reflect.Interface, false},
		},
		"move": spec{
			"id":    Prop{reflect.String, false},
			"after": Prop{reflect.String, true},
			"order": Prop{reflect.Slice, false},
		},
		"remove": spec{
			"id": Prop{reflect.String, false},
		},
	}
)
