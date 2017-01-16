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
	"testing"

	"github.com/dnaeon/gru/minion"
	"github.com/dnaeon/gru/task"

	"github.com/pborman/uuid"
)

func TestMinionTaskBacklog(t *testing.T) {
	tc := mustNewTestClient("fixtures/minion-task-backlog")
	defer tc.recorder.Stop()

	// Setup our minion
	minionName := "Kevin"
	cfg := &minion.EtcdMinionConfig{
		Name:       minionName,
		EtcdConfig: tc.config,
	}

	testMinion, err := minion.NewEtcdMinion(cfg)
	if err != nil {
		t.Fatal(err)
	}

	minionID := testMinion.ID()

	err = testMinion.SetName(minionName)
	if err != nil {
		t.Fatal(err)
	}

	// Create a test task and submit it
	wantTask := task.New("foo", "bar")
	wantTask.ID = uuid.Parse("e6d2bebd-2219-4a8c-9d30-a861097c147e")

	err = tc.client.MinionSubmitTask(minionID, wantTask)
	if err != nil {
		t.Fatal(err)
	}

	// Get pending tasks and verify the task we sent is the task we get
	backlog, err := tc.client.MinionTaskQueue(minionID)
	if err != nil {
		t.Fatal(err)
	}

	if len(backlog) != 1 {
		t.Errorf("want 1 backlog task, got %d backlog tasks", len(backlog))
	}

	gotTask := backlog[0]
	if !reflect.DeepEqual(wantTask, gotTask) {
		t.Errorf("want %q task, got %q task", wantTask.ID, gotTask.ID)
	}
}
