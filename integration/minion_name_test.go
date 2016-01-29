package integration

import (
	"testing"

	"github.com/dnaeon/gru/minion"
)

func TestMinionName(t *testing.T) {
	tc := mustNewTestClient("fixtures/minion-name")
	defer tc.recorder.Stop()

	wantName := "Kevin"
	m := minion.NewEtcdMinion(wantName, tc.config)
	minionId := m.ID()
	err := m.SetName(wantName)
	if err != nil {
		t.Fatal(err)
	}

	gotName, err := tc.client.MinionName(minionId)
	if err != nil {
		t.Fatal(err)
	}

	if wantName != gotName {
		t.Errorf("want %q, got %q", wantName, gotName)
	}
}
