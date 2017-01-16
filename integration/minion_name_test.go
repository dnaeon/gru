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

package integration

import (
	"testing"

	"github.com/dnaeon/gru/minion"
)

func TestMinionName(t *testing.T) {
	tc := mustNewTestClient("fixtures/minion-name")
	defer tc.recorder.Stop()

	wantName := "Kevin"
	cfg := &minion.EtcdMinionConfig{
		Name:       wantName,
		EtcdConfig: tc.config,
	}
	m, err := minion.NewEtcdMinion(cfg)
	if err != nil {
		t.Fatal(err)
	}

	minionID := m.ID()
	err = m.SetName(wantName)
	if err != nil {
		t.Fatal(err)
	}

	gotName, err := tc.client.MinionName(minionID)
	if err != nil {
		t.Fatal(err)
	}

	if wantName != gotName {
		t.Errorf("want %q, got %q", wantName, gotName)
	}
}
