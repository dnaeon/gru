package resource

import (
	"github.com/layeh/gopher-luar"
	"github.com/yuin/gopher-lua"
)

// LuaRegisterBuiltin registers resource providers in Lua
func LuaRegisterBuiltin(L *lua.LState) {
	for typ, provider := range providerRegistry {
		// Wrap resource providers, so that we can properly handle any
		// errors returned by providers during resource instantiation.
		// Since we don't want to return the error to Lua, this is the
		// place where we handle any errors returned by providers.
		wrapper := func(L *lua.LState) int {
			r, err := provider(L.CheckString(1))
			if err != nil {
				L.RaiseError(err.Error())
			}

			L.Push(luar.New(L, r))
			return 1
		}

		tbl := L.NewTable()
		tbl.RawSetH(lua.LString("new"), L.NewFunction(wrapper))
		L.SetGlobal(typ, tbl)
	}
}
