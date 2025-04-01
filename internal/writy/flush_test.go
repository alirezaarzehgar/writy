package writy_test

import (
	"fmt"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/alirezaarzehgar/writy/internal/writy"
)

func TestFlush(t *testing.T) {
	t.Fail()
	w, err := writy.New("./storage.db", time.Second)
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
		w.Set(fmt.Sprint("key-", i), "value")
	}
}
