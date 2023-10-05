package collector

import (
	"reflect"
	"testing"
	"time"
)

func TestNewRuntimeReader(t *testing.T) {
	type args struct {
		pollInterval time.Duration
	}
	tests := []struct {
		name string
		args args
		want Reader
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewRuntimeReader(tt.args.pollInterval); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRuntimeReader() = %v, want %v", got, tt.want)
			}
		})
	}
}
