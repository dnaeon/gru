package resource

import "testing"

func TestPacman(t *testing.T) {
	L := newLuaState()
	defer L.Close()

	const code = `
	tmux = pacman.new("tmux")
	`

	if err := L.DoString(code); err != nil {
		t.Fatal(err)
	}

	pkg := luaResource(L, "tmux").(*Pacman)
	errorIfNotEqual(t, "pkg", pkg.Type)
	errorIfNotEqual(t, "tmux", pkg.Name)
	errorIfNotEqual(t, StatePresent, pkg.State)
	errorIfNotEqual(t, []string{}, pkg.Before)
	errorIfNotEqual(t, []string{}, pkg.After)
	errorIfNotEqual(t, []string{"present", "installed"}, pkg.PresentStates)
	errorIfNotEqual(t, []string{"absent", "deinstalled"}, pkg.AbsentStates)
	errorIfNotEqual(t, "tmux", pkg.Package)
	errorIfNotEqual(t, "", pkg.Version)
}

func TestYum(t *testing.T) {
	L := newLuaState()
	defer L.Close()

	const code = `
	tmux = yum.new("tmux")
	`

	if err := L.DoString(code); err != nil {
		t.Fatal(err)
	}

	pkg := luaResource(L, "tmux").(*Yum)
	errorIfNotEqual(t, "pkg", pkg.Type)
	errorIfNotEqual(t, "tmux", pkg.Name)
	errorIfNotEqual(t, StatePresent, pkg.State)
	errorIfNotEqual(t, []string{}, pkg.Before)
	errorIfNotEqual(t, []string{}, pkg.After)
	errorIfNotEqual(t, []string{"present", "installed"}, pkg.PresentStates)
	errorIfNotEqual(t, []string{"absent", "deinstalled"}, pkg.AbsentStates)
	errorIfNotEqual(t, "tmux", pkg.Package)
	errorIfNotEqual(t, "", pkg.Version)
}
