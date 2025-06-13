package model

import "container/heap"

type Node struct {
	Char  byte  `json:"char"`
	Freq  int   `json:"freq"`
	Left  *Node `json:"left"`
	Right *Node `json:"right"`
}

type MinHeap []*Node

func NewMinHeapWithFreqMap(freqMap map[byte]int) *MinHeap {
	h := &MinHeap{}
	heap.Init(h)

	for char, freq := range freqMap {
		heap.Push(h, &Node{
			Char: char,
			Freq: freq,
		})
	}

	return h
}

func (h *MinHeap) Len() int {
	return len(*h)
}

func (h *MinHeap) Less(a, b int) bool {
	return (*h)[a].Freq < (*h)[b].Freq
}

func (h *MinHeap) Swap(a, b int) {
	(*h)[a], (*h)[b] = (*h)[b], (*h)[a]
}

func (h *MinHeap) Push(x any) {
	*h = append(*h, x.(*Node))
}

func (h *MinHeap) Pop() any {
	old := *h
	n := len(old)
	item := old[n-1]
	*h = old[0 : n-1]
	return item
}
