package client

import (
	"log"
	"path"
	"sync"
	"strconv"
	"encoding/json"
	"path/filepath"

	"github.com/dnaeon/gru/minion"
	"github.com/dnaeon/gru/utils"

	"code.google.com/p/go-uuid/uuid"
	"github.com/coreos/etcd/Godeps/_workspace/src/golang.org/x/net/context"
	etcdclient "github.com/coreos/etcd/client"
)

// Max number of concurrent requests to be
// performed against the etcd cluster at a time
const maxGoroutines = 4

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
	cm := utils.NewConcurrentMap()

	// We wait until all goroutines are complete
	// before returning the result to the client
	var wg sync.WaitGroup

	// A channel to which we send minion uuids to be
	// checked whether or not they have the given classifier
	minionChannel := make(chan uuid.UUID, 1024)

	// Get the minions from etcd
	resp, err := c.KAPI.Get(context.Background(), minion.EtcdMinionSpace, nil)
	if err != nil {
		return nil, err
	}

	// Producer sending uuids for processing over the channel
	producer := func() {
		for _, node := range resp.Node.Nodes {
			u := path.Base(node.Key)
			minionUUID := uuid.Parse(u)
			if minionUUID == nil {
				log.Printf("Bad minion uuid found: %s\n", u)
				continue
			}
			minionChannel <- minionUUID
		}

		close(minionChannel)
	}
	go producer()

	// Start our worker goroutines that will be
	// processing the minion uuids for the given classifiers
	for i := 0; i < maxGoroutines; i++ {
		wg.Add(1)
		worker := func() {
			defer wg.Done()
			for id := range minionChannel {
				klassifier, err := c.GetClassifier(id, key)
				if err != nil {
					continue
				}

				// Set uuid and classifier in the concurrent map
				cm.Set(id.String(), klassifier)
			}
		}
		go worker()
	}

	wg.Wait()

	// The result map should be of type[string]minion.MinionClassifier, so
	// perform any type assertions here
	result := make(map[string]minion.MinionClassifier)
	for item := range cm.Iter() {
		result[item.Key] = item.Value.(minion.MinionClassifier)
	}

	return result, nil
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
