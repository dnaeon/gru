package client

import (
	"log"
	"strconv"
	"encoding/json"
	"path/filepath"

	"github.com/dnaeon/gru/minion"

	"code.google.com/p/go-uuid/uuid"
	"github.com/coreos/etcd/Godeps/_workspace/src/golang.org/x/net/context"
	etcdclient "github.com/coreos/etcd/client"
)

type EtcdMinionClient struct {
	// KeysAPI client to etcd
	KAPI etcdclient.KeysAPI
}

func NewEtcdMinionClient(cfg etcdclient.Config) MinionClient {
	c, err := etcdclient.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	kapi := etcdclient.NewKeysAPI(c)
	klient := &EtcdMinionClient{
		KAPI: kapi,
	}

	return klient
}

// Gets the name of the minion
func (c *EtcdMinionClient) GetName(u uuid.UUID) (string, error) {
	nameKey := filepath.Join(minion.EtcdMinionSpace, u.String(), "name")
	resp, err := c.KAPI.Get(context.Background(), nameKey, nil)

	if err != nil {
		return "", err
	}

	return resp.Node.Value, nil
}

// Gets the time the minion was last seen
func (c *EtcdMinionClient) GetLastseen(u uuid.UUID) (int64, error) {
	lastseenKey := filepath.Join(minion.EtcdMinionSpace, u.String(), "lastseen")
	resp, err := c.KAPI.Get(context.Background(), lastseenKey, nil)

	if err != nil {
		return 0, err
	}

	lastseen, err := strconv.ParseInt(resp.Node.Value, 10, 64)

	if err != nil {
		return 0, err
	}

	return lastseen, nil
}

// Gets a classifier identified with key
func (c *EtcdMinionClient) GetClassifier(u uuid.UUID, key string) (minion.MinionClassifier, error) {
	// Classifier key in etcd
	classifierKey := filepath.Join(minion.EtcdMinionSpace, u.String(), "classifier", key)
	resp, err := c.KAPI.Get(context.Background(), classifierKey, nil)

	if err != nil {
		return nil, err
	}

	klassifier := new(minion.SimpleClassifier)
	err = json.Unmarshal([]byte(resp.Node.Value), &klassifier)

	return klassifier, err
}

// Submits a task to a minion
func (c *EtcdMinionClient) SubmitTask(u uuid.UUID, t minion.MinionTask) error {
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
