package writy

// TODO: Redesign this and find efficient way for writing may lines once.
// This implementation is not performant.
func writeIndex(w *Writy, k string, v any) {
	off := newStorageEncoder(w.storageWriter).Encode(k, v)
	err := newIndexEncoder(w.indexWriter).Encode(k, off)
	if err != nil {
		logger.Warn("unable to encode data", "error", err)
	}
}

func searchIndexByKey(w *Writy, k string) int64 {
	indDec := newIndexDecoder(w.indexReader)
	for indDec.Scan() {
		ind := indDec.Decode()
		if !ind.IsDeleted && k == ind.Key {
			logger.Debug("found offset", "fkey", ind.Key, "k", k, "isdel", ind.IsDeleted, "offset", ind.Offset)
			return ind.Offset
		}
	}

	return -1
}

func getValueByOffset(w *Writy, off int64) any {
	s := newStorageDecoder(w.storageReader)
	line, err := s.Decode(off)
	if err != nil {
		logger.Warn("failed to read storage", "offset", off, "error", err)
	}

	logger.Debug("desirable line found", "line", line, "error", err)
	return line
}
