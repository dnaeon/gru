package resource

import "testing"

func TestPacman(t *testing.T) {
	L := newLuaState()
	defer L.Close()

	const code = `
	tmux = resource.pacman.new("tmux")
	`

	if err := L.DoString(code); err != nil {
		t.Fatal(err)
	}

	pkg := luaResource(L, "tmux").(*Pacman)
	errorIfNotEqual(t, "package", pkg.Type)
	errorIfNotEqual(t, "tmux", pkg.Name)
	errorIfNotEqual(t, "installed", pkg.State)
	errorIfNotEqual(t, []string{}, pkg.Require)
	errorIfNotEqual(t, []string{"present", "installed"}, pkg.PresentStatesList)
	errorIfNotEqual(t, []string{"absent", "deinstalled"}, pkg.AbsentStatesList)
	errorIfNotEqual(t, false, pkg.Concurrent)
	errorIfNotEqual(t, "tmux", pkg.Package)
	errorIfNotEqual(t, "", pkg.Version)
}

func TestYum(t *testing.T) {
	L := newLuaState()
	defer L.Close()

	const code = `
	tmux = resource.yum.new("tmux")
	`

	if err := L.DoString(code); err != nil {
		t.Fatal(err)
	}

	pkg := luaResource(L, "tmux").(*Yum)
	errorIfNotEqual(t, "package", pkg.Type)
	errorIfNotEqual(t, "tmux", pkg.Name)
	errorIfNotEqual(t, "installed", pkg.State)
	errorIfNotEqual(t, []string{}, pkg.Require)
	errorIfNotEqual(t, []string{"present", "installed"}, pkg.PresentStatesList)
	errorIfNotEqual(t, []string{"absent", "deinstalled"}, pkg.AbsentStatesList)
	errorIfNotEqual(t, false, pkg.Concurrent)
	errorIfNotEqual(t, "tmux", pkg.Package)
	errorIfNotEqual(t, "", pkg.Version)
}
