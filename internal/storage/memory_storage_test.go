package storage

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMemoryStorage_PutGetListDeleteSimple(t *testing.T) {
	st := NewMemoryStorage[int, int]()
	ctx := context.TODO()
	for i := 0; i < 3; i++ {
		_, _ = st.Put(ctx, i, i+1)
		v1, ok := st.Get(ctx, i)
		assert.Equal(t, true, ok)
		assert.Equal(t, i+1, v1)
	}
	lst, _ := st.List(ctx)
	assert.Len(t, lst, 3)
	assert.ElementsMatch(t, []KeyValuer[int, int]{
		{0, 1},
		{1, 2},
		{2, 3},
	}, lst)
	_ = st.Delete(ctx, 0)
	lst, _ = st.List(ctx)
	assert.Len(t, lst, 2)
	_, ok := st.Get(ctx, 0)
	assert.False(t, ok)
	lst, _ = st.List(ctx)
	assert.ElementsMatch(t, []KeyValuer[int, int]{
		{1, 2},
		{2, 3},
	}, lst)
}

func TestMemoryStorage_PutGetListDeleteStruct(t *testing.T) {
	type KV struct {
		K1 int
		K2 int
	}
	st := NewMemoryStorage[KV, int]()
	ctx := context.TODO()
	for i := 0; i < 3; i++ {
		_, _ = st.Put(ctx, KV{i, i}, i)
		v, ok := st.Get(ctx, KV{i, i})
		assert.Equal(t, true, ok)
		assert.Equal(t, i, v)
	}
	lst, _ := st.List(ctx)
	assert.Len(t, lst, 3)
	assert.ElementsMatch(t, []KeyValuer[KV, int]{
		{KV{0, 0}, 0},
		{KV{1, 1}, 1},
		{KV{2, 2}, 2},
	}, lst)
	_ = st.Delete(ctx, KV{2, 2})
	lst, _ = st.List(ctx)
	assert.Len(t, lst, 2)
	_, ok := st.Get(ctx, KV{2, 2})
	assert.False(t, ok)
	lst, _ = st.List(ctx)
	assert.ElementsMatch(t, []KeyValuer[KV, int]{
		{KV{0, 0}, 0},
		{KV{1, 1}, 1},
	}, lst)
}
