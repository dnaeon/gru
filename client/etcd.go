package client

import (
	"context"
	"encoding/json"
	"log"
	"path"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/dnaeon/gru/classifier"
	"github.com/dnaeon/gru/minion"
	"github.com/dnaeon/gru/task"
	"github.com/dnaeon/gru/utils"

	etcdclient "github.com/coreos/etcd/client"
	"github.com/pborman/uuid"
)

// Max number of concurrent requests to be
// made at a time to the etcd cluster
const maxGoroutines = 4

type etcdMinionClient struct {
	// KeysAPI client to etcd
	kapi etcdclient.KeysAPI
}

// NewEtcdMinionClient creates a new client for interacting with
// minions using etcd as their interface implementation
func NewEtcdMinionClient(cfg etcdclient.Config) Client {
	c, err := etcdclient.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	kapi := etcdclient.NewKeysAPI(c)
	klient := &etcdMinionClient{
		kapi: kapi,
	}

	return klient
}

// MinionList returns a slice containing all registered minions
func (c *etcdMinionClient) MinionList() ([]uuid.UUID, error) {
	resp, err := c.kapi.Get(context.Background(), minion.EtcdMinionSpace, nil)

	if err != nil {
		return nil, err
	}

	var result []uuid.UUID
	for _, node := range resp.Node.Nodes {
		k := path.Base(node.Key)
		u := uuid.Parse(k)
		if u == nil {
			log.Printf("Bad minion uuid found: %s\n", k)
			continue
		}
		result = append(result, u)
	}

	return result, nil
}

// MinionName returns the name of a minion identified by its uuid
func (c *etcdMinionClient) MinionName(m uuid.UUID) (string, error) {
	nameKey := filepath.Join(minion.EtcdMinionSpace, m.String(), "name")
	resp, err := c.kapi.Get(context.Background(), nameKey, nil)

	if err != nil {
		return "", err
	}

	return resp.Node.Value, nil
}

// MinionLastSeen retrieves the time a minion was last seen
func (c *etcdMinionClient) MinionLastseen(m uuid.UUID) (int64, error) {
	lastseenKey := filepath.Join(minion.EtcdMinionSpace, m.String(), "lastseen")
	resp, err := c.kapi.Get(context.Background(), lastseenKey, nil)

	if err != nil {
		return 0, err
	}

	lastseen, err := strconv.ParseInt(resp.Node.Value, 10, 64)

	return lastseen, err
}

// MinionClassifier retrieves a classifier with the given key
func (c *etcdMinionClient) MinionClassifier(m uuid.UUID, key string) (*classifier.Classifier, error) {
	// Classifier key in etcd
	classifierKey := filepath.Join(minion.EtcdMinionSpace, m.String(), "classifier", key)
	resp, err := c.kapi.Get(context.Background(), classifierKey, nil)

	if err != nil {
		return nil, err
	}

	klassifier := new(classifier.Classifier)
	err = json.Unmarshal([]byte(resp.Node.Value), &klassifier)

	return klassifier, err
}

// MinionClassifierKeys returns all classifier keys that a minion has
func (c *etcdMinionClient) MinionClassifierKeys(m uuid.UUID) ([]string, error) {
	// Classifier directory in etcd
	classifierDir := filepath.Join(minion.EtcdMinionSpace, m.String(), "classifier")
	opts := &etcdclient.GetOptions{
		Recursive: true,
	}

	resp, err := c.kapi.Get(context.Background(), classifierDir, opts)
	if err != nil {
		return nil, err
	}

	var classifierKeys []string
	for _, node := range resp.Node.Nodes {
		klassifier := new(classifier.Classifier)
		err := json.Unmarshal([]byte(node.Value), &klassifier)
		if err != nil {
			return nil, err
		}

		classifierKeys = append(classifierKeys, klassifier.Key)
	}

	return classifierKeys, nil
}

// MinionWithClassifierKey returns all minions which are classified
// with a given classifier key
func (c *etcdMinionClient) MinionWithClassifierKey(key string) ([]uuid.UUID, error) {
	// Concurrent slice to hold the result
	cs := utils.NewConcurrentSlice()

	// We wait until all goroutines are complete
	// before returning the result to the client
	var wg sync.WaitGroup

	// A channel to which we send minion uuids to be
	// checked whether or not they have the given classifier
	queue := make(chan uuid.UUID, 1024)

	// Get the minions from etcd
	resp, err := c.kapi.Get(context.Background(), minion.EtcdMinionSpace, nil)
	if err != nil {
		return nil, err
	}

	// Producer sending uuids for processing over the channel
	producer := func() {
		for _, node := range resp.Node.Nodes {
			k := path.Base(node.Key)
			u := uuid.Parse(k)
			if u == nil {
				log.Printf("Bad minion uuid found: %s\n", k)
				continue
			}
			queue <- u
		}

		close(queue)
	}
	go producer()

	// Start our worker goroutines that will be
	// processing the minion uuids for the given classifiers
	for i := 0; i < maxGoroutines; i++ {
		wg.Add(1)
		worker := func() {
			defer wg.Done()
			for minionUUID := range queue {
				_, err := c.MinionClassifier(minionUUID, key)
				if err != nil {
					continue
				}

				cs.Append(minionUUID)
			}
		}
		go worker()
	}

	wg.Wait()

	// The result slice should be of []uuid.UUID, so
	// perform any type assertions here
	var result []uuid.UUID
	for item := range cs.Iter() {
		result = append(result, item.Value.(uuid.UUID))
	}

	return result, nil
}

// MinionTaskResult retrieves the result of a task for a given minion
func (c *etcdMinionClient) MinionTaskResult(m uuid.UUID, t uuid.UUID) (*task.Task, error) {
	// Task key in etcd
	taskKey := filepath.Join(minion.EtcdMinionSpace, m.String(), "log", t.String())

	// Get the task from etcd
	resp, err := c.kapi.Get(context.Background(), taskKey, nil)
	if err != nil {
		return nil, err
	}

	result := new(task.Task)
	err = json.Unmarshal([]byte(resp.Node.Value), &result)

	return result, err
}

// MinionWithTaskResult returns all minions which have a task result
// with the given task uuid
func (c *etcdMinionClient) MinionWithTaskResult(t uuid.UUID) ([]uuid.UUID, error) {
	// Concurrent slice to hold the result
	cs := utils.NewConcurrentSlice()

	// We wait until all goroutines are complete
	// before returning the result to the client
	var wg sync.WaitGroup

	// A channel to which we send minion uuids to be
	// checked whether or not they have the given task uuid
	queue := make(chan uuid.UUID, 1024)

	// Get the minions from etcd
	resp, err := c.kapi.Get(context.Background(), minion.EtcdMinionSpace, nil)
	if err != nil {
		return nil, err
	}

	// Producer sending uuids for processing over the channel
	producer := func() {
		for _, node := range resp.Node.Nodes {
			k := path.Base(node.Key)
			u := uuid.Parse(k)
			if u == nil {
				log.Printf("Bad minion uuid found: %s\n", k)
				continue
			}
			queue <- u
		}

		close(queue)
	}
	go producer()

	// Start our worker goroutines that will be
	// processing the minion uuids for the given task uuid
	for i := 0; i < maxGoroutines; i++ {
		wg.Add(1)
		worker := func() {
			defer wg.Done()
			for minionUUID := range queue {
				_, err := c.MinionTaskResult(minionUUID, t)
				if err != nil {
					continue
				}

				cs.Append(minionUUID)
			}
		}
		go worker()
	}

	wg.Wait()

	// The result slice should be of []uuid.UUID, so
	// perform any type assertions here
	var result []uuid.UUID
	for item := range cs.Iter() {
		result = append(result, item.Value.(uuid.UUID))
	}

	return result, nil
}

// MinionTaskQueue returns the tasks which are currently pending in the queue
func (c *etcdMinionClient) MinionTaskQueue(m uuid.UUID) ([]*task.Task, error) {
	queueDir := filepath.Join(minion.EtcdMinionSpace, m.String(), "queue")
	opts := &etcdclient.GetOptions{
		Recursive: true,
	}

	resp, err := c.kapi.Get(context.Background(), queueDir, opts)
	if err != nil {
		return nil, err
	}

	var tasks []*task.Task
	for _, node := range resp.Node.Nodes {
		t, err := minion.EtcdUnmarshalTask(node)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, t)
	}

	return tasks, nil
}

// MinionTaskLog returns the uuids of tasks which have already been
// processed by a minion
func (c *etcdMinionClient) MinionTaskLog(m uuid.UUID) ([]uuid.UUID, error) {
	logDir := filepath.Join(minion.EtcdMinionSpace, m.String(), "log")
	opts := &etcdclient.GetOptions{
		Recursive: true,
	}

	resp, err := c.kapi.Get(context.Background(), logDir, opts)
	if err != nil {
		return nil, err
	}

	var tasks []uuid.UUID
	for _, node := range resp.Node.Nodes {
		t, err := minion.EtcdUnmarshalTask(node)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, t.ID)
	}

	return tasks, nil
}

// MinionSubmitTask submits a new task to a given minion uuid
func (c *etcdMinionClient) MinionSubmitTask(m uuid.UUID, t *task.Task) error {
	rootDir := filepath.Join(minion.EtcdMinionSpace, m.String())
	queueDir := filepath.Join(rootDir, "queue")

	// Check if minion exists first
	_, err := c.kapi.Get(context.Background(), rootDir, nil)
	if err != nil {
		return err
	}

	// Serialize task and submit it to the minion
	data, err := json.Marshal(t)
	if err != nil {
		return err
	}

	_, err = c.kapi.CreateInOrder(context.Background(), queueDir, string(data), nil)

	return err
}
