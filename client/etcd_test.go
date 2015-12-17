// +build integration

package client

import (
	"reflect"
	"sort"
	"time"
	"testing"

	"github.com/dnaeon/gru/minion"
	"github.com/dnaeon/gru/classifier"

	"github.com/dnaeon/go-vcr/recorder"

	"github.com/pborman/uuid"
	"golang.org/x/net/context"
	etcdclient "github.com/coreos/etcd/client"
)

// Default config for etcd minions and clients
var defaultEtcdConfig = etcdclient.Config{
	Endpoints:               []string{"http://127.0.0.1:2379"},
	Transport:               etcdclient.DefaultTransport,
	HeaderTimeoutPerRequest: time.Second,
}

// Cleans up the minion space in etcd after tests
func cleanupAfterTest(t *testing.T) {
	c, err := etcdclient.New(defaultEtcdConfig)
	if err != nil {
		t.Fatal(err)
	}

	kapi := etcdclient.NewKeysAPI(c)
	deleteOpts := &etcdclient.DeleteOptions{
		Recursive: true,
		Dir: true,
	}

	_, err = kapi.Delete(context.Background(), minion.EtcdMinionSpace, deleteOpts)
	if err != nil {
		t.Fatal(err)
	}
}

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
	klient := NewEtcdMinionClient(cfg)
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

func TestMinionName(t *testing.T) {
	defer cleanupAfterTest(t)	

	wantName := "Kevin"
	m := minion.NewEtcdMinion(wantName, defaultEtcdConfig)
	minionId := m.ID()
	err := m.SetName(wantName)
	if err != nil {
		t.Fatal(err)
	}

	klient := NewEtcdMinionClient(defaultEtcdConfig)
	gotName, err := klient.MinionName(minionId)
	if err != nil {
		t.Fatal(err)
	}

	if wantName != gotName {
		t.Errorf("want %q, got %q", wantName, gotName)
	}
}

func TestMinionLastseen(t *testing.T) {
	defer cleanupAfterTest(t)

	m := minion.NewEtcdMinion("Kevin", defaultEtcdConfig)
	id := m.ID()
	want := time.Now().Unix()
	err := m.SetLastseen(want)
	if err != nil {
		t.Fatal(err)
	}

	klient := NewEtcdMinionClient(defaultEtcdConfig)
	got, err := klient.MinionLastseen(id)

	if want != got {
		t.Errorf("want %d lastseen, got %d lastseen", want, got)
	}
}

func TestMinionClassifiers(t *testing.T) {
	defer cleanupAfterTest(t)

	wantClassifierKeys := make([]string, 0)
	testClassifiers := []*classifier.Classifier{
		&classifier.Classifier{
			Key:   "foo",
			Value: "bar",
		},
		&classifier.Classifier{
			Key:   "baz",
			Value: "qux",
		},
	}

	m := minion.NewEtcdMinion("Kevin", defaultEtcdConfig)
	minionId := m.ID()

	// Set minion classifiers
	for _, tc := range testClassifiers {
		err := m.SetClassifier(tc)
		if err != nil {
			t.Error(err)
		}
		wantClassifierKeys = append(wantClassifierKeys, tc.Key)
	}

	// Get classifiers keys
	klient := NewEtcdMinionClient(defaultEtcdConfig)
	gotClassifierKeys, err := klient.MinionClassifierKeys(minionId)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(wantClassifierKeys, gotClassifierKeys) {
		t.Errorf("want %q classifier keys, got %q classifier keys", wantClassifierKeys, gotClassifierKeys)
	}

	// Get classifier values
	for _, tc := range testClassifiers {
		klassifier, err := klient.MinionClassifier(minionId, tc.Key)
		if err != nil {
			t.Fatal(err)
		}

		if tc.Value != klassifier.Value {
			t.Errorf("want %q classifier value, got %q classifier value", tc.Value, klassifier.Value)
		}
	}

	// Get minions which contain given classifier key
	for _, tc := range testClassifiers {
		minions, err := klient.MinionWithClassifierKey(tc.Key)
		if err != nil {
			t.Fatal(err)
		}

		// We expect a single minion with the test classifier keys
		if len(minions) != 1 {
			t.Errorf("want 1 minion, got %d minion(s)", len(minions))
		}
	}
}
