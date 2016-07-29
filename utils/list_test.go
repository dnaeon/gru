package utils

import "testing"

func TestListContains(t *testing.T) {
	l := NewList("foo", "bar", "qux")
	want := "foo"

	if !l.Contains(want) {
		t.Errorf("list does not contain %q", want)
	}
}

func TestStringInList(t *testing.T) {
	l := NewList("foo", "bar", "qux")
	s := NewString("foo")

	if !s.IsInList(l) {
		t.Errorf("string %q is not in list", s)
	}
}
