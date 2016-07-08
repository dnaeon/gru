package resource

import "testing"

func TestPacman(t *testing.T) {
	L := newLuaState()
	defer L.Close()

	const code = `
	pkg = pacman.new("tmux")
	`

	if err := L.DoString(code); err != nil {
		t.Fatal(err)
	}

	pkg := luaResource(L, "pkg").(*Pacman)
	errorIfNotEqual(t, "package", pkg.Type)
	errorIfNotEqual(t, "tmux", pkg.Name)
	errorIfNotEqual(t, StatePresent, pkg.State)
	errorIfNotEqual(t, []string{}, pkg.Before)
	errorIfNotEqual(t, []string{}, pkg.After)
	errorIfNotEqual(t, "tmux", pkg.Package)
	errorIfNotEqual(t, "", pkg.Version)
}

func TestYum(t *testing.T) {
	L := newLuaState()
	defer L.Close()

	const code = `
	pkg = yum.new("tmux")
	`

	if err := L.DoString(code); err != nil {
		t.Fatal(err)
	}

	pkg := luaResource(L, "pkg").(*Yum)
	errorIfNotEqual(t, "package", pkg.Type)
	errorIfNotEqual(t, "tmux", pkg.Name)
	errorIfNotEqual(t, StatePresent, pkg.State)
	errorIfNotEqual(t, []string{}, pkg.Before)
	errorIfNotEqual(t, []string{}, pkg.After)
	errorIfNotEqual(t, "tmux", pkg.Package)
	errorIfNotEqual(t, "", pkg.Version)
}
