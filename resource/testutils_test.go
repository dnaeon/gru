package resource

import (
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	"github.com/yuin/gopher-lua"
)

// newLuaState creates a new Lua state and registers the
// resource providers in Lua. It is up to the caller to
// close the Lua state once done with it.
func newLuaState() *lua.LState {
	L := lua.NewState()
	LuaRegisterBuiltin(L)

	return L
}

// luaResource retrieves a resource by it's name
func luaResource(L *lua.LState, name string) interface{} {
	return L.GetGlobal(name).(*lua.LUserData).Value
}

func positionString(level int) string {
	_, file, line, _ := runtime.Caller(level + 1)

	return fmt.Sprintf("%v:%v:", filepath.Base(file), line)
}

func errorIfNotEqual(t *testing.T, v1, v2 interface{}) {
	if !reflect.DeepEqual(v1, v2) {
		t.Errorf("%v want '%v', got '%v'", positionString(1), v1, v2)
	}
}
