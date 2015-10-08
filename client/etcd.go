package client

import (
	"log"
	"encoding/json"
	"path/filepath"

	"github.com/dnaeon/gru/minion"

	"code.google.com/p/go-uuid/uuid"
	"github.com/coreos/etcd/Godeps/_workspace/src/golang.org/x/net/context"
	etcdclient "github.com/coreos/etcd/client"
)

type EtcdClient struct {
	// KeysAPI client to etcd
	KAPI etcdclient.KeysAPI
}

func NewEtcdClient(cfg etcdclient.Config) Client {
	c, err := etcdclient.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	kapi := etcdclient.NewKeysAPI(c)
	klient := &EtcdClient{
		KAPI: kapi,
	}

	return klient
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
