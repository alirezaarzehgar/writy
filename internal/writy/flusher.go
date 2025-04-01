package writy

import (
	"log/slog"
	"time"
)

type Flusher struct {
	logger  *slog.Logger
	running bool
	cycle   time.Duration
	writy   *writy
}

func NewFlusher(cycle time.Duration) *Flusher {
	return &Flusher{cycle: cycle, running: true}
}

func (f *Flusher) SetCycle(c time.Duration) {
	f.cycle = c
}

func (f *Flusher) Run(w *writy) {
	f.writy = w

	go func() {
		f.logger.Debug("start filesystem flusher")
		for f.running {
			select {
			case <-time.NewTicker(f.cycle).C:
				f.flush()
			}
		}
	}()
}

func (f *Flusher) flush() {
	defer f.writy.cache.Clear()

	list, _ := f.writy.cache.List()
	for k, v := range list {
		if searchIndexByKey(f.writy, k) < 0 {
			writeIndex(f.writy, k, v)
		}
	}
}

func (f *Flusher) Stop() {
	f.running = false
}
