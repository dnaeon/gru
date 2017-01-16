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

package client

import (
	"github.com/dnaeon/gru/classifier"
	"github.com/dnaeon/gru/task"

	"github.com/pborman/uuid"
)

// Client interface for interacting with minions
type Client interface {
	// Gets all registered minions
	MinionList() ([]uuid.UUID, error)

	// Gets the name of a minion
	MinionName(m uuid.UUID) (string, error)

	// Gets the time a minion was last seen
	MinionLastseen(m uuid.UUID) (int64, error)

	// Gets a classifier of a minion
	MinionClassifier(m uuid.UUID, key string) (*classifier.Classifier, error)

	// Gets all classifier keys of a minion
	MinionClassifierKeys(m uuid.UUID) ([]string, error)

	// Gets minions which are classified with a given classifier key
	MinionWithClassifierKey(key string) ([]uuid.UUID, error)

	// Gets the result of a task for a minion
	MinionTaskResult(m uuid.UUID, t uuid.UUID) (*task.Task, error)

	// Gets the minions which have a task result with the given uuid
	MinionWithTaskResult(t uuid.UUID) ([]uuid.UUID, error)

	// Gets the tasks which are currently pending in the queue
	MinionTaskQueue(m uuid.UUID) ([]*task.Task, error)

	// Gets the uuids of tasks which have already been processed
	MinionTaskLog(m uuid.UUID) ([]uuid.UUID, error)

	// Submits a new task to a minion
	MinionSubmitTask(m uuid.UUID, t *task.Task) error
}
