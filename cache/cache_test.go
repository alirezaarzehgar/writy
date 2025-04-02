package cache_test

import (
	"log/slog"
	"os"
	"testing"

	"github.com/alirezaarzehgar/writy/cache"
)

func TestCache(t *testing.T) {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug})))

	c := cache.New()
	err := c.Set("name", "ali")
	if err != nil {
		t.Error("error in first write to cache:", err)
	}

	v := c.Get("name")
	if v == nil {
		t.Error("error in first read from cache:", err)
	}
	if v != "ali" {
		t.Error("key name is equal to ali not ", v)
	}

	err = c.Del("name")
	if err == nil {
		t.Error("false positive not found error:", err)
	}
	err = c.Del("name")
	if err == nil {
		t.Error("key does not exits")
	}
}
