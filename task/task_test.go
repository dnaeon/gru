package task

import "testing"

func TestTaskState(t *testing.T) {
	dummyTask := New("foo", "bar")
	got := dummyTask.State
	want := TaskStateUnknown
	if want != got {
		t.Errorf("Incorrect task state: want %q, got %q", want, got)
	}
}

func TestTaskTimeReceivedProcessed(t *testing.T) {
	dummyTask := New("foo", "bar")

	// Task time received and processed should be 0 when initially created
	var want int64

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
	dummyTask := New("foo", "bar")
	got := dummyTask.Result
	want := ""

	if want != got {
		t.Errorf("Incorrect task result: want %q, got %q", want, got)
	}
}
