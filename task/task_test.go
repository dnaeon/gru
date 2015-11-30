package task

import "testing"

func TestTaskState(t *testing.T) {
	dummyTask := NewTask("dummy", "foo", "bar")
	got := dummyTask.State
	want := TaskStateUnknown
	if got != want {
		t.Errorf("Incorrect task state: want %q, got %q", want, got)
	}
}
