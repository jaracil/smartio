package smartio

import (
	"io"
	"sync/atomic"
	"time"
)

const (
	defaultChunkSize = 4096
)

type SmartReader struct {
	r       io.Reader
	total   int64
	limit   int64
	last    int64
	limited uint32 //atomic bool
}

func NewSmartReader(reader io.Reader) *SmartReader {
	return &SmartReader{r: reader, last: time.Now().Unix()}
}

func (sr *SmartReader) SetTotal(total int64) int64 {
	return atomic.SwapInt64(&sr.total, total)
}

func (sr *SmartReader) GetTotal() int64 {
	return atomic.LoadInt64(&sr.total)
}

func (sr *SmartReader) SetLimit(limit int64) int64 {
	last := atomic.SwapInt64(&sr.limit, limit)
	if limit != 0 {
		atomic.StoreUint32(&sr.limited, 1)
	} else {
		atomic.StoreUint32(&sr.limited, 0)
	}
	return last
}

func (sr *SmartReader) GetLimit() int64 {
	return atomic.LoadInt64(&sr.limit)
}

func (sr *SmartReader) SetLast(last int64) int64 {
	return atomic.SwapInt64(&sr.last, last)
}

func (sr *SmartReader) GetLast() int64 {
	return atomic.LoadInt64(&sr.last)
}

func (sr *SmartReader) Read(p []byte) (n int, err error) {
	if atomic.LoadUint32(&sr.limited) != 0 && atomic.LoadInt64(&sr.limit) < 0 {
		err = io.EOF
		return
	}
	atomic.StoreInt64(&sr.last, time.Now().Unix())
	n, err = sr.r.Read(p)
	atomic.StoreInt64(&sr.last, time.Now().Unix())
	atomic.AddInt64(&sr.limit, int64(-n))
	atomic.AddInt64(&sr.total, int64(n))
	return
}

type SmartWriter struct {
	w     io.Writer
	total int64
	last  int64
}

func NewSmartWriter(writer io.Writer) *SmartWriter {
	return &SmartWriter{w: writer, last: time.Now().Unix()}
}

func (sw *SmartWriter) SetTotal(total int64) int64 {
	return atomic.SwapInt64(&sw.total, total)
}

func (sw *SmartWriter) GetTotal() int64 {
	return atomic.LoadInt64(&sw.total)
}

func (sw *SmartWriter) SetLast(last int64) int64 {
	return atomic.SwapInt64(&sw.last, last)
}

func (sw *SmartWriter) GetLast() int64 {
	return atomic.LoadInt64(&sw.last)
}

func (sw *SmartWriter) Write(p []byte) (n int, err error) {
	atomic.StoreInt64(&sw.last, time.Now().Unix())
	n, err = sw.w.Write(p)
	atomic.StoreInt64(&sw.last, time.Now().Unix())
	atomic.AddInt64(&sw.total, int64(n))

	return
}
