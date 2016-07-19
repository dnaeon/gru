package utils

import "sync"

// ConcurrentSlice type that can be safely shared between goroutines
type ConcurrentSlice struct {
	sync.RWMutex
	items []interface{}
}

// ConcurrentSliceItem contains the index/value pair of an item in a
// concurrent slice
type ConcurrentSliceItem struct {
	Index int
	Value interface{}
}

// NewConcurrentSlice creates a new concurrent slice
func NewConcurrentSlice() *ConcurrentSlice {
	cs := &ConcurrentSlice{
		items: make([]interface{}, 0),
	}

	return cs
}

// Append adds an item to the concurrent slice
func (cs *ConcurrentSlice) Append(item interface{}) {
	cs.Lock()
	defer cs.Unlock()

	cs.items = append(cs.items, item)
}

// Iter iterates over the items in the concurrent slice
// Each item is sent over a channel, so that
// we can iterate over the slice using the builin range keyword
func (cs *ConcurrentSlice) Iter() <-chan ConcurrentSliceItem {
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

// List type represents a slice of strings
type List []string

// NewList creates a new list with the given items
func NewList(s ...string) List {
	l := make(List, len(s))
	for _, v := range s {
		l = append(l, v)
	}

	return l
}

// Contains returns a boolean indicating whether the list
// contains the given string.
func (l List) Contains(x string) bool {
	for _, v := range l {
		if v == x {
			return true
		}
	}

	return false
}

// String type represents a string
type String struct {
	str string
}

// NewString creates a new string
func NewString(s string) String {
	return String{
		str: s,
	}
}

// IsInList returns a boolean indicating whether the string is
// contained within a given list
func (s String) IsInList(l List) bool {
	return l.Contains(s.str)
}
