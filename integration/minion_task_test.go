package integration

import (
	"reflect"
	"testing"

	"github.com/dnaeon/gru/client"
	"github.com/dnaeon/gru/minion"
	"github.com/dnaeon/gru/task"

	"github.com/dnaeon/go-vcr/recorder"

	etcdclient "github.com/coreos/etcd/client"
	"github.com/pborman/uuid"
)

func TestMinionTaskBacklog(t *testing.T) {
	// Start our recorder
	r, err := recorder.New("fixtures/minion-task-backlog")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()

	// Etcd config using our transport
	cfg := etcdclient.Config{
		Endpoints:               []string{"http://127.0.0.1:2379"},
		Transport:               r.Transport, // Inject our transport!
		HeaderTimeoutPerRequest: etcdclient.DefaultRequestTimeout,
	}

	klient := client.NewEtcdMinionClient(cfg)

	// Setup our minion
	minionName := "Kevin"
	testMinion := minion.NewEtcdMinion(minionName, cfg)
	minionId := testMinion.ID()

	err = testMinion.SetName(minionName)
	if err != nil {
		t.Fatal(err)
	}

	// Create a test task and submit it
	wantTask := task.NewTask("foo", "bar")
	wantTask.TaskID = uuid.Parse("e6d2bebd-2219-4a8c-9d30-a861097c147e")

	err = klient.MinionSubmitTask(minionId, wantTask)
	if err != nil {
		t.Fatal(err)
	}

	// Get pending tasks and verify the task we sent is the task we get
	backlog, err := klient.MinionTaskQueue(minionId)
	if err != nil {
		t.Fatal(err)
	}

	if len(backlog) != 1 {
		t.Errorf("want 1 backlog task, got %d backlog tasks", len(backlog))
	}

	gotTask := backlog[0]
	if !reflect.DeepEqual(wantTask, gotTask) {
		t.Errorf("want %q task, got %q task", wantTask.TaskID, gotTask.TaskID)
	}
}
