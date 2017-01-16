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
	"os"
	"testing"
)

func TestFile(t *testing.T) {
	L := newLuaState()
	defer L.Close()

	const code = `
	foo = resource.file.new("/tmp/foo")
	`

	if err := L.DoString(code); err != nil {
		t.Fatal(err)
	}

	foo := luaResource(L, "foo").(*File)
	errorIfNotEqual(t, "file", foo.Type)
	errorIfNotEqual(t, "/tmp/foo", foo.Name)
	errorIfNotEqual(t, "present", foo.State)
	errorIfNotEqual(t, []string{}, foo.Require)
	errorIfNotEqual(t, []string{"present"}, foo.PresentStatesList)
	errorIfNotEqual(t, []string{"absent"}, foo.AbsentStatesList)
	errorIfNotEqual(t, true, foo.Concurrent)
	errorIfNotEqual(t, "/tmp/foo", foo.Path)
	errorIfNotEqual(t, os.FileMode(0644), foo.Mode)
	errorIfNotEqual(t, "", foo.Source)
}

func TestDirectory(t *testing.T) {
	L := newLuaState()
	defer L.Close()

	const code = `
	bar = resource.directory.new("/tmp/bar")
	`

	if err := L.DoString(code); err != nil {
		t.Fatal(err)
	}

	bar := luaResource(L, "bar").(*Directory)
	errorIfNotEqual(t, "directory", bar.Type)
	errorIfNotEqual(t, "/tmp/bar", bar.Name)
	errorIfNotEqual(t, "present", bar.State)
	errorIfNotEqual(t, []string{}, bar.Require)
	errorIfNotEqual(t, []string{"present"}, bar.PresentStatesList)
	errorIfNotEqual(t, []string{"absent"}, bar.AbsentStatesList)
	errorIfNotEqual(t, true, bar.Concurrent)
	errorIfNotEqual(t, "/tmp/bar", bar.Path)
	errorIfNotEqual(t, os.FileMode(0755), bar.Mode)
	errorIfNotEqual(t, false, bar.Parents)
}

func TestLink(t *testing.T) {
	L := newLuaState()
	defer L.Close()

	const code = `
	qux = resource.link.new("/tmp/qux")
	qux.source = "/tmp/foo"
	`

	if err := L.DoString(code); err != nil {
		t.Fatal(err)
	}

	qux := luaResource(L, "qux").(*Link)
	errorIfNotEqual(t, "link", qux.Type)
	errorIfNotEqual(t, "/tmp/qux", qux.Name)
	errorIfNotEqual(t, "present", qux.State)
	errorIfNotEqual(t, []string{}, qux.Require)
	errorIfNotEqual(t, []string{"present"}, qux.PresentStatesList)
	errorIfNotEqual(t, []string{"absent"}, qux.AbsentStatesList)
	errorIfNotEqual(t, true, qux.Concurrent)
	errorIfNotEqual(t, "/tmp/foo", qux.Source)
	errorIfNotEqual(t, false, qux.Hard)
}
