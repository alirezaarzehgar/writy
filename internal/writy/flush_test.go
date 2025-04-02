package writy

import (
	"fmt"
	"log/slog"
	"os"
	"testing"
	"time"
)

var w *Writy

func TestFlush(t *testing.T) {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug})))

	var err error
	w, err = New(".", time.Second)
	if err != nil {
		t.Fatal("unable to open storage")
	}
	defer w.Flush()

	w.Set("name", "ali")
	t.Log("ali saved in name")

	v, _ := w.Get("name")
	if v != "ali" {
		t.Error("saved key not found")
	}
	t.Log("key name found:", v)

	for i := 0; i < 10; i++ {
		time.Sleep(time.Second / 3)
		w.Set(fmt.Sprint("key-", i), fmt.Sprint("vvv\nasd", i))
	}
}

func TestSearchIndexByKey(t *testing.T) {
	off := searchIndexByKey(w, "notfound")
	if off >= 0 {
		t.Error("false positive on search")
	}

	off = searchIndexByKey(w, "key-5")
	if off < 0 {
		t.Error("search by index is not correct", off)
	}
	t.Log("offset:", off)
}

func TestGetValueByOffset(t *testing.T) {
	for i := 0; i < 10; i++ {
		key := fmt.Sprint("key-", i)
		v := w.cache.Get(key)
		if v != nil {
			t.Error(key, "already exists in cache:", v)
			t.Log(w.cache.List())
		}

		v, err := w.Get(key)
		if err != nil {
			t.Error(key, " is not found:", err)
		}
		t.Log(key, " value:", v, ", error:", err)
	}
}

func BenchmarkWrityGet(b *testing.B) {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{})))

	for i := 0; i < b.N; i++ {
		_, err := w.Get("key-5")
		if err != nil {
			b.Fatal("benchmark failed. error:", err)
		}
	}
}
