package apk

type _errorSlice []error

func errorSlice(s []error) error {
	if len(s) == 0 {
		return nil
	}
	return _errorSlice(s)
}
func (e _errorSlice) Error() string {
	s := ""
	for _, v := range e {
		s += v.Error() + ";"
	}
	return s
}
