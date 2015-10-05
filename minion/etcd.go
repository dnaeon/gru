package minion

import (
	"os"
	"os/signal"
	"log"
	"time"
	"strings"
	"strconv"
	"path/filepath"
	"encoding/json"

	"code.google.com/p/go-uuid/uuid"
	"github.com/coreos/etcd/Godeps/_workspace/src/golang.org/x/net/context"
	"github.com/coreos/etcd/client"
)

// Root keyspace in etcd
const EtcdKeySpace = "/gru"

// Etcd Minion
type EtcdMinion struct {
	// Name of this minion
	Name string

	// Minion root node in etcd 
	MinionRootDir string

	// Minion queue node in etcd
	QueueDir string

	// Root node for results in etcd
	//ResultRootDir string

	// Root node for classifiers in etcd
	ClassifierDir string

	// Minion unique identifier
	UUID uuid.UUID

	// KeysAPI client to etcd
	KAPI etcdclient.KeysAPI
}

// Etcd minion task
type EtcdTask struct {
	// Command to be executed by the minion
	Command string

	// Time when the command was sent for processing
	Timestamp int64

	// Task unique identifier
	UUID uuid.UUID
}

func NewEtcdTask(command string, timestamp int64) MinionTask {
	t := &EtcdTask{
		Command: command,
		Timestamp: timestamp,
		UUID: uuid.NewRandom(),
	}

	return t
}

// Gets the task unique identifier
func (t *EtcdTask) GetUUID() uuid.UUID {
	return t.UUID
}

// Gets the task command to be executed
func (t *EtcdTask) GetCommand() (string, error) {
	return t.Command, nil
}

// Returns the time a task has been received for processing
func (t *EtcdTask) GetTimestamp() (int64, error) {
	return t.Timestamp, nil
}

// Processes a task
func (t *EtcdTask) Process() error {
	id := t.GetUUID()
	log.Printf("Processing task %s\n", id.String())

	// TODO: Actually process the command

	return nil
}

// Create a new minion
func NewEtcdMinion(name string, kapi etcdclient.KeysAPI) Minion {
	minionUUID := uuid.NewSHA1(uuid.NameSpace_DNS, []byte(name))
	minionRootDir := filepath.Join(EtcdKeySpace, "minion", minionUUID.String())
	queueDir := filepath.Join(minionRootDir, "queue")
//	resultDir := filepath.Join(...)
	classifierDir := filepath.Join(minionRootDir, "classifier")

	log.Printf("Created minion with uuid %s\n", minionUUID)

	m := &EtcdMinion{
		Name: name,
		MinionRootDir: minionRootDir,
		QueueDir: queueDir,
//		ResultRootDir: resultRootDir,
		ClassifierDir: classifierDir,
		UUID: minionUUID,
		KAPI: kapi,
	}

	return m
}

// Get the minion UUID
func (m *EtcdMinion) GetUUID() uuid.UUID {
	return m.UUID
}

// Get the human-readable name of the minion
func (m *EtcdMinion) GetName() (string, error) {
	return m.Name, nil
}

// Set the human-readable name of the minion
func (m *EtcdMinion) SetName(name string) error {
	key := filepath.Join(m.MinionRootDir, "name")
	opts := &etcdclient.SetOptions{
		PrevExist: etcdclient.PrevIgnore,
	}

	_, err := m.KAPI.Set(context.Background(), key, m.Name, opts)

	return err
}

// Set the time the minion was last seen in seconds since the Epoch
func (m *EtcdMinion) SetLastseen(s int64) error {
	key := filepath.Join(m.MinionRootDir, "lastseen")
	lastseen := strconv.FormatInt(s, 10)
	opts := &etcdclient.SetOptions{
		PrevExist: etcdclient.PrevIgnore,
	}

	_, err := m.KAPI.Set(context.Background(), key, lastseen, opts)

	return err
}

// Get classifier for a minion
func (m *EtcdMinion) GetClassifier(key string) (MinionClassifier, error) {
	klassifier := new(SimpleClassifier)
	klassifierNode := filepath.Join(m.ClassifierDir, key, "info")

	// Get classifier from etcd and deserialize
	resp, err := m.KAPI.Get(context.Background(), klassifierNode, nil)

	if err != nil {
		return klassifier, err
	}

	err = json.Unmarshal([]byte(resp.Node.Value), &klassifier)

	return klassifier, err
}

// Classify a minion  a given key and value
func (m *EtcdMinion) SetClassifier(c MinionClassifier) error {
	// Classifiers in etcd expire after an hour
	opts := &etcdclient.SetOptions{
		PrevExist: etcdclient.PrevIgnore,
		TTL: time.Hour,
	}

	// Get classifier values
	key, err := c.GetKey()
	description, err := c.GetDescription()
	value, err := c.GetValue(m)

	if err != nil {
		return err
	}

	// Create a simple classifier and serialize to JSON
	klassifier := NewSimpleClassifier(key, value, description)
	data, err := json.Marshal(klassifier)

	if err != nil {
		log.Printf("Failed to serialize classifier: %s\n", key)
		return err
	}

	// Set classifier in etcd
	klassifierNode := filepath.Join(m.ClassifierDir, key, "info")
	_, err = m.KAPI.Set(context.Background(), klassifierNode, string(data), opts)

	return err
}

// Runs periodic jobs such as refreshing classifiers and updating lastseen
func (m *EtcdMinion) Refresh(ticker *time.Ticker) error {
	for {
		// Update classifiers
		for _, classifier := range ClassifierRegistry {
			m.SetClassifier(classifier)
		}

		// Update lastseen time
		now := time.Now().Unix()
		err := m.SetLastseen(now)
		if err != nil {
			log.Printf("Failed to update lastseen time: %s\n", err)
		}

		<- ticker.C
	}

	return nil
}

// Unmarshals task from etcd and removes it from the queue
func UnmarshalEtcdTask(node *etcdclient.Node) (*EtcdTask, error) {
	task := new(EtcdTask)
	err := json.Unmarshal([]byte(node.Value), &task)

	if err != nil {
		log.Printf("Invalid task: key: %s\n", node.Key)
		log.Printf("Invalid task: value: %s\n", node.Value)
		log.Printf("Invalid task: error: %s\n", err)
	}

	return task, err
}

// Monitors etcd for new tasks for processing
func (m *EtcdMinion) TaskListener(c chan<- MinionTask) error {
	log.Printf("Starting task listener")

	watcherOpts := &etcdclient.WatcherOptions{
		Recursive: true,
	}
	watcher := m.KAPI.Watcher(m.QueueDir, watcherOpts)

	for {
		resp, err := watcher.Next(context.Background())
		if err != nil {
			log.Printf("Failed to read task: %s\n", err)
			continue
		}

		// Ignore "delete" events when removing a task from the queue
		action := strings.ToLower(resp.Action)
		if strings.EqualFold(action, "delete") {
			continue
		}

		task, err := UnmarshalEtcdTask(resp.Node)
		m.KAPI.Delete(context.Background(), resp.Node.Key, nil)

		if err != nil {
			continue
		}

		c <- task
	}

	return nil
}

// Processes new tasks
func (m *EtcdMinion) TaskRunner(c <-chan MinionTask) error {
	for {
		task := <-c
		task.Process()
	}

	return nil
}

// Checks for any tasks in backlog
func (m *EtcdMinion) CheckForBacklog(c chan<- MinionTask) error {
	opts := &etcdclient.GetOptions{
		Recursive: true,
		Sort: true,
	}

	// Get backlog tasks if any
	resp, err := m.KAPI.Get(context.Background(), m.QueueDir, opts)
	if err != nil {
		return nil
	}

	backlog := resp.Node.Nodes

	if len(backlog) == 0 {
		// No backlog tasks found
		return nil
	}

	log.Printf("Found %d tasks in backlog", len(backlog))
	for _, node := range backlog {
		task, err := UnmarshalEtcdTask(node)
		m.KAPI.Delete(context.Background(), node.Key, nil)

		if err != nil {
			continue
		}

		c <- task
	}

	return nil
}

// Main entry point of the minion
func (m *EtcdMinion) Serve() error {
	// Channel on which we send the quit signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	// Run any periodic tasks every hour
	ticker := time.NewTicker(time.Minute * 15)
	go m.Refresh(ticker)

	// Check for backlog tasks and start task listener and runner
	tasks := make(chan MinionTask)
	go m.TaskListener(tasks)
	go m.CheckForBacklog(tasks)
	go m.TaskRunner(tasks)

	// Block until a stop signal is received
	s := <-quit
	log.Printf("Received %s signal, shutting down", s)
	close(quit)
	close(tasks)

	return nil
}

