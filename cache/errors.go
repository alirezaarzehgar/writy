package cache

import "fmt"

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
	_, ok := err.(*notfoundError)
	return ok
}

func IsDuplicated(err error) bool {
	_, ok := err.(*duplicatedError)
	return ok
}
