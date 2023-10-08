package utils

import (
	"testing"
	"time"
)

func TestIntToDuration(t *testing.T) {
	type args struct {
		val int
	}
	tests := []struct {
		name string
		args args
		want time.Duration
	}{
		{
			"ok",
			args{
				10,
			},
			time.Duration(10) * time.Second,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IntToDuration(tt.args.val); got != tt.want {
				t.Errorf("IntToDuration() = %v, want %v", got, tt.want)
			}
		})
	}
}
