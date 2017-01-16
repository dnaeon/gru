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

package minion

import (
	"github.com/dnaeon/gru/classifier"
	"github.com/dnaeon/gru/task"

	"github.com/pborman/uuid"
)

// Minion interface type
type Minion interface {
	// ID returns the unique identifier of a minion
	ID() uuid.UUID

	// SetName sets the name of the minion
	SetName(string) error

	// SetLastseen sets the time the minion was last seen
	SetLastseen(int64) error

	// SetClassifier sets a classifier for the minion
	SetClassifier(*classifier.Classifier) error

	// TaskListener listens for new tasks and processes them
	TaskListener(c chan<- *task.Task) error

	// TaskRunner runs new tasks as received by the TaskListener
	TaskRunner(c <-chan *task.Task) error

	// SaveTaskResult saves the result of a task
	SaveTaskResult(t *task.Task) error

	// Sync syncs modules and data files
	Sync() error

	// Serve start the minion
	Serve() error

	// Stop stops the minion
	Stop() error
}
