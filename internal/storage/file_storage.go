package storage

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"os"
	"time"

	"github.com/ekubyshin/metrics_agent/internal/metrics"
)

type FileStorage[K any, V metrics.Keyable[K]] struct {
	st            *MemStorage[K, V]
	Filename      string
	file          *os.File
	storeInterval time.Duration
}

func (s FileStorage[K, V]) Put(ctx context.Context, k K, v V) error {
	err := s.st.Put(ctx, k, v)
	if err != nil {
		return err
	}
	if s.storeInterval == 0 {
		go func() {
			_ = s.flush(ctx)
		}()
	}
	return nil
}

func (s FileStorage[K, V]) Get(ctx context.Context, k K) (V, bool) {
	return s.st.Get(ctx, k)
}

func (s FileStorage[K, V]) Delete(ctx context.Context, k K) error {
	return s.st.Delete(ctx, k)
}

func (s FileStorage[K, V]) List(ctx context.Context) ([]KeyValuer[K, V], error) {
	return s.st.List(ctx)
}

func NewFileStorage[K any, V metrics.Keyable[K]](
	ctx context.Context,
	st *MemStorage[K, V],
	filename string,
	restore bool,
	interval time.Duration) (*FileStorage[K, V], error) {
	if st == nil {
		return nil, errors.New("storage is nil")
	}
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	fs := &FileStorage[K, V]{
		st:            st,
		Filename:      filename,
		file:          file,
		storeInterval: interval,
	}
	if restore {
		err = fs.restore(ctx)
	}
	if interval > 0 {
		fs.runInterval(ctx)
	}
	return fs, err
}

func (w *FileStorage[K, V]) Ping(ctx context.Context) error {
	return nil
}

func (w *FileStorage[K, V]) runInterval(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				_ = w.Close()
				return
			default:
				time.Sleep(w.storeInterval)
				_ = w.flush(ctx)
			}
		}
	}()
}

func (w *FileStorage[K, V]) Close() error {
	return w.file.Close()
}

func (w *FileStorage[K, V]) restore(ctx context.Context) error {
	reader := bufio.NewReader(w.file)
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanLines)
	for {
		if !scanner.Scan() {
			break
		}
		m := new(V)
		err := json.Unmarshal(scanner.Bytes(), &m)
		if err != nil {
			continue
		}
		key := V.Key(*m)
		err = w.st.Put(ctx, key, *m)
		if err != nil {
			return err
		}
	}
	return nil
}

func (w *FileStorage[K, V]) flush(ctx context.Context) error {
	elems, _ := w.st.List(ctx)
	if len(elems) == 0 {
		return nil
	}
	_ = w.file.Truncate(0)
	_, _ = w.file.Seek(0, 0)
	writer := bufio.NewWriter(w.file)
	for _, v := range elems {
		str, err := json.Marshal(v.Value)
		if err != nil {
			return err
		}
		_, err = writer.Write(str)
		if err != nil {
			return err
		}
		_, err = writer.Write([]byte("\n"))
		if err != nil {
			return err
		}
	}
	return writer.Flush()
}
