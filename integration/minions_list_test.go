package integration

import (
	"reflect"
	"sort"
	"testing"

	"github.com/dnaeon/gru/minion"
	"github.com/dnaeon/gru/client"

	"github.com/dnaeon/go-vcr/recorder"

	"github.com/pborman/uuid"
	etcdclient "github.com/coreos/etcd/client"
)

func TestMinionList(t *testing.T) {
	// Start our recorder
	r, err := recorder.New("fixtures/minion-list")
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

	minionNames := []string{
		"Bob", "Kevin", "Stuart",
	}

	wantMinions := []uuid.UUID{
		uuid.Parse("f827bffd-bd9e-5441-be36-a92a51d0b79e"), // Bob
		uuid.Parse("46ce0385-0e2b-5ede-8279-9cd98c268170"), // Kevin
		uuid.Parse("f87cf58e-1e19-57e1-bed3-9dff5064b86a"), // Stuart
	}

	// Convert minion uuids as strings for
	// sorting and equality testing
	var wantMinionsAsString []string
	for _, m := range wantMinions {
		wantMinionsAsString = append(wantMinionsAsString, m.String())
	}
	sort.Strings(wantMinionsAsString)

	// Register our minions in etcd
	for _, name := range minionNames {
		m := minion.NewEtcdMinion(name, cfg)
		err := m.SetName(name)
		if err != nil {
			t.Error(err)
		}
	}

	// Get minions from etcd
	klient := client.NewEtcdMinionClient(cfg)
	gotMinions, err := klient.MinionList()
	if err != nil {
		t.Fatal(err)
	}

	// Convert retrieved minion uuids as string for
	// sorting and equality testing
	var gotMinionsAsString []string
	for _, m := range gotMinions {
		gotMinionsAsString = append(gotMinionsAsString, m.String())
	}
	sort.Strings(gotMinionsAsString)

	if !reflect.DeepEqual(wantMinionsAsString, gotMinionsAsString) {
		t.Errorf("want %q minions, got %q minions", wantMinions, gotMinions)
	}
}
