package task

import "testing"

func TestTaskState(t *testing.T) {
	dummyTask := NewTask("dummy", "foo", "bar")
	got := dummyTask.State
	want := TaskStateUnknown
	if want != got {
		t.Errorf("Incorrect task state: want %q, got %q", want, got)
	}
}

func TestTaskCommand(t *testing.T) {
	dummyTask := NewTask("dummy", "foo", "bar")
	got := dummyTask.Command
	want := "dummy"
	if want != got {
		t.Errorf("Incorrect task command: want %q, got %q", want, got)
	}
}

