package utils

import (
	"fmt"
	"slices"
	"sync"
)

type ConcurrentSlice[T any] struct {
	sync.RWMutex
	items []T
}

func NewConcurrentSlice[T any]() *ConcurrentSlice[T] {
	return &ConcurrentSlice[T]{
		items: []T{},
	}
}

func (cs *ConcurrentSlice[T]) Append(item T) {
	cs.Lock()
	defer cs.Unlock()
	cs.items = append(cs.items, item)
}

func (cs *ConcurrentSlice[T]) Get(index int) (T, error) {
	cs.RLock()
	defer cs.RUnlock()
	if index < 0 || index >= len(cs.items) {
		var zero T
		return zero, fmt.Errorf("index out of range")
	}
	return cs.items[index], nil
}

func (cs *ConcurrentSlice[T]) Len() int {
	cs.RLock()
	defer cs.RUnlock()
	return len(cs.items)
}

func (cs *ConcurrentSlice[T]) IndexFunc(f func(T) bool) int {
	cs.RLock()
	defer cs.RUnlock()
	return slices.IndexFunc(cs.items, f)
}

func (cs *ConcurrentSlice[T]) Remove(index int) error {
	cs.Lock()
	defer cs.Unlock()
	if index < 0 || index >= len(cs.items) {
		return fmt.Errorf("index out of range")
	}
	cs.items = slices.Delete(cs.items, index, index+1)
	return nil
}

func (cs *ConcurrentSlice[T]) FindFunc(f func(T) bool) (T, error) {
	i := cs.IndexFunc(f)
	if i == -1 {
		var zero T
		return zero, fmt.Errorf("Could not find value")
	}
	return cs.items[i], nil
}

func (cs *ConcurrentSlice[T]) DeleteFunc(f func(T) bool) {
	cs.Lock()
	defer cs.Unlock()
	cs.items = slices.DeleteFunc(cs.items, f)
}

func (cs *ConcurrentSlice[T]) Items() []T {
	cs.RLock()
	defer cs.RUnlock()
	snapshot := make([]T, len(cs.items))
	copy(snapshot, cs.items)
	return snapshot
}

// Iterator type for the concurrent slice
type Iterator[T any] struct {
	slice    *ConcurrentSlice[T]
	index    int
	released bool
}

func (it *Iterator[T]) Next() bool {
	if it.released {
		return false
	}
	it.index++
	return it.index < len(it.slice.items)
}

func (it *Iterator[T]) Value() T {
	return it.slice.items[it.index]
}

func (it *Iterator[T]) Release() {
	if !it.released {
		it.slice.RUnlock()
		it.released = true
	}
}

func (cs *ConcurrentSlice[T]) Iterator() *Iterator[T] {
	cs.RLock() // Lock for reading
	return &Iterator[T]{
		slice:    cs,
		index:    -1,
		released: false,
	}
}
