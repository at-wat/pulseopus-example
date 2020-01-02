package main

import (
	"sync"
)

type buffer struct {
	b  []int16
	mu sync.Mutex
}

func (b *buffer) Read(p []int16) {
	b.mu.Lock()
	if len(b.b) > len(p) {
		copy(p, b.b[:len(p)])
		b.b = b.b[len(p):]
		b.mu.Unlock()
		return
	}
	b.mu.Unlock()
}

func (b *buffer) Write(p []int16) {
	b.mu.Lock()
	b.b = append(b.b, p...)
	b.mu.Unlock()
}

func (b *buffer) Len() int {
	b.mu.Lock()
	defer b.mu.Unlock()
	return len(b.b)
}
