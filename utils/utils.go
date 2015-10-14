package utils

// Map type that can be safely shared between
// goroutines that require read/write access to a map
type concurrentMap struct {
	sync.RWMutex
	items map[string]interface{}
}

// Concurrent map item
type ConcurrentMapItem struct {
	Key string
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
	cm.items[key] = value
	cm.Unlock()
}

// Gets a key from a concurrent map
func (cm *concurrentMap) Get(key string) (interface{}, bool) {
	defer cm.Unlock()

	cm.Lock()
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
		for k, v := range cm.items {
			c <- concurrentMapItem{k, v}
		}
		cm.Unlock()
		close(c)
	}
	go f()

	return c
}
