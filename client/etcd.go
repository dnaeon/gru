package client

import (
	"encoding/json"
	"path/filepath"

	"github.com/dnaeon/gru/minion"

	"code.google.com/p/go-uuid/uuid"
	"github.com/coreos/etcd/Godeps/_workspace/src/golang.org/x/net/context"
	etcdclient "github.com/coreos/etcd/client"
)

type EtcdClient struct {
	// etcd keys api client
	KAPI etcdclient.KeysAPI
}

func NewEtcdClient(kapi etcdclient.KeysAPI) *EtcdClient {
	c := &EtcdClient{
		KAPI: kapi,
	}

	return c
}

func (c *EtcdClient) SubmitTask(u uuid.UUID, t minion.MinionTask) error {
	minionRootDir := filepath.Join(minion.EtcdMinionSpace, u.String())
	queueDir := filepath.Join(minionRootDir, "queue")

	_, err := c.KAPI.Get(context.Background(), minionRootDir, nil)
	if err != nil {
		return err
	}

	data, err := json.Marshal(t)
	if err != nil {
		return err
	}

	_, err = c.KAPI.CreateInOrder(context.Background(), queueDir, string(data), nil)
	if err != nil {
		return err
	}

	return nil
}
