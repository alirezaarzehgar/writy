package writy

import (
	"fmt"
	"log/slog"
	"os"
	"testing"
	"time"
)

var w *writy

func TestFlush(t *testing.T) {
	var err error
	w, err = New(".", time.Second)
	if err != nil {
		t.Fatal("unable to open storage")
	}
	w.WithLogHandler(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug}))
	defer w.Close()

	w.Set("name", "ali")
	t.Log("ali saved in name")

	v, _ := w.Get("name")
	if v != "ali" {
		t.Error("saved key not found")
	}
	t.Log("key name found:", v)

	for i := 0; i < 10; i++ {
		time.Sleep(time.Second / 3)
		w.Set(fmt.Sprint("key-", i), "vvv")
	}
}

func TestSearchIndexByKey(t *testing.T) {
	off := searchIndexByKey(w, "notfound")
	if off >= 0 {
		t.Error("false positive on search")
	}

	off = searchIndexByKey(w, "key-5")
	if off != 119 {
		t.Error("search by index is not correct", off)
	}
	t.Log("offset:", off)
}

func TestGetValueByOffset(t *testing.T) {
	v, err := w.cache.Get("key-2")
	if err == nil {
		t.Error("key-2 already exists in cache:", v)
		t.Log(w.cache.List())
	}
	t.Log("key-2 is not in the cache:", err)

	v, err = w.Get("key-2")
	if err != nil {
		t.Error("key-2 is not found:", err)
	}
	t.Log("key-2 value:", v, ", error:", err)
}
