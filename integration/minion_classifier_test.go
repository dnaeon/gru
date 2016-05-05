package integration

import (
	"reflect"
	"sort"
	"testing"

	"github.com/dnaeon/gru/classifier"
	"github.com/dnaeon/gru/minion"
)

func TestMinionClassifiers(t *testing.T) {
	tc := mustNewTestClient("fixtures/minion-classifier")
	defer tc.recorder.Stop()

	// Classifiers to test
	var wantClassifierKeys []string
	testClassifiers := []*classifier.Classifier{
		{
			Key:   "foo",
			Value: "bar",
		},
		{
			Key:   "baz",
			Value: "qux",
		},
	}

	cfg := &minion.EtcdMinionConfig{
		Name:       "Kevin",
		EtcdConfig: tc.config,
	}
	m, err := minion.NewEtcdMinion(cfg)
	if err != nil {
		t.Fatal(err)
	}

	minionID := m.ID()

	// Set minion classifiers
	for _, c := range testClassifiers {
		err := m.SetClassifier(c)
		if err != nil {
			t.Error(err)
		}
		wantClassifierKeys = append(wantClassifierKeys, c.Key)
	}
	sort.Strings(wantClassifierKeys)

	// Get classifiers keys from etcd
	gotClassifierKeys, err := tc.client.MinionClassifierKeys(minionID)

	if err != nil {
		t.Fatal(err)
	}

	sort.Strings(gotClassifierKeys)
	if !reflect.DeepEqual(wantClassifierKeys, gotClassifierKeys) {
		t.Errorf("want %q classifier keys, got %q classifier keys", wantClassifierKeys, gotClassifierKeys)
	}

	// Get classifier values
	for _, c := range testClassifiers {
		klassifier, err := tc.client.MinionClassifier(minionID, c.Key)
		if err != nil {
			t.Fatal(err)
		}

		if c.Value != klassifier.Value {
			t.Errorf("want %q classifier value, got %q classifier value", c.Value, klassifier.Value)
		}
	}

	// Get minions which contain given classifier key
	for _, c := range testClassifiers {
		minions, err := tc.client.MinionWithClassifierKey(c.Key)
		if err != nil {
			t.Fatal(err)
		}

		// We expect a single minion with the test classifier keys
		if len(minions) != 1 {
			t.Errorf("want 1 minion, got %d minion(s)", len(minions))
		}
	}
}
