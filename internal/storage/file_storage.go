package storage

import (
	"bufio"
	"encoding/json"
	"errors"
	"os"
	"time"

	"github.com/ekubyshin/metrics_agent/internal/types"
)

type FileStorage[K any, V types.Keyable[K]] struct {
	st       *MemStorage[K, V]
	Filename string
	file     *os.File
	kType    K
	vType    V
}

func (s FileStorage[K, V]) Put(k K, v V) {
	s.st.Put(k, v)
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

func NewFileStorage[K any, V types.Keyable[K]](
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
		st:       st,
		Filename: filename,
		file:     file,
	}
	if restore {
		err = fs.restore()
	}
	return fs, err
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

// func (w *FileStorage) Flush() error {
// 	writer := bufio.NewWriter(w.file)
// 	elems := w.st.List()
// 	if len(elems) == 0 {
// 		return nil
// 	}
// 	for _, v := range elems {
// 		str, err := json.Marshal(v.Value)
// 		if err != nil {
// 			return err
// 		}
// 		_, err = writer.WriteString(string(str))
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	return writer.Flush()
// }
