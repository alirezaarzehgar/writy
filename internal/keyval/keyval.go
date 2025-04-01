package keyval

import "log/slog"

type StorageType map[string]any

type KeyVal interface {
	WithLogHandler(handler slog.Handler) KeyVal
	Set(key, value string) error
	ForceSet(key string, value any) error
	Get(key string) (any, error)
	Del(key string) error
	Clear() error
	List() (StorageType, error)
}
