package writy

import "log/slog"

var logger *slog.Logger = slog.Default()

func SetLogger(l *slog.Logger) {
	logger = l
}
