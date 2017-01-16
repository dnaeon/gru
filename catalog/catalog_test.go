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
