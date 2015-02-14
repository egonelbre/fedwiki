package spec

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
