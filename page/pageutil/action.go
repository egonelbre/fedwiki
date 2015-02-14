package pageutil

import (
	"encoding/json"
	"io"

	"github.com/egonelbre/wiki-go-server/page"
	"github.com/egonelbre/wiki-go-server/page/spec"
)

func ReadAction(r io.Reader) (page.Action, error) {
	dec := json.NewDecoder(r)
	action := make(page.Action)
	err := dec.Decode(&action)
	if err != nil {
		return nil, err
	}

	if validator, ok := spec.Actions[action.Type()]; ok {
		if err := validator.Validate(action); err != nil {
			return nil, err
		}
	}

	return action, nil
}
