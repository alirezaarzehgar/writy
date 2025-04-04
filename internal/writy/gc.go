package writy

import (
	"log/slog"
	"os"
	"time"
)

var (
	DefaultGarbageCollectorCycle = time.Minute * 2
)

type GarbageCollector struct {
	ir, itw *os.File
	sr, stw *os.File
	writy   *Writy
}

func newGarbageCollector() *GarbageCollector {
	return &GarbageCollector{}
}

func (gc *GarbageCollector) Run(w *Writy, flusher *Flusher) {
	gc.writy = w
	gc.ir = w.indexReader
	gc.sr = w.storageReader

	indTmpWriter, err := os.OpenFile(gc.ir.Name()+".tmp", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		slog.Warn("failed to open temporary file for index")
	}
	gc.itw = indTmpWriter

	storTmpWriter, err := os.OpenFile(gc.sr.Name()+".tmp", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		slog.Warn("failed to open temporary file for storage")
	}
	gc.stw = storTmpWriter

	go func() {
		slog.Debug("run garbage collector")
		for {
			select {
			case <-time.NewTicker(DefaultGarbageCollectorCycle).C:
				if len(w.cache.List()) > 0 {
					flusher.flush()
				}
				gc.collect()
				gc.writy.w8ForDaemons.Wait()
			}
		}
	}()
}

func (gc *GarbageCollector) collect() {
	gc.writy.w8ForDaemons.Add(1)
	defer gc.writy.w8ForDaemons.Done()
	glk.Lock()
	defer glk.Unlock()

	slog.Debug("start manageing deleted rows")

	gc.itw.Truncate(0)
	gc.stw.Truncate(0)

	ienc := newIndexEncoder(gc.itw)
	idec := newIndexDecoder(gc.ir)
	senc := newStorageEncoder(gc.stw)
	sdec := newStorageDecoder(gc.sr)

	for idec.Scan() {
		i := idec.Decode()
		if !i.IsDeleted {
			ienc.Encode(i.Key, i.ValueOffset)
			v, _ := sdec.Decode(i.ValueOffset)
			senc.Encode(i.Key, v)
		}
	}

	os.Rename(gc.itw.Name(), gc.ir.Name())
	os.Rename(gc.stw.Name(), gc.sr.Name())
	slog.Debug("managed garbage collector")
}
