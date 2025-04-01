package writy

import "log"

type Writy struct {
	logger log.Logger
}

func (w *Writy) Set(key, value string) error {
	// Write to cache

	return nil
}

func (w *Writy) Get(key string) (string, error) {
	return "", nil
}
