package resource

import "testing"

func TestParseSysRCOutput(t *testing.T) {
	table := []struct {
		Input string
		K     string
		V     string
	}{
		{
			Input: "keyrate: fast\n",
			K:     "keyrate",
			V:     "fast",
		},
		{
			Input: "dumpdev: \n",
			K:     "dumpdev",
			V:     "",
		},
	}

	for _, item := range table {
		t.Run(item.Input, func(t *testing.T) {
			k, v, err := parseSysRCOutput(item.Input)
			if err != nil {
				t.Error(err)
			}
			if k != item.K || v != item.V {
				t.Errorf("expected: k=%q, v=%q, got: k=%q, v=%q", item.K, item.V, k, v)
			}
		})
	}
}
