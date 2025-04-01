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

func IsNotFound(err error) bool {
	switch err.(type) {
	case notfoundError:
		return true
	case *notfoundError:
		return true
	default:
		return false
	}
}
