package main

import "github.com/gammazero/deque"

// TODO: Write tests for history.
type History struct {
	queue *deque.Deque[string]
	iter  int
}

const HISTORY_SIZE = 50

// Return if we inserted.
func (h *History) Push(in string) bool {
	if h.queue.Len() > 0 {
		if in == h.queue.Back() {
			return false
		}
	}

	h.queue.PushBack(in)
	h.iter = 0

	if h.queue.Len() > HISTORY_SIZE {
		h.queue.PopFront()
	}

	return true
}

func (h *History) Pop() {
	if h.queue.Len() == 0 {
		return
	}

	h.queue.PopBack()
}

func (h *History) Down() {
	h.iter -= 1
	if h.iter < -1 {
		h.iter = -1
	}
}

func (h *History) Up() {
	if h.queue.Len() == 0 {
		return
	}

	h.iter += 1
	if h.iter >= h.queue.Len() {
		h.iter = h.queue.Len() - 1
	}
}

func (h *History) Entry() string {
	if h.iter == -1 {
		return ""
	}
	index := h.queue.Len() - h.iter - 1
	return h.queue.At(index)
}
