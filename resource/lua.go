// Copyright (c) 2015-2017 Marin Atanasov Nikolov <dnaeon@gmail.com>
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions
// are met:
//
//  1. Redistributions of source code must retain the above copyright
//     notice, this list of conditions and the following disclaimer
//     in this position and unchanged.
//  2. Redistributions in binary form must reproduce the above copyright
//     notice, this list of conditions and the following disclaimer in the
//     documentation and/or other materials provided with the distribution.
//
// THIS SOFTWARE IS PROVIDED BY THE AUTHOR(S) ``AS IS'' AND ANY EXPRESS OR
// IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES
// OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED.
// IN NO EVENT SHALL THE AUTHOR(S) BE LIABLE FOR ANY DIRECT, INDIRECT,
// INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT
// NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF
// THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package resource

import (
	"github.com/yuin/gopher-lua"
	"layeh.com/gopher-luar"
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
