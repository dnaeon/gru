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
	"reflect"
	"sort"
	"testing"

	"github.com/dnaeon/gru/minion"
	"github.com/pborman/uuid"
)

func TestMinionList(t *testing.T) {
	tc := mustNewTestClient("fixtures/minion-list")
	defer tc.recorder.Stop()

	minionNames := []string{
		"Bob", "Kevin", "Stuart",
	}

	wantMinions := []uuid.UUID{
		uuid.Parse("f827bffd-bd9e-5441-be36-a92a51d0b79e"), // Bob
		uuid.Parse("46ce0385-0e2b-5ede-8279-9cd98c268170"), // Kevin
		uuid.Parse("f87cf58e-1e19-57e1-bed3-9dff5064b86a"), // Stuart
	}

	// Convert minion uuids as strings for
	// sorting and equality testing
	var wantMinionsAsString []string
	for _, m := range wantMinions {
		wantMinionsAsString = append(wantMinionsAsString, m.String())
	}
	sort.Strings(wantMinionsAsString)

	// Register our minions in etcd
	for _, name := range minionNames {
		cfg := &minion.EtcdMinionConfig{
			Name:       name,
			EtcdConfig: tc.config,
		}
		m, err := minion.NewEtcdMinion(cfg)
		if err != nil {
			t.Fatal(err)
		}

		err = m.SetName(name)
		if err != nil {
			t.Error(err)
		}
	}

	// Get minions from etcd
	gotMinions, err := tc.client.MinionList()
	if err != nil {
		t.Fatal(err)
	}

	// Convert retrieved minion uuids as string for
	// sorting and equality testing
	var gotMinionsAsString []string
	for _, m := range gotMinions {
		gotMinionsAsString = append(gotMinionsAsString, m.String())
	}
	sort.Strings(gotMinionsAsString)

	if !reflect.DeepEqual(wantMinionsAsString, gotMinionsAsString) {
		t.Errorf("want %q minions, got %q minions", wantMinions, gotMinions)
	}
}
