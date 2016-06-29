package resource

import (
	"os"
	"testing"
)

func TestFile(t *testing.T) {
	L := newLuaState()
	defer L.Close()

	const code = `
	foo = file.new("/tmp/foo")
	`

	if err := L.DoString(code); err != nil {
		t.Fatal(err)
	}

	foo := luaResource(L, "foo").(*File)
	errorIfNotEqual(t, "file", foo.Type)
	errorIfNotEqual(t, "/tmp/foo", foo.Name)
	errorIfNotEqual(t, StatePresent, foo.State)
	errorIfNotEqual(t, []string{}, foo.Before)
	errorIfNotEqual(t, []string{}, foo.After)
	errorIfNotEqual(t, "/tmp/foo", foo.Path)
	errorIfNotEqual(t, os.FileMode(0644), foo.Mode)
	errorIfNotEqual(t, "", foo.Source)
	errorIfNotEqual(t, fileTypeRegular, foo.FileType)
	errorIfNotEqual(t, false, foo.Recursive)
	errorIfNotEqual(t, false, foo.Purge)
}
