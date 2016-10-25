package resource

import (
	"os"
	"testing"
)

func TestFile(t *testing.T) {
	L := newLuaState()
	defer L.Close()

	const code = `
	foo = resource.file.new("/tmp/foo")
	`

	if err := L.DoString(code); err != nil {
		t.Fatal(err)
	}

	foo := luaResource(L, "foo").(*File)
	errorIfNotEqual(t, "file", foo.Type)
	errorIfNotEqual(t, "/tmp/foo", foo.Name)
	errorIfNotEqual(t, "present", foo.State)
	errorIfNotEqual(t, []string{}, foo.Require)
	errorIfNotEqual(t, []string{"present"}, foo.PresentStatesList)
	errorIfNotEqual(t, []string{"absent"}, foo.AbsentStatesList)
	errorIfNotEqual(t, true, foo.Concurrent)
	errorIfNotEqual(t, "/tmp/foo", foo.Path)
	errorIfNotEqual(t, os.FileMode(0644), foo.Mode)
	errorIfNotEqual(t, "", foo.Source)
}

func TestDirectory(t *testing.T) {
	L := newLuaState()
	defer L.Close()

	const code = `
	bar = resource.directory.new("/tmp/bar")
	`

	if err := L.DoString(code); err != nil {
		t.Fatal(err)
	}

	bar := luaResource(L, "bar").(*Directory)
	errorIfNotEqual(t, "directory", bar.Type)
	errorIfNotEqual(t, "/tmp/bar", bar.Name)
	errorIfNotEqual(t, "present", bar.State)
	errorIfNotEqual(t, []string{}, bar.Require)
	errorIfNotEqual(t, []string{"present"}, bar.PresentStatesList)
	errorIfNotEqual(t, []string{"absent"}, bar.AbsentStatesList)
	errorIfNotEqual(t, true, bar.Concurrent)
	errorIfNotEqual(t, "/tmp/bar", bar.Path)
	errorIfNotEqual(t, os.FileMode(0755), bar.Mode)
	errorIfNotEqual(t, false, bar.Parents)
}

func TestLink(t *testing.T) {
	L := newLuaState()
	defer L.Close()

	const code = `
	qux = resource.link.new("/tmp/qux")
	qux.source = "/tmp/foo"
	`

	if err := L.DoString(code); err != nil {
		t.Fatal(err)
	}

	qux := luaResource(L, "qux").(*Link)
	errorIfNotEqual(t, "link", qux.Type)
	errorIfNotEqual(t, "/tmp/qux", qux.Name)
	errorIfNotEqual(t, "present", qux.State)
	errorIfNotEqual(t, []string{}, qux.Require)
	errorIfNotEqual(t, []string{"present"}, qux.PresentStatesList)
	errorIfNotEqual(t, []string{"absent"}, qux.AbsentStatesList)
	errorIfNotEqual(t, true, qux.Concurrent)
	errorIfNotEqual(t, "/tmp/foo", qux.Source)
	errorIfNotEqual(t, false, qux.Hard)
}
