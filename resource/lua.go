package resource

import (
	"github.com/layeh/gopher-luar"
	"github.com/yuin/gopher-lua"
)

// LuaRegisterBuiltin registers resource providers in Lua
func LuaRegisterBuiltin(L *lua.LState) {
	// Go functions registered in Lua
	builtins := map[string]interface{}{
		"log": Log,
	}

	// Register functions in Lua
	for name, fn := range builtins {
		L.SetGlobal(name, luar.New(L, fn))
	}

	// The resource namespace in Lua
	namespace := L.NewTable()
	L.SetGlobal("resource", namespace)

	// Register resource providers in Lua
	for typ, provider := range providerRegistry {
		// Wrap resource providers, so that we can properly handle any
		// errors returned by providers during resource instantiation.
		// Since we don't want to return the error to Lua, this is the
		// place where we handle any errors returned by providers.
		wrapper := func(p Provider) lua.LGFunction {
			return func(L *lua.LState) int {
				// Create the resource by calling it's provider
				r, err := p(L.CheckString(1))
				if err != nil {
					L.RaiseError(err.Error())
				}

				L.Push(luar.New(L, r))
				return 1 // Number of arguments returned to Lua
			}
		}

		tbl := L.NewTable()
		tbl.RawSetH(lua.LString("new"), L.NewFunction(wrapper(provider)))
		L.SetField(namespace, typ, tbl)
	}
}
