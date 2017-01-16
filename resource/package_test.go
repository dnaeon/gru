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

import "testing"

func TestPacman(t *testing.T) {
	L := newLuaState()
	defer L.Close()

	const code = `
	tmux = resource.pacman.new("tmux")
	`

	if err := L.DoString(code); err != nil {
		t.Fatal(err)
	}

	pkg := luaResource(L, "tmux").(*Pacman)
	errorIfNotEqual(t, "package", pkg.Type)
	errorIfNotEqual(t, "tmux", pkg.Name)
	errorIfNotEqual(t, "installed", pkg.State)
	errorIfNotEqual(t, []string{}, pkg.Require)
	errorIfNotEqual(t, []string{"present", "installed"}, pkg.PresentStatesList)
	errorIfNotEqual(t, []string{"absent", "deinstalled"}, pkg.AbsentStatesList)
	errorIfNotEqual(t, false, pkg.Concurrent)
	errorIfNotEqual(t, "tmux", pkg.Package)
	errorIfNotEqual(t, "", pkg.Version)
}

func TestYum(t *testing.T) {
	L := newLuaState()
	defer L.Close()

	const code = `
	tmux = resource.yum.new("tmux")
	`

	if err := L.DoString(code); err != nil {
		t.Fatal(err)
	}

	pkg := luaResource(L, "tmux").(*Yum)
	errorIfNotEqual(t, "package", pkg.Type)
	errorIfNotEqual(t, "tmux", pkg.Name)
	errorIfNotEqual(t, "installed", pkg.State)
	errorIfNotEqual(t, []string{}, pkg.Require)
	errorIfNotEqual(t, []string{"present", "installed"}, pkg.PresentStatesList)
	errorIfNotEqual(t, []string{"absent", "deinstalled"}, pkg.AbsentStatesList)
	errorIfNotEqual(t, false, pkg.Concurrent)
	errorIfNotEqual(t, "tmux", pkg.Package)
	errorIfNotEqual(t, "", pkg.Version)
}
