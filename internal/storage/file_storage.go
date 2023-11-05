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

func (s FileStorage[K, V]) Put(k K, v V) {
	s.st.Put(k, v)
	if s.storeInterval == 0 {
		go func() {
			_ = s.flush()
		}()
	}
}

func (s FileStorage[K, V]) Get(k K) (V, bool) {
	return s.st.Get(k)
}

func (s FileStorage[K, V]) Delete(k K) {
	s.st.Delete(k)
}

func (s FileStorage[K, V]) List() []KeyValuer[K, V] {
	return s.st.List()
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
		err = fs.restore()
	}
	if interval > 0 {
		go func() {
			for {
				select {
				case <-ctx.Done():
					_ = fs.Close()
					return
				default:
					time.Sleep(interval)
					_ = fs.flush()
				}
			}
		}()
	}
	return fs, err
}

func (w *FileStorage[K, V]) Close() error {
	return w.file.Close()
}

func (w *FileStorage[K, V]) restore() error {
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
		w.st.Put(key, *m)
	}
	return nil
}

func (w *FileStorage[K, V]) flush() error {
	elems := w.st.List()
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
