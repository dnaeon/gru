package resource

import "testing"

func TestShell(t *testing.T) {
	L := newLuaState()
	defer L.Close()

	const code = `
	sh = shell.new("create /tmp/foo file")
	sh.command = "touch /tmp/foo"
	sh.creates = "/tmp/foo"
	`

	if err := L.DoString(code); err != nil {
		t.Fatal(err)
	}

	sh := luaResource(L, "sh").(*Shell)
	errorIfNotEqual(t, "shell", sh.Type)
	errorIfNotEqual(t, "create /tmp/foo file", sh.Name)
	errorIfNotEqual(t, "present", sh.State)
	errorIfNotEqual(t, []string{}, sh.After)
	errorIfNotEqual(t, []string{}, sh.Before)
	errorIfNotEqual(t, []string{"present"}, sh.PresentStates)
	errorIfNotEqual(t, []string{"absent"}, sh.AbsentStates)
	errorIfNotEqual(t, "touch /tmp/foo", sh.Command)
	errorIfNotEqual(t, "/tmp/foo", sh.Creates)
}
