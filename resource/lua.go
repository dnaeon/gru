package resource

import (
	"github.com/layeh/gopher-luar"
	"github.com/yuin/gopher-lua"
)

// functionRegistry contains the functions to be registered in Lua.
var functionRegistry = make([]FunctionItem, 0)

// DefaultResourceNamespace is the Lua table where resources are being
// registered to, when using the default namespace.
const DefaultResourceNamespace = "resource"

// DefaultFunctionNamespace is the Lua table where functions are being
// registered to, when using the default namespace.
const DefaultFunctionNamespace = "stdlib"

// FunctionItem type represents a single item from the function registry.
type FunctionItem struct {
	// Name of the function to register in Lua
	Name string

	// Namespace is the Lua table where the function will be registered to
	Namespace string

	// Function to execute when called from Lua
	Function interface{}
}

// RegisterFunction registers a function to the registry.
func RegisterFunction(items ...FunctionItem) {
	functionRegistry = append(functionRegistry, items...)
}

// LuaRegisterBuiltin registers resource providers and functions in Lua.
func LuaRegisterBuiltin(L *lua.LState) {
	// Register functions in Lua
	for _, item := range functionRegistry {
		namespace := L.GetGlobal(item.Namespace)
		if lua.LVIsFalse(namespace) {
			namespace = L.NewTable()
			L.SetGlobal(item.Namespace, namespace)
		}
		L.SetField(namespace, item.Name, luar.New(L, item.Function))
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

func init() {
	logf := FunctionItem{
		Name:      "logf",
		Namespace: "stdlib",
		Function:  Logf,
	}

	RegisterFunction(logf)
}
