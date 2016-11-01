package catalog

import (
	"log"
	"os"
	"testing"

	"github.com/dnaeon/gru/resource"
	"github.com/yuin/gopher-lua"
)

func TestCatalog(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	resource.LuaRegisterBuiltin(L)

	config := &Config{
		Module:   "",
		DryRun:   true,
		Logger:   log.New(os.Stdout, "", log.LstdFlags),
		SiteRepo: "",
		L:        L,
	}
	katalog := New(config)

	if len(katalog.Unsorted) != 0 {
		t.Errorf("want 0 resources, got %d\n", len(katalog.Unsorted))
	}

	code := `
	foo = resource.file.new("foo")
	bar = resource.file.new("bar")
	qux = resource.file.new("qux")
	catalog:add(foo, bar, qux)
	`

	if err := L.DoString(code); err != nil {
		t.Error(err)
	}

	if len(katalog.Unsorted) != 3 {
		t.Errorf("want 3 resources, got %d\n", len(katalog.Unsorted))
	}

	code = `
	if #catalog ~= 3 then
	   error("want 3 resources, got " .. #catalog)
	end
	`

	if err := L.DoString(code); err != nil {
		t.Error(err)
	}
}
