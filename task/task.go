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

package task

import "github.com/pborman/uuid"

// Task states
const (
	// Unknown state of the task
	// This is the default state of a task
	// when new task is initially created
	TaskStateUnknown = "unknown"

	// Task has been received by the
	// minion and is queued for execution
	TaskStateQueued = "queued"

	// Task is being processed
	TaskStateProcessing = "processing"

	// Task has been processed by the
	// minion and was flagged as successful
	TaskStateSuccess = "success"

	// Task has been processed by the
	// minion and was flagged as failed
	TaskStateFailed = "failed"

	// Task has been skipped
	TaskStateSkipped = "skipped"
)

// Task type represents a task that is processed by minions
type Task struct {
	// Do not take any actions, just report what would be done
	DryRun bool `json:"dryRun"`

	// Environment to use for this task
	Environment string `json:"environment"`

	// Command to be processed
	Command string `json:"command"`

	// Time when the command was sent for processing
	TimeReceived int64 `json:"timeReceived"`

	// Time when the command was processed
	TimeProcessed int64 `json:"timeProcessed"`

	// Task unique id
	ID uuid.UUID `json:"id"`

	// Result of task after processing
	Result string `json:"result"`

	// Task state
	State string `json:"state"`
}

// New creates a new task
func New(command, environment string) *Task {
	t := &Task{
		Command:     command,
		Environment: environment,
		ID:          uuid.NewRandom(),
		State:       TaskStateUnknown,
	}

	return t
}
