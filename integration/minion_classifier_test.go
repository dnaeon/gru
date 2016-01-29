package integration

import (
	"reflect"
	"sort"
	"testing"

	"github.com/dnaeon/gru/classifier"
	"github.com/dnaeon/gru/client"
	"github.com/dnaeon/gru/minion"

	"github.com/dnaeon/go-vcr/recorder"

	etcdclient "github.com/coreos/etcd/client"
)

func TestMinionClassifiers(t *testing.T) {
	// Start our recorder
	r, err := recorder.New("fixtures/minion-classifier")
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

	// Classifiers to test
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

	m := minion.NewEtcdMinion("Kevin", cfg)
	minionId := m.ID()

	// Set minion classifiers
	for _, tc := range testClassifiers {
		err := m.SetClassifier(tc)
		if err != nil {
			t.Error(err)
		}
		wantClassifierKeys = append(wantClassifierKeys, tc.Key)
	}
	sort.Strings(wantClassifierKeys)

	// Get classifiers keys from etcd
	klient := client.NewEtcdMinionClient(cfg)
	gotClassifierKeys, err := klient.MinionClassifierKeys(minionId)

	if err != nil {
		t.Fatal(err)
	}

	sort.Strings(gotClassifierKeys)

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
