package middlewares

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

func GzipReader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}
		oldBody := r.Body
		defer oldBody.Close()
		zr, err := gzip.NewReader(oldBody)
		if err != nil {
			io.WriteString(w, err.Error()) //nolint
			return
		}
		r.Body = zr
		next.ServeHTTP(w, r)
	})
}
