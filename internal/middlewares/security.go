package middlewares

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"io"
	"net/http"
	"strconv"

	"github.com/ekubyshin/metrics_agent/internal/crypto"
)

type cryptoWriter struct {
	http.ResponseWriter
	Writer io.Writer
	secret string
}

const hashHeader = "Hashsha256"

func (w cryptoWriter) Write(b []byte) (int, error) {
	h, err := crypto.HashData(b, w.secret)
	if err == nil {
		w.ResponseWriter.Header().Add(hashHeader, string(h))
	}
	return w.ResponseWriter.Write(b)
}

func (w cryptoWriter) WriteHeader(c int) {
	w.ResponseWriter.Header().Add("Status-Code", strconv.Itoa(c))
}

func NewSecurity(k string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sum := r.Header.Get(hashHeader)
			sw := cryptoWriter{ResponseWriter: w, secret: k}
			if sum == "" {
				next.ServeHTTP(sw, r)
				return
			}
			var buf bytes.Buffer
			_, err := buf.ReadFrom(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			h := hmac.New(sha256.New, []byte(k))
			_, err = h.Write(buf.Bytes())
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			if !hmac.Equal(h.Sum(nil), []byte(sum)) {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			next.ServeHTTP(sw, r)
		})
	}
}
