package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeleteEmpty(t *testing.T) {
	type args struct {
		s []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "empty array",
			args: args{
				s: []string{},
			},
			want: []string{},
		},
		{
			name: "not need to delete array",
			args: args{
				s: []string{"1", "2", "3"},
			},
			want: []string{"1", "2", "3"},
		},
		{
			name: "need to delete array",
			args: args{
				s: []string{"", "1", "", "2", "3", "", ""},
			},
			want: []string{"1", "2", "3"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, DeleteEmpty(tt.args.s), tt.want)
		})
	}
}
