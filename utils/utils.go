package utils

import (
	"sync"

	"github.com/pborman/uuid"
)

// Generates a uuid for a minion
func GenerateUUID(name string) uuid.UUID {
	u := uuid.NewSHA1(uuid.NameSpace_DNS, []byte(name))

	return u
}

// Map type that can be safely shared between
// goroutines that require read/write access to a map
type concurrentMap struct {
	sync.RWMutex
	items map[string]interface{}
}

// Concurrent map item
type ConcurrentMapItem struct {
	Key   string
	Value interface{}
}

// Creates a new concurrent map
func NewConcurrentMap() *concurrentMap {
	cm := &concurrentMap{
		items: make(map[string]interface{}),
	}

	return cm
}

// Sets a key in a concurrent map
func (cm *concurrentMap) Set(key string, value interface{}) {
	cm.Lock()
	defer cm.Unlock()

	cm.items[key] = value
}

// Gets a key from a concurrent map
func (cm *concurrentMap) Get(key string) (interface{}, bool) {
	cm.Lock()
	defer cm.Unlock()

	value, ok := cm.items[key]

	return value, ok
}

// Iterates over the items in a concurrent map
// Each item is sent over a channel, so that
// we can iterate over the map using the builtin range keyword
func (cm *concurrentMap) Iter() <-chan ConcurrentMapItem {
	c := make(chan ConcurrentMapItem)

	f := func() {
		cm.Lock()
		defer cm.Unlock()

		for k, v := range cm.items {
			c <- ConcurrentMapItem{k, v}
		}
		close(c)
	}
	go f()

	return c
}

// Slice type that can be safely shared between goroutines
type concurrentSlice struct {
	sync.RWMutex
	items []interface{}
}

// Concurrent slice item
type ConcurrentSliceItem struct {
	Index int
	Value interface{}
}

// Convenience function that creates a new concurrent slice
func NewConcurrentSlice() *concurrentSlice {
	cs := &concurrentSlice{
		items: make([]interface{}, 0),
	}

	return cs
}

// Append an item to the concurrent slice
func (cs *concurrentSlice) Append(item interface{}) {
	cs.Lock()
	defer cs.Unlock()

	cs.items = append(cs.items, item)
}

// Iterates over the items in the concurrent slice
// Each item is sent over a channel, so that
// we can iterate over the slice using the builin range keyword
func (cs *concurrentSlice) Iter() <-chan ConcurrentSliceItem {
	c := make(chan ConcurrentSliceItem)

	f := func() {
		cs.Lock()
		defer cs.Lock()
		for index, value := range cs.items {
			c <- ConcurrentSliceItem{index, value}
		}
		close(c)
	}
	go f()

	return c
}
