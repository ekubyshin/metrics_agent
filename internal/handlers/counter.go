package handlers

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type CounterHandler struct {
	route string
}

func NewCounterHandler() Handler {
	return &CounterHandler{
		route: "/counter/{name}/{value:^[0-9]+}",
	}
}

func (m *CounterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	paramName := chi.URLParam(r, "name")
	paramValue := chi.URLParam(r, "value")
	if paramName == "" || paramValue == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_, err := strconv.ParseInt(paramValue, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(paramValue))
	if err != nil {
		panic(err)
	}
}

func (m *CounterHandler) BaseURL() string {
	return m.route
}
