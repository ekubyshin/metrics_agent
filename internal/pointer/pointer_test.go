package pointer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToPointer(t *testing.T) {
	type args struct {
		val any
	}
	ival := 10
	fval := 10.1
	sval := "test"
	stval := args{}
	mval := map[string]string{}
	arrval := []int{123}
	tests := []struct {
		name string
		args args
	}{
		{
			"int",
			args{
				ival,
			},
		},
		{
			"float",
			args{
				fval,
			},
		},
		{
			"string",
			args{
				sval,
			},
		},
		{
			"struct",
			args{
				stval,
			},
		},
		{
			"map",
			args{
				mval,
			},
		},
		{
			"arr",
			args{
				arrval,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := From(tt.args.val)
			assert.Equal(t, tt.args.val, *got)
		})
	}
}
