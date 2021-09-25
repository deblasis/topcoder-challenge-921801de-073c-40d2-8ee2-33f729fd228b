package common_v1

func (e *Error) Error() string {
	if e == NilErr {
		return ""
	}
	return e.Message
}

var NilErr = (*Error)(nil)
