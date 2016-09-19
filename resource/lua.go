package resource

import (
	"github.com/layeh/gopher-luar"
	"github.com/yuin/gopher-lua"
)

// DefaultResourceNamespace is the Lua table where resources are being
// registered to, when using the default namespace.
const DefaultResourceNamespace = "resource"

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

	// Register resource providers in Lua
	for _, item := range providerRegistry {
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

		// Create the resource namespace
		namespace := L.GetGlobal(item.Namespace)
		if lua.LVIsFalse(namespace) {
			namespace = L.NewTable()
			L.SetGlobal(item.Namespace, namespace)
		}

		tbl := L.NewTable()
		tbl.RawSetH(lua.LString("new"), L.NewFunction(wrapper(item.Provider)))

		L.SetField(namespace, item.Type, tbl)
	}
}
