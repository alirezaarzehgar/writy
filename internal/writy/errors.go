package writy

import (
	"fmt"
)

type notfoundError struct {
	Key string
}

func (e notfoundError) Error() string {
	return fmt.Sprint("not found: ", e.Key)
}

func ErrIsNotFound(err error) bool {
	switch err.(type) {
	case notfoundError:
		return true
	case *notfoundError:
		return true
	default:
		return false
	}
}

func IsDuplicated(err error) bool {
	return !ErrIsNotFound(err)
}
