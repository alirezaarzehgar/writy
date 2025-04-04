package writy

import (
	"log/slog"
	"time"
)

var (
	DefaultFlushCycle time.Duration = time.Second * 5
)

type Flusher struct {
	running bool
	writy   *Writy
}

func newFlusher() *Flusher {
	return &Flusher{running: true}
}

func (f *Flusher) Run(w *Writy) {
	f.writy = w

	go func() {
		slog.Debug("start filesystem flusher")
		for f.running {
			select {
			case <-time.NewTicker(DefaultFlushCycle).C:
				f.flush()
			}
		}
	}()
}

func (f *Flusher) flush() {
	glk.Lock()
	f.writy.w8ForDaemons.Add(1)
	slog.Debug("flush cache to filesystem")

	for k, v := range f.writy.cache.List() {
		if searchIndexByKey(f.writy, k) < 0 {
			writeIndex(f.writy, k, v)
		}
	}

	f.writy.cache.Clear()
	glk.Unlock()
	f.writy.w8ForDaemons.Done()
}

func (f *Flusher) Stop() {
	f.running = false
}
