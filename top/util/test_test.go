package util

import (
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestHeap(t *testing.T) {
	myHeap := NewPairHeap(2)
	myHeap.TryPush(Pair{K: "1", V: 1})
	myHeap.TryPush(Pair{K: "2", V: 2})
	assert.Equal(t, myHeap.Len(), 2)
	assert.Equal(t, myHeap.Min(), Pair{K: "1", V: 1})

	myHeap.TryPush(Pair{K: "3", V: 3})
	assert.Equal(t, myHeap.Len(), 2)
	assert.Equal(t, myHeap.Min(), Pair{K: "2", V: 2})

	myHeap.TryPush(Pair{K: "1", V: 1})
	assert.Equal(t, myHeap.Len(), 2)
	assert.Equal(t, myHeap.Min(), Pair{K: "2", V: 2})

	myHeap.TryPush(Pair{K: "1", V: 2})
	assert.Equal(t, myHeap.Len(), 2)
	assert.Equal(t, myHeap.Min(), Pair{K: "1", V: 2})

	myHeap.TryPush(Pair{K: "1", V: 4})
	assert.Equal(t, myHeap.Len(), 2)
	assert.Equal(t, myHeap.Min(), Pair{K: "3", V: 3})

	for i := 0; i < 100; i++ {
		myHeap.TryPush(Pair{K: "0", V: int64(i)})
	}
	assert.Equal(t, myHeap.Len(), 2)
	assert.Equal(t, myHeap.Min(), Pair{K: "0", V: 98})
}

func TestHeap2(t *testing.T) {
	myHeap := NewPairHeap(100)
	myHeap.TryPush(Pair{K: "100", V: int64(100)})
	for i := 0; i < 100; i++ {
		myHeap.TryPush(Pair{K: "100", V: int64(i)})
	}
	myHeap.TryPush(Pair{K: "100", V: int64(101)})
	assert.Equal(t, myHeap.Min(), Pair{K: "100", V: int64(100)})
}
