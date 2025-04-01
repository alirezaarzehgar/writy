package cache

import (
	"fmt"
)

type notfoundError struct {
	Key string
}

func (e notfoundError) Error() string {
	return fmt.Sprint("not found: ", e.Key)
}

type duplicatedError struct {
	Key string
}

func (e duplicatedError) Error() string {
	return fmt.Sprint("duplicate key: ", e.Key)
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

func IsDuplicated(err error) bool {
	switch err.(type) {
	case duplicatedError:
		return true
	case *duplicatedError:
		return true
	default:
		return false
	}
}
