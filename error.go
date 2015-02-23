package fedwiki

import "errors"

var (
	ErrInvalid       = errors.New("invalid argument")
	ErrPermission    = errors.New("permission denied")
	ErrNotExist      = errors.New("page does not exist")
	ErrUnknownAction = errors.New("unknown action")
)
