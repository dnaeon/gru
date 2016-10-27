package integration

import (
	"github.com/dnaeon/go-vcr/recorder"
	"github.com/dnaeon/gru/client"

	etcdclient "github.com/coreos/etcd/client"
)

// Minion client used during integration testing
type testClient struct {
	client   client.Client
	config   etcdclient.Config
	recorder *recorder.Recorder
}

// Creates a new etcd client with recording enabled
func mustNewTestClient(cassette string) *testClient {
	// Start our recorder
	r, err := recorder.New(cassette)
	if err != nil {
		panic(err)
	}

	cfg := etcdclient.Config{
		Endpoints:               []string{"http://127.0.0.1:2379"},
		Transport:               r, // Inject our transport!
		HeaderTimeoutPerRequest: etcdclient.DefaultRequestTimeout,
	}

	klient := client.NewEtcdMinionClient(cfg)

	tc := &testClient{
		client:   klient,
		config:   cfg,
		recorder: r,
	}

	return tc
}
