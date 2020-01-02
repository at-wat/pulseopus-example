package main

import (
	"sync"
)

type buffer struct {
	b  []int16
	mu sync.RWMutex
}

func (b *buffer) Read(p []int16) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if len(b.b) > len(p) {
		copy(p, b.b[:len(p)])
		b.b = b.b[len(p):]
	}
}

func (b *buffer) Write(p []int16) {
	b.mu.Lock()
	b.b = append(b.b, p...)
	b.mu.Unlock()
}

func (b *buffer) Len() int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return len(b.b[:])
}
