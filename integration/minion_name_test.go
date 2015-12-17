package integration

import (
	"testing"

	"github.com/dnaeon/gru/minion"
	"github.com/dnaeon/gru/client"

	"github.com/dnaeon/go-vcr/recorder"

	etcdclient "github.com/coreos/etcd/client"
)

func TestMinionName(t *testing.T) {
	// Start our recorder
	r, err := recorder.New("fixtures/minion-name")
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

	wantName := "Kevin"
	m := minion.NewEtcdMinion(wantName, cfg)
	minionId := m.ID()
	err = m.SetName(wantName)
	if err != nil {
		t.Fatal(err)
	}

	klient := client.NewEtcdMinionClient(cfg)
	gotName, err := klient.MinionName(minionId)
	if err != nil {
		t.Fatal(err)
	}

	if wantName != gotName {
		t.Errorf("want %q, got %q", wantName, gotName)
	}
}
