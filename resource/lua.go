package resource

import (
	"github.com/layeh/gopher-luar"
	"github.com/yuin/gopher-lua"
)

// RegisterLuaBuiltins registers resource providers in Lua
func RegisterLuaBuiltins(L *lua.LState) error {
	// Register resource providers
	for typ, provider := range providerRegistry {
		tbl := L.NewTable()
		tbl.RawSetH(lua.LString("new"), luar.New(L, provider))
		L.SetGlobal(typ, tbl)
	}
}
