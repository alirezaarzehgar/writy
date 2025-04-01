package writy

import (
	"fmt"
	"log/slog"
	"os"
	"testing"
	"time"
)

func TestFlush(t *testing.T) {
	w, err := New(".", time.Second)
	if err != nil {
		t.Fatal("unable to open storage")
	}
	w.WithLogHandler(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug}))

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
	w, err := New(".", time.Second)
	if err != nil {
		t.Fatal("unable to open storage")
	}
	w.WithLogHandler(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug}))

	writyW := w.(*writy)
	_, err = searchIndexByKey(writyW, "notfound")
	if !IsNotFound(err) {
		t.Error("false positive on search")
	}

	off, err := searchIndexByKey(writyW, "key-5")
	if err != nil || off != 119 {
		t.Error("search by index is not correct", off, err)
	}
}
