package integration

import (
	"testing"
	"time"

	"github.com/dnaeon/gru/minion"
	"github.com/dnaeon/gru/client"

	"github.com/dnaeon/go-vcr/recorder"

	etcdclient "github.com/coreos/etcd/client"
)

func TestMinionLastseen(t *testing.T) {
	// Start our recorder
	r, err := recorder.New("fixtures/minion-lastseen")
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

	m := minion.NewEtcdMinion("Kevin", cfg)
	id := m.ID()
	want := time.Now().Unix()
	err = m.SetLastseen(want)
	if err != nil {
		t.Fatal(err)
	}

	klient := client.NewEtcdMinionClient(cfg)
	got, err := klient.MinionLastseen(id)

	if want != got {
		t.Errorf("want %d lastseen, got %d lastseen", want, got)
	}
}
