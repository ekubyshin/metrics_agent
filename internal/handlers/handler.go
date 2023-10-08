package handlers

import "net/http"

type Handler interface {
	http.Handler
	BaseURL() string
}

type Key struct {
	Type string
	Name string
}

const (
	GaugeActionKey   = "gauge"
	CounterActionKey = "counter"
	ParamTypeKey     = "type"
	ParamNameKey     = "name"
	ParamValueKey    = "value"
)
