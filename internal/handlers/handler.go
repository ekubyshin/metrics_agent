package handlers

import "net/http"

type Handler interface {
	http.Handler
	BaseURL() string
}
