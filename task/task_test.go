package task

import (
	"reflect"
	"testing"
)

func TestTaskState(t *testing.T) {
	dummyTask := NewTask("dummy")
	got := dummyTask.State
	want := TaskStateUnknown
	if want != got {
		t.Errorf("Incorrect task state: want %q, got %q", want, got)
	}
}

func TestTaskCommand(t *testing.T) {
	dummyTask := NewTask("dummy")
	got := dummyTask.Command
	want := "dummy"

	if want != got {
		t.Errorf("Incorrect task command: want %q, got %q", want, got)
	}
}

func TestTaskWithArgs(t *testing.T) {
	dummyTask := NewTask("dummy", "foo", "bar")
	got := dummyTask.Args
	want := []string{"foo", "bar"}

	if !reflect.DeepEqual(want, got) {
		t.Errorf("Incorrect task args: want %q, got %q", want, got)
	}
}

func TestTaskWithoutArgs(t *testing.T) {
	dummyTask := NewTask("dummy")
	got := dummyTask.Args
	var want []string

	if got != nil {
		t.Errorf("Incorrect task args: want %q, got %q", want, got)
	}
}

func TestTaskTimeReceivedProcessed(t *testing.T) {
	dummyTask := NewTask("dummy", "foo", "bar")

	// Task time received and processed should be 0 when initially created
	var want int64 = 0

	got := dummyTask.TimeReceived
	if want != got {
		t.Errorf("Incorrect task time received: want %q, got %q", want, got)
	}

	got = dummyTask.TimeProcessed
	if want != got {
		t.Errorf("Incorrect task time processed: want %q, got %q", want, got)
	}
}

func TestTaskResult(t *testing.T) {
	dummyTask := NewTask("dummy")
	got := dummyTask.Result
	want := ""

	if want != got {
		t.Errorf("Incorrect task result: want %q, got %q", want, got)
	}
}

func TestTaskIsConcurrent(t *testing.T) {
	dummyTask := NewTask("dummy")
	got := dummyTask.IsConcurrent
	want := false

	if want != got {
		t.Errorf("Incorrect task concurrency: want %q, got %q", want, got)
	}
}

func TestTaskError(t *testing.T) {
	dummyTask := NewTask("dummy")
	got := dummyTask.Result
	want := ""

	if want != got {
		t.Errorf("Incorrect task error: want %q, got %q", want, got)
	}
}
