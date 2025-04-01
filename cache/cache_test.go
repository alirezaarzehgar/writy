package cache_test

import (
	"testing"

	"github.com/alirezaarzehgar/writy/cache"
)

func TestCache(t *testing.T) {
	c := cache.New()

	err := c.Set("name", "ali")
	if err != nil {
		t.Error("error in first write to cache:", err)
	}

	v, err := c.Get("name")
	if err != nil {
		t.Error("error in first read from cache:", err)
	}
	if v != "ali" {
		t.Error("key name is equal to ali not ", v)
	}

	err = c.Del("name")
	if cache.IsNotFound(err) {
		t.Error("false positive not found error:", err)
	}
	err = c.Del("name")
	if !cache.IsNotFound(err) {
		t.Error("key does not exits")
	}
}
