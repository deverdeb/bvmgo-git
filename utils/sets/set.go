package sets

import (
	"iter"
	"maps"
)

type Set[T comparable] struct {
	values map[T]bool
}

func (set Set[T]) Add(value T) {
	set.values[value] = true
}

func (set Set[T]) AddAll(values ...T) {
	for _, value := range values {
		set.Add(value)
	}
}

func (set Set[T]) Remove(value T) {
	delete(set.values, value)
}

func (set Set[T]) RemoveAll(values ...T) {
	for _, value := range values {
		set.Remove(value)
	}
}

func (set Set[T]) Contains(value T) bool {
	_, contains := set.values[value]
	return contains
}

func (set Set[T]) Size() int {
	return len(set.values)
}

func (set Set[T]) Iter() iter.Seq[T] {
	return maps.Keys(set.values)
}

func (set Set[T]) ToSlice() []T {
	result := make([]T, len(set.values))
	idx := 0
	for value := range set.Iter() {
		result[idx] = value
		idx++
	}
	return result
}

func (set Set[T]) Clone() Set[T] {
	return Clone(set)
}
