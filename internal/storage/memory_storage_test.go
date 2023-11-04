package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMemoryStorage_PutGetListDeleteSimple(t *testing.T) {
	st := NewMemoryStorage[int, int]()
	for i := 0; i < 3; i++ {
		st.Put(i, i+1)
		v1, ok := st.Get(i)
		assert.Equal(t, true, ok)
		assert.Equal(t, i+1, v1)
	}
	lst := st.List()
	assert.Len(t, lst, 3)
	assert.ElementsMatch(t, []KeyValuer[int, int]{
		{0, 1},
		{1, 2},
		{2, 3},
	}, lst)
	st.Delete(0)
	assert.Len(t, st.List(), 2)
	_, ok := st.Get(0)
	assert.False(t, ok)
	assert.ElementsMatch(t, []KeyValuer[int, int]{
		{1, 2},
		{2, 3},
	}, st.List())
}

func TestMemoryStorage_PutGetListDeleteStruct(t *testing.T) {
	type KV struct {
		K1 int
		K2 int
	}
	st := NewMemoryStorage[KV, int]()
	for i := 0; i < 3; i++ {
		st.Put(KV{i, i}, i)
		v, ok := st.Get(KV{i, i})
		assert.Equal(t, true, ok)
		assert.Equal(t, i, v)
	}
	lst := st.List()
	assert.Len(t, lst, 3)
	assert.ElementsMatch(t, []KeyValuer[KV, int]{
		{KV{0, 0}, 0},
		{KV{1, 1}, 1},
		{KV{2, 2}, 2},
	}, lst)
	st.Delete(KV{2, 2})
	assert.Len(t, st.List(), 2)
	_, ok := st.Get(KV{2, 2})
	assert.False(t, ok)
	assert.ElementsMatch(t, []KeyValuer[KV, int]{
		{KV{0, 0}, 0},
		{KV{1, 1}, 1},
	}, st.List())
}

// func TestRestoreStorage(t *testing.T) {
// 	type args struct {
// 		filename string
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		wantErr bool
// 	}{
// 		{
// 			"ok",
// 			args{
// 				"./test/test.json",
// 			},
// 			false,
// 		},
// 		{
// 			"false",
// 			args{
// 				"./test/test2.json",
// 			},
// 			true,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			st := storage.NewMemoryStorage[types.MetricsKey, types.Metrics]()
// 			err := RestoreStorage(st, tt.args.filename)
// 			if tt.wantErr {
// 				assert.Error(t, err)
// 			} else {
// 				assert.NoError(t, err)
// 			}
// 			if !tt.wantErr {
// 				elems := st.List()
// 				assert.Equal(t, 2, len(elems))
// 			}
// 		})
// 	}
// }
