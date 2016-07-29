package utils

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

// String implements the fmt.Stringer interface
func (s String) String() string {
	return s.str
}

// IsInList returns a boolean indicating whether the string is
// contained within a given list
func (s String) IsInList(l List) bool {
	return l.Contains(s.str)
}
