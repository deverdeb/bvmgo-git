package sets

func New[T comparable]() Set[T] {
	return Set[T]{
		values: make(map[T]bool),
	}
}

func NewFromSlice[T comparable](data []T) Set[T] {
	return NewFromData(data...)
}

func NewFromData[T comparable](data ...T) Set[T] {
	set := Set[T]{
		values: make(map[T]bool),
	}
	set.AddAll(data...)
	return set
}

func Clone[T comparable](set Set[T]) Set[T] {
	clone := New[T]()
	for key, value := range set.values {
		clone.values[key] = value
	}
	return clone
}
