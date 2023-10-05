package reporter

import (
	"net/http"
	"time"
)

func NewReporterClient() *http.Client {
	return &http.Client{
		Timeout: time.Second * 1,
	}
}
