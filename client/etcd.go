package client

import (
	"log"
	"path"
	"sync"
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

// Gets all classifiers for a minion
func (c *EtcdMinionClient) GetAllClassifiers(u uuid.UUID) ([]minion.MinionClassifier, error) {
	var classifiers []minion.MinionClassifier

	// Classifier directory key in etcd
	classifierDirKey := filepath.Join(minion.EtcdMinionSpace, u.String(), "classifier")
	opts := &etcdclient.GetOptions{
		Recursive: true,
	}

	resp, err := c.KAPI.Get(context.Background(), classifierDirKey, opts)
	if err != nil {
		return classifiers, err
	}

	for _, node := range resp.Node.Nodes {
		klassifier := new(minion.SimpleClassifier)
		err := json.Unmarshal([]byte(node.Value), &klassifier)
		if err != nil {
			return classifiers, err
		}

		classifiers = append(classifiers, klassifier)
	}

	return classifiers, nil
}

// Gets all minions which are classified with a given key
// The keys of the result map are the minion uuids
// represented as a string
func (c *EtcdMinionClient) GetClassifiedMinions(key string) (map[string]minion.MinionClassifier, error) {
	// Searching for classified minions is performed in a
	// concurrent way, so we need to make sure any writes to
	// the result map are synchronized between each goroutine
	var wg sync.WaitGroup
	type concurrentMap struct {
		sync.RWMutex
		m map[string]minion.MinionClassifier
	}

	// A channel to which we send minion uuids to be
	// checked whether or not they have the given classifier
	tasks := make(chan uuid.UUID, 1024)

	// Concurrent map which we use to store the minions
	// which have the given clasifier key
	minionsMap := concurrentMap{m: make(map[string]minion.MinionClassifier)}

	// Get the minions from etcd
	resp, err := c.KAPI.Get(context.Background(), minion.EtcdMinionSpace, nil)
	if err != nil {
		return minionsMap.m, err
	}

	// Start four worker goroutines that will be
	// processing the uuids from the tasks channel
	for i := 0; i < 4; i++ {
		wg.Add(1)
		worker := func() {
			defer wg.Done()
			for u := range tasks {
				klassifier, err := c.GetClassifier(u, key)
				if err != nil {
					continue
				}

				minionsMap.Lock()
				minionsMap.m[u.String()] = klassifier
				minionsMap.Unlock()
			}
		}

		go worker()
	}

	// Send the minion uuids to our workers for processing
	for _, node := range resp.Node.Nodes {
		u := path.Base(node.Key)
		minionUUID := uuid.Parse(u)
		if minionUUID == nil {
			log.Printf("Bad minion uuid found: %s\n", u)
			continue
		}

		tasks <- minionUUID
	}
	close(tasks)

	wg.Wait()

	return minionsMap.m, nil
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
