package spec

import "reflect"

var (
	Actions = map[string]Action{
		"add": Action{
			"id":    Prop{reflect.String, false},
			"item":  Prop{reflect.Interface, false},
			"after": Prop{reflect.String, true},
		},
		"edit": Action{
			"id":   Prop{reflect.String, false},
			"item": Prop{reflect.Interface, false},
		},
		"move": Action{
			"id":    Prop{reflect.String, false},
			"after": Prop{reflect.String, true},
			"order": Prop{reflect.Slice, false},
		},
		"remove": Action{
			"id": Prop{reflect.String, false},
		},
	}
)
