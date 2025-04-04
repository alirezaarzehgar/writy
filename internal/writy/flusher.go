package writy

import (
	"log/slog"
	"time"
)

type Flusher struct {
	running bool
	cycle   time.Duration
	writy   *Writy
}

func newFlusher(cycle time.Duration) *Flusher {
	return &Flusher{cycle: cycle, running: true}
}

func (f *Flusher) SetCycle(c time.Duration) {
	f.cycle = c
}

func (f *Flusher) Run(w *Writy) {
	f.writy = w

	go func() {
		slog.Debug("start filesystem flusher")
		for f.running {
			select {
			case <-time.NewTicker(f.cycle).C:
				f.flush()
			}
		}
	}()
}

func (f *Flusher) flush() {
	glk.Lock()
	f.writy.w8ForDaemons.Add(1)

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
