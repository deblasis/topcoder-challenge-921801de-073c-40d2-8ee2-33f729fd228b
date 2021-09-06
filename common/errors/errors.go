package errors

import "errors"

var (
	ErrUnknown = errors.New("unknown argument passed")

	ErrInvalidArgument = errors.New("invalid argument passed")
)

func Str2err(s string) error {
	if s == "" {
		return nil
	}
	return errors.New(s)
}

func Err2str(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}
