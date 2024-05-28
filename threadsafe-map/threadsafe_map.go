package main

import (
	"fmt"
	"sync"
)

type ThreadSafeMap[K comparable, V any] struct {
	internalMap map[K]V
	mutex      *sync.RWMutex
}

func NewThreadSafeMap[K comparable, V any]() *ThreadSafeMap[K ,V] {
	return &ThreadSafeMap[K, V]{
		internalMap: make(map[K]V),
		mutex:      &sync.RWMutex{},
	}
}

func (m *ThreadSafeMap[K, V]) Get(key K) (V, bool) {
	m.mutex.RLock()
	val, ok := m.internalMap[key]
	m.mutex.RUnlock()
	return val, ok
}

func (m *ThreadSafeMap[K, V]) Set(key K, val V) {
	m.mutex.Lock()
	m.internalMap[key] = val
	m.mutex.Unlock()
}

func (m *ThreadSafeMap[K, V]) Len() int {
	m.mutex.RLock()
	length := len(m.internalMap)
	m.mutex.RUnlock()
	return length
}

func main() {
	m := NewThreadSafeMap[int, int]()
	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		for i := 0; i < 1000; i++ {
			m.Set(i, i)
			fmt.Println("Set key:", i, "Len:", m.Len())
		}
		fmt.Println("Done setting")
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		for i := 0; i < 1000; i++ {
			val, ok := m.Get(i)
			if ok {
				fmt.Println("Key:", i, "Value:", val, "Len:", m.Len())
			} else {
				fmt.Println("Key not found:", i, "Len:", m.Len())
			}
		}
		fmt.Println("Done getting")
		wg.Done()
	}()

	wg.Wait()
}
