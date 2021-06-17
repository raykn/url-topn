package util

import "container/heap"

type PairHeap struct {
	d    []Pair
	topN int
}

func NewPairHeap(topN int) *PairHeap {
	h := &PairHeap{
		d:    make([]Pair, 0),
		topN: topN,
	}
	return h
}

func (h PairHeap) Len() int           { return len(h.d) }
func (h PairHeap) Less(i, j int) bool { return !h.d[i].Greater(h.d[j]) }
func (h PairHeap) Swap(i, j int)      { h.d[i], h.d[j] = h.d[j], h.d[i] }
func (h *PairHeap) Push(x interface{}) {
	h.d = append(h.d, x.(Pair))
}
func (h *PairHeap) Pop() interface{} {
	old := h.d
	n := len(old)
	x := old[n-1]
	h.d = old[:n-1]
	return x
}
func (h *PairHeap) Min() Pair {
	return h.d[0]
}
func (h *PairHeap) Data() []Pair {
	return h.d
}
func (h *PairHeap) TryPush(p Pair) (push, pop bool) {
	if h.Len() > 0 && h.Min().Greater(p) {
		return
	}
	if h.topN > 0 && h.Len() >= h.topN {
		heap.Pop(h)
		pop = true
	}
	heap.Push(h, p)
	push = true
	return
}
