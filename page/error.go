package page

import "errors"

var (
	ErrInvalid    = errors.New("invalid argument")
	ErrPermission = errors.New("permission denied")
	ErrNotExist   = errors.New("page does not exist")
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
