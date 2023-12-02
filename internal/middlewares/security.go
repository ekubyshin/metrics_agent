package middlewares

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"io"
	"net/http"

	"github.com/ekubyshin/metrics_agent/internal/crypto"
)

type cryptoWriter struct {
	http.ResponseWriter
	Writer  io.Writer
	secret  string
	status  int
	UseHash bool
}

func (w *cryptoWriter) Write(b []byte) (int, error) {
	if !w.UseHash {
		return w.ResponseWriter.Write(b)
	}
	h, err := crypto.HashData(b, w.secret)
	if err == nil {
		w.ResponseWriter.Header().Set(crypto.HashHeader, string(h))
	}
	if w.status != 0 {
		w.ResponseWriter.WriteHeader(w.status)
	}
	return w.ResponseWriter.Write(b)
}

func (w *cryptoWriter) WriteHeader(c int) {
	w.status = c
}

func NewSecurity(k string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sum := r.Header.Get(crypto.HashHeader)
			sw := &cryptoWriter{ResponseWriter: w, secret: k}
			if sum == "" {
				next.ServeHTTP(sw, r)
				return
			}
			var buf bytes.Buffer
			_, err := buf.ReadFrom(r.Body)
			if err != nil {
				sw.WriteHeader(http.StatusInternalServerError)
				next.ServeHTTP(sw, r)
				return
			}
			h := hmac.New(sha256.New, []byte(k))
			_, err = h.Write(buf.Bytes())
			if err != nil {
				sw.WriteHeader(http.StatusInternalServerError)
				next.ServeHTTP(sw, r)
				return
			}
			if !hmac.Equal(h.Sum(nil), []byte(sum)) {
				sw.WriteHeader(http.StatusBadRequest)
				next.ServeHTTP(sw, r)
				return
			}
			sw.UseHash = true
			next.ServeHTTP(sw, r)
		})
	}
}
