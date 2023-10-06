package handlers

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type GaugeHandler struct {
	route string
}

func NewGaugeHandler() Handler {
	return &GaugeHandler{
		route: "/gauge/{name}/{value:^[0-9]\\.[0-9]$}",
	}
}

func (m *GaugeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_ = chi.URLParam(r, "name")
	paramValue := chi.URLParam(r, "value")
	_, err := strconv.ParseFloat(paramValue, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("{}"))
	if err != nil {
		panic(err)
	}
}

func (m *GaugeHandler) BaseURL() string {
	return m.route
}
