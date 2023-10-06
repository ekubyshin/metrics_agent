package handlers

import (
	"net/http"
	"testing"
)

func TestGaugeHandler_ServeHTTP(t *testing.T) {
	type fields struct {
		route string
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &GaugeHandler{
				route: tt.fields.route,
			}
			m.ServeHTTP(tt.args.w, tt.args.r)
		})
	}
}
