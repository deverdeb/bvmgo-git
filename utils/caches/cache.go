package caches

type Cache[K comparable, T any] interface {
	Put(key K, value T)
	Get(key K) T
	GetOrCompute(key K, provider func(K) T) T
	Size() int
	Clear()
}

func BuildCache[K comparable, T any]() Cache[K, T] {
	return &cacheImpl[K, T]{data: make(map[K]T)}
}

type cacheImpl[K comparable, T any] struct {
	data map[K]T
}

func (cache *cacheImpl[K, T]) Put(key K, value T) {
	cache.data[key] = value
}

func (cache *cacheImpl[K, T]) Get(key K) T {
	value, _ := cache.data[key]
	return value
}

func (cache *cacheImpl[K, T]) GetOrCompute(key K, provider func(K) T) T {
	value, ok := cache.data[key]
	if !ok {
		value = provider(key)
		cache.Put(key, value)
	}
	return value
}

func (cache *cacheImpl[K, T]) Size() int {
	return len(cache.data)
}

func (cache *cacheImpl[K, T]) Clear() {
	cache.data = make(map[K]T)
}
